package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/fast01/rpcx"
	uidmap_pb "github.com/fast01/rpcx/_examples/uidmap/pb"
	codec "github.com/fast01/rpcx/codec/brpc"
	"github.com/fast01/rpcx/log"
	proto "github.com/gogo/protobuf/proto"
	"sync"
	"time"
)

var (
	NUM_GO_ROUTINUES = 1
	NUM_REQUEST      = 1
)

func main() {
	var address = flag.String("address", "127.0.0.1:8320", "set the listening address:port")
	var methodName = flag.String("method", "Query_uid", "full name of rpc service method")
	var serviceName = flag.String("service", "demo.uidmap.UidmapService", "service name")
	var numGoRoutinues = flag.Int("r", NUM_GO_ROUTINUES, "num of go routinues")
	var numRequests = flag.Int("c", NUM_REQUEST, "num of requests per go routinue")
	var loglevelVar int
	flag.IntVar(&loglevelVar, "loglevel", int(log.LOG_LEVEL_INFO),
		"set loglevel: 0-trace, 1-debug, 2-info, 3-notice, 4-wran, 5-error, 6-fatal, 7-panic")
	flag.Parse()
	log.SetLogLevel(log.LogLevel(loglevelVar))

	fmt.Println("loglevel: ", log.GetLogLevel())

	s := &rpcx.DirectClientSelector{Network: "tcp", Address: *address, DialTimeout: 10 * time.Second}
	log.Infof("client select address:%s", *address)
	wg := sync.WaitGroup{}
	tsum := time.Now().UnixNano()
	for i := 0; i < *numGoRoutinues; i++ {
		wg.Add(1)
		go func(id int) {
			client := rpcx.NewClient(s)
			client.ClientCodecFunc = codec.NewClientCodec
			args := &uidmap_pb.QueryUidRequest{
				Accinfo: &uidmap_pb.AccountInfo{
					Type: uidmap_pb.AccountType_NORMAL.Enum(),
					Id:   proto.Uint64(1000100),
					//Id:   proto.Uint64(30000000100),
					Flag: proto.Int32(int32(uidmap_pb.AccountFlag_INIT)),
				},
			}
			var reply uidmap_pb.QueryUidResponse
			for j := 0; j < *numRequests; j++ {
				reply.Reset()
				t := time.Now().UnixNano()
				err := client.Call(context.Background(),
					fmt.Sprintf("%s.%s", *serviceName, *methodName),
					args, &reply)
				t = time.Now().UnixNano() - t
				if err != nil {
					log.Errorf("query_uid failed: req=%+v res=%+v cost=%d(us) err=%v", args, reply, t/1000, err)
				} else {
					log.Infof("query_uid success: req=%+v res=%+v %q cost=%d(us)", args, reply, reply.GetMsg(), t/1000)
				}
				//time.Sleep(time.Millisecond * 500)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	tsum = time.Now().UnixNano() - tsum
	log.Infof("total time cost=%d(us)", tsum/1000)

}
