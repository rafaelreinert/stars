package planet

// Planet struct represents a planet
type Planet struct {
	ID                          string `json:"id"`
	Name                        string `json:"name"`
	Climate                     string `json:"climate"`
	Terrain                     string `json:"terrain"`
	NumberOfAppearancesOnMovies int    `json:"numberOfAppearancesOnMovies"`
}
