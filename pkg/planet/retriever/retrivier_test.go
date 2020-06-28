package retriever

import (
	"context"
	"testing"

	"github.com/rafaelreinert/starts/pkg/planet"
	"github.com/stretchr/testify/assert"
)

func TestRetrivePlanet(t *testing.T) {
	p, err := RetrivePlanet(context.Background(), "id", counterMock{}, finderMock{})
	assert.NoError(t, err)
	assert.Equal(t, "Tatooine", p.Name)
	assert.Equal(t, 6, p.NumberOfAppearancesOnMovies)
}

func TestRetrivePlanetByName(t *testing.T) {
	p, err := RetrivePlanetByName(context.Background(), "Alderaan", counterMock{}, finderMock{})
	assert.NoError(t, err)
	assert.Equal(t, "Alderaan", p.Name)
	assert.Equal(t, 1, p.NumberOfAppearancesOnMovies)
}

func TestRetriveAllPlanet(t *testing.T) {
	p, err := RetriveAllPlanets(context.Background(), counterMock{}, finderMock{})
	assert.NoError(t, err)
	assert.Contains(t, p, planet.Planet{
		Name:                        "Tatooine",
		Climate:                     "arid",
		Terrain:                     "desert",
		NumberOfAppearancesOnMovies: 6,
	})
	assert.Contains(t, p, planet.Planet{
		Name:                        "Alderaan",
		Climate:                     "temperate",
		Terrain:                     "grasslands, mountains",
		NumberOfAppearancesOnMovies: 1,
	})
}

func TestRetriveAllPlanetEmpty(t *testing.T) {
	p, err := RetriveAllPlanets(context.Background(), counterMock{}, finderMock{Empty: true})
	assert.NoError(t, err)
	assert.Empty(t, p)
}

type counterMock struct {
}

func (c counterMock) CountPlanetAppearancesOnMovies(ctx context.Context, name string) (int, error) {
	if name == "Tatooine" {
		return 6, nil
	}
	return 1, nil
}

type finderMock struct {
	Empty bool
}

func (r finderMock) FindByID(ctx context.Context, id string) (planet.Planet, error) {

	return planet.Planet{
		Name:    "Tatooine",
		Climate: "arid",
		Terrain: "desert",
	}, nil
}

func (r finderMock) FindByName(ctx context.Context, name string) (planet.Planet, error) {

	return planet.Planet{
		Name:    "Alderaan",
		Climate: "arid",
		Terrain: "desert",
	}, nil
}

func (r finderMock) FindAll(ctx context.Context) ([]planet.Planet, error) {
	if r.Empty {
		return []planet.Planet{}, nil
	}

	planetOneToCreate := planet.Planet{
		Name:    "Tatooine",
		Climate: "arid",
		Terrain: "desert",
	}

	planetTwoToCreate := planet.Planet{
		Name:    "Alderaan",
		Climate: "temperate",
		Terrain: "grasslands, mountains",
	}

	return []planet.Planet{planetOneToCreate, planetTwoToCreate}, nil
}
