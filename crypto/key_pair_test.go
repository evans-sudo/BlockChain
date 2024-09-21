package crypto

import (
	
	"testing"

	"github.com/stretchr/testify/assert"
)




func TestKeyPairSignVerifySuccess(t *testing.T) {
	privkey := GeneratePrivateKey()
	pubKey := privkey.PublicKey()
	msg := []byte("Hello world")


	sig, err := privkey.Sign(msg)
	assert.Nil(t, err)

	assert.True(t, sig.Verify(pubKey, msg))
}

func TestKeyPairSignVerifyFalse(t *testing.T) {
	privkey := GeneratePrivateKey()
	pubKey := privkey.PublicKey()
	msg := []byte("Hello world")


	sig, err := privkey.Sign(msg)
	assert.Nil(t, err)

	otherPrivKey := GeneratePrivateKey()
	otherPubKey := otherPrivKey.PublicKey()


	assert.False(t, sig.Verify(otherPubKey, msg))
	assert.False(t, sig.Verify(pubKey, []byte("xxxxxxxxx")))

}