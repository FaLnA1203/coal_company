package equipment

type EquipmentType string

const (
	EquipmentTypePickaxe     EquipmentType = "pickaxe"
	EquipmentTypeVentilation EquipmentType = "ventilation"
	EquipmentTypeTrolleys    EquipmentType = "trolleys"
)

type Equipment struct {
	pickaxe     bool
	ventilation bool
	trolleys    bool
}

func NewEquipment() Equipment {
	return Equipment{
		pickaxe:     false,
		ventilation: false,
		trolleys:    false,
	}
}

func (e *Equipment) BuyPickaxe() {
	e.pickaxe = true
	println("buying pickaxe")
}

func (e *Equipment) BuyVentilation() {
	e.ventilation = true
	println("buying ventilation")
}

func (e *Equipment) BuyTrolleys() {
	e.trolleys = true
	println("buying trolleys")
}

func (e *Equipment) PickaxesPurchased() bool {
	return e.pickaxe
}

func (e *Equipment) VentilationPurchased() bool {
	return e.ventilation
}

func (e *Equipment) TrolleysPurchased() bool {
	return e.trolleys
}

func (e *Equipment) AllBought() bool {
	return e.PickaxesPurchased() && e.VentilationPurchased() && e.TrolleysPurchased()
}
