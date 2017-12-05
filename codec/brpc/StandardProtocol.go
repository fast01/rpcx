package brpc

import (

	lg  "github.com/fast01/rpcx/log"
	br  "github.com/fast01/rpcx/codec/brpc/proto/baidu/rpc"
	brp "github.com/fast01/rpcx/codec/brpc/proto/baidu/rpc/policy"
	pb "github.com/golang/protobuf/proto"
	"github.com/fast01/rpcx/core"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
)


// Configs
const (
	//bufferSize = 4096
	MaxBodySize = 20 * 1024 * 1024
)

// RawMessage stands for the initial layer of protocol data to unpack
type RawMessage struct {
	meta    []byte  // meta data
	payLoad []byte  // payload length includes attachments
}

// PackRPCHeader serialize the message header
// Notes:
// 1. 12-byte header [PRPC][body_size][meta_size]
// 2. body_size and meta_size are in network byte order
// 3. Use service->full_name() + method_name to specify the method to call
// 4. `attachment_size' is set iff request/response has attachment
// 5. Not supported: chunk_info
func PackRPCHeader(metaSize, payLoadSize int32) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write([]byte("PRPC"))
	//lg.Infof("packing, metaSize=%d, payLoadSize=%d", metaSize, payLoadSize)
	ts := metaSize + payLoadSize
	binary.Write(&buf, binary.BigEndian, ts)
	binary.Write(&buf, binary.BigEndian, metaSize)

	ret := buf.Bytes()
	if len(ret) != 12 {
		return nil, fmt.Errorf("The Head length is not 12!")
	}

	return ret, nil
}

// ParseRPCMessage parse the received bytes from socket
// Notes:
// 1. 12-byte header [PRPC][body_size][meta_size]
// 2. body_size and meta_size are in network byte order
//
func ParseRPCMessage(rd io.Reader) (msg *RawMessage, reterr error) {
	backingBuf := make([]byte, 12)
	lth, err := rd.Read(backingBuf)
	if err != nil {
		// may be eof
		return nil, err
	}

	if lth < len(backingBuf) {
		return nil, fmt.Errorf("Not Enough Data to decode message! %d/%d", lth, 12)
	}

	//reterr = fmt.Errorf("Not enough data")
	if string(backingBuf[0:4]) != "PRPC" {
		return nil, fmt.Errorf("Bad request schema")
	}

	bodySize := binary.BigEndian.Uint32(backingBuf[4:8])
	metaSize := binary.BigEndian.Uint32(backingBuf[8:12])
	//lg.Debugf("backingbuf = %v, bodySize = %d, metaSize = %d", backingBuf, bodySize, metaSize)

	if bodySize > MaxBodySize {
		return nil, fmt.Errorf("Data too big! %d/%d", bodySize, MaxBodySize)
	}
	if uint32(metaSize) > bodySize {
		return nil, fmt.Errorf("metaSize > bodySize")
	}

	restData := make([]byte, bodySize)
	restDataLen, err := rd.Read(restData)
	if err != nil {
		// may be eof
		return nil, err
	}

	if uint32(restDataLen) < bodySize {
		return nil, fmt.Errorf("Not enough data in body: %d/%d", restDataLen, bodySize)
	}
	// payload length includes attachments length
	meta := restData[0: metaSize]
	payLoad := restData[metaSize: bodySize]

	//lg.Infof("metaSize = %d, payLoadSize = %d", len(meta), len(payLoad))
	msg = &RawMessage{meta: meta, payLoad: payLoad }
	return
}


func Header2Meta(header *core.Header, meta *MetaData) {

	CorrelationId := header.Get("CorrelationId")
	if len(CorrelationId) > 0 {
		meta.CorrelationId, _ = strconv.ParseInt(CorrelationId, 10, 64)
	}

	if LogId := header.Get("LogId"); len(LogId) > 0 {
		meta.LogId, _ = strconv.ParseInt(LogId, 10 ,64)
	}

	if SpanId := header.Get("SpanId"); len(SpanId)> 0 {
		meta.SpanId, _ = strconv.ParseInt(SpanId, 10 ,64)
	}

	if PSpanId := header.Get("PSpanId"); len(PSpanId) > 0  {
		meta.PSpanId, _ = strconv.ParseInt(PSpanId, 10 ,64)
	}

	if TraceId := header.Get("TraceId"); len(TraceId)  > 0 {
		meta.TraceId, _ = strconv.ParseInt(TraceId, 10 ,64)
	}
}

func Meta2Header(meta *MetaData, header *core.Header) {
	if meta.CorrelationId > 0 {
		header.Set("CorrelationId",strconv.FormatInt(meta.CorrelationId,10))
	}

	if meta.LogId > 0 {
		header.Set("LogId",strconv.FormatInt(meta.LogId,10))
	}

	if meta.SpanId > 0 {
		header.Set("SpanId",strconv.FormatInt(meta.SpanId,10))
	}

	if meta.PSpanId > 0 {
		header.Set("PSpanId",strconv.FormatInt(meta.PSpanId,10))
	}

	if meta.TraceId > 0 {
		header.Set("TraceId",strconv.FormatInt(meta.TraceId,10))
	}
}


//StandardProtocol stands for baidu rpc standard protocol
type StandardProtocol struct {
}

// Name gives the name of this protocol
func (p StandardProtocol) Name() string {
	return "Standard"
}

// SendRequest serialize and give sends the given request
func (p StandardProtocol) SendRequest(wt io.Writer, req *RPCRequest) (reterr error) {
	defer func() {
		err := recover()
		if err != nil {
			reterr = fmt.Errorf("Error in send request, err = %v", err)
		}
	}()

	requestMeta := &brp.RpcRequestMeta{}
	meta := &brp.RpcMeta{Request: requestMeta}
	meta.Request = requestMeta
	meta.CompressType = pb.Int32(req.Conf.CompressType)
	requestMeta.ServiceName = pb.String(req.MethodDesc.ServiceName)
	requestMeta.MethodName = pb.String(req.MethodDesc.MethodName)

	//TODO: setting requestMeta logid traceid span_id parrent_span_id
	// meta.AuthenticationData
	meta.CorrelationId = &req.Meta.CorrelationId
	requestMeta.LogId  = &req.Meta.LogId
	requestMeta.SpanId = &req.Meta.SpanId
	requestMeta.ParentSpanId = &req.Meta.PSpanId
	requestMeta.TraceId = &req.Meta.TraceId

	if req.Attachment != nil {
		meta.AttachmentSize = pb.Int32(int32(len(req.Attachment)))
	}

	if req.ParamData == nil {
		req.ParamData = make([]byte, 0)
	}

	metaSize := int32(pb.Size(meta))
	var metaData []byte
	var reqBody  []byte
	var head 	 []byte
	var err error

	defer func() {
		if head == nil {
			lg.Infof("Invalid head, cannot send response!")
		} else if metaData == nil {
			lg.Infof("Invalid meta serialization, cannot send request!")
		} else {
			attLen := 0
			sendData := append(head, metaData...)
			sendData = append(sendData, reqBody...)
			if req.Attachment != nil {
				sendData = append(sendData, req.Attachment...)
				attLen = len(req.Attachment)
			}

			var lth int
			lth, reterr = wt.Write(sendData)
			lg.Debugf("%d bytes wrote to request[%d/%d/%d/%d|%d] req=%v err=%v",
					lth, len(head), len(metaData), len(reqBody), attLen, len(sendData), meta, reterr)
		}
	}()

	reqBody, err = CompressData(req.ParamData,  br.CompressType(req.Conf.CompressType))
	if err != nil {
		reterr = err
		reqBody = make([]byte, 0)
	}

	respSize := int32(len(reqBody))
	attSize := int32(len(req.Attachment))
	head, err = PackRPCHeader(metaSize, respSize+attSize)
	if err != nil {
		reterr = err
		return
	}

	metaData, err = pb.Marshal(meta)
	if err != nil {
		reterr = err
		return
	}
	return
}

// ParseRequest receives the request and parse it from binary data
func (p StandardProtocol) ParseRequest(rd io.Reader) (res *RPCRequest, reterr error) {
	defer func() {
		err := recover()
		if err != nil {
			reterr = fmt.Errorf("Error in parse request, err = %v", err)
		}
	}()

	rawMsg, err := ParseRPCMessage(rd)
	if err != nil {
		return nil, err
	}

	metaMsg := &brp.RpcMeta{}
	err = pb.Unmarshal(rawMsg.meta, metaMsg)
	if err != nil {
		lg.Infof("Cannot convert RPC meta")
		return nil, err
	}

	reqMeta := metaMsg.Request
	if reqMeta == nil {
		return nil, fmt.Errorf("request meta is nil")
	}

	ret := &RPCRequest{}
	ret.MethodDesc = MethodDescriptor{MethodName: reqMeta.GetMethodName(),
		ServiceName: reqMeta.GetServiceName()}
	ret.Meta.LogId = reqMeta.GetLogId()
	ret.Meta.TraceId = reqMeta.GetTraceId()
	ret.Meta.SpanId = reqMeta.GetSpanId()
	ret.Meta.PSpanId = reqMeta.GetParentSpanId()
	ret.Meta.CorrelationId = metaMsg.GetCorrelationId()

	requestSize := int32(len(rawMsg.payLoad))
	bodySize := requestSize - metaMsg.GetAttachmentSize()
	requestBuf := rawMsg.payLoad[0:bodySize]
	var attachmentData []byte
	if metaMsg.AttachmentSize != nil {
		if requestSize < metaMsg.GetAttachmentSize() {
			return nil, fmt.Errorf("attachment size [%d] is larger than request size [%d]",
				metaMsg.GetAttachmentSize(), requestSize)
		}
		attachmentData = rawMsg.payLoad[bodySize:]
	}

	dcmpBuf, err := DecompressData(requestBuf, br.CompressType(metaMsg.GetCompressType()))
	if err != nil {
		return nil, err
	}
	ret.ParamData = dcmpBuf
	ret.Attachment = attachmentData

	lg.Debugf("recv from request[12/%d/%d/%d|%d] resp=%v err=%v",
		len(rawMsg.meta), bodySize, len(attachmentData),
		12 + len(rawMsg.meta) + len(rawMsg.payLoad), metaMsg, reterr)

	return ret, nil
}

// ParseResponse receives and deserialize the response
func (p StandardProtocol) ParseResponse(rd io.Reader) (res *RPCResponse, reterr error) {
	defer func() {
		err := recover()
		if err != nil {
			reterr = fmt.Errorf("Error in parse response, err = %v", err)
		}
	}()

	rawMsg, err := ParseRPCMessage(rd)
	if err != nil {
		return nil, err
	}

	metaMsg := &brp.RpcMeta{}
	err = pb.Unmarshal(rawMsg.meta, metaMsg)
	if err != nil {
		lg.Infof("Cannot convert RPC meta")
		return nil, err
	}

	respMeta := metaMsg.Response
	if respMeta == nil {
		reterr = fmt.Errorf("Response meta is nil")
	}

	res = &RPCResponse{}
	res.RetCode = respMeta.GetErrorCode()
	res.FailReason = respMeta.GetErrorText()
	res.Conf.CompressType = metaMsg.GetCompressType()
	res.Meta.CorrelationId = metaMsg.GetCorrelationId()

	requestSize := int32(len(rawMsg.payLoad))
	bodySize := requestSize - metaMsg.GetAttachmentSize()
	respBuf := rawMsg.payLoad[0:bodySize]
	var attachmentData []byte
	if metaMsg.AttachmentSize != nil {
		if requestSize < metaMsg.GetAttachmentSize() {
			return nil, fmt.Errorf("attachment size [%d] is larger than request size [%d]",
				metaMsg.GetAttachmentSize(), requestSize)
		}
		attachmentData = rawMsg.payLoad[bodySize:]
	}

	res.Attachment = attachmentData
	res.ReturnData, reterr = DecompressData(respBuf, br.CompressType(res.Conf.CompressType))

	lg.Debugf("recv from response[12/%d/%d/%d|%d] resp=%v err=%v",
		len(rawMsg.meta), bodySize, len(attachmentData),
			12 + len(rawMsg.meta) + len(rawMsg.payLoad), metaMsg, reterr)
	return res, nil
}


// SendResponse serialize and send the response
func (p StandardProtocol) SendResponse(res *RPCResponse, wt io.Writer) (reterr error) {
	// TODO attachment
	errorCode := res.RetCode
	responseMeta := &brp.RpcResponseMeta{}
	meta := &brp.RpcMeta{Response: responseMeta}
	responseMeta.ErrorCode = pb.Int32(errorCode)
	if res.FailReason != "" {
		responseMeta.ErrorText = pb.String(res.FailReason)
	}

	meta.CorrelationId = pb.Int64(res.Meta.CorrelationId)
	meta.CompressType = pb.Int32(int32(res.Conf.CompressType))
	if res.Attachment != nil {
		meta.AttachmentSize = pb.Int32(int32(len(res.Attachment)))
	}

	metaSize := int32(pb.Size(meta))
	var metaData []byte
	var resBody []byte
	var head []byte
	var err error

	if res.ReturnData == nil {
		res.ReturnData = make([]byte, 0)
	}

	defer func() {
		if head == nil {
			lg.Infof("Invalid head, cannot send response!")
		} else if metaData == nil {
			lg.Infof("Invalid meta serialization, cannot send response!")
		} else {
			sendData := append(head, metaData...)
			sendData = append(sendData, resBody...)
			attLen := 0
			if res.Attachment != nil {
				attLen = len(res.Attachment)
				sendData = append(sendData, res.Attachment...)
			}

			var lth int
			lth, reterr = wt.Write(sendData)
			lg.Debugf("%d bytes wrote to response[%d/%d/%d/%d|%d] resp=%v err=%v",
				lth, len(head), len(metaData), len(resBody), attLen, len(sendData), meta, reterr)
		}
	}()

	resBody, err = CompressData(res.ReturnData, br.CompressType(res.Conf.CompressType))
	if err != nil {
		reterr = err
		resBody = make([]byte, 0)
	}

	respSize := int32(len(resBody))
	attSize := int32(len(res.Attachment))
	head, err = PackRPCHeader(metaSize, respSize+attSize)
	if err != nil {
		reterr = err
		return
	}

	metaData, err = pb.Marshal(meta)
	if err != nil {
		reterr = err
		return
	}

	return
}





