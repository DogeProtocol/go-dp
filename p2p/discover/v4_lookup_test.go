// Copyright 2019 The go-ethereum Authors
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

package discover

import (
	"encoding/hex"
	"fmt"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"net"
	"sort"
	"testing"

	"github.com/DogeProtocol/dp/p2p/discover/v4wire"
	"github.com/DogeProtocol/dp/p2p/enode"
	"github.com/DogeProtocol/dp/p2p/enr"
)

func TestUDPv4_Lookup(t *testing.T) {
	t.Parallel()
	test := newUDPTest(t)

	// Lookup on empty table returns no nodes.
	targetKey, _ := decodePubkey(lookupTestnet.target.PubBytes)
	if results := test.udp.LookupPubkey(targetKey); len(results) > 0 {
		t.Fatalf("lookup on empty table returned %d results: %#v", len(results), results)
	}

	// Seed table with initial node.
	fillTable(test.table, []*node{wrapNode(lookupTestnet.node(256, 0))})

	// Start the lookup.
	resultC := make(chan []*enode.Node, 1)
	go func() {
		resultC <- test.udp.LookupPubkey(targetKey)
		test.close()
	}()

	// Answer lookup packets.
	serveTestnet(test, lookupTestnet)

	// Verify result nodes.
	results := <-resultC
	t.Logf("results:")
	for _, e := range results {
		t.Logf("  ld=%d, %x", enode.LogDist(lookupTestnet.target.id(), e.ID()), e.ID().Bytes())
	}
	if len(results) != bucketSize {
		t.Errorf("wrong number of results: got %d, want %d", len(results), bucketSize)
	}
	checkLookupResults(t, lookupTestnet, results)
}

func TestUDPv4_LookupIterator(t *testing.T) {
	t.Parallel()
	test := newUDPTest(t)
	defer test.close()

	// Seed table with initial nodes.
	bootnodes := make([]*node, len(lookupTestnet.dists[256]))
	for i := range lookupTestnet.dists[256] {
		bootnodes[i] = wrapNode(lookupTestnet.node(256, i))
	}
	fillTable(test.table, bootnodes)
	go serveTestnet(test, lookupTestnet)

	// Create the iterator and collect the nodes it yields.
	iter := test.udp.RandomNodes()
	seen := make(map[enode.ID]*enode.Node)
	for limit := lookupTestnet.len(); iter.Next() && len(seen) < limit; {
		seen[iter.Node().ID()] = iter.Node()
	}
	iter.Close()

	// Check that all nodes in lookupTestnet were seen by the iterator.
	results := make([]*enode.Node, 0, len(seen))
	for _, n := range seen {
		results = append(results, n)
	}
	sortByID(results)
	want := lookupTestnet.nodes()
	if err := checkNodesEqual(results, want); err != nil {
		t.Fatal(err)
	}
}

// TestUDPv4_LookupIteratorClose checks that lookupIterator ends when its Close
// method is called.
func TestUDPv4_LookupIteratorClose(t *testing.T) {
	t.Parallel()
	test := newUDPTest(t)
	defer test.close()

	// Seed table with initial nodes.
	bootnodes := make([]*node, len(lookupTestnet.dists[256]))
	for i := range lookupTestnet.dists[256] {
		bootnodes[i] = wrapNode(lookupTestnet.node(256, i))
	}
	fillTable(test.table, bootnodes)
	go serveTestnet(test, lookupTestnet)

	it := test.udp.RandomNodes()
	if ok := it.Next(); !ok || it.Node() == nil {
		t.Fatalf("iterator didn't return any node")
	}

	it.Close()

	ncalls := 0
	for ; ncalls < 100 && it.Next(); ncalls++ {
		if it.Node() == nil {
			t.Error("iterator returned Node() == nil node after Next() == true")
		}
	}
	t.Logf("iterator returned %d nodes after close", ncalls)
	if it.Next() {
		t.Errorf("Next() == true after close and %d more calls", ncalls)
	}
	if n := it.Node(); n != nil {
		t.Errorf("iterator returned non-nil node after close and %d more calls", ncalls)
	}
}

func serveTestnet(test *udpTest, testnet *preminedTestnet) {
	for done := false; !done; {
		done = test.waitPacketOut(func(p v4wire.Packet, to *net.UDPAddr, hash []byte) {
			n, key := testnet.nodeByAddr(to)
			switch p.(type) {
			case *v4wire.Ping:
				test.packetInFrom(nil, key, to, &v4wire.Pong{Expiration: futureExp, ReplyTok: hash})
			case *v4wire.Findnode:
				dist := enode.LogDist(n.ID(), testnet.target.id())
				nodes := testnet.nodesAtDistance(dist - 1)
				test.packetInFrom(nil, key, to, &v4wire.Neighbors{Expiration: futureExp, Nodes: nodes})
			}
		})
	}
}

// checkLookupResults verifies that the results of a lookup are the closest nodes to
// the testnet's target.
func checkLookupResults(t *testing.T, tn *preminedTestnet, results []*enode.Node) {
	t.Helper()
	t.Logf("results:")
	for _, e := range results {
		t.Logf("  ld=%d, %x", enode.LogDist(tn.target.id(), e.ID()), e.ID().Bytes())
	}
	if hasDuplicates(wrapNodes(results)) {
		t.Errorf("result set contains duplicate entries")
	}
	if !sortedByDistanceTo(tn.target.id(), wrapNodes(results)) {
		t.Errorf("result set not sorted by distance to target")
	}
	wantNodes := tn.closest(len(results))
	if err := checkNodesEqual(results, wantNodes); err != nil {
		t.Error(err)
	}
}

// This is the test network for the Lookup test.
// The nodes were obtained by running lookupTestnet.mine with a random NodeID as target.
var (
	key1, _    = cryptobase.SigAlg.GenerateKey()
	hexkey1, _ = cryptobase.SigAlg.PrivateKeyToHex(key1)

	key2, _    = cryptobase.SigAlg.GenerateKey()
	hexkey2, _ = cryptobase.SigAlg.PrivateKeyToHex(key2)

	key3, _    = cryptobase.SigAlg.GenerateKey()
	hexkey3, _ = cryptobase.SigAlg.PrivateKeyToHex(key3)

	key4, _    = cryptobase.SigAlg.GenerateKey()
	hexkey4, _ = cryptobase.SigAlg.PrivateKeyToHex(key4)

	key5, _    = cryptobase.SigAlg.GenerateKey()
	hexkey5, _ = cryptobase.SigAlg.PrivateKeyToHex(key5)

	key6, _    = cryptobase.SigAlg.GenerateKey()
	hexkey6, _ = cryptobase.SigAlg.PrivateKeyToHex(key6)

	key7, _    = cryptobase.SigAlg.GenerateKey()
	hexkey7, _ = cryptobase.SigAlg.PrivateKeyToHex(key7)

	key8, _    = cryptobase.SigAlg.GenerateKey()
	hexkey8, _ = cryptobase.SigAlg.PrivateKeyToHex(key8)

	key9, _    = cryptobase.SigAlg.GenerateKey()
	hexkey9, _ = cryptobase.SigAlg.PrivateKeyToHex(key9)

	key10, _    = cryptobase.SigAlg.GenerateKey()
	hexkey10, _ = cryptobase.SigAlg.PrivateKeyToHex(key10)
)

var lookupTestnet = &preminedTestnet{
	target: hexEncPubkey(hex.EncodeToString(key1.PriData)),
	dists: [257][]*signaturealgorithm.PrivateKey{
		251: {
			hexEncPrivkey(hexkey1),
			hexEncPrivkey(hexkey2),
			hexEncPrivkey(hexkey3),
			hexEncPrivkey(hexkey9),
			hexEncPrivkey(hexkey10),
		},
		252: {
			hexEncPrivkey(hexkey4),
			hexEncPrivkey(hexkey5),
			hexEncPrivkey(hexkey6),
		},
		253: {
			hexEncPrivkey(hexkey1),
			hexEncPrivkey(hexkey2),
			hexEncPrivkey(hexkey3),
			hexEncPrivkey(hexkey4),
			hexEncPrivkey(hexkey5),
			hexEncPrivkey(hexkey6),
			hexEncPrivkey(hexkey7),
			hexEncPrivkey(hexkey8),
			hexEncPrivkey(hexkey9),
			hexEncPrivkey(hexkey10),
		},
		254: {
			hexEncPrivkey(hexkey1),
			hexEncPrivkey(hexkey2),
			hexEncPrivkey(hexkey3),
			hexEncPrivkey(hexkey4),
			hexEncPrivkey(hexkey5),
			hexEncPrivkey(hexkey6),
			hexEncPrivkey(hexkey7),
			hexEncPrivkey(hexkey8),
		},
		255: {
			hexEncPrivkey(hexkey1),
			hexEncPrivkey(hexkey2),
			hexEncPrivkey(hexkey3),
			hexEncPrivkey(hexkey4),
			hexEncPrivkey(hexkey5),
			hexEncPrivkey(hexkey6),
			hexEncPrivkey(hexkey7),
			hexEncPrivkey(hexkey8),
			hexEncPrivkey(hexkey9),
			hexEncPrivkey(hexkey10),
		},
		256: {
			hexEncPrivkey(hexkey1),
			hexEncPrivkey(hexkey2),
			hexEncPrivkey(hexkey3),
			hexEncPrivkey(hexkey4),
			hexEncPrivkey(hexkey5),
			hexEncPrivkey(hexkey6),
			hexEncPrivkey(hexkey7),
			hexEncPrivkey(hexkey8),
		},
	},
}

type preminedTestnet struct {
	target encPubkey
	dists  [hashBits + 1][]*signaturealgorithm.PrivateKey
}

func (tn *preminedTestnet) len() int {
	n := 0
	for _, keys := range tn.dists {
		n += len(keys)
	}
	return n
}

func (tn *preminedTestnet) nodes() []*enode.Node {
	result := make([]*enode.Node, 0, tn.len())
	for dist, keys := range tn.dists {
		for index := range keys {
			result = append(result, tn.node(dist, index))
		}
	}
	sortByID(result)
	return result
}

func (tn *preminedTestnet) node(dist, index int) *enode.Node {
	key := tn.dists[dist][index]
	rec := new(enr.Record)
	rec.Set(enr.IP{127, byte(dist >> 8), byte(dist), byte(index)})
	rec.Set(enr.UDP(5000))
	enode.SignV4(rec, key)
	n, _ := enode.New(enode.ValidSchemes, rec)
	return n
}

func (tn *preminedTestnet) nodeByAddr(addr *net.UDPAddr) (*enode.Node, *signaturealgorithm.PrivateKey) {
	dist := int(addr.IP[1])<<8 + int(addr.IP[2])
	index := int(addr.IP[3])
	key := tn.dists[dist][index]
	return tn.node(dist, index), key
}

func (tn *preminedTestnet) nodesAtDistance(dist int) []v4wire.Node {
	result := make([]v4wire.Node, len(tn.dists[dist]))
	for i := range result {
		result[i] = nodeToRPC(wrapNode(tn.node(dist, i)))
	}
	return result
}

func (tn *preminedTestnet) neighborsAtDistances(base *enode.Node, distances []uint, elems int) []*enode.Node {
	var result []*enode.Node
	for d := range lookupTestnet.dists {
		for i := range lookupTestnet.dists[d] {
			n := lookupTestnet.node(d, i)
			d := enode.LogDist(base.ID(), n.ID())
			if containsUint(uint(d), distances) {
				result = append(result, n)
				if len(result) >= elems {
					return result
				}
			}
		}
	}
	return result
}

func (tn *preminedTestnet) closest(n int) (nodes []*enode.Node) {
	for d := range tn.dists {
		for i := range tn.dists[d] {
			nodes = append(nodes, tn.node(d, i))
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return enode.DistCmp(tn.target.id(), nodes[i].ID(), nodes[j].ID()) < 0
	})
	return nodes[:n]
}

var _ = (*preminedTestnet).mine // avoid linter warning about mine being dead code.

// mine generates a testnet struct literal with nodes at
// various distances to the network's target.
func (tn *preminedTestnet) mine() {
	// Clear existing slices first (useful when re-mining).
	for i := range tn.dists {
		tn.dists[i] = nil
	}

	targetSha := tn.target.id()
	found, need := 0, 40
	for found < need {
		k := newkey()
		ld := enode.LogDist(targetSha, encodePubkey(&k.PublicKey).id())
		if len(tn.dists[ld]) < 8 {
			tn.dists[ld] = append(tn.dists[ld], k)
			found++
			fmt.Printf("found ID with ld %d (%d/%d)\n", ld, found, need)
		}
	}
	fmt.Printf("&preminedTestnet{\n")
	fmt.Printf("	target: hexEncPubkey(\"%x\"),\n", tn.target.PubBytes)
	fmt.Printf("	dists: [%d][]*oqs.PrivateKey{\n", len(tn.dists))
	for ld, ns := range tn.dists {
		if len(ns) == 0 {
			continue
		}
		fmt.Printf("		%d: {\n", ld)
		for _, key := range ns {
			privKey, err := cryptobase.SigAlg.SerializePrivateKey(key)
			if err != nil {
				panic(err)
			}
			fmt.Printf("			hexEncPrivkey(\"%x\"),\n", privKey)
		}
		fmt.Printf("		},\n")
	}
	fmt.Printf("	},\n")
	fmt.Printf("}\n")
}
