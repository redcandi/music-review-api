package api

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"music-review-api/internal/models"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (env *ApiEnv) getOrCreateUser(username string) (int, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		username = "anonymous"
	}

	var userID int
	tx, err := env.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	query := "SELECT user_id FROM users WHERE username = ?"
	err = tx.QueryRow(query, username).Scan(&userID)

	if err == sql.ErrNoRows {
		email := username + "@app.com"
		if username == "anonymous" {
			email = "anonymous@app.com"
		}

		insertQuery := "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)"
		res, err := tx.Exec(insertQuery, username, email, "dummy_hash")
		if err != nil {
			return 0, err
		}
		newID, _ := res.LastInsertId()
		userID = int(newID)
	} else if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return userID, nil
}

func (env *ApiEnv) SignUp(c *gin.Context) {
	var input models.SignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	query := "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)"
	_, err = env.DB.Exec(query, input.Username, input.Email, hash)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email or username already exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (env *ApiEnv) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	query := "SELECT user_id, username, email, password_hash FROM users WHERE email = ?"
	err := env.DB.QueryRow(query, input.Email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !CheckPasswordHash(input.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{Username: user.Username})
}

func (env *ApiEnv) CreateArtist(c *gin.Context) {
	var input models.CreateArtistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO artists (name, bio, formed_year) VALUES (?, ?, ?)"
	res, err := env.DB.Exec(query, input.Name, input.Bio, input.FormedYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create artist"})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"message": "Artist created", "artist_id": id})
}

func (env *ApiEnv) GetArtists(c *gin.Context) {
	rows, err := env.DB.Query("SELECT artist_id, name, bio, formed_year FROM artists")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch artists"})
		return
	}
	defer rows.Close()

	artists := []models.Artist{}
	for rows.Next() {
		var artist models.Artist
		if err := rows.Scan(&artist.ID, &artist.Name, &artist.Bio, &artist.FormedYear); err != nil {
			log.Printf("Error scanning artist: %v", err)
			continue
		}
		artists = append(artists, artist)
	}
	c.JSON(http.StatusOK, artists)
}

func (env *ApiEnv) CreateAlbum(c *gin.Context) {
	var input models.CreateAlbumInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	releaseDate, err := time.Parse("2006-01-02", input.ReleaseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	query := "INSERT INTO albums (title, release_date, cover_image_url, artist_id) VALUES (?, ?, ?, ?)"
	res, err := env.DB.Exec(query, input.Title, releaseDate, input.CoverImageURL, input.ArtistID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create album"})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"message": "Album created", "album_id": id})
}

func (env *ApiEnv) GetAlbumsByRating(c *gin.Context) {
	query := "SELECT album_id, title, artist_name, cover_image_url, average_rating, total_comments FROM v_album_avg_rating ORDER BY total_comments DESC, average_rating DESC"
	
	rows, err := env.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query albums"})
		return
	}
	defer rows.Close()

	albums := []models.AlbumRating{}
	for rows.Next() {
		var album models.AlbumRating
		if err := rows.Scan(&album.AlbumID, &album.Title, &album.ArtistName, &album.CoverImageURL, &album.AverageRating, &album.TotalComments); err != nil {
			log.Printf("Error scanning album rating: %v", err)
			continue
		}
		albums = append(albums, album)
	}
	c.JSON(http.StatusOK, albums)
}

func (env *ApiEnv) SearchAlbums(c *gin.Context) {
	searchTerm := c.Query("q")
	query := `
		SELECT album_id, title, artist_name, cover_image_url, average_rating, total_comments 
		FROM v_album_avg_rating 
		WHERE title LIKE ? OR artist_name LIKE ?
		ORDER BY total_comments DESC, average_rating DESC
	`
	rows, err := env.DB.Query(query, "%"+searchTerm+"%", "%"+searchTerm+"%")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search albums"})
		return
	}
	defer rows.Close()

	albums := []models.AlbumRating{}
	for rows.Next() {
		var album models.AlbumRating
		if err := rows.Scan(&album.AlbumID, &album.Title, &album.ArtistName, &album.CoverImageURL, &album.AverageRating, &album.TotalComments); err != nil {
			log.Printf("Error scanning album rating: %v", err)
			continue
		}
		albums = append(albums, album)
	}
	c.JSON(http.StatusOK, albums)
}

func (env *ApiEnv) GetAlbumDetails(c *gin.Context) {
	albumID := c.Param("id")
	var album models.Album
	albumSQL := `
		SELECT a.album_id, a.title, a.release_date, a.cover_image_url, a.artist_id, ar.name 
		FROM albums a
		JOIN artists ar ON a.artist_id = ar.artist_id
		WHERE a.album_id = ?
	`
	err := env.DB.QueryRow(albumSQL, albumID).Scan(
		&album.ID, &album.Title, &album.ReleaseDate, &album.CoverImageURL, &album.ArtistID, &album.ArtistName,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	comments, _ := env.fetchCommentsForAlbum(albumID)
	genres, _ := env.fetchGenresForAlbum(albumID)

	c.JSON(http.StatusOK, gin.H{
		"album_details": album,
		"comments":    comments,
		"genres":      genres,
	})
}

func (env *ApiEnv) CreateComment(c *gin.Context) {
	albumID := c.Param("id")
	var input models.CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := env.getOrCreateUser(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user"})
		return
	}

	query := "INSERT INTO comments (album_id, user_id, rating, comment_text) VALUES (?, ?, ?, ?)"
	_, err = env.DB.Exec(query, albumID, userID, input.Rating, input.CommentText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post comment"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Comment created"})
}

func (env *ApiEnv) CreateGenre(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := "INSERT INTO genres (name) VALUES (?)"
	res, err := env.DB.Exec(query, input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create genre"})
		return
	}
	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"message": "Genre created", "genre_id": id})
}

func (env *ApiEnv) GetGenres(c *gin.Context) {
	rows, err := env.DB.Query("SELECT genre_id, name FROM genres ORDER BY name")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch genres"})
		return
	}
	defer rows.Close()

	genres := []models.Genre{}
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			continue
		}
		genres = append(genres, genre)
	}
	c.JSON(http.StatusOK, genres)
}

func (env *ApiEnv) AddGenreToAlbum(c *gin.Context) {
	albumID := c.Param("id")
	var input struct {
		GenreID int `json:"genre_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO album_genres (album_id, genre_id) VALUES (?, ?)"
	_, err := env.DB.Exec(query, albumID, input.GenreID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "This genre is already on the album"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Genre added to album"})
}

func (env *ApiEnv) GetCommentsByUser(c *gin.Context) {
	username := c.Param("username")
	query := `
		SELECT 
			c.comment_id, c.album_id, c.rating, c.comment_text, c.created_at,
			u.username,
			a.title as album_title, a.cover_image_url
		FROM comments c
		JOIN users u ON c.user_id = u.user_id
		JOIN albums a ON c.album_id = a.album_id
		WHERE u.username = ?
		ORDER BY c.created_at DESC
	`
	rows, err := env.DB.Query(query, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user comments"})
		return
	}
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID, &comment.AlbumID, &comment.Rating, &comment.CommentText, &comment.CreatedAt,
			&comment.Username, &comment.AlbumTitle, &comment.CoverImageURL,
		)
		if err != nil {
			log.Printf("Error scanning comment: %v", err)
			continue
		}
		comments = append(comments, comment)
	}
	c.JSON(http.StatusOK, comments)
}

func (env *ApiEnv) fetchCommentsForAlbum(albumID string) ([]models.Comment, error) {
	query := `
		SELECT c.comment_id, c.album_id, c.user_id, c.rating, c.comment_text, c.created_at, u.username
		FROM comments c
		JOIN users u ON c.user_id = u.user_id
		WHERE c.album_id = ?
		ORDER BY c.created_at DESC
	`
	rows, err := env.DB.Query(query, albumID)
	if err != nil { return nil, err }
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.AlbumID, &comment.UserID, &comment.Rating, &comment.CommentText, &comment.CreatedAt, &comment.Username); err != nil {
			log.Printf("Error scanning comment: %v", err)
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (env *ApiEnv) fetchGenresForAlbum(albumID string) ([]models.Genre, error) {
	query := `
		SELECT g.genre_id, g.name
		FROM genres g
		JOIN album_genres ag ON g.genre_id = ag.genre_id
		WHERE ag.album_id = ?
	`
	rows, err := env.DB.Query(query, albumID)
	if err != nil { return nil, err }
	defer rows.Close()

	genres := []models.Genre{}
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			log.Printf("Error scanning genre: %v", err)
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func (env *ApiEnv) DeleteUser(c *gin.Context) {
	username := c.Param("username")

	query := "DELETE FROM users WHERE username = ?"

	res, err := env.DB.Exec(query, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check deletion status"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User and all their comments deleted"})
}
