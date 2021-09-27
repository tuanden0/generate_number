package main

import (
	"flag"

	"github.com/golang/glog"
	v1 "github.com/tuanden0/generate_number/internal/numberapis/v1"
)

var (
	addr = flag.String("ip", "0.0.0.0:8000", "server address:port")
)

func init() {
	// Default setting log for glob
	flag.Lookup("v").Value.Set("2")
	flag.Lookup("logtostderr").Value.Set("true")
	flag.Parse()
}

func main() {

	// Flush all glog to file
	defer glog.Flush()

	// Create number repo
	repo := v1.NewRepoManager()

	// Create number validator
	validator := v1.NewValidate()
	if err := validator.Init(); err != nil {
		glog.Fatalf("failed to create number validator %v\n", err)
	}

	// Create number service
	srv := v1.NewService(repo, validator)

	if err := v1.RunServer(srv, *addr); err != nil {
		glog.Fatal(err)
	}
}
