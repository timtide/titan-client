package util

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"github.com/linguohua/titan/api"
	"github.com/linguohua/titan/api/client"
	http2 "github.com/timtide/titan-client/util/http"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"time"
)

// todo: there is no domain name at present. Use IP first
const defaultLocatorAddress = "http://192.168.0.132:5000/rpc/v0"
const sdkName = "go-titan-client"

var logger = logging.Logger("titan-client/util")

type FetcherOption func(*fetcher)

func WithLocatorAddressOption(locatorUrl string) FetcherOption {
	return func(dg *fetcher) {
		dg.locatorAddr = locatorUrl
	}
}

// Fetcher from titan or common gateway or local gateway to get data
type Fetcher interface {
	GetBlockData(ctx context.Context, c cid.Cid) ([]byte, error)
	GetBlockDataFromTitanOrGateway(ctx context.Context, customGatewayURL string, c cid.Cid) ([]byte, error)
	GetBlocksFromTitanOrGateway(ctx context.Context, customGatewayURL string, ks []cid.Cid) <-chan blocks.Block
	GetBlocksFromTitan(ctx context.Context, ks []cid.Cid) <-chan blocks.Block
}

type fetcher struct {
	// store edge node information of all cid
	pool        []*api.DownloadInfoResult
	locatorAddr string
	// for carfile, root cid load failure record
	err error
}

func NewFetcher(option ...FetcherOption) Fetcher {
	dg := &fetcher{}
	for _, v := range option {
		v(dg)
	}
	if dg.locatorAddr == "" {
		dg.locatorAddr = defaultLocatorAddress
	}
	if !strings.HasSuffix(dg.locatorAddr, "/rpc/v0") {
		dg.locatorAddr = fmt.Sprintf("%s%s", dg.locatorAddr, "/rpc/v0")
	}
	return dg
}

func (d *fetcher) getDownloadInfosByRootCid(ctx context.Context, c cid.Cid) error {
	locator, closer, err := client.NewLocator(ctx, d.locatorAddr, nil)
	if err != nil {
		logger.Error("create schedule fail : ", err.Error())
		return err
	}
	defer closer()
	publicKey := GetSigner().GetPublicKey()
	X509PublicKey := x509.MarshalPKCS1PublicKey(&publicKey)
	publicKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: X509PublicKey,
		})
	downloadInfos, err := locator.GetDownloadInfosWithCarfile(ctx, c.String(), string(publicKeyPem))
	if err != nil {
		logger.Error("get download info fail : ", err.Error())
		return err
	}

	if downloadInfos == nil || len(downloadInfos) == 0 {
		return fmt.Errorf("%s%s", "titan does not cache the carfile cid : ", c.String())
	}

	d.pool = downloadInfos

	return nil
}

func (d *fetcher) GetBlockData(ctx context.Context, c cid.Cid) ([]byte, error) {
	if d.err != nil {
		return nil, d.err
	}
	if d.pool == nil || len(d.pool) == 0 {
		err := d.getDownloadInfosByRootCid(ctx, c)
		if err != nil {
			d.err = err
			return nil, err
		}
	}
	df, err := d.allotDownloadInfo()
	if err != nil {
		return nil, err
	}
	data, err := d.getDataFromEdgeNode(df, c)
	if err != nil {
		logger.Error("fail get data from edge node : ", err.Error())
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil, err
		}
		go d.callback(c, df.SN, false)
		return nil, err
	}
	go d.callback(c, df.SN, true)
	return data, nil
}

func (d *fetcher) allotDownloadInfo() (*api.DownloadInfoResult, error) {
	if len(d.pool) == 1 {
		return d.pool[0], nil
	}
	weightAllot := false
	for _, v := range d.pool {
		if v.Weight != 0 {
			weightAllot = true
			break
		}
	}
	if weightAllot {
		cs, err := NewChooser(d.pool...)
		if err != nil {
			return nil, err
		}
		return cs.Pick(), nil
	}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(d.pool))
	return d.pool[index], nil
}

func (d *fetcher) GetBlockDataFromTitanOrGateway(ctx context.Context, customGatewayAddr string, c cid.Cid) ([]byte, error) {
	data, err := d.GetBlockData(ctx, c)
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return nil, err
	}
	if data == nil {
		data, err = d.getDataFromCommonGateway(customGatewayAddr, c)
		if err != nil {
			logger.Error("fail get data from gateway : ", err.Error())
			return nil, err
		}
	}
	return data, nil
}

func (d *fetcher) GetBlocksFromTitanOrGateway(ctx context.Context, customGatewayAddr string, ks []cid.Cid) <-chan blocks.Block {
	ch := make(chan blocks.Block)
	go func() {
		defer close(ch)

		var wg sync.WaitGroup
		for _, v := range ks {
			value := v

			wg.Add(1)
			go func(cc context.Context, c cid.Cid) {
				defer wg.Done()
				data, err := d.GetBlockData(cc, c)
				if data == nil {
					data, err = d.getDataFromCommonGateway(customGatewayAddr, c)
					if err != nil {
						logger.Error("fail get data from gateway : ", err.Error())
						return
					}
				}
				block, err := blocks.NewBlockWithCid(data, c)
				if err != nil {
					logger.Error("create block fail : ", err.Error())
					return
				}
				select {
				case ch <- block:
					return
				case <-cc.Done():
					return
				}
			}(ctx, value)
		}
		wg.Wait()
	}()

	return ch
}

func (d *fetcher) GetBlocksFromTitan(ctx context.Context, ks []cid.Cid) <-chan blocks.Block {
	ch := make(chan blocks.Block)
	go func() {
		defer close(ch)

		var wg sync.WaitGroup
		for _, v := range ks {
			value := v
			wg.Add(1)
			go func(cc context.Context, c cid.Cid) {
				defer wg.Done()

				data, err := d.GetBlockData(cc, c)
				if err != nil {
					logger.Error("fail get data from edge node : ", err.Error())
					return
				}
				block, err := blocks.NewBlockWithCid(data, c)
				if err != nil {
					logger.Error("create block fail : ", err.Error())
					return
				}
				select {
				case ch <- block:
					return
				case <-cc.Done():
					return
				}
			}(ctx, value)
		}
		wg.Wait()
	}()

	return ch
}

// getDataFromEdgeNode connect Titan edge node by http get method
func (d *fetcher) getDataFromEdgeNode(di *api.DownloadInfoResult, cid cid.Cid) ([]byte, error) {
	if di.URL == "" {
		return nil, fmt.Errorf("not found target host")
	}
	if di.Sign == "" {
		return nil, fmt.Errorf("sign data is null")
	}
	url := fmt.Sprintf("%s?cid=%s&sign=%s&sn=%d&signTime=%d&timeout=%d",
		di.URL,
		cid.String(),
		di.Sign,
		di.SN,
		di.SignTime,
		di.TimeOut)
	return http2.Get(url, sdkName)
}

func (d *fetcher) getDataFromCommonGateway(customGatewayAddr string, c cid.Cid) ([]byte, error) {
	if customGatewayAddr == "" {
		return nil, fmt.Errorf("not found target host")
	}
	logger.Debugf("got data from common gateway with cid [%s]", c.String())
	url := fmt.Sprintf("%s%s", customGatewayAddr, c.String())
	return http2.PostFromGateway(url)
}

func (d *fetcher) callback(c cid.Cid, sn int64, downloadSuccess bool) {
	// give up the CPU, download first
	runtime.Gosched()

	locator, closer, err := client.NewLocator(context.TODO(), d.locatorAddr, nil)
	if err != nil {
		logger.Error("create schedule fail : ", err.Error())
		return
	}
	defer closer()

	logger.Debugf("[%s] downlaod state : %v", c.String(), downloadSuccess)

	cidSign, err := GetSigner().Sign([]byte(c.String()))
	if err != nil {
		logger.Warn("cid sign fail : ", err.Error())
		return
	}
	bdResult := []api.UserBlockDownloadResult{
		{
			SN:     sn,
			Sign:   cidSign,
			Result: downloadSuccess,
		},
	}

	err = locator.UserDownloadBlockResults(context.TODO(), bdResult)
	if err != nil {
		logger.Warnf("[%s] download callback fail : %s", c.String(), err.Error())
		return
	}

	logger.Debugf("[%s] downlaod callback success", c.String())
}
