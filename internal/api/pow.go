package api

import (
	pb "github.com/skillcoder/pow-protect/internal/api/sdk/pow"
)

type PowServer struct {
	pb.UnimplementedPowServer
	challenger challengeGenerator
}

func NewPowServer(challenger challengeGenerator) *PowServer {
	return &PowServer{
		challenger: challenger,
	}
}
