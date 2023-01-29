package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"github.com/ipfs/go-cid"
	"testing"
)

var data = "hello world !!!"

func TestRsaSign_Sign(t *testing.T) {
	t.Log("source data :", data)
	msg := []byte(data)
	signer := GetSigner()
	signData, err := signer.Sign(msg)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("sign data :", signData)
	b := signer.VerifySign(msg, signData)
	if !b {
		t.Error("verify sign fail")
	}
	t.Log("verify sign success")
}

func TestVerifySignByPublicKey(t *testing.T) {
	cd, err := cid.Decode("QmUbaDBz6YKn3dVzoKrLDyupMmyWk5am2QSdgfKsU1RN3N")
	if err != nil {
		t.Error(err)
		return
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Error(err)
		return
	}

	publicKey := privateKey.PublicKey

	hash := sha256.New()
	hash.Write([]byte(cd.String()))
	bytes := hash.Sum(nil)
	signData, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, bytes)
	if err != nil {
		t.Error(err)
		return
	}
	// public key to string
	X509PublicKey := x509.MarshalPKCS1PublicKey(&publicKey)
	publicKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: X509PublicKey,
		})

	t.Log(string(publicKeyPem))

	// string to public key
	block, _ := pem.Decode(publicKeyPem)
	pk, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		t.Error(err)
		return
	}

	hashv := sha256.New()
	hashv.Write([]byte(cd.String()))
	vbytes := hashv.Sum(nil)
	err = rsa.VerifyPKCS1v15(pk, crypto.SHA256, vbytes, signData)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}
