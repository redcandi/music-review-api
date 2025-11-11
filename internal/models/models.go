package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Artist struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Bio        string `json:"bio,omitempty"`
	FormedYear int    `json:"formed_year,omitempty"`
}

type Album struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	ReleaseDate   time.Time `json:"release_date"`
	CoverImageURL string    `json:"cover_image_url,omitempty"`
	ArtistID      int       `json:"artist_id"`
	ArtistName    string    `json:"artist_name,omitempty"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Comment struct {
	ID          int       `json:"id"`
	AlbumID     int       `json:"album_id"`
	UserID      int       `json:"user_id"`
	Rating      int       `json:"rating"`
	CommentText string    `json:"comment_text,omitempty"`
	CreatedAt   time.Time `json:"created_at"`

	Username      string `json:"username,omitempty"`
	AlbumTitle    string `json:"album_title,omitempty"`
	CoverImageURL string `json:"cover_image_url,omitempty"`
}

type AlbumRating struct {
	AlbumID int `json:"album_id"`

	Title         string  `json:"title"`
	ArtistName    string  `json:"artist_name"`
	CoverImageURL string  `json:"cover_image_url"`
	AverageRating float64 `json:"average_rating"`
	TotalComments int     `json:"total_comments"`
}

type SignUpInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Username string `json:"username"`
}

type CreateCommentInput struct {
	Username    string `json:"username"`
	Rating      int    `json:"rating" binding:"required"`
	CommentText string `json:"comment_text"`
}

type CreateArtistInput struct {
	Name       string `json:"name" binding:"required"`
	Bio        string `json:"bio"`
	FormedYear int    `json:"formed_year"`
}

type CreateAlbumInput struct {
	Title         string `json:"title" binding:"required"`
	ReleaseDate   string `json:"release_date" binding:"required"`
	CoverImageURL string `json:"cover_image_url"`
	ArtistID      int    `json:"artist_id" binding:"required"`
}
