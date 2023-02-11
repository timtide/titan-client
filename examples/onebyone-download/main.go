package main

import (
	"context"
	"fmt"
	"github.com/ipfs/go-cid"
	"github.com/timtide/titan-client/util"
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

func main() {
	// one by one to download
	bg := util.NewFetcher()

	ctx := context.Background()
	for _, v := range keys {
		c, err := cid.Decode(v)
		if err != nil {
			panic(err.Error())
		}
		_, err = bg.GetBlockData(ctx, c)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("download success")
	}

	/// ===

	// batch to download
	ks := make([]cid.Cid, 0, len(keys))
	for _, v := range keys {
		c, err := cid.Decode(v)
		if err != nil {
			panic(err.Error())
		}
		ks = append(ks, c)
	}
	ch := bg.GetBlocksFromTitan(ctx, ks)
	var count int
	for {
		select {
		case b, ok := <-ch:
			if !ok {
				fmt.Println("channel is not ok, and download block is : ", count)
				return
			}
			count++
			fmt.Println("batch download success, with cid :", b.Cid())
		}
	}
}
