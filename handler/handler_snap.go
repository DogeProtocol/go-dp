// Copyright 2020 The go-ethereum Authors
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
	"github.com/DogeProtocol/dp/core"
	"github.com/DogeProtocol/dp/eth/protocols/snap"
	"github.com/DogeProtocol/dp/p2p/enode"
)

// SnapHandler implements the snap.Backend interface to handle the various network
// packets that are sent as replies or broadcasts.
type SnapHandler P2PHandler

func (h *SnapHandler) Chain() *core.BlockChain { return h.chain }

// RunPeer is invoked when a peer joins on the `snap` protocol.
func (h *SnapHandler) RunPeer(peer *snap.Peer, hand snap.Handler) error {
	return (*P2PHandler)(h).runSnapExtension(peer, hand)
}

// PeerInfo retrieves all known `snap` information about a peer.
func (h *SnapHandler) PeerInfo(id enode.ID) interface{} {
	if p := h.peers.peer(id.String()); p != nil {
		if p.snapExt != nil {
			return p.snapExt.info()
		}
	}
	return nil
}

// Handle is invoked from a peer's message P2PHandler when it receives a new remote
// message that the P2PHandler couldn't consume and serve itself.
func (h *SnapHandler) Handle(peer *snap.Peer, packet snap.Packet) error {
	return h.Downloader.DeliverSnapPacket(peer, packet)
}
