package miners

import (
	"context"

	"github.com/google/uuid"
)

type (
	Coal       int
	MinersType string
)

type Miner interface {
	Run(ctx context.Context) <-chan Coal
	Info() MinerInfo
}

type MinerInfo struct {
	ID          uuid.UUID
	MinerType   MinersType
	MinerEnergy int64
	MinerPower  int64
}
