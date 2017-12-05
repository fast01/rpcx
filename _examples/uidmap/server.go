package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/fast01/rpcx"
	pb "github.com/fast01/rpcx/_examples/uidmap/pb"
	codec "github.com/fast01/rpcx/codec/brpc"
	"github.com/fast01/rpcx/log"
	"github.com/fast01/rpcx/plugin"
	proto "github.com/golang/protobuf/proto"
	"time"
)

type UidmapService struct {
}

func (t *UidmapService) Query_uid(ctx context.Context, args *pb.QueryUidRequest, reply *pb.QueryUidResponse) error {
	time_cost := time.Now().UnixNano()
	accinfo := args.GetAccinfo()
	if accinfo.GetId() < 2000000000 {
		reply.Msg = proto.String("OK")
		reply.Status = pb.ECode_SUCCESS.Enum()
		reply.Uid = proto.Clone(accinfo).(*pb.AccountInfo)
		reply.Uid.Type = pb.AccountType_NORMAL.Enum()
		reply.Uid.Flag = proto.Int32(int32(pb.AccountFlag_ACTIVE))
		return nil
	} else {
		reply.Msg = proto.String("Diabled")
		reply.Status = pb.ECode_SUCCESS.Enum()
		reply.Uid = proto.Clone(accinfo).(*pb.AccountInfo)
		reply.Uid.Type = pb.AccountType_THIRD.Enum()
		reply.Uid.Flag = proto.Int32(int32(pb.AccountFlag_INIT))
	}
	time_cost = time.Now().UnixNano() - time_cost
	log.Logf("query_uid done: req=%+v res=%+v %q cost=%d(us)", args, reply, reply.GetMsg(), time_cost/1000)
	return nil
}

func main() {
	var address = flag.String("listen", "127.0.0.1:8320", "set the listening address:port")
	var serviceName = flag.String("service", "demo.uidmap.UidmapService", " service name")
	var methodName = flag.String("method", "Query_uid", "full name of rpc service method")
	var loglevelVar int
	flag.IntVar(&loglevelVar, "loglevel", int(log.LOG_LEVEL_INFO),
		"set loglevel: 0-trace, 1-debug, 2-info, 3-notice, 4-wran, 5-error, 6-fatal, 7-panic")
	flag.Parse()
	log.SetLogLevel(log.LogLevel(loglevelVar))
	fmt.Println("loglevel: ", log.GetLogLevel())

	server := rpcx.NewServer()
	server.ServerCodecFunc = codec.NewServerCodec
	server.RegisterName(*serviceName, new(UidmapService))

	{
		p := plugin.NewRateLimitingPlugin(time.Second, 2000)
		server.PluginContainer.Add(p)
	}
	{
		p := plugin.NewAliasPlugin()
		server.PluginContainer.Add(p)
		//set alias
		p.Alias(*methodName, fmt.Sprintf("%s.%s", *serviceName, "Query_uid"))
		p.Alias(fmt.Sprintf("%s.%s", *serviceName, *methodName), fmt.Sprintf("%s.%s", *serviceName, "Query_uid"))
	}

	log.Infof("start server at :%s, accept method: %s.%s", *address, *serviceName, *methodName)
	err := server.Serve("tcp", *address)
	log.Error(err)
}
