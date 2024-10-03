package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(connection string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), connection)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Задача.
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Tasks возвращает список задач из БД.
func (s *Storage) TasksByAuthor(authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id, 
			opened, 
			closed, 
			author_id, 
			assigned_id, 
			title, 
			content 
		FROM tasks
		WHERE author_id = $1`,
		authorID)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)
	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// Tasks возвращает список задач из БД с заданным тегом.
func (s *Storage) TasksByTag(tag string) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id, 
			opened, 
			closed, 
			author_id, 
			assigned_id, 
			title, 
			content 
		FROM tasks t
		WHERE EXISTS (
			SELECT * FROM tasks_labels tl 
			WHERE tl.task_id = t.id
			AND tl.label_id IN (
				SELECT l.id FROM labels l
				WHERE l.name = $1
			)
		)`,
		tag)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)
	}

	return tasks, rows.Err()
}

// UpdateTask обновляет задачу в БД.
func (s *Storage) UpdateTask(t Task) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE tasks
		SET opened = $1, closed = $2, author_id = $3, assigned_id = $4, title = $5, content = $6
		WHERE id = $7;
		`,
		t.Opened,
		t.Closed,
		t.AuthorID,
		t.AssignedID,
		t.Title,
		t.Content,
		t.ID,
	)
	return err
}

// NewTask создаёт новую задачу и возвращает её id.
func (s *Storage) NewTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content)
		VALUES ($1, $2) RETURNING id;
		`,
		t.Title,
		t.Content,
	).Scan(&id)
	return id, err
}

// TaskByID возвращает задачу по её id.
func (s *Storage) TaskByID(id int) (Task, error) {
	var t Task
	err := s.db.QueryRow(context.Background(), `
		SELECT id, opened, closed, author_id, assigned_id, title, content
		FROM tasks
		WHERE id = $1;
		`,
		id,
	).Scan(
		&t.ID,
		&t.Opened,
		&t.Closed,
		&t.AuthorID,
		&t.AssignedID,
		&t.Title,
		&t.Content,
	)
	return t, err
}

// DeleteTaskByID удаляет задачу по её id.
func (s *Storage) DeleteTaskByID(id int) error {
	_, err := s.db.Exec(context.Background(), `
		DELETE FROM tasks
		WHERE id = $1;
		`,
		id,
	)
	return err
}

//
