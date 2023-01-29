package titan_client

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"io"
	"os"
	"testing"
)

func TestNewDownloader(t *testing.T) {
	//downloader := NewDownloader(WithCustomGatewayAddressOption("http://127.0.0.1:5001"), WithLocatorAddressOption(""))
	// t.Log(downloader)
	f, err := os.Create("./data.txt")
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()

	for i := 0; i < 5100; i++ {
		_, err = f.WriteString(fmt.Sprintf("%d %s\n", i, "QmTp2hEo8eXRp6wg7jXv1BLCMh5a4F3B7buAUZNZUu772j"))
		if err != nil {
			t.Error(err.Error())
			return
		}
	}
}

func TestTitanDownloader_Download(t *testing.T) {
	err := logging.SetLogLevel("titan-client/util", "DEBUG")
	if err != nil {
		t.Error(err.Error())
		return
	}
	ctx := context.Background()
	c, err := cid.Decode("QmXRrLjxgHd2Ls8jFZby2fx2wQuuqBkamQE8ibY6TnREA4")
	if err != nil {
		t.Error(err)
		return
	}
	downloader := NewDownloader(WithCustomGatewayAddressOption("http://127.0.0.1:5001"), WithLocatorAddressOption("http://192.168.0.132:5000"))
	err = downloader.Download(ctx, c, false, gzip.NoCompression, "./titan.txt")
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log("download success")
}

func TestTitanDownloader_GetReader(t *testing.T) {
	/*	err := os.Remove("./download.txt")
		if err != nil {
			t.Error(err)
			return
		}*/
	c, err := cid.Decode("QmYLniRF9EL5CCV5hY5z9KEYqnRdhvjpCGZgvNxL8JCc2E")
	if err != nil {
		t.Error(err)
		return
	}
	reader, err := NewDownloader(WithLocatorAddressOption("http://192.168.0.132:5000"), WithCustomGatewayAddressOption("http://127.0.0.1:5001")).GetReader(context.Background(), c, false, gzip.NoCompression)
	if err != nil {
		t.Error(err)
		return
	}
	defer reader.Close()
	f, err := os.Create("./download.txt")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = io.Copy(f, reader)
	if err != nil {
		t.Error(err.Error())
		return
	}

	/*ow := util.Writer{
		Archive:     false,
		Compression: gzip.NoCompression,
	}

	err = ow.Write(reader, "./download.txt")
	if err != nil {
		t.Error(err)
		return
	}
	*/
	t.Log("download success")
}
