package asymmetric

import (
	"crypto/rand"
	"io"
	"log"

	"golang.org/x/crypto/nacl/box"
)

type Nonce = *[24]byte
type AsymmetricBox []byte
type PublicKey = *[32]byte
type PrivateKey = *[32]byte

func generateNonce() Nonce {
	var nonce Nonce
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		log.Panicln("Failed to generate nonce ", err)
	}
	return nonce
}

// Generates a nacl box keypair
func GenerateKeypair() (PublicKey, PrivateKey, error) {
	var PubKey PublicKey
	var PrivKey PrivateKey
	PubKey, PrivKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		log.Println("Failed to generate keypair: ", err)
	}
	return PubKey, PrivKey, err
}

// Returns the output from an Asymmetric box; provided the box, the peer pubkey, and the private key
func Decrypt(encBox AsymmetricBox, peersPublicKey PublicKey, privateKey PrivateKey) ([]byte, bool) {
	var output []byte
	nonce := generateNonce()
	output, boolSuccess := box.Open(output, encBox, nonce, peersPublicKey, privateKey)
	if !boolSuccess {
		log.Println("Failed to open secret box, received: ", boolSuccess)
	}
	return output, boolSuccess
}

func Encrypt(message AsymmetricBox, peersPublicKey PublicKey, privateKey PrivateKey) AsymmetricBox {
	nonce := generateNonce()
	var output []byte
	output = box.Seal(output, message, nonce, peersPublicKey, privateKey)
	return output
}