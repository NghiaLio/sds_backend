package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"golang.org/x/crypto/bcrypt"
)

// InitDB initializes the SQLite database and returns the *sql.DB handle.
func InitDB(dbPath string) (*sql.DB, error) {
	// Open connection to SQLite
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables
	if err = migrate(db); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	// Seed data if empty
	if err = seed(db); err != nil {
		return nil, fmt.Errorf("seeding failed: %w", err)
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	// Enable foreign key constraints in SQLite
	_, err := db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}

	usersQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);`
	if _, err = db.Exec(usersQuery); err != nil {
		return err
	}

	categoriesQuery := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);`
	if _, err = db.Exec(categoriesQuery); err != nil {
		return err
	}

	productsQuery := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		code TEXT UNIQUE NOT NULL,
		price REAL NOT NULL,
		stock INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		description TEXT,
		image TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
	);`
	if _, err = db.Exec(productsQuery); err != nil {
		return err
	}

	return nil
}

func seed(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Seed user 'cuongpc10' (with password '123456' from Postman payload)
	var userID int64
	err = tx.QueryRow("SELECT id FROM users WHERE username = 'cuongpc10'").Scan(&userID)
	if err == sql.ErrNoRows {
		hashed, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		res, err := tx.Exec("INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)", "cuongpc10", string(hashed), time.Now())
		if err != nil {
			return err
		}
		userID, err = res.LastInsertId()
		if err != nil {
			return err
		}
		log.Println("Seeded default user 'cuongpc10' (password: '123456')")
	} else if err != nil {
		return err
	}

	// 2. Seed Categories if empty
	var catCount int
	err = tx.QueryRow("SELECT COUNT(1) FROM categories").Scan(&catCount)
	if err != nil {
		return err
	}

	if catCount == 0 {
		log.Println("Seeding product categories...")
		categories := []string{"Chai lọ", "Hộp nhựa", "Bao bì"}
		stmt, err := tx.Prepare("INSERT INTO categories (name) VALUES (?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, name := range categories {
			if _, err = stmt.Exec(name); err != nil {
				return err
			}
		}
	}

	// 3. Seed Products if empty
	var prodCount int
	err = tx.QueryRow("SELECT COUNT(1) FROM products").Scan(&prodCount)
	if err != nil {
		return err
	}

	if prodCount == 0 {
		log.Println("Seeding sample products...")
		// Fetch category IDs
		var catChaiLo, catHopNhua, catBaoBi int64
		_ = tx.QueryRow("SELECT id FROM categories WHERE name = 'Chai lọ'").Scan(&catChaiLo)
		_ = tx.QueryRow("SELECT id FROM categories WHERE name = 'Hộp nhựa'").Scan(&catHopNhua)
		_ = tx.QueryRow("SELECT id FROM categories WHERE name = 'Bao bì'").Scan(&catBaoBi)

		// Set default fallbacks if query failed for some reason
		if catChaiLo == 0 {
			catChaiLo = 1
		}
		if catHopNhua == 0 {
			catHopNhua = 2
		}
		if catBaoBi == 0 {
			catBaoBi = 3
		}

		products := []struct {
			Name        string
			Code        string
			Price       float64
			Stock       int
			CategoryID  int64
			Description string
			Image       string
		}{
			{"Chai thủy tinh 500ml", "SP01", 10.0, 100, catChaiLo, "Chai thủy tinh cao cấp đựng nước hoa quả", "https://example.com/images/chai-tt-500ml.png"},
			{"Chai nhựa PET 250ml", "SP02", 5.0, 200, catChaiLo, "Chai nhựa PET tiện lợi đựng sữa chua uống", "https://example.com/images/chai-pet-250ml.png"},
			{"Hộp nhựa tròn 1000ml", "SP03", 15.0, 150, catHopNhua, "Hộp nhựa đựng thực phẩm quay được lò vi sóng", "https://example.com/images/hop-tron-1000ml.png"},
			{"Hộp nhựa vuông 750ml", "SP04", 12.0, 180, catHopNhua, "Hộp nhựa đựng thực phẩm chia ngăn", "https://example.com/images/hop-vuong-750ml.png"},
			{"Túi giấy kraft quai xoắn", "SP05", 3.0, 500, catBaoBi, "Túi giấy kraft bảo vệ môi trường", "https://example.com/images/tui-kraft.png"},
			{"Túi ni lông sinh học", "SP06", 2.0, 1000, catBaoBi, "Túi ni lông phân hủy sinh học tự nhiên", "https://example.com/images/tui-nilon-sh.png"},
			{"Chai thủy tinh lùn 330ml", "SP07", 9.5, 120, catChaiLo, "Chai thủy tinh đựng trà sữa, nước ép", "https://example.com/images/chai-tt-330ml.png"},
			{"Example Product", "SP22", 12.5, 100, catChaiLo, "An example product with optional description.", "https://example.com/image.png"},
		}

		stmt, err := tx.Prepare(`
			INSERT INTO products (name, code, price, stock, category_id, description, image, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`)
		if err != nil {
			return err
		}
		defer stmt.Close()

		now := time.Now()
		for _, p := range products {
			if _, err = stmt.Exec(p.Name, p.Code, p.Price, p.Stock, p.CategoryID, p.Description, p.Image, now, now); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
