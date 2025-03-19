package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/memdb"
	"GoNews/pkg/storage/mongo"
	"GoNews/pkg/storage/postgres"
	"fmt"
	"log"
	"net/http"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера.
	var srv server

	// Создаём объекты баз данных.
	//
	// БД в памяти.
	db := memdb.New()

	// Реляционная БД PostgreSQL.
	db2, err := postgres.New("postgres://postgres:admin@192.168.1.165/posts")
	if err != nil {
		log.Fatal(err)
	}

	// Документная БД MongoDB.
	db3, err := mongo.New("mongodb://192.168.1.165:27017/")
	if err != nil {
		log.Fatal(err)
	}

	_, _, _ = db, db2, db3

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db3

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Получаем список постов
	posts, err := db2.Posts()
	if err != nil {
		log.Fatalf("Ошибка при получении постов: %v", err)
	}

	// Выводим посты в лог
	for _, post := range posts {
		fmt.Printf("ID: %d, AuthorID: %d, Title: %s, Content: %s, CreatedAt: %d",
			post.ID, post.AuthorID, post.Title, post.Content, post.CreatedAt)
	}

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	http.ListenAndServe(":8080", srv.api.Router())
}
