package swapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountPlanetAppearancesOnMoviesWithAValidPlanet(t *testing.T) {
	ts := initTestServer()
	defer ts.Close()

	count, err := SWAPI{ApiURL: ts.URL}.CountPlanetAppearancesOnMovies(context.Background(), "Tatooine")

	assert.NoError(t, err)
	assert.Equal(t, 5, count, "The Number of Movies should be 5")
}

func TestCountPlanetAppearancesOnMoviesWithAnInvalidPlanet(t *testing.T) {
	ts := initTestServer()
	defer ts.Close()

	count, err := SWAPI{ApiURL: ts.URL}.CountPlanetAppearancesOnMovies(context.Background(), "Tatoo")

	assert.NoError(t, err)
	assert.Equal(t, 0, count, "The Number of Movies should be 0")
}

func TestCountPlanetAppearancesOnMoviesWithAnInexistentPlanet(t *testing.T) {
	ts := initTestServer()
	defer ts.Close()

	count, err := SWAPI{ApiURL: ts.URL}.CountPlanetAppearancesOnMovies(context.Background(), "Pluto")

	assert.NoError(t, err)
	assert.Equal(t, 0, count, "The Number of Movies should be 0")
}

func TestCountPlanetAppearancesOnMoviesWhenAPIReturnMultiplesPlanets(t *testing.T) {
	ts := initTestServer()
	defer ts.Close()

	count, err := SWAPI{ApiURL: ts.URL}.CountPlanetAppearancesOnMovies(context.Background(), "Yavin")

	assert.NoError(t, err)
	assert.Equal(t, 1, count, "The Number of Movies should be 1")
}

func initTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/planets/" {
			if r.URL.RawQuery == "search=Tatooine" || r.URL.RawQuery == "search=Tatoo" {
				fmt.Fprint(w, `{ 
					"count": 1,
					"next": null,
					"previous": null,
					"results": [
						{
							"name": "Tatooine",
							"rotation_period": "23",
							"orbital_period": "304",
							"diameter": "10465",
							"climate": "arid",
							"gravity": "1 standard",
							"terrain": "desert",
							"surface_water": "1",
							"population": "200000",
							"residents": [
								"http://swapi.dev/api/people/1/",
								"http://swapi.dev/api/people/2/",
								"http://swapi.dev/api/people/4/",
								"http://swapi.dev/api/people/6/",
								"http://swapi.dev/api/people/7/",
								"http://swapi.dev/api/people/8/",
								"http://swapi.dev/api/people/9/",
								"http://swapi.dev/api/people/11/",
								"http://swapi.dev/api/people/43/",
								"http://swapi.dev/api/people/62/"
							],
							"films": [
								"http://swapi.dev/api/films/1/",
								"http://swapi.dev/api/films/3/",
								"http://swapi.dev/api/films/4/",
								"http://swapi.dev/api/films/5/",
								"http://swapi.dev/api/films/6/"
							],
							"created": "2014-12-09T13:50:49.641000Z",
							"edited": "2014-12-20T20:58:18.411000Z",
							"url": "http://swapi.dev/api/planets/1/"
						}
					]
				}`)
			} else if r.URL.RawQuery == "search=Pluto" {
				fmt.Fprint(w, `{ 
					"count": 0,
					"next": null,
					"previous": null,
					"results": []
				}`)
			} else if r.URL.RawQuery == "search=Yavin" {
				fmt.Fprint(w, `{ 
					"count": 2,
					"next": null,
					"previous": null,
					"results": [
						{
							"name": "Yavin IV",
							"rotation_period": "24",
							"orbital_period": "4818",
							"diameter": "10200",
							"climate": "temperate, tropical",
							"gravity": "1 standard",
							"terrain": "jungle, rainforests",
							"surface_water": "8",
							"population": "1000",
							"residents": [],
							"films": [
								"http://swapi.dev/api/films/2/",
								"http://swapi.dev/api/films/1/"
							],
							"created": "2014-12-10T11:37:19.144000Z",
							"edited": "2014-12-20T20:58:18.421000Z",
							"url": "http://swapi.dev/api/planets/3/"
						},
						{
							"name": "Yavin",
							"rotation_period": "24",
							"orbital_period": "4818",
							"diameter": "10200",
							"climate": "temperate, tropical",
							"gravity": "1 standard",
							"terrain": "jungle, rainforests",
							"surface_water": "8",
							"population": "1000",
							"residents": [],
							"films": [
								"http://swapi.dev/api/films/1/"
							],
							"created": "2014-12-10T11:37:19.144000Z",
							"edited": "2014-12-20T20:58:18.421000Z",
							"url": "http://swapi.dev/api/planets/3/"
						}
					]
				}`)
			}
		}
	}))

}
