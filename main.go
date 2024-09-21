package main

import (
	"blockchainsystem/network"
	"time"
)

//"blockchainsystem/network"

// server
// Transport ==> tcp, udp
// block
// Transacation
// key pairs,


func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func ()  {
	for{
		trRemote.SendMessage(network.NetAddr(trLocal.Addr()), []byte("Hello world"))
		time.Sleep(1 * time.Second)
	}	
	}()

	opts := network.Serveropts{
		Transport: []network.Transport{trLocal},
	}

	s := network.NewServer(opts)
	s.Start()
}