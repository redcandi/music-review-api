package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"music-review-api/internal/models"
)

// --- Artist Handlers ---

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

func (env *ApiEnv) GetArtistByID(c *gin.Context) {
	id := c.Param("id")
	var artist models.Artist
	query := "SELECT artist_id, name, bio, formed_year FROM artists WHERE artist_id = ?"

	err := env.DB.QueryRow(query, id).Scan(&artist.ID, &artist.Name, &artist.Bio, &artist.FormedYear)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artist not found"})
		return
	}
	c.JSON(http.StatusOK, artist)
}

// --- Album Handlers ---

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

func (env *ApiEnv) GetAlbums(c *gin.Context) {
	query := `
		SELECT a.album_id, a.title, a.release_date, a.cover_image_url, a.artist_id, ar.name 
		FROM albums a
		JOIN artists ar ON a.artist_id = ar.artist_id
		ORDER BY a.release_date DESC
	`
	rows, err := env.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch albums"})
		return
	}
	defer rows.Close()

	albums := []models.Album{}
	for rows.Next() {
		var album models.Album
		if err := rows.Scan(&album.ID, &album.Title, &album.ReleaseDate, &album.CoverImageURL, &album.ArtistID, &album.ArtistName); err != nil {
			log.Printf("Error scanning album: %v", err)
			continue
		}
		albums = append(albums, album)
	}
	c.JSON(http.StatusOK, albums)
}

func (env *ApiEnv) GetAlbumByID(c *gin.Context) {
	albumID := c.Param("id")

	// 1. Get Album Details
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

	// 2. Get Comments
	comments, err := env.fetchCommentsForAlbum(albumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	// 3. Combine and Return
	c.JSON(http.StatusOK, gin.H{
		"album_details": album,
		"comments":    comments,
	})
}

// --- Aggregate Query Handler (Using the VIEW) ---

func (env *ApiEnv) GetTopRatedAlbums(c *gin.Context) {
	query := `
		SELECT album_id, title, artist_name, cover_image_url, average_rating, total_comments
		FROM v_album_avg_rating
		WHERE total_comments > 0
		ORDER BY average_rating DESC
		LIMIT 10
	`
	rows, err := env.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query top albums"})
		return
	}
	defer rows.Close()

	albums := []models.AlbumRating{}
	for rows.Next() {
		var album models.AlbumRating
		if err := rows.Scan(&album.AlbumID, &album.Title, &album.ArtistName, &album.CoverImageURL, &album.AverageRating, &album.TotalComments); err != nil {
			log.Printf("Error scanning top album: %v", err)
			continue
		}
		albums = append(albums, album)
	}
	c.JSON(http.StatusOK, albums)
}

// --- Comment Handlers ---

func (env *ApiEnv) CreateComment(c *gin.Context) {
	albumID := c.Param("id")

	var input models.CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO comments (album_id, username, rating, comment_text) VALUES (?, ?, ?, ?)"
	_, err := env.DB.Exec(query, albumID, input.Username, input.Rating, input.CommentText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post comment"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Comment created"})
}

func (env *ApiEnv) GetCommentsForAlbum(c *gin.Context) {
	albumID := c.Param("id")
	comments, err := env.fetchCommentsForAlbum(albumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

// --- Helper Function ---

// fetchCommentsForAlbum is now simpler (no JOIN)
func (env *ApiEnv) fetchCommentsForAlbum(albumID string) ([]models.Comment, error) {
	query := `
		SELECT comment_id, album_id, username, rating, comment_text, created_at
		FROM comments
		WHERE album_id = ?
		ORDER BY created_at DESC
	`
	rows, err := env.DB.Query(query, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.AlbumID, &comment.Username, &comment.Rating, &comment.CommentText, &comment.CreatedAt); err != nil {
			log.Printf("Error scanning comment: %v", err)
			continue
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
