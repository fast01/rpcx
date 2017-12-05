package brpc

import (
    "log"
    "testing"
    "bytes"
)

func SampleRequest() *RPCRequest {
    return &RPCRequest{
        MethodDesc: MethodDescriptor{ServiceName:"Test", MethodName:"Test"},
        ParamData:[]byte{0x00, 0x01, 0x02, 0x03, 0x04},
        Attachment:[]byte{0x00, 0x01, 0x02, 0x03, 0x04},
        Conf:RPCDataConf{CompressType:int32(TypeNone)},
    }
}

func TestReqeustString(t *testing.T) {
    log.SetFlags(log.Lshortfile | log.LstdFlags)
    req := SampleRequest()
    log.Printf("req info: %+v", req)
}

func TestSendParseReqeust(t *testing.T) {
    log.Printf("Testing Parse request")
    
    var buf bytes.Buffer
    std := StandardProtocol{}
    err := std.SendRequest(&buf, SampleRequest())
    if err != nil {
        t.Fatalf("send request error = %v", err)
    }
    
    recReq, err := std.ParseRequest(&buf)
    if err != nil {
        t.Errorf("Parse returns error = %v", err)    
    }
    
    if recReq == nil {
        t.Errorf("Parse result is nil")
    }
    
    if recReq.MethodDesc.ServiceName != "Test" {
        t.Errorf("Parse service name error")
    }
    
    if recReq.MethodDesc.MethodName != "Test" {
        t.Errorf("Parse method name error")
    }
    
    if len(recReq.ParamData) != 5 {
        t.Errorf("Param data length error [%d]", len(recReq.ParamData))
    }
    
    if len(recReq.Attachment) != 5 {
        t.Errorf("Attachment data length error [%d]", len(recReq.Attachment))
    }
}

func SampleResponse() *RPCResponse {
    return &RPCResponse {
        Meta: MetaData{CorrelationId:1},
        RetCode:-1,
        FailReason:"Does not fail",
        ReturnData:[]byte{0x00, 0x01, 0x02, 0x03, 0x04},
        Attachment:[]byte{0x00, 0x01, 0x02, 0x03, 0x04},
        Conf:RPCDataConf{CompressType:int32(TypeNone)},
    }
}

func TestSendParseResponse(t *testing.T) {
    log.Printf("Testing Parse response")
    
    var buf bytes.Buffer
    std := StandardProtocol{}
    err := std.SendResponse(SampleResponse(), &buf)
    if err != nil {
        t.Fatalf("send request error = %v", err)
    }
    
    recRes, err := std.ParseResponse(&buf)
    if err != nil {
        t.Errorf("Parse returns error = %v", err)    
    }
    if recRes.RetCode != -1 {
        t.Errorf("RetCode not equal!")
    }
    if recRes.Meta.CorrelationId != 1 {
        t.Errorf("Correlation id is not 1")
    }
    
    if recRes.FailReason != "Does not fail" {
        t.Errorf("Fail reason is inconsistent")
    }
    
    if len(recRes.ReturnData) != 5 {
        t.Errorf("Return data length is incorrect!")
    }
    
    if len(recRes.Attachment) != 5 {
        t.Errorf("Attachment length is incorrect!")
    }
}
