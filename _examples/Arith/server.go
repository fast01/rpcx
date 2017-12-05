package main

import (
	"context"
	"time"

	"github.com/fast01/rpcx"
	"github.com/fast01/rpcx/plugin"
	codec "github.com/fast01/rpcx/codec/brpc"
	_ "github.com/fast01/rpcx/log"
	arith_pb "github.com/fast01/rpcx/_examples/Arith/pb"
	proto "github.com/golang/protobuf/proto"
)
/*

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}
*/

type Arith int

func (t *Arith) Mul(ctx context.Context, args *arith_pb.Args, reply *arith_pb.Reply) error {
	reply.C = proto.Int32(args.GetA() * args.GetB())
	return nil
}

func main() {
	server := rpcx.NewServer()
	server.ServerCodecFunc = codec.NewServerCodec
	server.RegisterName("Arith", new(Arith))

	p := plugin.NewRateLimitingPlugin(time.Second, 1000)
	server.PluginContainer.Add(p)

	server.Serve("tcp", "127.0.0.1:8972")

}
