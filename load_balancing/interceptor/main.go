package main
import (
	//"context"
	"flag"
	"log"
	"golang.org/x/net/context"
	//pb "github.com/smallnest/grpc/a/pb"
	ecpb "google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc"
)
var (
	address = flag.String("addr", "localhost:50052", "address")
	name    = flag.String("n", "world", "name")
)
func main() {
	flag.Parse()
	// 连接服务器
	chain := chainInterceptor(UnaryClientInterceptor1, UnaryClientInterceptor2)
	conn, err := grpc.Dial(*address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(chain),
		grpc.WithStreamInterceptor(StreamClientInterceptor))
	if err != nil {
		log.Fatalf("faild to connect: %v", err)
	}
	defer conn.Close()
	c := ecpb.NewEchoClient(conn)
	r, err := c.UnaryEcho(context.Background(), &ecpb.EchoRequest{Message: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("before invoker. method: %+v, request:%+v", method, req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("after invoker. reply: %+v", reply)
	return err
}

func UnaryClientInterceptor1(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("1before invoker. method: %+v, request:%+v", method, req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("1after invoker. reply: %+v", reply)
	return err
}

func UnaryClientInterceptor2(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("2before invoker. method: %+v, request:%+v", method, req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("2after invoker. reply: %+v", reply)
	return err
}

func chainInterceptor(first grpc.UnaryClientInterceptor, second grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	return func(ctxOuter context.Context, methodOuter string, reqOuter, replyOuter interface{}, ccOuter *grpc.ClientConn, invokerOuter grpc.UnaryInvoker, optsOuter ...grpc.CallOption) error {
		var inner grpc.UnaryInvoker
		inner = func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			return second(ctx, method, req, reply, cc, invokerOuter, opts...)
		}
		return first(ctxOuter, methodOuter, reqOuter, replyOuter, ccOuter, inner, optsOuter...)
	}
}

func StreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Printf("before invoker. method: %+v, StreamDesc:%+v", method, desc)
	clientStream, err := streamer(ctx, desc, cc, method, opts...)
	log.Printf("before invoker. method: %+v", method)
	return clientStream, err
}
