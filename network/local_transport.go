package network

import (
	"bytes"
	"fmt"
	"sync"
)


type LocalTransport struct {
	addr NetAddr
	consumeCh chan RPC
	lock sync.RWMutex
	peers map[NetAddr] *LocalTransport
}

func NewLocalTransport(addr NetAddr) Transport {
	return &LocalTransport{
		addr: addr,
		consumeCh: make(chan RPC, 1024),
		peers: make(map[NetAddr]*LocalTransport),
	}
}

func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeCh
}

func (t *LocalTransport) Connect(tr Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.peers[NetAddr(tr.Addr())] = tr.(*LocalTransport)

	return nil
}

func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: could not send message to %s", t.addr, to)
	}

	peer.consumeCh <- RPC {
		From: NetAddr(t.Addr()),
		Payload: bytes.NewReader(payload),
	}

	return nil
}

func (t *LocalTransport) Broadcast(payload []byte) error {
	for _, peer := range t.peers {
		if err := t.SendMessage(NetAddr(peer.Addr()), payload); err != nil {
			return err
		}
	}

	return nil 
}


func (t *LocalTransport) Addr() string {
	return string(t.addr)
}