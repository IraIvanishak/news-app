package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	data_errors "github.com/IraIvanishak/news-app/internal/adapter/http/errors"
	"github.com/IraIvanishak/news-app/internal/core/domain"
	"github.com/IraIvanishak/news-app/internal/core/ports"
)

type PostRepository struct {
	collection *mongo.Collection
}

func NewPostRepository(db *mongo.Database) *PostRepository {
	return &PostRepository{
		collection: db.Collection("posts"),
	}
}

var _ ports.PostRepository = (*PostRepository)(nil)

func (r *PostRepository) Store(ctx context.Context, post *domain.Post) error {
	post.ID = primitive.NewObjectID().Hex()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, post)
	return err
}

func (r *PostRepository) Find(ctx context.Context, id string) (*domain.Post, error) {
	var post domain.Post
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&post); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, data_errors.ErrPostNotFound
		}
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) List(ctx context.Context, filter domain.ListFilter) ([]*domain.Post, error) {
	findOptions := options.Find().
		SetSkip(filter.Offset).
		SetLimit(filter.Limit).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	query := bson.D{}
	if filter.Query != "" {
		// Case-insensitive search in title or content
		query = bson.D{{Key: "$or", Value: bson.A{
			bson.D{{Key: "title", Value: bson.D{{Key: "$regex", Value: filter.Query}, {Key: "$options", Value: "i"}}}},
			bson.D{{Key: "content", Value: bson.D{{Key: "$regex", Value: filter.Query}, {Key: "$options", Value: "i"}}}},
		}}}
	}

	cursor, err := r.collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*domain.Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.D{})
}

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	post.UpdatedAt = time.Now()

	result, err := r.collection.UpdateByID(ctx, post.ID, bson.M{
		"$set": bson.M{
			"title":     post.Title,
			"content":   post.Content,
			"updatedAt": post.UpdatedAt,
		},
	})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return data_errors.ErrPostNotFound
	}

	return nil
}

func (r *PostRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return data_errors.ErrPostNotFound
	}

	return nil
}
