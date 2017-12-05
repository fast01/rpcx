package brpc

type RPC  struct {
	Head  []byte

	Body  []byte
}

type MethodDescriptor struct {
	ServiceName string
	MethodName  string
}

type MetaData struct {
	LogId         int64
	TraceId       int64
	SpanId    	  int64
	PSpanId       int64
	CorrelationId int64
	Header        string
}

type RPCRequest struct {
	Meta  		MetaData
	MethodDesc  MethodDescriptor
	Conf RPCDataConf
	ParamData  []byte
    Attachment []byte
}

type RPCDataConf struct {
    CompressType int32
}

type RPCResponse struct {
	Meta  	MetaData
	RetCode    int32
	FailReason string
	Conf RPCDataConf
	ReturnData []byte
    Attachment []byte
}





