CREATE TABLE users(
    id uuid PRIMARY KEY,
    username varchar(32) NOT NULL UNIQUE, 
    password_hash varchar(64) NOT NULL
);

CREATE TABLE films(
    title varchar(32) PRIMARY KEY,
    cover_path varchar(64) not null
);

CREATE TABLE reviews(
    id uuid PRIMARY KEY,
    author_id uuid REFERENCES users(id) NOT NULL,
    film_title varchar(32) REFERENCES films(title) NOT NULL,
    rating int NOT NULL,
    body text NOT NULL
    CHECK (rating >= 1 AND rating <= 10)
);

CREATE TABLE tokens(
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) NOT NULL,
    token varchar(64) NOT NULL,
    exp_date Timestamp NOT NULL
);