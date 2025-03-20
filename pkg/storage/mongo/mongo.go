package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func New(url string) (*Storage, error) {
	// Опции подключения
	clientOptions := options.Client().ApplyURI(url)

	// Создаём клиент MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к MongoDB: %w", err)
	}

	// Проверяем подключение
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("MongoDB недоступна: %w", err)
	}

	// Выбираем базу данных и коллекцию
	db := client.Database("data")            // Имя базы данных
	collection := db.Collection("languages") // Имя коллекции

	return &Storage{
		client:     client,
		collection: collection,
	}, nil
}

// Posts получает все публикации из MongoDB.
func (s *Storage) Posts() ([]storage.Post, error) {
	// Создаём контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Выполняем запрос
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к MongoDB: %w", err)
	}
	defer cursor.Close(ctx)

	// Читаем результаты
	var posts []storage.Post
	for cursor.Next(ctx) {
		var post storage.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, fmt.Errorf("ошибка декодирования документа: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// AddPost добавляет новую публикацию.
func (s *Storage) AddPost(post storage.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.collection.InsertOne(ctx, post)
	if err != nil {
		return fmt.Errorf("ошибка добавления публикации в MongoDB: %w", err)
	}

	return nil
}

// UpdatePost обновляет публикацию.
func (s *Storage) UpdatePost(post storage.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": post.ID}
	update := bson.M{"$set": post}

	_, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("ошибка обновления публикации: %w", err)
	}

	return nil
}

// DeletePost удаляет публикацию по ID.
func (s *Storage) DeletePost(post storage.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": post.ID}

	_, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("ошибка удаления публикации: %w", err)
	}

	return nil
}
