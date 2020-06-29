package retriever

import (
	"context"
	"sync"

	"github.com/rafaelreinert/stars/pkg/planet"
	"github.com/rafaelreinert/stars/pkg/planet/repository"
)

type PlanetAppearancesOnMoviesCounter interface {
	CountPlanetAppearancesOnMovies(context.Context, string) (int, error)
}

func RetrivePlanet(ctx context.Context, id string, counter PlanetAppearancesOnMoviesCounter, rep repository.PlanetFinder) (planet.Planet, error) {
	p, err := rep.FindByID(ctx, id)
	if err != nil {
		return planet.Planet{}, err
	}
	return fillNumberOfAppearancesOnMovies(ctx, p, counter)
}

func RetrivePlanetByName(ctx context.Context, name string, counter PlanetAppearancesOnMoviesCounter, rep repository.PlanetFinder) (planet.Planet, error) {
	p, err := rep.FindByName(ctx, name)
	if err != nil {
		return planet.Planet{}, err
	}
	return fillNumberOfAppearancesOnMovies(ctx, p, counter)
}

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
