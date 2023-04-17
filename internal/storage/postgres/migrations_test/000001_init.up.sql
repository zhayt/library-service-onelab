CREATE TABLE IF NOT EXISTS "user" (
    id SERIAL PRIMARY KEY,
    fio VARCHAR(70) NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL,
    password char(60) NOT NULL
);

CREATE TABLE IF NOT EXISTS book (
    id serial PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    author VARCHAR(70) NOT NULL
);

CREATE TABLE IF NOT EXISTS book_issue_history (
    id SERIAL PRIMARY KEY,
    book_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    return_date TIMESTAMP,
    FOREIGN KEY (book_id) REFERENCES book (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
);

INSERT INTO "user" (fio, email, password)
VALUES ('Test Fio Test', 'existsemail@mail.ru', 'hashed_password');