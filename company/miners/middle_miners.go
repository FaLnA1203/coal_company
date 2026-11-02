package miners

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type MiddleMiner struct {
	id       uuid.UUID
	price    int
	energy   *atomic.Int64
	power    *atomic.Int64
	interval time.Duration
}

func NewMiddleMiner() *MiddleMiner {
	const (
		normalMinerEnergy    = 45
		normalMinerPower     = 3
		normalMinerlInterval = 2 * time.Second
	)

	energy := &atomic.Int64{}
	power := &atomic.Int64{}

	energy.Add(normalMinerEnergy)
	power.Add(normalMinerPower)

	return &MiddleMiner{
		id:     uuid.New(),
		energy: energy,
		power:  power,
	}
}

func (m *MiddleMiner) Run(ctx context.Context) <-chan Coal {
	transferPoint := make(chan Coal)

	go func() {
		defer func() {
			close(transferPoint)
			fmt.Println("Middle miner stopped:", m.id)
		}()

		fmt.Println("Middle miner started:", m.id)

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(m.interval):
			}
			select {
			case <-ctx.Done():
				return
			case transferPoint <- Coal(m.power.Load()):
				if new := m.energy.Add(-1); new <= 0 {
					return
				}
			}
		}
	}()
	return transferPoint
}

func (m *MiddleMiner) Info() MinerInfo {
	return MinerInfo{
		ID:          m.id,
		MinerType:   "Middle",
		MinerEnergy: m.energy.Load(),
		MinerPower:  m.power.Load(),
	}
}
