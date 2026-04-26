package models

// Place represents a normalized place
// ready to be inserted into the database.
// This is your internal shape — not the
// raw shape from any external API.
type Place struct {
	CityID       int
	CategoryID   int
	Name         string
	Subcategory  string
	Address      string
	Neighborhood string
	Latitude     float64
	Longitude    float64
	Description  string
	Website      string
	Phone        string
	Source       string
	SourceID     string
}
