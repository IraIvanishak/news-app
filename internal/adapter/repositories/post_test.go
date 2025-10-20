package repositories_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/IraIvanishak/news-app/internal/adapter/repositories"
	"github.com/IraIvanishak/news-app/internal/core/domain"
)

var (
	mongoContainer *dockertest.Resource
	mongoPool      *dockertest.Pool
	mongoClient    *mongo.Client
	testDB         *mongo.Database
)

func TestMain(m *testing.M) {
	var err error
	mongoPool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	mongoContainer, err = mongoPool.RunWithOptions(&dockertest.RunOptions{
		Name:       "mongodb",
		Repository: "mongo",
		Tag:        "latest",
	})
	if err != nil {
		fmt.Printf("Could not start Mongodb: %v \n", err)
		return
	}

	if err = mongoPool.Retry(func() error {
		var err error
		mongoClient, err = mongo.Connect(
			context.Background(),
			options.Client().ApplyURI(fmt.Sprintf("mongodb://localhost:%s", mongoContainer.GetPort("27017/tcp"))),
		)
		if err != nil {
			return err
		}
		return mongoClient.Ping(context.Background(), nil)
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	testDB = mongoClient.Database("testdb")

	code := m.Run()

	if err = mongoPool.Purge(mongoContainer); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}

func TestPostRepository_Store(t *testing.T) {
	repo := repositories.NewPostRepository(testDB)

	post := &domain.Post{
		Title:   "Test Post",
		Content: "This is a test post content",
	}

	ctx := context.Background()
	err := repo.Store(ctx, post)
	assert.NoError(t, err)
	assert.NotZero(t, post.ID)
}

func TestPostRepository_Find(t *testing.T) {
	repo := repositories.NewPostRepository(testDB)

	post := &domain.Post{
		Title:   "Find Test Post",
		Content: "Content for finding test",
	}

	ctx := context.Background()
	err := repo.Store(ctx, post)
	assert.NoError(t, err)

	foundPost, err := repo.Find(ctx, post.ID)
	assert.NoError(t, err)
	assert.Equal(t, post.Title, foundPost.Title)
	assert.Equal(t, post.Content, foundPost.Content)
}

func TestPostRepository_List(t *testing.T) {
	repo := repositories.NewPostRepository(testDB)

	ctx := context.Background()
	posts := []*domain.Post{
		{
			Title:     "Post 1",
			Content:   "Content 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Title:     "Post 2",
			Content:   "Content 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, post := range posts {
		err := repo.Store(ctx, post)
		assert.NoError(t, err)
	}

	filter := domain.ListFilter{
		Limit:  10,
		Offset: 0,
	}
	listedPosts, err := repo.List(ctx, filter)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(listedPosts), 2)
}

func TestPostRepository_Update(t *testing.T) {
	repo := repositories.NewPostRepository(testDB)

	ctx := context.Background()
	post := &domain.Post{
		Title:     "Original Title",
		Content:   "Original Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Store(ctx, post)
	assert.NoError(t, err)

	post.Title = "Updated Title"
	post.Content = "Updated Content"
	err = repo.Update(ctx, post)
	assert.NoError(t, err)

	updatedPost, err := repo.Find(ctx, post.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedPost.Title)
	assert.Equal(t, "Updated Content", updatedPost.Content)
}

func TestPostRepository_Delete(t *testing.T) {
	repo := repositories.NewPostRepository(testDB)

	ctx := context.Background()
	post := &domain.Post{
		Title:     "Post to Delete",
		Content:   "This post will be deleted",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Store(ctx, post)
	assert.NoError(t, err)

	err = repo.Delete(ctx, post.ID)
	assert.NoError(t, err)

	_, err = repo.Find(ctx, post.ID)
	assert.Error(t, err)
}

func TestPostRepository_Count(t *testing.T) {
	repo := repositories.NewPostRepository(testDB)

	ctx := context.Background()
	initialCount, err := repo.Count(ctx)
	assert.NoError(t, err)

	posts := []*domain.Post{
		{
			Title:     "Count Test 1",
			Content:   "Content 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Title:     "Count Test 2",
			Content:   "Content 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, post := range posts {
		err := repo.Store(ctx, post)
		assert.NoError(t, err)
	}

	finalCount, err := repo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, initialCount+int64(len(posts)), finalCount)
}

// Negative test cases
func TestPostRepository_Find_NonExistent(t *testing.T) {
	repo := repositories.NewPostRepository(testDB)

	ctx := context.Background()
	nonExistentID := primitive.NewObjectID().Hex()

	_, err := repo.Find(ctx, nonExistentID)
	assert.Error(t, err)
}

func TestPostRepository_Update_NonExistent(t *testing.T) {
	repo := repositories.NewPostRepository(testDB)

	ctx := context.Background()
	nonExistentPost := &domain.Post{
		ID:        primitive.NewObjectID().Hex(),
		Title:     "Non-existent Post",
		Content:   "This post doesn't exist",
		UpdatedAt: time.Now(),
	}

	err := repo.Update(ctx, nonExistentPost)
	assert.Error(t, err)
}

func TestPostRepository_Update_ZeroID(t *testing.T) {
	repo := repositories.NewPostRepository(testDB)

	ctx := context.Background()
	post := &domain.Post{
		ID:        primitive.NilObjectID.Hex(),
		Title:     "Post with zero ID",
		Content:   "This post has a zero ID",
		UpdatedAt: time.Now(),
	}

	err := repo.Update(ctx, post)
	assert.Error(t, err)
}

func TestPostRepository_Delete_NonExistent(t *testing.T) {
	repo := repositories.NewPostRepository(testDB)

	ctx := context.Background()
	nonExistentID := primitive.NewObjectID().Hex()

	err := repo.Delete(ctx, nonExistentID)
	assert.Error(t, err)
}
func TestPostRepository_Store_VeryLongContent(t *testing.T) {
	repo := repositories.NewPostRepository(testDB)

	ctx := context.Background()

	longContent := strings.Repeat("This is a very long content. ", 1000)

	post := &domain.Post{
		Title:     "Post with very long content",
		Content:   longContent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Store(ctx, post)
	assert.NoError(t, err)
	assert.NotZero(t, post.ID)

	foundPost, err := repo.Find(ctx, post.ID)
	assert.NoError(t, err)
	assert.Equal(t, longContent, foundPost.Content)
}

func TestPostRepository_Store_VeryLongTitle(t *testing.T) {
	repo := repositories.NewPostRepository(testDB)

	ctx := context.Background()

	longTitle := strings.Repeat("Very long title. ", 100)

	post := &domain.Post{
		Title:     longTitle,
		Content:   "Normal content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Store(ctx, post)
	assert.NoError(t, err)
	assert.NotZero(t, post.ID)

	foundPost, err := repo.Find(ctx, post.ID)
	assert.NoError(t, err)
	assert.Equal(t, longTitle, foundPost.Title)
}
