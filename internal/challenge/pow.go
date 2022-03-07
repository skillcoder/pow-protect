package challenge

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	vdf "github.com/harmony-one/vdf/src/vdf_go"

	pb "github.com/skillcoder/pow-protect/internal/api/sdk/pow"
)

const (
	nonceLen                = 32
	defaultDifficulty int32 = 5000
	stampTimeout            = 10 * time.Second
)

var (
	errShortNonce       = errors.New("short nonce")
	errChallengeTimeout = errors.New("challenge timeout")
	errChallengeInvalid = errors.New("challenge invalid")
)

type Challenge struct {
	rnd io.Reader
}

func New(rnd io.Reader) *Challenge {
	return &Challenge{
		rnd: rnd,
	}
}

func (c *Challenge) Generate() (*pb.Challenge, error) {
	nonce := make([]byte, nonceLen)
	n, err := c.rnd.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("generate %w", err)
	}

	if n != nonceLen {
		return nil, errShortNonce
	}

	difficulty := defaultDifficulty
	stamp := time.Now().UnixMilli()

	return &pb.Challenge{
		Nonce:      nonce,
		Difficulty: difficulty,
		Stamp:      stamp,
		Sign:       signChallenge(nonce, difficulty, stamp),
	}, nil
}

func (c *Challenge) Verify(challenge *pb.Challenge, proof []byte) (bool, error) {
	// challenge timeout
	if time.Since(time.UnixMilli(challenge.GetStamp())) > stampTimeout {
		return false, errChallengeTimeout
	}

	// check challenge consistency
	sign := signChallenge(challenge.GetNonce(), challenge.GetDifficulty(), challenge.GetStamp())
	res := bytes.Compare(sign, challenge.GetSign())
	if res != 0 {
		return false, errChallengeInvalid
	}

	difficulty := int(challenge.GetDifficulty())
	input := *(*[32]byte)(challenge.GetNonce())

	return vdf.New(difficulty, input).Verify(*(*[516]byte)(proof)), nil
}

func signChallenge(nonce []byte, difficulty int32, stamp int64) []byte {
	sign := sha256.Sum256(serialiseChallenge(nonce, difficulty, stamp))
	return sign[:]
}

func serialiseChallenge(nonce []byte, difficulty int32, stamp int64) []byte {
	buf := make([]byte, 0, len(nonce)+4+8)
	buf = append(buf, nonce...)

	b32 := make([]byte, 4)
	binary.LittleEndian.PutUint32(b32, uint32(difficulty))
	buf = append(buf, b32...)

	b64 := make([]byte, 8)
	binary.LittleEndian.PutUint64(b64, uint64(stamp))
	buf = append(buf, b64...)

	return buf
}
