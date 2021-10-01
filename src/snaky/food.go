package snaky

// Food
type Food struct {
	cell Cell
}

// NewFood
func NewFood() *Food {
	return &Food{
		cell: Cell{5, 5},
	}
}
