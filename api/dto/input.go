package dto

import (
	"coal_company/company/equipment"
	"coal_company/company/miners"
)

type InputDTOMiner struct {
	MinerType miners.MinersType `json:"miner_type"`
}

type InputDTOEquipment struct {
	EquipmentType equipment.EquipmentType `json:"equipment_type"`
}
