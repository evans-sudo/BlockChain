package core

type Storage interface {
	Put(*Block) error
}

type Memorystore struct {
}

func NewMemoryStore() *Memorystore {
	return &Memorystore{}
}

func (s *Memorystore) Put(b *Block) error {
	return nil
}