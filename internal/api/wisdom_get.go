package api

import (
	"context"
	"errors"
	"fmt"
	"log"

	pb "github.com/skillcoder/pow-protect/internal/api/sdk/wisdom"
)

var errInvalidRequest = errors.New("invalid request")

func (s *WisdomServer) Get(_ context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("wisdom request")

	ok, err := s.challenger.Verify(req.GetChallenge(), req.GetProof())
	if err != nil {
		return nil, fmt.Errorf("verify %w", err)
	}

	if !ok {
		return nil, errInvalidRequest
	}

	// TODO: real wisdom
	return &pb.GetResponse{
		WordOfWisdom: "test",
	}, nil
}
