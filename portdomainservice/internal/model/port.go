package model

type Port struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	City        string    `db:"city"`
	Country     string    `db:"country"`
	Alias       []string  `db:"alias"`
	Regions     []string  `db:"regions"`
	Coordinates []float64 `db:"coordinates"`
	Province    string    `db:"province"`
	Timezone    string    `db:"timezone"`
	Unlocs      []string  `db:"unlocs"`
	Code        string    `db:"code"`
}
