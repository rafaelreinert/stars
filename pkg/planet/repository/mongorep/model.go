package mongorep

import (
	"github.com/rafaelreinert/stars/pkg/planet"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlanetMongoModel struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name,omitempty"`
	Climate string             `bson:"climate,omitempty"`
	Terrain string             `bson:"terrain,omitempty"`
}

func (p PlanetMongoModel) ToPlanet() planet.Planet {
	return planet.Planet{
		ID:      p.ID.Hex(),
		Name:    p.Name,
		Climate: p.Climate,
		Terrain: p.Terrain,
	}
}
