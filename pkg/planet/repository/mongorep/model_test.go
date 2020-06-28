package mongorep

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestToPlanet(t *testing.T) {
	model := PlanetMongoModel{
		ID:      primitive.NewObjectID(),
		Name:    "Tatooine",
		Climate: "arid",
		Terrain: "desert",
	}

	planet := model.ToPlanet()

	assert.Equal(t, model.ID.Hex(), planet.ID, "The ids should be equals.")
	assert.Equal(t, model.Name, planet.Name, "The names should be equals.")
	assert.Equal(t, model.Climate, planet.Climate, "The climates should be equals.")
	assert.Equal(t, model.Terrain, planet.Terrain, "The terrains should be equals.")
}
