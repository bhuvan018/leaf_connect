package utils

import (
        "database/sql"
        "log"
        "os"
        "time"

        _ "github.com/lib/pq"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB() {
        var err error
        connStr := os.Getenv("DATABASE_URL")
        if connStr == "" {
                log.Fatal("DATABASE_URL environment variable not set")
        }

        // Connect to the database
        db, err = sql.Open("postgres", connStr)
        if err != nil {
                log.Fatalf("Failed to open database: %v", err)
        }

        // Check the connection
        err = db.Ping()
        if err != nil {
                log.Fatalf("Failed to ping database: %v", err)
        }

        log.Println("Successfully connected to the database")

        // Create the required tables if they don't exist
        createTables()
}

// CreateTables creates all the necessary tables if they don't exist
func createTables() {
        // Create users table
        _, err := db.Exec(`
                CREATE TABLE IF NOT EXISTS users (
                        id SERIAL PRIMARY KEY,
                        email VARCHAR(100) UNIQUE NOT NULL,
                        username VARCHAR(50) UNIQUE NOT NULL,
                        password VARCHAR(255) NOT NULL,
                        name VARCHAR(100) NOT NULL,
                        location VARCHAR(100),
                        bio TEXT,
                        profile_pic TEXT,
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                        last_login_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                )
        `)
        if err != nil {
                log.Fatalf("Failed to create users table: %v", err)
        }

        // Create listings table
        _, err = db.Exec(`
                CREATE TABLE IF NOT EXISTS listings (
                        id SERIAL PRIMARY KEY,
                        user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                        title VARCHAR(200) NOT NULL,
                        description TEXT NOT NULL,
                        type VARCHAR(20) NOT NULL,
                        plant_type VARCHAR(50) NOT NULL,
                        price NUMERIC(10, 2) NOT NULL,
                        trade_for TEXT,
                        location VARCHAR(100) NOT NULL,
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                        status VARCHAR(20) DEFAULT 'available'
                )
        `)
        if err != nil {
                log.Fatalf("Failed to create listings table: %v", err)
        }

        // Create images table (for listing images)
        _, err = db.Exec(`
                CREATE TABLE IF NOT EXISTS listing_images (
                        id SERIAL PRIMARY KEY,
                        listing_id INTEGER REFERENCES listings(id) ON DELETE CASCADE,
                        image_url TEXT NOT NULL,
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                )
        `)
        if err != nil {
                log.Fatalf("Failed to create listing_images table: %v", err)
        }

        // Create messages table
        _, err = db.Exec(`
                CREATE TABLE IF NOT EXISTS messages (
                        id SERIAL PRIMARY KEY,
                        from_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                        to_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                        listing_id INTEGER REFERENCES listings(id) ON DELETE SET NULL,
                        content TEXT NOT NULL,
                        read BOOLEAN DEFAULT FALSE,
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                )
        `)
        if err != nil {
                log.Fatalf("Failed to create messages table: %v", err)
        }

        // Create favorites table
        _, err = db.Exec(`
                CREATE TABLE IF NOT EXISTS favorites (
                        id SERIAL PRIMARY KEY,
                        user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                        listing_id INTEGER REFERENCES listings(id) ON DELETE CASCADE,
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                        UNIQUE(user_id, listing_id)
                )
        `)
        if err != nil {
                log.Fatalf("Failed to create favorites table: %v", err)
        }

        log.Println("Database tables created successfully")
}

// GetDB returns the database connection
func GetDB() *sql.DB {
        return db
}

// CloseDB closes the database connection
func CloseDB() {
        if db != nil {
                db.Close()
        }
}

// StringToNullString converts a string to sql.NullString
func StringToNullString(s string) sql.NullString {
        if s == "" {
                return sql.NullString{}
        }
        return sql.NullString{
                String: s,
                Valid:  true,
        }
}

// TimeToNullTime converts a time.Time to sql.NullTime
func TimeToNullTime(t time.Time) sql.NullTime {
        if t.IsZero() {
                return sql.NullTime{}
        }
        return sql.NullTime{
                Time:  t,
                Valid: true,
        }
}