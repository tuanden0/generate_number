package main

import (
	"context"
	"flag"
	"time"

	"github.com/golang/glog"
	numv1 "github.com/tuanden0/generate_number/proto/gen/go/v1/number"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	addr = flag.String("ip", "0.0.0.0:8000", "server address:port")
)

func init() {
	flag.Lookup("v").Value.Set("2")
	flag.Lookup("logtostderr").Value.Set("true")
	flag.Parse()
}

func main() {

	// Make sure to write log to file
	defer glog.Flush()

	// Create client connection
	ctx, cancer := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancer()
	cc, err := grpc.DialContext(ctx, *addr, setupClientDialOpts()...)
	if err != nil {
		glog.Fatal(err)
	}
	defer cc.Close()

	// Create client
	client := numv1.NewNumberServiceClient(cc)

	// Call server
	getPing(client)
	getNumber(client)

}

func getPing(client numv1.NumberServiceClient) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	if _, err := client.Ping(ctx, &numv1.NumberServicePingRequest{}); err != nil {
		glog.Error("server is down")
	} else {
		glog.Info("server is alive")
	}
}

func getNumber(client numv1.NumberServiceClient) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	val, err := client.Get(ctx, &numv1.NumberServiceGetRequest{
		NumA: 10,
		NumB: -10,
		NumC: 100,
		NumD: 20,
	})

	errDetails := handleServerError(err)
	if len(errDetails) != 0 {
		glog.Errorf("failed to get number %v", errDetails)
		return
	}

	glog.Infof("got a random number %v", val.GetValue())
}

func setupClientDialOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}
}

func handleServerError(err error) map[string]string {

	if err != nil {
		res := make(map[string]string)
		st, ok := status.FromError(err)
		if ok {
			if len(st.Details()) != 0 {
				for _, detail := range st.Details() {
					switch t := detail.(type) {
					case *errdetails.BadRequest:
						for _, violation := range t.GetFieldViolations() {
							res[violation.GetField()] = violation.GetDescription()
						}
					}
				}
			} else {
				res[st.Code().String()] = st.Message()
			}
		}

		return res
	}

	return nil
}
