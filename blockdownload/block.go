package blockdownload

import (
	"context"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log/v2"
	"github.com/timtide/titan-client/util"
)

var logger = logging.Logger("titan-client/blockdownload")

// BlockGetter only from titan get block data
type BlockGetter interface {
	// GetBlock gets the requested block.
	GetBlock(ctx context.Context, c cid.Cid) (blocks.Block, error)

	// GetBlocks The scheduler queries the corresponding
	// edge node information according to the incoming value. Each value
	// is assigned to the corresponding edge node for global optimization.
	// schedule service mapping cid to edge node.
	GetBlocks(ctx context.Context, ks []cid.Cid) <-chan blocks.Block
}

func NewBlockGetter(option ...Option) BlockGetter {
	bg := &blockGetter{}
	for _, v := range option {
		v(bg)
	}
	bg.ds = util.NewDataService(util.WithLocatorAddressOption(bg.locatorAddr))
	return bg
}

type blockGetter struct {
	ds          util.DataService
	locatorAddr string
}

func (b *blockGetter) GetBlock(ctx context.Context, c cid.Cid) (blocks.Block, error) {
	logger.Debugf("get block with cid [%s]", c.String())
	if !c.Defined() {
		return nil, ipld.ErrNotFound{Cid: c}
	}
	data, err := b.ds.GetDataFromTitanByCid(ctx, c)
	if err != nil {
		return nil, err
	}
	logger.Debug("block data download success")
	block, err := blocks.NewBlockWithCid(data, c)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (b *blockGetter) GetBlocks(ctx context.Context, ks []cid.Cid) <-chan blocks.Block {
	logger.Debug("start batch download block")
	return b.ds.GetBlockFromTitanByCids(ctx, ks)
}
