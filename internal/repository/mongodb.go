package repository

import (
	"context"
	"errors"

	"github.com/0xSangeet/user-api/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoDBRepo struct {
	col *mongo.Collection
}

func NewMongoDBRepo(col *mongo.Collection) *MongoDBRepo {
	return &MongoDBRepo{
		col: col,
	}
}

func (m *MongoDBRepo) Create(ctx context.Context, u *domain.User) error {
	_, err := m.col.InsertOne(ctx, u)

	if mongo.IsDuplicateKeyError(err) {
		return domain.ErrUserAlreadyExists
	}

	return err
}

func (m *MongoDBRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	err := m.col.FindOne(ctx, bson.M{"_id": id}).Decode(&u)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (m *MongoDBRepo) GetAll(ctx context.Context) ([]domain.User, error) {
	cursor, err := m.col.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	results := make([]domain.User, 0)
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *MongoDBRepo) Delete(ctx context.Context, id string) error {
	res, err := m.col.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
