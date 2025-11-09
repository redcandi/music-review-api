package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func InitDB() (*sql.DB, error){
	
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	dsn := os.Getenv("DB_DSN")

	if dsn == "" {
		return nil, fmt.Errorf("DB_DSN variable error")
	}

	db,err := sql.Open("mysql",dsn)
	if err != nil{
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	fmt.Println("Successfully connected to database!")

	return db,nil

}
