CREATE TABLE IF NOT EXISTS transaction (
                                           id SERIAL PRIMARY KEY,
                                           user_name VARCHAR(255) NOT NULL,
                                           amount FLOAT NOT NULL,
                                           created_at TIMESTAMP NOT NULL  DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transaction_item (
                                                id SERIAL PRIMARY KEY,
                                                transaction_id INTEGER NOT NULL,
                                                book_title VARCHAR(255) NOT NULL,
                                                book_author VARCHAR(255) NOT NULL,
                                                price FLOAT NOT NULL,
                                                FOREIGN KEY (transaction_id) REFERENCES transaction(id) ON DELETE CASCADE
);