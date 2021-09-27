package v1

import (
	"context"

	numv1 "github.com/tuanden0/generate_number/proto/gen/go/v1/number"
)

type Repository interface {
	Get(ctx context.Context, in *numv1.NumberServiceGetRequest) (int64, error)
}

type repoManager struct{}

func NewRepoManager() Repository {
	return &repoManager{}
}

func (r *repoManager) Get(ctx context.Context, in *numv1.NumberServiceGetRequest) (int64, error) {

	// Get data
	num_a := in.GetNumA()
	num_b := in.GetNumB()
	num_c := in.GetNumC()
	num_d := in.GetNumD()

	ch := make(chan int64)

	go func() {
		ch <- generateRangeNumber(num_a, num_b)
	}()

	go func() {
		ch <- generateRangeNumber(num_c, num_d)
	}()

	select {
	case <-ctx.Done():
		return int64(0), ctx.Err()
	case n := <-ch:
		return n, nil
	}
}
