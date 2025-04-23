package mongo

import (
	"7-solutions-test-backend/internal/core/user"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongoRepo struct {
	col *mongo.Collection
}

func NewUserMongoRepo(db *mongo.Database) *UserMongoRepo {
	return &UserMongoRepo{col: db.Collection("users")}
}

func (r *UserMongoRepo) Create(ctx context.Context, u *user.User) error {
	_, err := r.col.InsertOne(ctx, u)
	return err
}

func (r *UserMongoRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	var u user.User
	err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&u)
	return &u, err
}

func (r *UserMongoRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &u, err
}

func (r *UserMongoRepo) List(ctx context.Context) ([]*user.User, error) {
	cursor, _ := r.col.Find(ctx, bson.M{})
	var users []*user.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserMongoRepo) Update(ctx context.Context, u *user.User) error {
	oid, _ := primitive.ObjectIDFromHex(u.ID)
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"name": u.Name, "email": u.Email}})
	return err
}

func (r *UserMongoRepo) Delete(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *UserMongoRepo) Count(ctx context.Context) (int64, error) {
	return r.col.CountDocuments(ctx, bson.M{})
}
