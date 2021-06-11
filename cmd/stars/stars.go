package main

import (
	"context"
	"log"

	"github.com/rafaelreinert/stars/pkg/api"
	"github.com/rafaelreinert/stars/pkg/config"
	"github.com/rafaelreinert/stars/pkg/planet/repository/mongorep"
	"github.com/rafaelreinert/stars/pkg/swapi"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.Println("Initiating stars...")
	log.Println("Initiating Config...")
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Config OK")
	log.Println("Initiating Mongo Client...")
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Mongo Client OK")

	s := api.Server{
		PlanetRepository: mongorep.NewMongoRepository(client.Database("starwars")),
		CountRetriever:   swapi.SWAPI{APIURL: cfg.SWAPIURL},
		Cfg:              cfg,
	}
	log.Println("Stars OK")
	s.ListenAndServe()
}
