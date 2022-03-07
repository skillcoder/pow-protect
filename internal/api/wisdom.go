package api

import pb "github.com/skillcoder/pow-protect/internal/api/sdk/wisdom"

type WisdomServer struct {
	pb.UnimplementedWisdomServer
	challenger challengeVerifier
}

func NewWisdomServer(challenger challengeVerifier) *WisdomServer {
	return &WisdomServer{
		challenger: challenger,
	}
}
