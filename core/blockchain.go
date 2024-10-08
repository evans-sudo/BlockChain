package core

import (
	"fmt"
	"sync"
	"github.com/go-kit/log"

)

type Blockchain struct {
	logger log.Logger
	store Storage
	lock  sync.RWMutex
	headers []*Header
	validator Validator
}

func NewBlockchain(l log.Logger, genesis *Block) (*Blockchain, error) {
	bc :=  &Blockchain{
		headers: []*Header{},
		store: NewMemoryStore(),
		logger: l,
	}

	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)  // Capture both values
  
	return bc, err
}



func (bc *Blockchain) Addblock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	for _, tx := range b.Transactions {
		bc.logger.Log("msg", "excecuting code", "len", len(tx.Data), "hash", tx.Hash(&TxHasher{}))
		vm := NewVM(tx.Data)
		if err := vm.Run(); err != nil {
			return err 
		}

		bc.logger.Log("vm result", vm.stack.data[vm.stack.sp])
	}

	bc.addBlockWithoutValidation(b)

	return nil
}

func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given (%d) height too high", height)
	}

	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[height], nil
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return  height <= bc.Height()
}


func (bc *Blockchain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return uint32(len(bc.headers)-1)
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.lock.Unlock()

	bc.logger.Log(
		"msg", "new block",
		"hash", b.Hash(BlockHasher{}),
		"height", b.Height,
		"transactions", len(b.Transactions),
	)

	return bc.store.Put(b)
}

// func (bc *Blockchain) addGenesisBlock(b *Block) {

// }