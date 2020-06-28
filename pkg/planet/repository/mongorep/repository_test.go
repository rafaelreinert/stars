package mongorep

import (
	"context"
	"testing"
	"time"

	"github.com/rafaelreinert/starts/pkg/planet"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreate(t *testing.T) {
	client, err := ConnectMongoClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	defer client.Database("starwars").Drop(context.Background())

	ctx := context.Background()
	repo := NewMongoRepository(client.Database("starwars"))
	planetToCreate := planet.Planet{
		Name:    "Tatooine",
		Climate: "arid",
		Terrain: "desert",
	}

	planetCreated, err := repo.Create(ctx, planetToCreate)

	assert.NoError(t, err)
	assert.NotNil(t, planetCreated.ID, "The id should not be nil.")
	assert.Equal(t, planetToCreate.Name, planetCreated.Name, "The names should be equals.")
	assert.Equal(t, planetToCreate.Climate, planetCreated.Climate, "The climates should be equals.")
	assert.Equal(t, planetToCreate.Terrain, planetCreated.Terrain, "The terrains should be equals.")
}

func TestUpdate(t *testing.T) {
	client, err := ConnectMongoClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	defer client.Database("starwars").Drop(context.Background())

	ctx := context.Background()
	repo := NewMongoRepository(client.Database("starwars"))
	planetToCreate := planet.Planet{
		Name:    "Tatooine",
		Climate: "arid",
		Terrain: "desert",
	}

	planetCreated, err := repo.Create(ctx, planetToCreate)

	planetToUpdate := planetCreated
	planetToUpdate.Climate = "cold"
	_, err = repo.Update(ctx, planetToUpdate)
	planetUpdated, _ := repo.FindByID(ctx, planetToUpdate.ID)

	assert.NoError(t, err)
	assert.Equal(t, planetCreated.ID, planetUpdated.ID, "The IDs should be equals.")
	assert.Equal(t, planetCreated.Name, planetUpdated.Name, "The names should be equals.")
	assert.Equal(t, "cold", planetUpdated.Climate, "The climate should be cold.")
	assert.Equal(t, planetCreated.Terrain, planetUpdated.Terrain, "The terrains should be equals.")
}

func TestUpdateWhenDocumentDoesNotExists(t *testing.T) {
	client, err := ConnectMongoClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	defer client.Database("starwars").Drop(context.Background())

	ctx := context.Background()
	repo := NewMongoRepository(client.Database("starwars"))
	planetToUpdate := planet.Planet{
		ID:      primitive.NewObjectID().Hex(),
		Name:    "Tatooine",
		Climate: "arid",
		Terrain: "desert",
	}
	_, err = repo.Update(ctx, planetToUpdate)
	planetCreated, errFind := repo.FindByID(ctx, planetToUpdate.ID)

	assert.NoError(t, err)
	assert.NoError(t, errFind)

	assert.Equal(t, planetToUpdate.ID, planetCreated.ID, "The IDs should be equals.")
	assert.Equal(t, planetToUpdate.Name, planetCreated.Name, "The names should be equals.")
	assert.Equal(t, planetToUpdate.Climate, planetCreated.Climate, "The climate should be cold.")
	assert.Equal(t, planetToUpdate.Terrain, planetCreated.Terrain, "The terrains should be equals.")
}

func TestUpdateWithAnInvalidId(t *testing.T) {
	client, err := ConnectMongoClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	defer client.Database("starwars").Drop(context.Background())

	ctx := context.Background()
	repo := NewMongoRepository(client.Database("starwars"))
	planetToUpdate := planet.Planet{
		ID:      "asdsdas",
		Name:    "Tatooine",
		Climate: "arid",
		Terrain: "desert",
	}
	_, err = repo.Update(ctx, planetToUpdate)

	assert.Error(t, err)
	assert.Equal(t, "encoding/hex: invalid byte: U+0073 's'", err.Error())
}

func TestDelete(t *testing.T) {
	client, err := ConnectMongoClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	defer client.Database("starwars").Drop(context.Background())

	ctx := context.Background()
	repo := NewMongoRepository(client.Database("starwars"))
	planetToCreate := planet.Planet{
		Name:    "Tatooine",
		Climate: "arid",
		Terrain: "desert",
	}

	planetCreated, err := repo.Create(ctx, planetToCreate)

	err = repo.Delete(ctx, planetCreated.ID)
	_, findErr := repo.FindByID(ctx, planetCreated.ID)

	assert.NoError(t, err)
	assert.Error(t, findErr)
}

func TestDeleteWithAnInvalidId(t *testing.T) {
	client, err := ConnectMongoClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	defer client.Database("starwars").Drop(context.Background())

	ctx := context.Background()
	repo := NewMongoRepository(client.Database("starwars"))

	err = repo.Delete(ctx, "sdfsd")

	assert.Error(t, err)
	assert.Equal(t, "encoding/hex: invalid byte: U+0073 's'", err.Error())
}

func TestFindById(t *testing.T) {
	client, err := ConnectMongoClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	defer client.Database("starwars").Drop(context.Background())

	ctx := context.Background()
	repo := NewMongoRepository(client.Database("starwars"))
	planetToCreate := planet.Planet{
		Name:    "Tatooine",
		Climate: "arid",
		Terrain: "desert",
	}

	planetCreated, _ := repo.Create(ctx, planetToCreate)

	planetFound, err := repo.FindByID(ctx, planetCreated.ID)

	assert.NoError(t, err)
	assert.Equal(t, planetCreated.ID, planetFound.ID, "The IDs should be equals.")
	assert.Equal(t, planetCreated.Name, planetFound.Name, "The names should be equals.")
	assert.Equal(t, planetCreated.Climate, planetFound.Climate, "The climates should be equals.")
	assert.Equal(t, planetCreated.Terrain, planetFound.Terrain, "The terrains should be equals.")
}

func TestFindByIdWhenDocumentDoesNotExists(t *testing.T) {
	client, err := ConnectMongoClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	defer client.Database("starwars").Drop(context.Background())

	ctx := context.Background()
	repo := NewMongoRepository(client.Database("starwars"))

	_, err = repo.FindByID(ctx, primitive.NewObjectID().Hex())

	assert.Error(t, err)
	assert.Equal(t, "mongo: no documents in result", err.Error())
}

func TestFindByIdWithAnInvalidID(t *testing.T) {
	client, err := ConnectMongoClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	defer client.Database("starwars").Drop(context.Background())

	ctx := context.Background()
	repo := NewMongoRepository(client.Database("starwars"))

	_, err = repo.FindByID(ctx, "KJSDFNSDJKNFDJKS")

	assert.Error(t, err)
	assert.Equal(t, "encoding/hex: invalid byte: U+004B 'K'", err.Error())
}

func TestFindByName(t *testing.T) {
	client, err := ConnectMongoClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	defer client.Database("starwars").Drop(context.Background())

	ctx := context.Background()
	repo := NewMongoRepository(client.Database("starwars"))
	planetToCreate := planet.Planet{
		Name:    "Tatooine",
		Climate: "arid",
		Terrain: "desert",
	}

	planetCreated, _ := repo.Create(ctx, planetToCreate)

	planetFound, err := repo.FindByName(ctx, planetCreated.Name)

	assert.NoError(t, err)
	assert.Equal(t, planetCreated.ID, planetFound.ID, "The IDs should be equals.")
	assert.Equal(t, planetCreated.Name, planetFound.Name, "The names should be equals.")
	assert.Equal(t, planetCreated.Climate, planetFound.Climate, "The climates should be equals.")
	assert.Equal(t, planetCreated.Terrain, planetFound.Terrain, "The terrains should be equals.")
}

func TestFindByNameWhenDocumentDoesNotExists(t *testing.T) {
	client, err := ConnectMongoClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	defer client.Database("starwars").Drop(context.Background())

	ctx := context.Background()
	repo := NewMongoRepository(client.Database("starwars"))

	_, err = repo.FindByName(ctx, "Pluto")

	assert.Error(t, err)
}

func TestFindAll(t *testing.T) {
	client, err := ConnectMongoClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	defer client.Database("starwars").Drop(context.Background())

	ctx := context.Background()
	repo := NewMongoRepository(client.Database("starwars"))
	planetOneToCreate := planet.Planet{
		Name:    "Tatooine",
		Climate: "arid",
		Terrain: "desert",
	}

	planetOneCreated, _ := repo.Create(ctx, planetOneToCreate)

	planetTwoToCreate := planet.Planet{
		Name:    "Tatooine 2",
		Climate: "arid",
		Terrain: "desert",
	}

	planetTwoCreated, _ := repo.Create(ctx, planetTwoToCreate)

	planets, err := repo.FindAll(ctx)

	assert.NoError(t, err)
	assert.Contains(t, planets, planetOneCreated)
	assert.Contains(t, planets, planetTwoCreated)
}

func TestFindAllEmpty(t *testing.T) {
	client, err := ConnectMongoClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	defer client.Database("starwars").Drop(context.Background())

	ctx := context.Background()
	repo := NewMongoRepository(client.Database("starwars"))

	planets, err := repo.FindAll(ctx)

	assert.NoError(t, err)
	assert.Empty(t, planets)
}

func ConnectMongoClient() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	client.Database("starwars").Drop(context.Background())
	return client, nil

}
