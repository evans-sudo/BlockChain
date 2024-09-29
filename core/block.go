package core

import (
	"blockchainsystem/crypto"
	"blockchainsystem/types"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)



type Header struct {
	Version uint32
	DataHash types.Hash
	PrevBlockHash types.Hash
	Timestamp uint64
	Height uint32
	
}


func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)
	
	return buf.Bytes()
}



type Block struct {
	*Header
	Transactions []*Transaction
	Validator crypto.PublicKey
	Signature *crypto.Signature

	// cached version of the Header hash
	hash types.Hash
}

func NewBlock(h *Header, tx []*Transaction) (*Block, error) {
	return	&Block{
		Header: h,
		Transactions: tx,
	}, nil
}


func NewBlockFromPrevHeader(prevHeader *Header, txx []*Transaction) (*Block, error) {
	datahash, err := CalculateDataHash(txx)
	if err != nil {
		return nil, err
	}
	header := &Header{
		Version: 1,
		Height: prevHeader.Height + 1,
		DataHash: datahash,
		PrevBlockHash: BlockHasher{}.Hash(prevHeader),
		Timestamp: uint64(time.Now().UnixNano()),
	}

	return NewBlock(header, txx)


}


func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)
}

func (b *Block) Sign(privkey crypto.PrivateKey) error {
	sig, err := privkey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}

	b.Validator = privkey.PublicKey()
	b.Signature = sig

	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has not signature")
	}

	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("block has invalid signature")
	}


	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}

	 dataHash, err := CalculateDataHash(b.Transactions)

	 if err != nil {
		return err
	 }
	 if dataHash != b.DataHash {
		return fmt.Errorf("block (%s) has an invalid data hash", b.Hash(BlockHasher{}))
	 }


	return nil
}

func (b *Block) Decode( dec Decoder[*Block]) error {
	return dec.Decode(b)
}

func (b *Block) Encode( enc Encoder[*Block]) error {
	return enc.Encode(b)
}

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}

	return b.hash
}


func  CalculateDataHash(txx []*Transaction) (hash types.Hash, err error) {
	buf := &bytes.Buffer{}

	for _, tx := range txx {
		if err = tx.Encode(NewGobTxEncoder(buf)); err != nil {
			return 
		}
	}

	hash = sha256.Sum256(buf.Bytes())
	
	return
}



