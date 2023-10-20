package raf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/jpeg"
	"io"

	"github.com/koud-fi/pkg/blob"
)

// based on https://libopenraw.freedesktop.org/formats/raf/

const (
	RAFMime = "image/x-fuji-raf"

	RAFMetaSensorDimensions      RAFMetaTag = 0x100
	RAFMetaActiveAreaTopLeft     RAFMetaTag = 0x110
	RAFMetaActiveAreaHeightWidth RAFMetaTag = 0x111
	RAFMetaOutputHeightWidth     RAFMetaTag = 0x121
	RAFMetaRawInfo               RAFMetaTag = 0x130
)

var rafMagic = []byte("FUJIFILMCCD-RAW ")

type RAF struct {
	Header      RAFHeader
	JPEGConfig  image.Config
	JPEG        []byte
	MetaHeader  RAFMetaHeader
	MetaRecords []RAFMetaRecord
	CFA         []byte
}

type RAFHeader struct {
	Magic         [16]byte
	FormatVersion [4]byte
	CameraID      [8]byte
	Camera        [32]byte
	Dir           struct {
		Version  [4]byte
		_        [20]byte
		JPEG     RAFOffset
		Metadata RAFOffset
		CFA      RAFOffset
	}
}

type RAFOffset struct {
	Offset int32
	Len    int32
}

type RAFMetaHeader struct {
	RecordCount int32
}

type RAFMetaRecord struct {
	Tag  RAFMetaTag
	Size int16
	Data []byte
}

type RAFMetaTag int16

func DecodeRAF(b blob.Blob) (raf RAF, _ error) {
	rc, err := b.Open()
	if err != nil {
		return raf, fmt.Errorf("open: %w", err)
	}
	type bufType interface {
		io.Reader
		io.ReaderAt
	}
	var buf bufType
	if b, ok := rc.(bufType); ok {
		buf = b
		defer rc.Close()

	} else {
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return raf, fmt.Errorf("read-all: %w", err)
		}
		buf = bytes.NewReader(data)
	}
	if err := binary.Read(buf, binary.BigEndian, &raf.Header); err != nil {
		return raf, fmt.Errorf("read header: %w", err)
	}
	if !bytes.Equal(raf.Header.Magic[:], rafMagic) {
		return raf, fmt.Errorf("bad magic: %s", string(raf.Header.Magic[:]))
	}
	if raf.JPEG, err = readRAFData(buf, raf.Header.Dir.JPEG); err != nil {
		return raf, fmt.Errorf("read jpg: %w", err)
	}
	if raf.JPEGConfig, err = jpeg.DecodeConfig(bytes.NewReader(raf.JPEG)); err != nil {
		return raf, fmt.Errorf("decode jpg config: %w", err)
	}

	// TODO: extract exif data from the JPEG

	// TODO: parse metadata records

	if raf.CFA, err = readRAFData(buf, raf.Header.Dir.CFA); err != nil {
		return raf, fmt.Errorf("read cfa: %w", err)
	}
	return
}

func readRAFData(buf io.ReaderAt, ro RAFOffset) ([]byte, error) {
	data := make([]byte, ro.Len)
	_, err := buf.ReadAt(data, int64(ro.Offset))
	return data, err
}
