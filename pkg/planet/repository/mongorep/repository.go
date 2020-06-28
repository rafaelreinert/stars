package mongorep

import (
	"context"

	"github.com/rafaelreinert/starts/pkg/planet"
	"github.com/rafaelreinert/starts/pkg/planet/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type planetMongoRepositoryImpl struct {
	Collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) repository.PlanetRepository {
	return planetMongoRepositoryImpl{Collection: db.Collection("planet")}
}

func (r planetMongoRepositoryImpl) Create(ctx context.Context, p planet.Planet) (planet.Planet, error) {
	model := PlanetMongoModel{
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

func (r planetMongoRepositoryImpl) FindByID(ctx context.Context, id string) (planet.Planet, error) {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return planet.Planet{}, err
	}
	result := r.Collection.FindOne(ctx, bson.M{"_id": oID})

	var model PlanetMongoModel
	err = result.Decode(&model)
	if err != nil {
		return planet.Planet{}, err
	}

	return model.ToPlanet(), nil
}

func (r planetMongoRepositoryImpl) FindByName(ctx context.Context, name string) (planet.Planet, error) {
	result := r.Collection.FindOne(ctx, bson.M{"name": name})

	var model PlanetMongoModel
	err := result.Decode(&model)
	if err != nil {
		return planet.Planet{}, err
	}

	return model.ToPlanet(), nil
}

func (r planetMongoRepositoryImpl) FindAll(ctx context.Context) ([]planet.Planet, error) {
	result, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var models []PlanetMongoModel
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

func (r planetMongoRepositoryImpl) Update(ctx context.Context, p planet.Planet) (planet.Planet, error) {
	oID, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		return planet.Planet{}, err
	}
	model := PlanetMongoModel{
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

func (r planetMongoRepositoryImpl) Delete(ctx context.Context, id string) error {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Collection.DeleteOne(ctx, bson.M{"_id": oID})
	return err
}
