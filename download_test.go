package titan_client

import (
	"compress/gzip"
	"context"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"io"
	"os"
	"testing"
)

func TestNewDownloader(t *testing.T) {
	downloader := NewDownloader(WithCustomGatewayAddressOption("http://127.0.0.1:5001"), WithLocatorAddressOption(""))
	t.Log(downloader)
}

func TestTitanDownloader_Download(t *testing.T) {
	carfiles := []string{
		"QmT2dwc94QJypTuACcdBGLdzJGLz7m1LCvPvH43HrZdTWn",
		"QmXmPUA8CGNDebNZZfhg6MYhfq6hRvRLuMf7sRWDsBUnU9",
		"QmXRrLjxgHd2Ls8jFZby2fx2wQuuqBkamQE8ibY6TnREA4",
		//"QmQztiWG1rvim9d8HgcK34UEJbTkih8wvfxaycpW2Zwccc",
		"QmUbaDBz6YKn3dVzoKrLDyupMmyWk5am2QSdgfKsU1RN3N",
	}
	err := logging.SetLogLevel("titan-client/util", "DEBUG")
	if err != nil {
		t.Error(err.Error())
		return
	}
	ctx := context.Background()
	downloader := NewDownloader(WithCustomGatewayAddressOption("http://127.0.0.1:5001"), WithLocatorAddressOption("http://39.108.143.56:5000"))
	for {
		for _, v := range carfiles {
			t.Logf("================>> start download carfile[%s] <<================", v)
			c, err := cid.Decode(v)
			if err != nil {
				t.Error(err)
				return
			}
			err = downloader.Download(ctx, c, false, gzip.NoCompression, "./titan03.txt")
			if err != nil {
				t.Error(err.Error())
				return
			}
			t.Logf("================>> carfile[%s] download success <<================", v)
		}
	}
}

func TestTitanDownloader_GetReader(t *testing.T) {
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
	f, err := os.Create("./titan02.txt")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = io.Copy(f, reader)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log("download success")
}
