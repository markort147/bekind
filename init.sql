\c bekind;

CREATE TABLE people
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
);

CREATE TABLE roles
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
);

CREATE TABLE sagas
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
);

CREATE TABLE studios
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
);

CREATE TABLE genres
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
);

CREATE TABLE countries
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
);

CREATE TABLE movies
(
    id    SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    rate  INT          NOT NULL,
    year  INT          NOT NULL,
    seen  INT          NOT NULL,
);

CREATE TABLE people_movies
(
    id        SERIAL PRIMARY KEY,
    person_id INT NOT NULL,
    movie_id  INT NOT NULL,
    role_id   INT NOT NULL,

    FOREIGN KEY (person_id) REFERENCES people (id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
);

CREATE TABLE sagas_movies
(
    id        SERIAL PRIMARY KEY,
    saga_id   INT NOT NULL,
    movie_id  INT NOT NULL,

    FOREIGN KEY (saga_id) REFERENCES sagas (id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE CASCADE,
);

CREATE TABLE studios_movies
(
    id        SERIAL PRIMARY KEY,
    studio_id INT NOT NULL,
    movie_id  INT NOT NULL,

    FOREIGN KEY (studio_id) REFERENCES studios (id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE CASCADE,
);

CREATE TABLE genres_movies
(
    id        SERIAL PRIMARY KEY,
    genre_id  INT NOT NULL,
    movie_id  INT NOT NULL,

    FOREIGN KEY (genre_id) REFERENCES genres (id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE CASCADE,
);

CREATE TABLE countries_movies
(
    id        SERIAL PRIMARY KEY,
    country_id INT NOT NULL,
    movie_id   INT NOT NULL,

    FOREIGN KEY (country_id) REFERENCES countries (id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE CASCADE,
);

INSERT INTO roles (name) VALUES
('director'),
('writer'),
('cinematographer'),
('composer'),
('editor'),
('producer');