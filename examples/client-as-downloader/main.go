package main

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/ipfs/go-cid"
	titan_client "github.com/timtide/titan-client"
)

func main() {
	ctx := context.Background()
	c, err := cid.Decode("QmUbaDBz6YKn3dVzoKrLDyupMmyWk5am2QSdgfKsU1RN3N")
	if err != nil {
		panic(err.Error())
	}
	d := titan_client.NewDownloader(titan_client.WithCustomGatewayAddressOption("http://127.0.0.1:5001"))
	err = d.Download(ctx, c, false, gzip.NoCompression, "./titan.mp4")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("download success")
}
