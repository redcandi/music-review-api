CREATE DATABASE IF NOT EXISTS music;
USE music;

DROP VIEW IF EXISTS v_album_avg_rating;
DROP TABLE IF EXISTS album_genres;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS genres;
DROP TABLE IF EXISTS albums;
DROP TABLE IF EXISTS artists;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE artists (
    artist_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    bio TEXT,
    formed_year INT
);

CREATE TABLE albums (
    album_id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    release_date DATE,
    cover_image_url TEXT,
    artist_id INT NOT NULL,
    FOREIGN KEY (artist_id) REFERENCES artists(artist_id) ON DELETE CASCADE
);

CREATE TABLE genres (
    genre_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE comments (
    comment_id INT AUTO_INCREMENT PRIMARY KEY,
    album_id INT NOT NULL,
    user_id INT NOT NULL,
    rating SMALLINT NOT NULL CHECK (rating >= 1 AND rating <= 10),
    comment_text TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (album_id) REFERENCES albums(album_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE album_genres (
    album_id INT NOT NULL,
    genre_id INT NOT NULL,
    
    PRIMARY KEY (album_id, genre_id),
    FOREIGN KEY (album_id) REFERENCES albums(album_id) ON DELETE CASCADE,
    FOREIGN KEY (genre_id) REFERENCES genres(genre_id) ON DELETE CASCADE
);

CREATE VIEW v_album_avg_rating AS
SELECT
    a.album_id,
    a.title,
    ar.name AS artist_name,
    a.cover_image_url,
    COALESCE(AVG(c.rating), 0) AS average_rating,
    COUNT(c.comment_id) AS total_comments
FROM
    albums a
JOIN
    artists ar ON a.artist_id = ar.artist_id
LEFT JOIN
    comments c ON a.album_id = c.album_id
GROUP BY
    a.album_id, a.title, ar.name, a.cover_image_url;

INSERT INTO users (username, email, password_hash) 
VALUES ('anonymous', 'anonymous@app.com', 'dummy_hash_for_anonymous_user');


--dummy data
INSERT INTO artists (name, formed_year) VALUES ('Radiohead', 1985);
INSERT INTO artists (name, formed_year) VALUES ('Pink Floyd', 1965);

INSERT INTO albums (title, release_date, artist_id, cover_image_url) 
VALUES ('OK Computer', '1997-05-21', 1, 'https://placehold.co/300x300/E8DED8/333?text=OK+Computer');
INSERT INTO albums (title, release_date, artist_id, cover_image_url) 
VALUES ('The Dark Side of the Moon', '1973-03-01', 2, 'https://placehold.co/300x300/000/FFF?text=Dark+Side');

INSERT INTO genres (name) VALUES ('Rock'), ('Electronic'), ('Psychedelic');
INSERT INTO album_genres (album_id, genre_id) VALUES (1, 1), (1, 2), (2, 1), (2, 3);

INSERT INTO comments (album_id, user_id, rating, comment_text) 
VALUES (1, 1, 10, 'A masterpiece.');
INSERT INTO comments (album_id, user_id, rating, comment_text) 
VALUES (2, 1, 9, 'Classic album.');
