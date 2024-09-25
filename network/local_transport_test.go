package network

import (
	"io/ioutil"
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
	b, err := ioutil.ReadAll(rpc.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)
	assert.Equal(t, rpc.From, NetAddr(tra.Addr())) 

}



func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	trc := NewLocalTransport("C")

	tra.Connect(trb)
	tra.Connect(trc)

	msg := []byte("foo")
	assert.Nil(t, tra.Broadcast(msg))


	rpcb := <- trb.Consume()
	b, err := ioutil.ReadAll(rpcb.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)


	rpcC := <- trc.Consume()
	c, err := ioutil.ReadAll(rpcC.Payload)
	assert.Nil(t, err)
	assert.Equal(t, c, msg)
}


