package main

import (
	"context"
	"time"
	"sync"
	"github.com/fast01/rpcx"
	codec "github.com/fast01/rpcx/codec/brpc"
	"github.com/fast01/rpcx/log"
	arith_pb "github.com/fast01/rpcx/_examples/Arith/pb"
	proto "github.com/gogo/protobuf/proto"
)
/*
type Args struct {
	A int `msg:"a"`
	B int `msg:"b"`
}

type Reply struct {
	C int `msg:"c"`
}*/

func main() {
	s := &rpcx.DirectClientSelector{Network: "tcp", Address: "127.0.0.1:8972", DialTimeout: 10 * time.Second}
	wg := sync.WaitGroup{}
	tsum := time.Now().UnixNano()
	for i:=0; i< 10; i++ {
		wg.Add(1)
		go func(i int) {
			client := rpcx.NewClient(s)
			client.ClientCodecFunc = codec.NewClientCodec
			defer client.Close()
			args := &arith_pb.Args{A: proto.Int32(7), B: proto.Int32(8)}
			var reply arith_pb.Reply
			for j:=0; j< 2; j++ {
				reply.Reset()
				t := time.Now().UnixNano()
				divCall := client.Go(context.Background(), "Arith.Mul", args, &reply, nil)
				replyCall := <-divCall.Done // will be equal to divCall
				t = time.Now().UnixNano() - t
				if replyCall.Error != nil {
					log.Errorf("error for Arith: %d*%d, %v cost=%d(ns)", args.GetA(), args.GetB(), replyCall.Error, t)
				} else {
					log.Infof("Arith: %d*%d=%d cost=%d(ns)", args.GetA(), args.GetB(), reply.GetC(), t)
				}
				time.Sleep(time.Millisecond * 500)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	tsum = time.Now().UnixNano() - tsum
	log.Infof("total time cost=%d(ns)", tsum)
}
