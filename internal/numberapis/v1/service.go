package v1

import (
	"context"

	"github.com/golang/glog"
	numv1 "github.com/tuanden0/generate_number/proto/gen/go/v1/number"
)

type Service interface {
	numv1.NumberServiceServer
}

type service struct {
	numv1.UnimplementedNumberServiceServer
	repo      Repository
	validator validate
}

func NewService(repo Repository, validator validate) Service {
	return &service{repo: repo, validator: validator}
}

func (s *service) Ping(ctx context.Context, in *numv1.NumberServicePingRequest) (*numv1.NumberServicePingResponse, error) {
	return &numv1.NumberServicePingResponse{}, nil
}

func (s *service) Get(ctx context.Context, in *numv1.NumberServiceGetRequest) (*numv1.NumberServiceGetResponse, error) {

	// Validate user input
	if err := s.validator.Get(ctx, in); err != nil {
		glog.Error("user input error %v", err)
		return nil, err
	}

	// Get random number
	num, err := s.repo.Get(ctx, in)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	return &numv1.NumberServiceGetResponse{
		Value: num,
	}, nil
}
