/* Схема для информационной системы отслеживания выполнения задач */

DROP TABLE IF EXISTS tasks_labels, task, users, labels;

CREATE TABLE IF NOT EXIST users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXIST labels (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXIST tasks (
    id SERIAL PRIMARY KEY,
    opened BIGINT,
    closed BIGINT,
    author_id INTEGER REFERENCES users (id),
    assigned_id INTEGER REFERENCES users (id),
    title TEXT NOT NULL,
    content TEXT
);

CREATE TABLE IF NOT EXIST tasks_labels (
    task_id INTEGER REFERENCES tasks (id),
    label_id INTEGER REFERENCES labels (id)
);