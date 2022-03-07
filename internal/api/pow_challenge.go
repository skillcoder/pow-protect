package api

import (
	"context"
	"fmt"
	"log"

	pb "github.com/skillcoder/pow-protect/internal/api/sdk/pow"
)

func (s *PowServer) Challenge(_ context.Context, _ *pb.ChallengeRequest) (*pb.ChallengeResponse, error) {
	log.Printf("challenge request")

	challenge, err := s.challenger.Generate()
	if err != nil {
		return nil, fmt.Errorf("generate %w", err)
	}

	return &pb.ChallengeResponse{
		Challenge: challenge,
	}, nil
}
