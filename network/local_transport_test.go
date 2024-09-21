package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {

	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	assert.Equal(t, tra.(*LocalTransport).peers[NetAddr(trb.Addr())], trb)
	assert.Equal(t, trb.(*LocalTransport).peers[NetAddr(tra.Addr())], tra)

}

func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	msg := []byte("Hello world")
	//assert.Nil(t, tra.SendMessage(trb.addr, msg))
	assert.Nil(t, tra.SendMessage(NetAddr(trb.Addr()), msg))


	rpc := <-trb.Consume()
	assert.Equal(t, rpc.Payload, msg)
	assert.Equal(t, rpc.From, NetAddr(tra.Addr())) 

}
