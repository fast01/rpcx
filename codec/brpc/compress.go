package brpc

//package protocols

import (
	br  "github.com/fast01/rpcx/codec/brpc/proto/baidu/rpc"
	//brp  "github.com/fast01/rpcx/codec/brpc/proto/baidu/rpc/policy"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"github.com/golang/snappy"
	"io/ioutil"
)

// CompressType
const (
	TypeNone   = br.CompressType_COMPRESS_TYPE_NONE
	TypeSnappy = br.CompressType_COMPRESS_TYPE_SNAPPY
	TypeGZIP   = br.CompressType_COMPRESS_TYPE_GZIP
	TypeZLIB   = br.CompressType_COMPRESS_TYPE_ZLIB
)

func CompressTypeFromString(t string) (ctp br.CompressType, suc bool) {
	ctpN, suc := br.CompressType_value[t]
	ctp = br.CompressType(ctpN)
	return
}

// CompressData compresses the data with given algorithm
func CompressData(src []byte, ctp br.CompressType) (ret []byte, err error) {
	dst := make([]byte, len(src))
	var buf bytes.Buffer
	switch ctp {
	case br.CompressType_COMPRESS_TYPE_NONE:
		ret = src
	case br.CompressType_COMPRESS_TYPE_SNAPPY:
		ret = snappy.Encode(dst, src)
		if ret == nil {
			err = fmt.Errorf("cannot encode with snappy")
		}
	case br.CompressType_COMPRESS_TYPE_GZIP:
		wt := gzip.NewWriter(&buf)
		wt.Write(src)
		ret, err = ioutil.ReadAll(&buf)
	case br.CompressType_COMPRESS_TYPE_ZLIB:
		wt := zlib.NewWriter(&buf)
		wt.Write(src)
		ret, err = ioutil.ReadAll(&buf)
	case br.CompressType_COMPRESS_TYPE_LZ4:
		err = fmt.Errorf("LZ4 is currently not supported")
	}

	return
}

// DecompressData decompresses the data with given algorithm
func DecompressData(src []byte, ctp br.CompressType) (ret []byte, err error) {
	dst := make([]byte, len(src))
	var buf bytes.Buffer
	switch ctp {
	case br.CompressType_COMPRESS_TYPE_NONE:
		ret = src
	case br.CompressType_COMPRESS_TYPE_SNAPPY:
		ret, err = snappy.Decode(dst, src)
	case br.CompressType_COMPRESS_TYPE_GZIP:
		buf.Write(src)
		rd, err := gzip.NewReader(&buf)
		if err == nil {
			ret, err = ioutil.ReadAll(rd)
		}
	case br.CompressType_COMPRESS_TYPE_ZLIB:
		wt := zlib.NewWriter(&buf)
		wt.Write(src)
		ret, err = ioutil.ReadAll(&buf)
	case br.CompressType_COMPRESS_TYPE_LZ4:
		err = fmt.Errorf("LZ4 is currently not supported")
	}

	return
}
