package api

import pb "github.com/skillcoder/pow-protect/internal/api/sdk/pow"

type challengeGenerator interface {
	Generate() (*pb.Challenge, error)
}

type challengeVerifier interface {
	Verify(challenge *pb.Challenge, proof []byte) (bool, error)
}

type challenger interface {
	challengeGenerator
	challengeVerifier
}
