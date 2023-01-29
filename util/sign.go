package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"sync"
)

var rs *rsaSign
var one sync.Once

func init() {
	one.Do(
		func() {
			rs = new(rsaSign)
			err := rs.generateRSAKey(1024)
			if err != nil {
				panic(err)
			}
		},
	)
}

func GetSigner() Signer {
	return rs
}

type Signer interface {
	Sign(msg []byte) ([]byte, error)
	VerifySign(msg []byte, sign []byte) bool
	GetPublicKey() rsa.PublicKey
}

type rsaSign struct {
	privateKey *rsa.PrivateKey
}

func (r *rsaSign) GetPublicKey() rsa.PublicKey {
	return r.privateKey.PublicKey
}

func (r *rsaSign) Sign(msg []byte) ([]byte, error) {
	hash := sha256.New()
	hash.Write(msg)
	bytes := hash.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.privateKey, crypto.SHA256, bytes)
}

func (r *rsaSign) VerifySign(msg []byte, sign []byte) bool {
	hash := sha256.New()
	hash.Write(msg)
	bytes := hash.Sum(nil)
	err := rsa.VerifyPKCS1v15(&r.privateKey.PublicKey, crypto.SHA256, bytes, sign)
	return err == nil
}

func (r *rsaSign) generateRSAKey(bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	r.privateKey = privateKey
	// publicKey := privateKey.PublicKey
	return nil
}
