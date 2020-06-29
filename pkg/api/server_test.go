package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/rafaelreinert/stars/pkg/config"
	"github.com/rafaelreinert/stars/pkg/planet"
	"github.com/rafaelreinert/stars/pkg/planet/repository/mongorep"
	"github.com/rafaelreinert/stars/pkg/swapi"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestServer(t *testing.T) {
	mongoClient, _ := initMongoClient()
	defer mongoClient.Disconnect(context.Background())
	mongoClient.Database("starwars").Drop(context.Background())
	defer mongoClient.Database("starwars").Drop(context.Background())

	swapiServer := initTestSWAPIServer()
	defer swapiServer.Close()

	s := Server{
		PlanetRepository: mongorep.NewMongoRepository(mongoClient.Database("starwars")),
		CountRetriever:   swapi.SWAPI{APIURL: swapiServer.URL},
		Cfg:              config.Config{Port: 8080},
	}

	go s.ListenAndServe()

	t.Run("A=1,GetEmptyList", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/planets")
		assert.NoError(t, err)
		var planets []planet.Planet
		err = json.NewDecoder(resp.Body).Decode(&planets)
		assert.NoError(t, err)
		assert.Empty(t, planets)
	})

	t.Run("A=2,CreateAPlanet", func(t *testing.T) {
		planetToCreate := planet.Planet{
			Name:    "Tatooine",
			Climate: "arid",
			Terrain: "desert",
		}
		planetJSON, _ := json.Marshal(planetToCreate)
		resp, err := http.Post("http://localhost:8080/planets", "application/json", bytes.NewReader(planetJSON))
		assert.NoError(t, err)
		var planetResponse planet.Planet
		err = json.NewDecoder(resp.Body).Decode(&planetResponse)
		assert.NoError(t, err)
		assert.NotNil(t, planetResponse.Name)
		assert.Equal(t, planetToCreate.Name, planetResponse.Name)
		assert.Equal(t, planetToCreate.Climate, planetResponse.Climate)
		assert.Equal(t, planetToCreate.Terrain, planetResponse.Terrain)
	})

	t.Run("A=3,CreateAPlanetWithPut", func(t *testing.T) {
		planetToCreate := planet.Planet{
			ID:      "507f1f77bcf86cd799439011",
			Name:    "Yavin",
			Climate: "temperate, tropical",
			Terrain: "jungle, rainforests",
		}
		planetJSON, _ := json.Marshal(planetToCreate)
		req, _ := http.NewRequest(http.MethodPut, "http://localhost:8080/planets/507f1f77bcf86cd799439011", bytes.NewReader(planetJSON))
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		var planetResponse planet.Planet
		err = json.NewDecoder(resp.Body).Decode(&planetResponse)
		assert.NoError(t, err)
		assert.NotNil(t, planetResponse.Name)
		assert.Equal(t, planetToCreate.Name, planetResponse.Name)
		assert.Equal(t, planetToCreate.Climate, planetResponse.Climate)
		assert.Equal(t, planetToCreate.Terrain, planetResponse.Terrain)
	})

	t.Run("A=4,GetFilledList", func(t *testing.T) {

		planetOne := planet.Planet{
			ID:                          "507f1f77bcf86cd799439011",
			Name:                        "Yavin",
			Climate:                     "temperate, tropical",
			Terrain:                     "jungle, rainforests",
			NumberOfAppearancesOnMovies: 1,
		}

		resp, err := http.Get("http://localhost:8080/planets")
		assert.NoError(t, err)
		var planets []planet.Planet
		err = json.NewDecoder(resp.Body).Decode(&planets)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(planets))
		assert.Contains(t, planets, planetOne)
	})

	t.Run("A=5,GetPlanetByID", func(t *testing.T) {

		planetOne := planet.Planet{
			ID:                          "507f1f77bcf86cd799439011",
			Name:                        "Yavin",
			Climate:                     "temperate, tropical",
			Terrain:                     "jungle, rainforests",
			NumberOfAppearancesOnMovies: 1,
		}

		resp, err := http.Get("http://localhost:8080/planets/507f1f77bcf86cd799439011")
		assert.NoError(t, err)
		var planet planet.Planet
		err = json.NewDecoder(resp.Body).Decode(&planet)
		assert.NoError(t, err)
		assert.Equal(t, planetOne, planet)
	})

	t.Run("A=6,GetPlanetByName", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/planets?name=Tatooine")
		assert.NoError(t, err)
		var planet planet.Planet
		err = json.NewDecoder(resp.Body).Decode(&planet)
		assert.NoError(t, err)
		assert.Equal(t, "Tatooine", planet.Name)
		assert.Equal(t, "arid", planet.Climate)
		assert.Equal(t, "desert", planet.Terrain)
		assert.Equal(t, 5, planet.NumberOfAppearancesOnMovies)
	})

	t.Run("A=7,UpdateAPlanet", func(t *testing.T) {
		planetToUpdate := planet.Planet{
			ID:      "507f1f77bcf86cd799439011",
			Name:    "Yavin IV",
			Climate: "temperate, tropical",
			Terrain: "jungle, rainforests",
		}
		planetJSON, _ := json.Marshal(planetToUpdate)
		req, _ := http.NewRequest(http.MethodPut, "http://localhost:8080/planets/507f1f77bcf86cd799439011", bytes.NewReader(planetJSON))
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		var planetResponse planet.Planet
		err = json.NewDecoder(resp.Body).Decode(&planetResponse)
		assert.NoError(t, err)
		assert.NotNil(t, planetResponse.Name)
		assert.Equal(t, planetToUpdate.Name, planetResponse.Name)
		assert.Equal(t, planetToUpdate.Climate, planetResponse.Climate)
		assert.Equal(t, planetToUpdate.Terrain, planetResponse.Terrain)
	})

	t.Run("A=8,GetUpdatedPlanet", func(t *testing.T) {

		planetOne := planet.Planet{
			ID:                          "507f1f77bcf86cd799439011",
			Name:                        "Yavin IV",
			Climate:                     "temperate, tropical",
			Terrain:                     "jungle, rainforests",
			NumberOfAppearancesOnMovies: 2,
		}

		resp, err := http.Get("http://localhost:8080/planets?name=" + url.QueryEscape("Yavin IV"))
		assert.NoError(t, err)
		var planet planet.Planet
		err = json.NewDecoder(resp.Body).Decode(&planet)
		assert.NoError(t, err)
		assert.Equal(t, planetOne, planet)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("A=9,DeleteAPlanet", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "http://localhost:8080/planets/507f1f77bcf86cd799439011", nil)
		_, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
	})

	t.Run("A=10,GetDeletedPlanetByID", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/planets/507f1f77bcf86cd799439011")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("A=11,GetDeletedPlanetByName", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/planets?name=" + url.QueryEscape("Yavin IV"))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("A=12,CreateAPlanetWithNoSWAPIInformations", func(t *testing.T) {
		planetToCreate := planet.Planet{
			ID:      "507f1f77bcf86cd799439019",
			Name:    "Yavin Not",
			Climate: "temperate, tropical",
			Terrain: "jungle, rainforests",
		}
		planetJSON, _ := json.Marshal(planetToCreate)
		req, _ := http.NewRequest(http.MethodPut, "http://localhost:8080/planets/507f1f77bcf86cd799439019", bytes.NewReader(planetJSON))
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		var planetResponse planet.Planet
		err = json.NewDecoder(resp.Body).Decode(&planetResponse)
		assert.NoError(t, err)
		assert.NotNil(t, planetResponse.Name)
		assert.Equal(t, planetToCreate.Name, planetResponse.Name)
		assert.Equal(t, planetToCreate.Climate, planetResponse.Climate)
		assert.Equal(t, planetToCreate.Terrain, planetResponse.Terrain)
	})

	t.Run("A=13,GetNoSWAPIInformationsPlanet", func(t *testing.T) {

		planetOne := planet.Planet{
			ID:                          "507f1f77bcf86cd799439019",
			Name:                        "Yavin Not",
			Climate:                     "temperate, tropical",
			Terrain:                     "jungle, rainforests",
			NumberOfAppearancesOnMovies: 0,
		}

		resp, err := http.Get("http://localhost:8080/planets?name=" + url.QueryEscape("Yavin Not"))
		assert.NoError(t, err)
		var planet planet.Planet
		err = json.NewDecoder(resp.Body).Decode(&planet)
		assert.NoError(t, err)
		assert.Equal(t, planetOne, planet)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func initMongoClient() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil

}

func initTestSWAPIServer() *httptest.Server {
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
			} else if r.URL.RawQuery == "search=Yavin" || r.URL.RawQuery == "search=Yavin+IV" {
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
			} else {
				fmt.Fprint(w, `{ 
					"count": 0,
					"next": null,
					"previous": null,
					"results": []
				}`)
			}
		}
	}))

}
