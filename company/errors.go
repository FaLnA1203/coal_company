package company

import "errors"

var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrUnknokwnMinerType = errors.New("unkown miner type")
var ErrUnknokwnEquipmentType = errors.New("unkown equipment type")
var ErrNotAllEquipmentPurchased = errors.New("not all equipment purchased")
