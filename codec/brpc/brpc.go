package brpc

import (
	"fmt"
	"io"
	"strings"
	"context"
	"sync"
	"github.com/fast01/rpcx/core"
	//"github.com/fast01/rpcx/log"
	br  "github.com/fast01/rpcx/codec/brpc/proto/baidu/rpc"
	//brp "github.com/fast01/rpcx/codec/brpc/proto/baidu/rpc/policy"
	proto "github.com/golang/protobuf/proto"
)


const defaultBufferSize = 4 * 1024

type serverCodec struct {
	mu   sync.Mutex // exclusive writer lock
	rwc  io.ReadWriteCloser
	c    io.Closer
	req  *RPCRequest   //  store req
}

// NewServerCodec returns a new core.ServerCodec.
//
// A ServerCodec implements reading of RPC requests and writing of RPC
// responses for the server side of an RPC session. The server calls
// ReadRequestHeader and ReadRequestBody in pairs to read requests from the
// connection, and it calls WriteResponse to write a response back. The
// server calls Close when finished with the connection. ReadRequestBody
// may be called with a nil argument to force the body of the request to be
// read and discarded.
func NewServerCodec(rwc io.ReadWriteCloser) core.ServerCodec {
	//w := bufio.NewWriterSize(rwc, defaultBufferSize)
	//r := bufio.NewReaderSize(rwc, defaultBufferSize)
	return &serverCodec{
		rwc:  rwc,
		c:   rwc,
	}
}

func (s *serverCodec) WriteResponse(ctx context.Context, res *core.Response, body interface{}) (err error) {
	/*msg, ok := body.(proto.Message)
	if !ok {
		return fmt.Errorf("%T does not implement proto.Message", body)
	}
	*/
	var msg proto.Message
	var resCode  int32 = 0
	var resError = ""
	switch body.(type) {
		case core.InvalidRequest: {
			resCode = int32(br.Errno_EREQUEST)
			if res.Error != "" {
				resError = res.Error
			}
		}
	case  proto.Message : {
			resCode = 0
			resError = ""
			msg = body.(proto.Message)
		}
	default:
		resCode  = int32(br.Errno_SYS_EPROTO)
		resError =  fmt.Sprintf("%T does not implement proto.Message", body)
	}

	rpc_resp :=&RPCResponse{}
	rpc_resp.Conf.CompressType = int32(br.CompressType_COMPRESS_TYPE_NONE)
	rpc_resp.Meta.CorrelationId = int64(res.Seq)
	rpc_resp.FailReason = resError
	rpc_resp.RetCode = resCode
	rpc_resp.Attachment = []byte(res.Header)
	if msg != nil {
		rpc_resp.ReturnData, err = proto.Marshal(msg)
		if err != nil {
			return err
		}
	}

	sp := StandardProtocol{}
	s.mu.Lock()
	err = sp.SendResponse(rpc_resp, s.rwc)
	s.mu.Unlock()
	return err
}

func (s *serverCodec) ReadRequestHeader(ctx context.Context, req *core.Request) (err error) {
	sp := StandardProtocol{}
	if s.req, err = sp.ParseRequest(s.rwc); err != nil  {
		return
	}
	req.ServiceMethod = s.req.MethodDesc.ServiceName + "." + s.req.MethodDesc.MethodName
	req.Seq = uint64(s.req.Meta.CorrelationId)
	if len(s.req.Attachment) > 0 {
		req.Header = string(s.req.Attachment)
	}

	return nil
}

func (s *serverCodec) ReadRequestBody(ctx context.Context, body interface{}) error {
	if body == nil {
		return nil
	}
	if pb, ok := body.(proto.Message); ok {
		return proto.Unmarshal(s.req.ParamData, pb)
	}
	return fmt.Errorf("%T does not implement proto.Message", body)
}

func (s *serverCodec) Close() error { return s.c.Close() }


type clientCodec struct {
	mu  sync.Mutex // exclusive writer lock
	rwc     io.ReadWriteCloser
	c    	io.Closer
	resp 	*RPCResponse   // store resp
	req 	*RPCRequest
}

// NewClientCodec returns a new core.Client.
//
// A ClientCodec implements writing of RPC requests and reading of RPC
// responses for the client side of an RPC session. The client calls
// WriteRequest to write a request to the connection and calls
// ReadResponseHeader and ReadResponseBody in pairs to read responses. The
// client calls Close when finished with the connection. ReadResponseBody
// may be called with a nil argument to force the body of the response to
// be read and then discarded.
func NewClientCodec(rwc io.ReadWriteCloser) core.ClientCodec {
	//w := bufio.NewWriterSize(rwc, defaultBufferSize)
	//r := bufio.NewReaderSize(rwc, defaultBufferSize)
	return &clientCodec{
		rwc: rwc,
		c:   rwc,
	}
}

func (c *clientCodec) WriteRequest(ctx context.Context, req *core.Request, body interface{}) (err error) {
	msg, ok := body.(proto.Message)
	if !ok {
		return fmt.Errorf("%T does not implement proto.Message", body)
	}
	// here req already has some info now ,like seq
	rpc_req :=  &RPCRequest{}
	sidx := strings.LastIndex(req.ServiceMethod, ".")
	if sidx < 0 {
		return fmt.Errorf("wrong format SericeMethod param, %v", req.ServiceMethod)
	}

	rpc_req.MethodDesc.ServiceName = req.ServiceMethod[:sidx]
	rpc_req.MethodDesc.MethodName = req.ServiceMethod[sidx+1:]
	rpc_req.Meta.CorrelationId = int64(req.Seq)
	rpc_req.Conf.CompressType = int32(br.CompressType_COMPRESS_TYPE_NONE)
	rpc_req.Attachment = []byte(req.Header)
	rpc_req.ParamData , err = proto.Marshal(msg)
	if err != nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	sp := StandardProtocol{}
	err =  sp.SendRequest(c.rwc, rpc_req)
	return err
}

func (c *clientCodec) ReadResponseHeader(resp *core.Response) (err error) {
	sp := StandardProtocol{}
	if c.resp, err = sp.ParseResponse(c.rwc); err != nil  {
		return
	}
	//resp.ServiceMethod = c.resp.Meta.Header
	resp.Seq = uint64(c.resp.Meta.CorrelationId)
	resp.Error = c.resp.FailReason
	return nil
}

func (c *clientCodec) ReadResponseBody(body interface{}) (err error) {
	if body == nil {
		// discard body data
		return nil
	}
	if pb, ok := body.(proto.Message); ok {
		return proto.Unmarshal(c.resp.ReturnData, pb)
	}

	return fmt.Errorf("%T does not implement proto.Message", body)
}


func (c *clientCodec) Close() error { return c.c.Close() }



// NewBrpcServerCodec creates a protobuf ServerCodec by https://github.com/mars9/codec
func NewBrpcServerCodec(conn io.ReadWriteCloser) core.ServerCodec {
	return NewServerCodec(conn)
}

// NewBrpcClientCodec creates a protobuf ClientCodec by https://github.com/mars9/codec
func NewBrpcClientCodec(conn io.ReadWriteCloser) core.ClientCodec {
	return NewClientCodec(conn)
}
