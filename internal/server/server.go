package server

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/skillcoder/pow-protect/internal/api"
	pbPow "github.com/skillcoder/pow-protect/internal/api/sdk/pow"
	pbWisdom "github.com/skillcoder/pow-protect/internal/api/sdk/wisdom"
	"github.com/skillcoder/pow-protect/internal/challenge"
	"github.com/skillcoder/pow-protect/internal/random"
)

const defaultPort = ":8088"

func Run() error {
	rndSource := random.New(rand.NewSource(time.Now().UnixNano()))
	challenger := challenge.New(rndSource)

	lis, err := net.Listen("tcp", defaultPort)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	s := grpc.NewServer()

	pbPow.RegisterPowServer(s, api.NewPowServer(challenger))
	pbWisdom.RegisterWisdomServer(s, api.NewWisdomServer(challenger))

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("serve: %w", err)
	}

	return nil
}
