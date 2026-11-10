package plan

var DemoPlans = []Plan{
	{
		WarehouseCode: "LYS-01", Width: 600, Height: 300,
		Zones: []Zone{
			{Label: "Réception", X: 20, Y: 20, Width: 160, Height: 120},
			{Label: "Stockage", X: 200, Y: 20, Width: 380, Height: 260},
			{Label: "Expédition", X: 20, Y: 160, Width: 160, Height: 120},
		},
	},
	{
		WarehouseCode: "NTE-02", Width: 500, Height: 300,
		Zones: []Zone{
			{Label: "Quai", X: 20, Y: 20, Width: 460, Height: 80},
			{Label: "Stockage froid", X: 20, Y: 120, Width: 220, Height: 160},
			{Label: "Stockage sec", X: 260, Y: 120, Width: 220, Height: 160},
		},
	},
	{
		WarehouseCode: "LIL-03", Width: 400, Height: 250,
		Zones: []Zone{
			{Label: "Réception", X: 20, Y: 20, Width: 170, Height: 90},
			{Label: "Préparation", X: 210, Y: 20, Width: 170, Height: 90},
			{Label: "Stockage", X: 20, Y: 130, Width: 360, Height: 100},
		},
	},
}
