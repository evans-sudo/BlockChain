package network

import (
	"fmt"
	"time"
)

type Serveropts struct {
	Transport []Transport
}

type Server struct {
	Serveropts
	rpcCh chan RPC
	quitCh chan struct{}
}

func NewServer(opts Serveropts) *Server {
	return &Server{
		Serveropts: opts,
		rpcCh: make(chan RPC),
		quitCh: make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransports()
	ticker :=  time.NewTicker(5 * time.Second)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			fmt.Printf("%+v\n", rpc)
		case <-s.quitCh:
			break free
		case <- ticker.C:
			fmt.Println("do stuff every x seconds")
		}
	}

	fmt.Println("Server shutdown")
}

func (s *Server) initTransports() {
	for _, tr  := range s.Transport {
		go func (tr Transport)  {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc  
			}
		}(tr)
	}
	
}