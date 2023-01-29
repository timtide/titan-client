package titan_client

import (
	"bufio"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"github.com/ipfs/go-cid"
	files "github.com/ipfs/go-ipfs-files"
	logging "github.com/ipfs/go-log/v2"
	md "github.com/ipfs/go-merkledag"
	unixFile "github.com/ipfs/go-unixfs/file"
	"io"
	gopath "path"
	"strings"
)

// need scientific Internet access
const defaultGatewayAddress = "https://ipfs.io/ipfs/"
const RouteProtocol = "/api/v0/block/get"

var logger = logging.Logger("go-titan-client")

// defaultBufSize is the buffer size for gets. for now, 1MiB, which is ~4 blocks.
const defaultBufSize = 1048576

type Downloader interface {
	// GetReader returns a read pipe
	// note: remember to close after using
	// eg: defer reader.close()
	GetReader(ctx context.Context, cid cid.Cid, archive bool, compressLevel int) (io.ReadCloser, error)

	// Download data from titan to the specified directory according to the cid
	// archive: compress to tar file
	// compressLevel: compress level, eg: gzip.NoCompression
	Download(ctx context.Context, cid cid.Cid, archive bool, compressLevel int, outPath string) error
}

func NewDownloader(option ...Option) Downloader {
	td := &titanDownloader{}
	for _, v := range option {
		v(td)
	}
	if td.customGatewayAddr == "" {
		td.customGatewayAddr = defaultGatewayAddress
	}
	if strings.Contains(td.customGatewayAddr, ":") {
		td.customGatewayAddr = fmt.Sprintf("%s%s%s", strings.TrimRight(td.customGatewayAddr, "/"), RouteProtocol, "?arg=")
	}
	return td
}

type titanDownloader struct {
	customGatewayAddr string
	locatorAddr       string
}

// GetReader returns a read pipe
// note: remember to close after using
// eg: defer reader.close()
func (t *titanDownloader) GetReader(ctx context.Context, cid cid.Cid, archive bool, compressLevel int) (io.ReadCloser, error) {
	logger.Info("begin get reader with cid : ", cid.String())
	bs := newBlockService(t.customGatewayAddr, t.locatorAddr)
	ds := md.NewDAGService(bs)
	nd, err := ds.Get(ctx, cid)
	if err != nil {
		return nil, err
	}

	file, err := unixFile.NewUnixfsFile(ctx, ds, nd)
	if err != nil {
		return nil, err
	}

	return fileArchive(file, cid.String(), archive, compressLevel)
}

// Download data from titan to the specified directory according to the cid
// archive: compress to tar file
// compressLevel: compress level, eg: gzip.NoCompression
func (t *titanDownloader) Download(ctx context.Context, cid cid.Cid, archive bool, compressLevel int, outPath string) error {
	reader, err := t.GetReader(ctx, cid, archive, compressLevel)
	if err != nil {
		return err
	}
	defer reader.Close()

	ow := Writer{
		Archive:     archive,
		Compression: compressLevel,
	}
	logger.Debugf("%s%s", "download data to ", outPath)
	return ow.Write(reader, outPath)
}

func fileArchive(f files.Node, name string, archive bool, compression int) (io.ReadCloser, error) {
	cleaned := gopath.Clean(name)
	_, filename := gopath.Split(cleaned)

	// need to connect a writer to a reader
	pipeReader, pipeWriter := io.Pipe()
	checkErrAndClosePipe := func(err error) bool {
		if err != nil {
			_ = pipeWriter.CloseWithError(err)
			return true
		}
		return false
	}

	// use a buffered writer to parallelize task
	bufWriter := bufio.NewWriterSize(pipeWriter, defaultBufSize)

	// compression determines whether to use gzip compression.
	maybeGzw, err := newMaybeGzWriter(bufWriter, compression)
	if checkErrAndClosePipe(err) {
		return nil, err
	}

	closeGzwAndPipe := func() {
		if err := maybeGzw.Close(); checkErrAndClosePipe(err) {
			return
		}
		if err := bufWriter.Flush(); checkErrAndClosePipe(err) {
			return
		}
		_ = pipeWriter.Close() // everything seems to be ok.
	}

	if !archive && compression != gzip.NoCompression {
		// the case when the node is a file
		r := files.ToFile(f)
		if r == nil {
			return nil, errors.New("file is not regular")
		}

		go func() {
			if _, err := io.Copy(maybeGzw, r); checkErrAndClosePipe(err) {
				return
			}
			closeGzwAndPipe() // everything seems to be ok
		}()
	} else {
		// the case for 1. archive, and 2. not archived and not compressed, in which tar is used anyway as a transport format

		// construct the tar writer
		w, err := files.NewTarWriter(maybeGzw)
		if checkErrAndClosePipe(err) {
			return nil, err
		}

		go func() {
			// write all the nodes recursively
			if err := w.WriteFile(f, filename); checkErrAndClosePipe(err) {
				return
			}
			_ = w.Close()     // close tar writer
			closeGzwAndPipe() // everything seems to be ok
		}()
	}

	return pipeReader, nil
}

type identityWriteCloser struct {
	w io.Writer
}

func (i *identityWriteCloser) Write(p []byte) (int, error) {
	return i.w.Write(p)
}

func (i *identityWriteCloser) Close() error {
	return nil
}

func newMaybeGzWriter(w io.Writer, compression int) (io.WriteCloser, error) {
	if compression != gzip.NoCompression {
		return gzip.NewWriterLevel(w, compression)
	}
	return &identityWriteCloser{w}, nil
}
