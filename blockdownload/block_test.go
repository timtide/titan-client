package blockdownload

import (
	"context"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"testing"
)

var keys = []string{
	"QmbRnNGU6pLBDaSRsjWpwHniNV1qMNt1rYEAqtazXTypDj",
	"QmQ7xvThkG7N32imK7Xpud5mWCbk2ZiggPpoYQTL91SdTr",
	"QmbvurDWFtRRHDLehR3fWDhou7zvziAeeaJJb2kakWu2SG",
	"QmPgaP4SiadmrtFzEVY5aGTCRou5vbMDJCgEaJwuN9Lk4H",
	"QmczVAtYCdNH6WVNZCCLV7dsfVVZgWZopuTwbAaUTk3927",
	"QmY3XHXx4RZJn4XiX38F9PhsBfJNcVH9dPEjEAK6iumHm3",
	"QmXqW6eQqj3Ng3To4LRz5nXh7cdPLMfFbvQM7poDqEQgir",
	"QmdrDA5q9mMqJ5NeVTU4JJ83Ao96RR8Y8JWPdYcu25AR59",
	"QmcHcWLuy9ki2i2XALhLryXAHT7AxVLdE7E8v6s7yFfmmq",
}

func TestBlockGetter_GetBlock(t *testing.T) {
	err := logging.SetLogLevel("titan-client/blockdownload", "DEBUG")
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = logging.SetLogLevel("titan-client/util", "DEBUG")
	if err != nil {
		t.Error(err.Error())
		return
	}
	c, err := cid.Decode("QmPgaP4SiadmrtFzEVY5aGTCRou5vbMDJCgEaJwuN9Lk4H")
	if err != nil {
		t.Error(err)
		return
	}
	block, err := NewBlockGetter(WithLocatorAddressOption("http://221.4.187.172:3456")).GetBlock(context.Background(), c)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(block.Cid())
}

func TestBlock_GetBlocks(t *testing.T) {
	err := logging.SetLogLevel("titan-client/blockdownload", "DEBUG")
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = logging.SetLogLevel("titan-client/util", "DEBUG")
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log("keys length : ", len(keys))
	ks := make([]cid.Cid, 0, len(keys))
	for _, v := range keys {
		c, err := cid.Decode(v)
		if err != nil {
			t.Error(err)
			continue
		}
		ks = append(ks, c)
	}
	t.Log("ks length : ", len(ks))
	ch := NewBlockGetter().GetBlocks(context.Background(), ks)
	var count int
	for {
		select {
		case b, ok := <-ch:
			if !ok {
				t.Log("channel is not ok, and download block is : ", count)
				return
			}
			count++
			t.Log(b.Cid())
		}
	}
}
