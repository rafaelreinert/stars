package repository

import (
	"context"

	"github.com/rafaelreinert/starts/pkg/planet"
)

type PlanetRepository interface {
	Create(ctx context.Context, p planet.Planet) (planet.Planet, error)
	FindByID(ctx context.Context, id string) (planet.Planet, error)
	FindByName(ctx context.Context, name string) (planet.Planet, error)
	FindAll(ctx context.Context) ([]planet.Planet, error)
	Update(ctx context.Context, p planet.Planet) (planet.Planet, error)
	Delete(ctx context.Context, id string) error
}

type PlanetFinder interface {
	FindByID(ctx context.Context, id string) (planet.Planet, error)
	FindByName(ctx context.Context, name string) (planet.Planet, error)
	FindAll(ctx context.Context) ([]planet.Planet, error)
}
