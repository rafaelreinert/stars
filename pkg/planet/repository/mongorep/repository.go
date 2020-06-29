package mongorep

import (
	"context"

	"github.com/rafaelreinert/stars/pkg/planet"
	"github.com/rafaelreinert/stars/pkg/planet/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type planetMongoRepositoryImpl struct {
	Collection *mongo.Collection
}

// NewMongoRepository creates an Repository instace to maniputate planets on MongoDB
func NewMongoRepository(db *mongo.Database) repository.PlanetRepository {
	return planetMongoRepositoryImpl{Collection: db.Collection("planet")}
}

// Create a new planet on Mongo
func (r planetMongoRepositoryImpl) Create(ctx context.Context, p planet.Planet) (planet.Planet, error) {
	model := planetMongoModel{
		ID:      primitive.NewObjectID(),
		Name:    p.Name,
		Climate: p.Climate,
		Terrain: p.Terrain,
	}
	result, err := r.Collection.InsertOne(ctx, model)
	if err != nil {
		return planet.Planet{}, err
	}
	model.ID = result.InsertedID.(primitive.ObjectID)
	return model.ToPlanet(), nil
}

// FindByID finds a planet on Mongo using the id
func (r planetMongoRepositoryImpl) FindByID(ctx context.Context, id string) (planet.Planet, error) {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return planet.Planet{}, err
	}
	result := r.Collection.FindOne(ctx, bson.M{"_id": oID})

	var model planetMongoModel
	err = result.Decode(&model)
	if err != nil {
		return planet.Planet{}, err
	}

	return model.ToPlanet(), nil
}

// FindByName finds a planet on Mongo using the planet name
func (r planetMongoRepositoryImpl) FindByName(ctx context.Context, name string) (planet.Planet, error) {
	result := r.Collection.FindOne(ctx, bson.M{"name": name})

	var model planetMongoModel
	err := result.Decode(&model)
	if err != nil {
		return planet.Planet{}, err
	}

	return model.ToPlanet(), nil
}

// FindAll finds all planets on Mongo
func (r planetMongoRepositoryImpl) FindAll(ctx context.Context) ([]planet.Planet, error) {
	result, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var models []planetMongoModel
	err = result.All(ctx, &models)
	if err != nil {
		return nil, err
	}

	planets := make([]planet.Planet, len(models))
	for i, m := range models {
		planets[i] = m.ToPlanet()
	}
	return planets, nil
}

// Update a planet on mongo
func (r planetMongoRepositoryImpl) Update(ctx context.Context, p planet.Planet) (planet.Planet, error) {
	oID, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		return planet.Planet{}, err
	}
	model := planetMongoModel{
		ID:      oID,
		Name:    p.Name,
		Climate: p.Climate,
		Terrain: p.Terrain,
	}

	opts := options.Update().SetUpsert(true)
	_, err = r.Collection.UpdateOne(ctx, bson.M{"_id": model.ID}, bson.D{{"$set", model}}, opts)
	if err != nil {
		return planet.Planet{}, err
	}

	return model.ToPlanet(), nil
}

// Delete a planet on mongo
func (r planetMongoRepositoryImpl) Delete(ctx context.Context, id string) error {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Collection.DeleteOne(ctx, bson.M{"_id": oID})
	return err
}
