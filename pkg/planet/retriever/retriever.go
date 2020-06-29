package retriever

import (
	"context"
	"sync"

	"github.com/rafaelreinert/stars/pkg/planet"
	"github.com/rafaelreinert/stars/pkg/planet/repository"
)

// PlanetAppearancesOnMoviesCounter defines the interface to retive the planet appearances on StarWars movies
type PlanetAppearancesOnMoviesCounter interface {
	CountPlanetAppearancesOnMovies(context.Context, string) (int, error)
}

// RetrivePlanet finds a planet on database using id, then fill the planet with the appearances on movies
func RetrivePlanet(ctx context.Context, id string, counter PlanetAppearancesOnMoviesCounter, rep repository.PlanetFinder) (planet.Planet, error) {
	p, err := rep.FindByID(ctx, id)
	if err != nil {
		return planet.Planet{}, err
	}
	return fillNumberOfAppearancesOnMovies(ctx, p, counter)
}

// RetrivePlanetByName finds a planet on database using name, then fills the planet with the appearances on movies
func RetrivePlanetByName(ctx context.Context, name string, counter PlanetAppearancesOnMoviesCounter, rep repository.PlanetFinder) (planet.Planet, error) {
	p, err := rep.FindByName(ctx, name)
	if err != nil {
		return planet.Planet{}, err
	}
	return fillNumberOfAppearancesOnMovies(ctx, p, counter)
}

// RetriveAllPlanets finds all planets on database then fills the planets with the appearances on movies
func RetriveAllPlanets(ctx context.Context, counter PlanetAppearancesOnMoviesCounter, rep repository.PlanetFinder) ([]planet.Planet, error) {
	planets, err := rep.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	planetInputChannel := make(chan *planet.Planet)

	numberOfConnection := len(planets)
	if len(planets) > 10 {
		numberOfConnection = 10
	}

	for i := 0; i < numberOfConnection; i++ {
		go func() {
			for p := range planetInputChannel {
				*p, _ = fillNumberOfAppearancesOnMovies(ctx, *p, counter)
				wg.Done()
			}

		}()
	}
	for i := 0; i < len(planets); i++ {
		wg.Add(1)
		planetInputChannel <- &planets[i]
	}
	wg.Wait()
	close(planetInputChannel)
	return planets, nil
}

func fillNumberOfAppearancesOnMovies(ctx context.Context, p planet.Planet, counter PlanetAppearancesOnMoviesCounter) (planet.Planet, error) {
	n, err := counter.CountPlanetAppearancesOnMovies(ctx, p.Name)
	if err != nil {
		return planet.Planet{}, err
	}
	p.NumberOfAppearancesOnMovies = n
	return p, nil
}
