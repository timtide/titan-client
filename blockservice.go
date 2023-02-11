package titan_client

import (
	"context"
	"fmt"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/timtide/titan-client/util"
)

type blockService struct {
	ds                util.Fetcher
	customGatewayAddr string
	locatorAddr       string
}

// newBlockService creates a BlockService with given datastore instance.
func newBlockService(customGatewayAddr, locatorAddr string) *blockService {
	return &blockService{
		ds:                util.NewFetcher(util.WithLocatorAddressOption(locatorAddr)),
		customGatewayAddr: customGatewayAddr,
		locatorAddr:       locatorAddr,
	}
}

// Blockstore returns the blockstore behind this blockservice.
func (s *blockService) Blockstore() blockstore.Blockstore {
	logger.Error("not implemented")
	return nil
}

// Exchange returns the exchange behind this blockservice.
func (s *blockService) Exchange() exchange.Interface {
	logger.Error("not implemented")
	return nil
}

// AddBlock adds a particular block to the service, Putting it into the datastore.
func (s *blockService) AddBlock(ctx context.Context, o blocks.Block) error {
	return fmt.Errorf("%s", "not implemented")
}

func (s *blockService) AddBlocks(ctx context.Context, bs []blocks.Block) error {
	return fmt.Errorf("%s", "not implemented")
}

// GetBlock retrieves a particular block from the service,
// Getting it from the datastore using the key (hash).
func (s *blockService) GetBlock(ctx context.Context, c cid.Cid) (blocks.Block, error) {
	if !c.Defined() {
		return nil, ipld.ErrNotFound{Cid: c}
	}
	data, err := s.ds.GetBlockDataFromTitanOrGateway(ctx, s.customGatewayAddr, c)
	if err != nil {
		return nil, err
	}
	logger.Debug("block data download success")
	block, err := blocks.NewBlockWithCid(data, c)
	if err != nil {
		logger.Error("create block fail : ", err.Error())
		return nil, err
	}
	return block, nil
}

// GetBlocks gets a list of blocks asynchronously and returns through
// the returned channel.
func (s *blockService) GetBlocks(ctx context.Context, ks []cid.Cid) <-chan blocks.Block {
	return s.ds.GetBlocksFromTitanOrGateway(ctx, s.customGatewayAddr, ks)
}

// DeleteBlock deletes a block in the blockservice from the datastore
func (s *blockService) DeleteBlock(ctx context.Context, c cid.Cid) error {
	return fmt.Errorf("%s", "not implemented")
}

func (s *blockService) Close() error {
	return fmt.Errorf("%s", "not implemented")
}
