package network

import (
	"blockchainsystem/core"
	"blockchainsystem/crypto"
	"blockchainsystem/types"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/log"
)

var defaultBlockTime = 5 * time.Second

type Serveropts struct {
	ID string
	Logger			 log.Logger
	RPCDecodeFunc   RPCDecodeFunc
	RPCProcessor    RPCProcessor
	Transport []Transport
	BlockTime time.Duration
	PrivateKey  *crypto.PrivateKey
}

type Server struct {
	Serveropts
	memPool *TxPool
	chain  *core.Blockchain
	isValidator bool
	rpcCh chan RPC
	quitCh chan struct{}
}

func NewServer(opts Serveropts) (*Server, error) {

	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}

	
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}

	chain, err := core.NewBlockchain(opts.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}

	
	s:= &Server{
		Serveropts: opts,
		memPool: NewTxPool(),
		chain: chain,  
		isValidator: opts.PrivateKey != nil,
		rpcCh: make(chan RPC),
		quitCh: make(chan struct{}, 1),
	}

// if we don't have any processor as the server opt, we assume the server default
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	
	if s.isValidator {
		go s.validatorLoop()
	}


	for _, tr := range s.Transport {
		if err := s.sendGetStatusMessage(tr); err != nil {
			s.Logger.Log("send get status error", err)
		}
	}


	return s, nil 

}

func (s *Server) Start() {
	s.initTransports()		

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc) 
			if err != nil {
				s.Logger.Log("error", err)
			}

			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				s.Logger.Log("error", err)
			}


		case <-s.quitCh:
			break free
	}
			
	}

	s.Logger.Log("msg", "Server is shutting down")
	
}


func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.BlockTime)

	s.Logger.Log("msg", "Starting validator loop","blocktime", s.BlockTime)

	for {
		<- ticker.C
		s.createNewBlock()
	}
}


func (s *Server) ProcessMessage(msg *DecodeMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	case *core.Block:
		return s.processBlock(t)

	case *GetStatusMessage:
		return s.processGetStatusMessage(msg.From, t)
	}

	return nil
}


func (s *Server) sendGetStatusMessage(tr Transport) error {
	var (
	getStatusMsg = new(GetStatusMessage)
	buf = new(bytes.Buffer)
	)

	if err := gob.NewEncoder(buf).Encode(getStatusMsg); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeGetStatus, buf.Bytes())


	if err := tr.SendMessage(tr.Addr(), msg.Bytes()); err != nil {

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

func (s *Server) processGetStatusMessage(from NetAddr, msg *GetStatusMessage) error {
	fmt.Printf("received status msg from %s => %+v\n", from, msg)
	
	return nil 
}

func (s *Server) processBlock(b *core.Block) error {
	if err := s.chain.Addblock(b); err != nil {
		return err
	}

	go s.broadcastBlock(b)

	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		return nil 
	} 
	
	if err := tx.Verify(); err != nil {
		return nil
	}

		// s.Logger.Log(
		// "msg", "adding new tx to mempool", 
		// "hash", hash,
		//  "mempoollength", s.memPool.Len(),
		// )

	go s.broadcastTx(tx)

	 s.memPool.Add(tx)

	 return nil 
}

func (s *Server) broadcastBlock(b *core.Block) error {
	buf := &bytes.Buffer{}
	if err := b.Encode(core.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeBlock, buf.Bytes())
	
	return s.broadcast(msg.Bytes())
}


func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())
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


func (s *Server) createNewBlock() error {
	currentHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}

	txx := s.memPool.Transaction()

	block, err := core.NewBlockFromPrevHeader(currentHeader, txx)
	if err != nil {
		return err
	}

	if err := block.Sign(*s.PrivateKey); err != nil {
		return err
	}

	if err := s.chain.Addblock(block); err != nil {
		return err
	}

	s.memPool.Flush()


	return nil
}

func genesisBlock() *core.Block {
	Header := &core.Header{
		Version: 1,
		DataHash: types.Hash{},
		Height: 0,
		Timestamp: 000000,
	}


	b, _ := core.NewBlock(Header, nil)

	return b 

}