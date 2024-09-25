package network

import (
	"blockchainsystem/core"
	"blockchainsystem/crypto"
	"bytes"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var defaultBlockTime = 5 * time.Second

type Serveropts struct {
	RPCDecodeFunc   RPCDecodeFunc
	RPCProcessor    RPCProcessor
	Transport []Transport
	BlockTime time.Duration
	privKey  *crypto.PrivateKey
}

type Server struct {
	Serveropts
	memPool *TxPool
	isValidator bool
	rpcCh chan RPC
	quitCh chan struct{}
}

func NewServer(opts Serveropts) *Server {

	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}

	
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	
	s:= &Server{
		Serveropts: opts,
		memPool: NewTxPool(),  
		isValidator: opts.privKey != nil,
		rpcCh: make(chan RPC),
		quitCh: make(chan struct{}, 1),
	}

// if we don't have any processor as the server opt, we assume the server default
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}


	return s
}

func (s *Server) Start() {
	s.initTransports()
	ticker :=  time.NewTicker(s.BlockTime)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc) 
			if err != nil {
				logrus.Error(err)
			}

			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				logrus.Error(err)
			}


		case <-s.quitCh:
			break free
		case <- ticker.C:
			if s.isValidator {
				s.createNewBlock()
			}
	}
			
	}

	fmt.Println("Server shutdown")
}


func (s *Server) ProcessMessage(msg *DecodeMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	}

	return nil
}


func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transport {
		if err := tr.Broadcast(payload); err != nil {
			return err
		}
	}

	return nil 
}

func (s *Server) processTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		
		logrus.WithFields(logrus.Fields{
			"hash" : hash,
		}).Info("transaction already in mempool")

		return nil 
		
	} 
	
	if err := tx.Verify(); err != nil {
		return nil
	}

	tx.SetFirstseen(time.Now().UnixNano())



	logrus.WithFields(logrus.Fields{
		"hash" : hash,
		"mempool length": s.memPool.Len(),
	}).Info("adding new tx to the mempool")

	go s.broadcastTx(tx)

	return s.memPool.Add(tx)
}


func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())
}


func (s *Server) createNewBlock() error {
	fmt.Println("creating a new block")
	return nil
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