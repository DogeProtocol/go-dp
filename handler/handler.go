// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package handler

import (
	"errors"
	"math"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/core"
	"github.com/DogeProtocol/dp/core/forkid"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/eth/downloader"
	"github.com/DogeProtocol/dp/eth/fetcher"
	"github.com/DogeProtocol/dp/eth/protocols/eth"
	"github.com/DogeProtocol/dp/ethdb"
	"github.com/DogeProtocol/dp/event"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/p2p"
	"github.com/DogeProtocol/dp/params"
	"github.com/DogeProtocol/dp/trie"
)

const (
	// txChanSize is the size of channel listening to NewTxsEvent.
	// The number is referenced from the size of tx pool.
	txChanSize = 4096
)

var (
	syncChallengeTimeout = 15 * time.Second // Time allowance for a node to reply to the sync progress challenge
)

// txPool defines the methods needed from a transaction pool implementation to
// support all the operations needed by the Ethereum chain protocols.
type txPool interface {
	// Has returns an indicator whether txpool has a transaction
	// cached with the given hash.
	Has(hash common.Hash) bool

	// Get retrieves the transaction from local txpool with given
	// tx hash.
	Get(hash common.Hash) *types.Transaction

	// AddRemotes should add the given transactions to the pool.
	AddRemotes([]*types.Transaction) []error

	// Pending should return pending transactions.
	// The slice should be modifiable by the caller.
	Pending(enforceTips bool) (map[common.Address]types.Transactions, error)

	// SubscribeNewTxsEvent should return an event subscription of
	// NewTxsEvent and send events to the given channel.
	SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
}

// HandlerConfig is the collection of initialization parameters to create a full
// node network P2PHandler.
type HandlerConfig struct {
	Database               ethdb.Database            // Database for direct sync insertions
	Chain                  *core.BlockChain          // Blockchain to serve data from
	TxPool                 txPool                    // Transaction pool to propagate from
	Network                uint64                    // Network identifier to adfvertise
	Sync                   downloader.SyncMode       // Whether to fast or full sync
	BloomCache             uint64                    // Megabytes to alloc for fast sync bloom
	EventMux               *event.TypeMux            // Legacy event mux, deprecate for `feed`
	Checkpoint             *params.TrustedCheckpoint // Hard coded checkpoint for sync challenges
	Whitelist              map[uint64]common.Hash    // Hard coded whitelist for sync challenged
	ConsensusPacketHandler *ConsensusPacketHandler
	RebroadcastCount       int
}

type ConsensusHandler interface {
	HandleConsensusPacket(packet *eth.ConsensusPacket) error
	HandleRequestConsensusDataPacket(packet *eth.RequestConsensusDataPacket) ([]*eth.ConsensusPacket, error)
	OnPeerConnected(peerId string) error
	OnPeerDisconnected(peerId string) error
}

type ConsensusPacketHandler struct {
	Handler ConsensusHandler
}

type HandlePeerList func(peerList []string) error

type P2PHandler struct {
	networkID  uint64
	forkFilter forkid.Filter // Fork ID filter, constant across the lifetime of the node

	fastSync   uint32 // Flag whether fast sync is enabled (gets disabled if we already have blocks)
	snapSync   uint32 // Flag whether fast sync should operate on top of the snap protocol
	AcceptTxns uint32 // Flag whether we're considered synchronised (enables transaction processing)

	checkpointNumber uint64      // Block number for the sync progress validator to cross reference
	checkpointHash   common.Hash // Block hash for the sync progress validator to cross reference

	database ethdb.Database
	txpool   txPool
	chain    *core.BlockChain
	maxPeers int

	Downloader       *downloader.Downloader
	stateBloom       *trie.SyncBloom
	blockFetcher     *fetcher.BlockFetcher
	txFetcher        *fetcher.TxFetcher
	consensusHandler *ConsensusPacketHandler
	peers            *peerSet

	eventMux      *event.TypeMux
	txsCh         chan core.NewTxsEvent
	txsSub        event.Subscription
	minedBlockSub *event.TypeMuxSubscription

	whitelist map[uint64]common.Hash

	// channels for fetcher, syncer, txsyncLoop
	txsyncCh chan *txsync
	quitSync chan struct{}

	chainSync *chainSyncer
	wg        sync.WaitGroup
	peerWG    sync.WaitGroup

	handlePeerListFn HandlePeerList

	rebroadcastCount int

	rebroadcastMap map[common.Hash]int64 //packetHash to unixnano of packet first-time-received-time hash

	rebroadcastLock            sync.Mutex
	rebroadcastLastCleanupTime time.Time
}

var lock = &sync.Mutex{}
var p2phandler *P2PHandler

func GetPacketHandler() *P2PHandler {
	lock.Lock()
	defer lock.Unlock()
	return p2phandler
}

func (h *P2PHandler) SetConsensusHandler(handler ConsensusHandler) {
	lock.Lock()
	defer lock.Unlock()
	consensusHandler := &ConsensusPacketHandler{
		Handler: handler,
	}
	h.consensusHandler = consensusHandler
}

func (h *P2PHandler) SetPeerHandler(handleFn HandlePeerList) {
	lock.Lock()
	defer lock.Unlock()
	h.handlePeerListFn = handleFn
}

// NewHandler returns a P2PHandler for all Ethereum chain management protocol.
func NewHandler(config *HandlerConfig) (*P2PHandler, error) {
	lock.Lock()
	defer lock.Unlock()
	if p2phandler != nil {
		return p2phandler, nil
	}

	// Create the protocol manager with the base fields
	if config.EventMux == nil {
		config.EventMux = new(event.TypeMux) // Nicety initialization for tests
	}
	h := &P2PHandler{
		networkID:                  config.Network,
		forkFilter:                 forkid.NewFilter(config.Chain),
		eventMux:                   config.EventMux,
		database:                   config.Database,
		txpool:                     config.TxPool,
		chain:                      config.Chain,
		peers:                      newPeerSet(),
		whitelist:                  config.Whitelist,
		txsyncCh:                   make(chan *txsync),
		quitSync:                   make(chan struct{}),
		rebroadcastCount:           config.RebroadcastCount,
		rebroadcastMap:             make(map[common.Hash]int64),
		rebroadcastLastCleanupTime: time.Now(),
	}
	if config.Sync == downloader.FullSync {
		// The database seems empty as the current block is the genesis. Yet the fast
		// block is ahead, so fast sync was enabled for this node at a certain point.
		// The scenarios where this can happen is
		// * if the user manually (or via a bad block) rolled back a fast sync node
		//   below the sync point.
		// * the last fast sync is not finished while user specifies a full sync this
		//   time. But we don't have any recent state for full sync.
		// In these cases however it's safe to reenable fast sync.
		fullBlock, fastBlock := h.chain.CurrentBlock(), h.chain.CurrentFastBlock()
		if fullBlock.NumberU64() == 0 && fastBlock.NumberU64() > 0 {
			h.fastSync = uint32(1)
			log.Warn("Switch sync mode from full sync to fast sync")
		}
	} else {
		if h.chain.CurrentBlock().NumberU64() > 0 {
			// Print warning log if database is not empty to run fast sync.
			log.Warn("Switch sync mode from fast sync to full sync")
		} else {
			// If fast sync was requested and our database is empty, grant it
			h.fastSync = uint32(1)
		}
	}
	// If we have trusted checkpoints, enforce them on the chain
	if config.Checkpoint != nil {
		h.checkpointNumber = (config.Checkpoint.SectionIndex+1)*params.CHTFrequency - 1
		h.checkpointHash = config.Checkpoint.SectionHead
	}
	// Construct the Downloader (long sync) and its backing state bloom if fast
	// sync is requested. The Downloader is responsible for deallocating the state
	// bloom when it's done.
	// Note: we don't enable it if snap-sync is performed, since it's very heavy
	// and the heal-portion of the snap sync is much lighter than fast. What we particularly
	// want to avoid, is a 90%-finished (but restarted) snap-sync to begin
	// indexing the entire trie
	if atomic.LoadUint32(&h.fastSync) == 1 && atomic.LoadUint32(&h.snapSync) == 0 {
		h.stateBloom = trie.NewSyncBloom(config.BloomCache, config.Database)
	}
	heighter := func() uint64 {
		return h.chain.CurrentBlock().NumberU64()
	}

	h.Downloader = downloader.New(h.checkpointNumber, config.Database, h.stateBloom, h.eventMux, h.chain, nil, h.removePeer)
	h.Downloader.SetChainHeighter(heighter)

	// Construct the fetcher (short sync)
	validator := func(header *types.Header) error {
		return h.chain.Engine().VerifyHeader(h.chain, header, true)
	}

	inserter := func(blocks types.Blocks) (int, error) {
		// If sync hasn't reached the checkpoint yet, deny importing weird blocks.
		//
		// Ideally we would also compare the head block's timestamp and similarly reject
		// the propagated block if the head is too old. Unfortunately there is a corner
		// case when starting new networks, where the genesis might be ancient (0 unix)
		// which would prevent full nodes from accepting it.
		if h.chain.CurrentBlock().NumberU64() < h.checkpointNumber {
			log.Warn("Unsynced yet, discarded propagated block", "number", blocks[0].Number(), "hash", blocks[0].Hash())
			return 0, nil
		}
		// If fast sync is running, deny importing weird blocks. This is a problematic
		// clause when starting up a new network, because fast-syncing miners might not
		// accept each others' blocks until a restart. Unfortunately we haven't figured
		// out a way yet where nodes can decide unilaterally whether the network is new
		// or not. This should be fixed if we figure out a solution.
		if atomic.LoadUint32(&h.fastSync) == 1 {
			log.Warn("Fast syncing, discarded propagated block", "number", blocks[0].Number(), "hash", blocks[0].Hash())
			return 0, nil
		}
		n, err := h.chain.InsertChain(blocks)
		if err == nil {
			atomic.StoreUint32(&h.AcceptTxns, 1) // Mark initial sync done on any fetcher import
		}
		return n, err
	}
	h.blockFetcher = fetcher.NewBlockFetcher(false, nil, h.chain.GetBlockByHash, validator, h.BroadcastBlock, heighter, nil, inserter, h.removePeer)

	fetchTx := func(peer string, hashes []common.Hash) error {
		p := h.peers.peer(peer)
		if p == nil {
			return errors.New("unknown peer")
		}
		return p.RequestTxs(hashes)
	}
	h.txFetcher = fetcher.NewTxFetcher(h.txpool.Has, h.txpool.AddRemotes, fetchTx)
	h.chainSync = newChainSyncer(h)
	p2phandler = h
	return h, nil
}

// runEthPeer registers an eth peer into the joint eth/snap peerset, adds it to
// various subsystems and starts handling messages.
func (h *P2PHandler) runEthPeer(peer *eth.Peer, handler eth.Handler) error {
	// TODO(karalabe): Not sure why this is needed
	if !h.chainSync.handlePeerEvent(peer) {
		return p2p.DiscQuitting
	}
	h.peerWG.Add(1)
	defer h.peerWG.Done()

	// Execute the Ethereum handshake
	var (
		genesis = h.chain.Genesis()
		head    = h.chain.CurrentHeader()
		hash    = head.Hash()
		number  = head.Number.Uint64()
		td      = h.chain.GetTd(hash, number)
	)
	forkID := forkid.NewID(h.chain.Config(), h.chain.Genesis().Hash(), h.chain.CurrentHeader().Number.Uint64())
	if err := peer.Handshake(h.networkID, td, hash, genesis.Hash(), forkID, h.forkFilter); err != nil {
		peer.Log().Debug("Ethereum handshake failed", "err", err)
		return err
	}
	reject := false // reserved peer slots

	// Ignore maxPeers if this is a trusted peer
	if !peer.Peer.Info().Network.Trusted {
		if reject || h.peers.len() >= h.maxPeers {
			return p2p.DiscTooManyPeers
		}
	}
	peer.Log().Debug("Ethereum peer connected", "name", peer.Name())

	// Register the peer locally
	if err := h.peers.registerPeer(peer); err != nil {
		peer.Log().Error("Ethereum peer registration failed", "err", err)
		return err
	}
	err := h.consensusHandler.Handler.OnPeerConnected(peer.ID())
	if err != nil {
		log.Debug("OnPeerConnected", "error", err)
	}
	defer h.unregisterPeer(peer.ID())

	p := h.peers.peer(peer.ID())
	if p == nil {
		return errors.New("peer dropped during handling")
	}
	// Register the peer in the Downloader. If the Downloader considers it banned, we disconnect
	if err := h.Downloader.RegisterPeer(peer.ID(), peer.Version(), peer); err != nil {
		peer.Log().Error("Failed to register peer in eth syncer", "err", err)
		return err
	}
	h.chainSync.handlePeerEvent(peer)

	// Propagate existing transactions. new transactions appearing
	// after this will be sent via broadcasts.
	h.syncTransactions(peer)

	// If we have a trusted CHT, reject all peers below that (avoid fast sync eclipse)
	if h.checkpointHash != (common.Hash{}) {
		// Request the peer's checkpoint header for chain height/weight validation
		if err := peer.RequestHeadersByNumber(h.checkpointNumber, 1, 0, false); err != nil {
			return err
		}
		// Start a timer to disconnect if the peer doesn't reply in time
		p.syncDrop = time.AfterFunc(syncChallengeTimeout, func() {
			peer.Log().Warn("Checkpoint challenge timed out, dropping", "addr", peer.RemoteAddr(), "type", peer.Name())
			h.removePeer(peer.ID())
		})
		// Make sure it's cleaned up if the peer dies off
		defer func() {
			if p.syncDrop != nil {
				p.syncDrop.Stop()
				p.syncDrop = nil
			}
		}()
	}
	// If we have any explicit whitelist block hashes, request them
	for number := range h.whitelist {
		if err := peer.RequestHeadersByNumber(number, 1, 0, false); err != nil {
			return err
		}
	}
	// Handle incoming messages until the connection is torn down
	return handler(peer)
}

// removePeer requests disconnection of a peer.
func (h *P2PHandler) removePeer(id string) {
	peer := h.peers.peer(id)
	if peer != nil {
		peer.Peer.Disconnect(p2p.DiscUselessPeer)
	}
}

// unregisterPeer removes a peer from the Downloader, fetchers and main peer set.
func (h *P2PHandler) unregisterPeer(id string) {
	// Create a custom logger to avoid printing the entire id
	var logger log.Logger
	if len(id) < 16 {
		// Tests use short IDs, don't choke on them
		logger = log.New("peer", id)
	} else {
		logger = log.New("peer", id[:8])
	}
	// Abort if the peer does not exist
	peer := h.peers.peer(id)
	if peer == nil {
		logger.Error("Ethereum peer removal failed", "err", errPeerNotRegistered)
		return
	}
	// Remove the `eth` peer if it exists
	logger.Debug("Removing Ethereum peer")

	h.Downloader.UnregisterPeer(id)
	h.txFetcher.Drop(id)

	if err := h.peers.unregisterPeer(id); err != nil {
		logger.Error("Ethereum peer removal failed", "err", err)
	}

	err := h.consensusHandler.Handler.OnPeerDisconnected(id)
	if err != nil {
		log.Debug("OnPeerDisconnected", "error", err)
	}
}

func (h *P2PHandler) Start(maxPeers int) {
	h.maxPeers = maxPeers

	// broadcast transactions
	h.wg.Add(1)
	h.txsCh = make(chan core.NewTxsEvent, txChanSize)
	h.txsSub = h.txpool.SubscribeNewTxsEvent(h.txsCh)
	go h.txBroadcastLoop()

	// broadcast mined blocks
	h.wg.Add(1)
	h.minedBlockSub = h.eventMux.Subscribe(core.NewMinedBlockEvent{})
	go h.minedBroadcastLoop()

	// start sync handlers
	h.wg.Add(2)
	go h.chainSync.loop()
	go h.txsyncLoop64() // TODO(karalabe): Legacy initial tx echange, drop with eth/64.
}

func (h *P2PHandler) Stop() {
	h.txsSub.Unsubscribe()        // quits txBroadcastLoop
	h.minedBlockSub.Unsubscribe() // quits blockBroadcastLoop

	// Quit chainSync and txsync64.
	// After this is done, no new peers will be accepted.
	close(h.quitSync)
	h.wg.Wait()

	// Disconnect existing sessions.
	// This also closes the gate for any new registrations on the peer set.
	// sessions which are already established but not added to h.peers yet
	// will exit when they try to register.
	h.peers.close()
	h.peerWG.Wait()

	log.Info("Ethereum protocol stopped")
}

// BroadcastBlock will either propagate a block to a subset of its peers, or
// will only announce its availability (depending what's requested).
func (h *P2PHandler) BroadcastBlock(block *types.Block, propagate bool) {
	hash := block.Hash()
	peers := h.peers.peersWithoutBlock(hash)

	// If propagation is requested, send to a subset of the peer
	if propagate {
		// Calculate the TD of the block (it's not imported yet, so block.Td is not valid)
		var td *big.Int
		if parent := h.chain.GetBlock(block.ParentHash(), block.NumberU64()-1); parent != nil {
			td = new(big.Int).Add(block.Difficulty(), h.chain.GetTd(block.ParentHash(), block.NumberU64()-1))
		} else {
			log.Error("Propagating dangling block", "number", block.Number(), "hash", hash)
			return
		}
		// Send the block to a subset of our peers
		transfer := peers[:int(math.Sqrt(float64(len(peers))))]
		for _, peer := range transfer {
			peer.AsyncSendNewBlock(block, td)
		}
		log.Trace("Propagated block", "hash", hash, "recipients", len(transfer), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))
		return
	}
	// Otherwise if the block is indeed in out own chain, announce it
	if h.chain.HasBlock(hash, block.NumberU64()) {
		for _, peer := range peers {
			peer.AsyncSendNewBlockHash(block)
		}
		log.Trace("Announced block", "hash", hash, "recipients", len(peers), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))
	}
}

// BroadcastTransactions will propagate a batch of transactions
// - To a square root of all peers
// - And, separately, as announcements to all peers which are not known to
// already have the given transaction.
func (h *P2PHandler) BroadcastTransactions(txs types.Transactions) {
	var (
		annoCount   int // Count of announcements made
		annoPeers   int
		directCount int // Count of the txs sent directly to peers
		directPeers int // Count of the peers that were sent transactions directly

		txset = make(map[*ethPeer][]common.Hash) // Set peer->hash to transfer directly
		annos = make(map[*ethPeer][]common.Hash) // Set peer->hash to announce

	)
	// Broadcast transactions to a batch of peers not knowing about it
	for _, tx := range txs {
		peers := h.peers.peersWithoutTransaction(tx.Hash())
		// Send the tx unconditionally to a subset of our peers
		numDirect := int(math.Sqrt(float64(len(peers))))
		for _, peer := range peers[:numDirect] {
			txset[peer] = append(txset[peer], tx.Hash())
		}
		// For the remaining peers, send announcement only
		for _, peer := range peers[numDirect:] {
			annos[peer] = append(annos[peer], tx.Hash())
		}
	}
	for peer, hashes := range txset {
		directPeers++
		directCount += len(hashes)
		peer.AsyncSendTransactions(hashes)
	}
	for peer, hashes := range annos {
		annoPeers++
		annoCount += len(hashes)
		peer.AsyncSendPooledTransactionHashes(hashes)
	}
	log.Debug("Transaction broadcast", "txs", len(txs),
		"announce packs", annoPeers, "announced hashes", annoCount,
		"tx packs", directPeers, "broadcast txs", directCount)
}

func (h *P2PHandler) RequestTransactions(txns []common.Hash) error {
	peers := h.peers.allPeers()

	numDirect := int(math.Sqrt(float64(len(peers))))
	for _, peer := range peers[:numDirect] {
		peer.RequestTxs(txns)
	}

	return nil
}

func (h *P2PHandler) RequestConsensusData(packet *eth.RequestConsensusDataPacket) error {
	peers := h.peers.allPeers()

	numDirect := int(math.Sqrt(float64(len(peers))))
	for _, peer := range peers[:numDirect] {
		peer.AsyncSendRequestConsensusDataPacket(packet)
	}

	return nil
}

func (h *P2PHandler) BroadcastConsensusData(packet *eth.ConsensusPacket) error {
	peers := h.peers.allPeers()
	for _, peer := range peers {
		_, err := peer.Node().Address()
		if err != nil {
			log.Trace("BroadcastConsensusData", "err", err, "peer", peer.ID())
		}
		peer.AsyncSendConsensusPacket(packet)
	}

	return nil
}

func (h *P2PHandler) RequestPeerList() error {
	packet := &eth.RequestPeerListPacket{
		MaxPeers: 10,
	}
	peers := h.peers.allPeers()
	for _, peer := range peers {
		_, err := peer.Node().Address()
		if err != nil {
			log.Trace("RequestPeerList", "err", err, "peer", peer.ID())
		}
		peer.AsyncSendRequestPeerListPacket(packet)
	}
	return nil
}

// minedBroadcastLoop sends mined blocks to connected peers.
func (h *P2PHandler) minedBroadcastLoop() {
	defer h.wg.Done()

	for obj := range h.minedBlockSub.Chan() {
		if ev, ok := obj.Data.(core.NewMinedBlockEvent); ok {
			h.BroadcastBlock(ev.Block, true)  // First propagate block to peers
			h.BroadcastBlock(ev.Block, false) // Only then announce to the rest
		}
	}
}

// txBroadcastLoop announces new transactions to connected peers.
func (h *P2PHandler) txBroadcastLoop() {
	defer h.wg.Done()
	for {
		select {
		case event := <-h.txsCh:
			h.BroadcastTransactions(event.Txs)
		case <-h.txsSub.Err():
			return
		}
	}
}
