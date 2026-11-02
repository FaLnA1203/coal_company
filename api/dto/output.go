package dto

import (
	"coal_company/company"
	"coal_company/company/equipment"
	"coal_company/company/miners"

	"github.com/google/uuid"
)

type OutputDTOMiner struct {
	ID          uuid.UUID
	MinerType   miners.MinersType
	MinerEnergy int64
	MinerPower  int64
}

func NewOutputDTOMiner(miner miners.MinerInfo) *OutputDTOMiner {
	return &OutputDTOMiner{
		ID:          miner.ID,
		MinerType:   miner.MinerType,
		MinerEnergy: miner.MinerEnergy,
		MinerPower:  miner.MinerPower,
	}
}

type DTOPriceMiner struct {
	LittleMinerPrice int64 `json:"little_miner_price"`
	MiddleMinerPrice int64 `json:"middle_miner_price"`
	StrongMinerPrice int64 `json:"strong_miner_price"`
}

func NewDTOPriceMiner() *DTOPriceMiner {
	return &DTOPriceMiner{
		LittleMinerPrice: 5,
		MiddleMinerPrice: 50,
		StrongMinerPrice: 450,
	}
}

type EquipmentPricesDTO struct {
	PickaxePrice     int64 `json:"pickaxe_price"`
	VentilationPrice int64 `json:"ventilation_price"`
	TrolleysPrice    int64 `json:"trolleys_price"`
}

func NewEquipmentPricesDTO() *EquipmentPricesDTO {
	return &EquipmentPricesDTO{
		PickaxePrice:     3000,
		VentilationPrice: 15000,
		TrolleysPrice:    50000,
	}
}

type EquipmentDTO struct {
	Pickaxe     bool `json:"pickaxe"`
	Ventilation bool `json:"ventilation"`
	Trolleys    bool `json:"trolleys"`
}

func NewEquipmentDTO(eq equipment.Equipment) *EquipmentDTO {
	return &EquipmentDTO{
		Pickaxe:     eq.PickaxesPurchased(),
		Ventilation: eq.VentilationPurchased(),
		Trolleys:    eq.TrolleysPurchased(),
	}
}

type CompanyStatisticsDTO struct {
	Balance               int64
	TotalEarned           int64
	TotalMinersStatistics map[miners.MinersType]int
	TotalTime             string
}

func NewCompanyStatisticsDTO(cm company.CompanyStatistics) *CompanyStatisticsDTO {
	return &CompanyStatisticsDTO{
		Balance:               cm.Ballance(),
		TotalEarned:           cm.TotalEarned(),
		TotalMinersStatistics: cm.TotalMinersStatistics(),
		TotalTime:             cm.TimeToComple(),
	}
}
