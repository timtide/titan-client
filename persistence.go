package titan_client

import (
	"compress/gzip"
	"github.com/ipfs/tar-utils"
	"io"
	"os"
	"strings"
)

type Writer struct {
	Archive     bool
	Compression int
}

func (gw *Writer) Write(r io.Reader, fpath string) error {
	if gw.Archive || gw.Compression != gzip.NoCompression {
		return gw.writeArchive(r, fpath)
	}
	return gw.writeExtracted(r, fpath)
}

func (gw *Writer) writeArchive(r io.Reader, fpath string) error {
	// adjust file name if tar
	if gw.Archive {
		if !strings.HasSuffix(fpath, ".tar") && !strings.HasSuffix(fpath, ".tar.gz") {
			fpath += ".tar"
		}
	}

	// adjust file name if gz
	if gw.Compression != gzip.NoCompression {
		if !strings.HasSuffix(fpath, ".gz") {
			fpath += ".gz"
		}
	}

	// create file
	file, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return err
}

func (gw *Writer) writeExtracted(r io.Reader, fpath string) error {
	extractor := &tar.Extractor{Path: fpath, Progress: nil}
	return extractor.Extract(r)
}
