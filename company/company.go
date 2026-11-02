package company

import (
	"coal_company/company/equipment"
	"coal_company/company/miners"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Company struct {
	incomeCh    chan miners.Coal
	mtx         sync.RWMutex
	miners      map[miners.MinersType]map[uuid.UUID]miners.Miner
	equipment   equipment.Equipment
	companyCtx  context.Context
	companyStop context.CancelFunc
	statistics  *CompanyStatistics
}

func NewCompany(ctx context.Context) *Company {
	companyCtx, companyStop := context.WithCancel(ctx)
	c := &Company{
		incomeCh:    make(chan miners.Coal),
		mtx:         sync.RWMutex{},
		miners:      make(map[miners.MinersType]map[uuid.UUID]miners.Miner),
		equipment:   equipment.NewEquipment(),
		companyCtx:  companyCtx,
		companyStop: companyStop,
		statistics:  NewCompanyStatistics(),
	}

	go c.baseIncome()
	go c.collectIncome()

	return c
}

func (c *Company) HireMiner(minerType miners.MinersType) (miners.Miner, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	var miner miners.Miner

	switch minerType {
	case "little":
		if c.statistics.ballance.Load() >= 5 {
			miner = miners.NewLittleMiner()
			c.statistics.ballance.Add(-5)
		} else {
			return nil, ErrInsufficientFunds
		}
	case "middle":
		if c.statistics.ballance.Load() >= 50 {
			miner = miners.NewMiddleMiner()
			c.statistics.ballance.Add(-50)
		} else {
			return nil, ErrInsufficientFunds
		}
	case "strong":
		if c.statistics.ballance.Load() >= 450 {
			miner = miners.NewStrongMiner()
			c.statistics.ballance.Add(-450)
		} else {
			return nil, ErrInsufficientFunds
		}
	default:
		return nil, ErrUnknokwnMinerType
	}
	info := miner.Info()
	if c.miners[info.MinerType] == nil {
		c.miners[info.MinerType] = make(map[uuid.UUID]miners.Miner)
	}
	c.miners[info.MinerType][info.ID] = miner

	coalCh := miner.Run(c.companyCtx)

	go func() {
		for v := range coalCh {
			c.incomeCh <- v
		}

		c.mtx.Lock()
		delete(c.miners[info.MinerType], info.ID)

		if len(c.miners[info.MinerType]) == 0 {
			delete(c.miners, info.MinerType)
		}
		c.mtx.Unlock()
	}()

	c.statistics.totalMinersStatistics[minerType]++

	return miner, nil
}

func (c *Company) GetAllMiners() map[miners.MinersType]map[uuid.UUID]miners.Miner {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	tmp := make(map[miners.MinersType]map[uuid.UUID]miners.Miner)
	for minerType, minersMap := range c.miners {
		tmp[minerType] = make(map[uuid.UUID]miners.Miner)
		for minerId, miner := range minersMap {
			tmp[minerType][minerId] = miner
		}
	}
	return tmp
}

func (c *Company) GetMinersByType(minerType miners.MinersType) map[uuid.UUID]miners.Miner {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	tmp := make(map[uuid.UUID]miners.Miner)
	for id, miner := range c.miners[minerType] {
		tmp[id] = miner
	}
	return tmp
}

func (c *Company) BuyEquipment(equipmentType equipment.EquipmentType) (equipment.Equipment, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	switch equipmentType {
	case "pickaxe":
		if !c.equipment.PickaxesPurchased() {
			if c.statistics.ballance.Load() >= 3000 {
				c.equipment.BuyPickaxe()
				c.statistics.ballance.Add(-3000)
			} else {
				return equipment.Equipment{}, ErrInsufficientFunds
			}
		}
	case "ventilation":
		if !c.equipment.VentilationPurchased() {
			if c.statistics.ballance.Load() >= 15000 {
				c.equipment.BuyVentilation()
				c.statistics.ballance.Add(-15000)
			} else {
				return equipment.Equipment{}, ErrInsufficientFunds
			}
		}
	case "trolleys":
		if !c.equipment.TrolleysPurchased() {
			if c.statistics.ballance.Load() >= 50000 {
				c.equipment.BuyTrolleys()
				c.statistics.ballance.Add(-50000)
			} else {
				return equipment.Equipment{}, ErrInsufficientFunds
			}
		}
	default:
		return equipment.Equipment{}, ErrUnknokwnEquipmentType
	}
	return c.equipment, nil
}

func (c *Company) GetEquipment() equipment.Equipment {
	return c.equipment
}

func (c *Company) GetStatistics() CompanyStatistics {
	return *c.statistics
}

func (c *Company) Complete() (CompanyStatistics, error) {
	if !c.equipment.AllBought() {
		return CompanyStatistics{}, ErrNotAllEquipmentPurchased
	}

	c.companyStop()

	now := time.Now()
	c.statistics.completedTime = &now

	return *c.statistics, nil
}

func (c *Company) collectIncome() {
	for {
		select {
		case <-c.companyCtx.Done():
			return
		case income := <-c.incomeCh:
			c.statistics.ballance.Add(int64(income))
			fmt.Println("ballance:", c.statistics.ballance.Load())

			c.statistics.totalEarned.Add(int64(income))
		}
	}
}

func (c *Company) baseIncome() {
	for {
		select {
		case <-c.companyCtx.Done():
			return

		case <-time.After(1 * time.Second):
		}

		select {
		case <-c.companyCtx.Done():
			return
		case c.incomeCh <- 1:
		}
	}
}
