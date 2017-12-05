WORK_ROOT=`pwd`/../../../../../..
PROTO_PATH="github.com/fast01/rpcx/codec/brpc/proto"

IMPORT_PREFIX="github.com/fast01/rpcx/codec/brpc/proto"
#protoc --go_out=. --proto_path .  google/protobuf/*.proto
#protoc --go_out=import_path=$IMPORT_PREFIX,Mbaidu/rpc=$IMPORT_PREFIX/baidu/rpc:. --proto_path=.  baidu/rpc/policy/*.proto

protoc --go_out=import_path=$IMPORT_PREFIX,Mbaidu/rpc/options.proto=$IMPORT_PREFIX/baidu/rpc:. -I. baidu/rpc/policy/*.proto
protoc --go_out=import_path=$IMPORT_PREFIX,Mbaidu/rpc/options.proto=$IMPORT_PREFIX/baidu/rpc,Mgoogle/protobuf/descriptor.proto=$IMPORT_PREFIX/google/protobuf:. -I. baidu/rpc/*.proto
protoc --go_out=import_path=$IMPORT_PREFIX,Mbaidu/rpc/options.proto=$IMPORT_PREFIX/baidu/rpc,Mgoogle/protobuf/descriptor.proto=$IMPORT_PREFIX/google/protobuf:. -I. google/protobuf/*.proto


