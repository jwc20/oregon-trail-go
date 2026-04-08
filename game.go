package oregontrail

const TotalRequiredMileage = 2040

type Player struct {
	Cash          int
	ShootingLevel int
}

type Inventory struct {
	Oxen          int
	Food          int
	Ammo          int
	Clothing      int
	Miscellaneous int
}

type TripState struct {
	Mileage       int
	TurnNumber    int
	FortAvailable bool
}

type Flags struct {
	Injured          bool
	Ill              bool
	ClearedSouthPass bool
	ClearedBlueMtns  bool
	SouthPassMileage bool
}

type GameState struct {
	Player    Player
	Inventory Inventory
	Trip      TripState
	Flags     Flags
}

//func InitSVT(state *GameState) {
//	state.Trip.FortAvailable = true
//	state.Flags.Injured = false
//	state.Flags.Ill = false
//	state.Flags.ClearedSouthPass = false
//	state.Flags.ClearedBlueMtns = false
//	state.Trip.Mileage = 0
//	state.Flags.SouthPassMileage = false
//	state.Trip.TurnNumber = 0
//}
