package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	pbPow "github.com/skillcoder/pow-protect/internal/api/sdk/pow"
	pbWisdom "github.com/skillcoder/pow-protect/internal/api/sdk/wisdom"

	vdf "github.com/harmony-one/vdf/src/vdf_go"
	"google.golang.org/grpc"
)

const (
	address    = "localhost:8088"
	sizeInBits = 2048
)

func main() {
	word, err := getWordOfWisdom()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	log.Printf("Word Of Wisdom: %s\n", word)
}

func getWordOfWisdom() (string, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return "", fmt.Errorf("dial %w", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close grpc conn: %v\n", err)
		}
	}()

	cPow := pbPow.NewPowClient(conn)
	cWisdom := pbWisdom.NewWisdomClient(conn)

	log.Println("get challenge")
	challenge, err := getChallenge(cPow)
	if err != nil {
		return "", fmt.Errorf("dial %w", err)
	}

	log.Println("pow")
	proof := pow(challenge)

	log.Println("get wisdom")
	word, err := getWisdom(cWisdom, challenge, proof)
	if err != nil {
		return "", fmt.Errorf("dial %w", err)
	}

	return word, nil
}

func getChallenge(c pbPow.PowClient) (*pbPow.Challenge, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r, err := c.Challenge(ctx, &pbPow.ChallengeRequest{})
	if err != nil {
		return nil, fmt.Errorf("challenge %w", err)
	}

	return r.GetChallenge(), nil
}

func pow(c *pbPow.Challenge) []byte {
	y, proof := vdf.GenerateVDF(c.GetNonce(), int(c.GetDifficulty()), sizeInBits)
	return append(y, proof...)
}

func getWisdom(c pbWisdom.WisdomClient, challenge *pbPow.Challenge, proof []byte) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &pbWisdom.GetRequest{
		Challenge: challenge,
		Proof:     proof,
	})
	if err != nil {
		return "", fmt.Errorf("get %w", err)
	}

	return r.GetWordOfWisdom(), nil
}
