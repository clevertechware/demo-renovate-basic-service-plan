package plan

type Plan struct {
	WarehouseCode string
	Width         int
	Height        int
	Zones         []Zone
}

type Zone struct {
	Label  string
	X      int
	Y      int
	Width  int
	Height int
}
