package core

import (
	"blockchainsystem/crypto"
	"blockchainsystem/types"
	"bytes"
	"encoding/gob"
	"fmt"
	
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
	Transacations []Transaction
	Validator crypto.PublicKey
	Signature *crypto.Signature

	// cached version of the Header hash
	hash types.Hash
}

func NewBlock(h *Header, tx []Transaction) *Block {
	return	&Block{
		Header: h,
		Transacations: tx,
	}
}


func (b *Block) AddTransaction(tx *Transaction) {
	b.Transacations = append(b.Transacations, *tx)
}

func (b *Block) Sign(privkey crypto.PrivateKey) error {
	sig, err :=privkey.Sign(b.Header.Bytes())
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


	for _, tx := range b.Transacations {
		if err := tx.Verify(); err != nil {
			return err
		}
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



