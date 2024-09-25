package main

import (
	"blockchainsystem/core"
	"blockchainsystem/crypto"
	"blockchainsystem/network"
	"bytes"
	"math/rand"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
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
	//	trRemote.SendMessage(network.NetAddr(trLocal.Addr()), []byte("Hello world"))
		if err := SendTransaction(trRemote, network.NetAddr(trLocal.Addr())); err != nil {
			logrus.Error(err)
		}
		time.Sleep(1 * time.Second)
	}	
	}()

	opts := network.Serveropts{
		Transport: []network.Transport{trLocal},
	}

	s := network.NewServer(opts)
	s.Start()
}


func SendTransaction(tr network.Transport, to network.NetAddr) error {
	privkey := crypto.GeneratePrivateKey()
	data :=  []byte(strconv.FormatInt(int64(rand.Intn(1000000000)),10))
	tx := core.NewTransaction(data)
	tx.Sign(privkey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err 
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())

}