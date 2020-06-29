package repository

import (
	"context"

	"github.com/rafaelreinert/stars/pkg/planet"
)

// PlanetRepository is the interface used to access the CRUD methods on database
type PlanetRepository interface {
	Create(ctx context.Context, p planet.Planet) (planet.Planet, error)
	FindByID(ctx context.Context, id string) (planet.Planet, error)
	FindByName(ctx context.Context, name string) (planet.Planet, error)
	FindAll(ctx context.Context) ([]planet.Planet, error)
	Update(ctx context.Context, p planet.Planet) (planet.Planet, error)
	Delete(ctx context.Context, id string) error
}

// PlanetFinder is the interface used to access the Finder methods on database
type PlanetFinder interface {
	FindByID(ctx context.Context, id string) (planet.Planet, error)
	FindByName(ctx context.Context, name string) (planet.Planet, error)
	FindAll(ctx context.Context) ([]planet.Planet, error)
}
