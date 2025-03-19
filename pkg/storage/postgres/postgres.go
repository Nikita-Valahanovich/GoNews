package postgres

import (
	"GoNews/pkg/storage"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

// NewStorage создаёт новый экземпляр Storage
func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{db: db}
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

func (s *Storage) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT id, author_id, title, content,created_at
		FROM posts
		`)
	if err != nil {
		return nil, err
	}

	defer rows.Close() // закрываем rows после использования

	posts := make([]storage.Post, 0)
	for rows.Next() {
		var post storage.Post
		err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *Storage) AddPost(post storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		INSERT INTO posts(id, author_id, title, content, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
		post.ID, post.AuthorID, post.Title, post.Content, post.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdatePost(post storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE posts
		SET author_id = $1, title = $2, content = $3, created_at = $4
		WHERE id = $5
	`, post.AuthorID, post.Title, post.Content, post.CreatedAt, post.ID)

	return err
}

func (s *Storage) DeletePost(post storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		DELETE FROM posts WHERE id = $1
	`, post.ID)

	return err
}
