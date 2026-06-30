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
			{"Chai thủy tinh dẹt 500ml", "SP01", 10.0, 100, catChaiLo, "Chai thủy tinh dẹt cao cấp đựng nước ép, cold brew", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRNN7NWhhYSJ_B5ha8Vk0Bfc1eZU60gYaGrnsOvd9H9cA&s=10"},
			{"Chai nhựa PET tròn 250ml", "SP02", 5.0, 200, catChaiLo, "Chai nhựa PET tròn trong suốt đựng sữa chua uống, trà sữa", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSpZjbVQOhXfiakBP3kpwJezzqfxZRsCKcVG-AXEuc4_g&s=10"},
			{"Chai thủy tinh tròn 330ml", "SP03", 9.5, 120, catChaiLo, "Chai thủy tinh tròn thích hợp đựng mật ong, nước ép", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQdCbJhdbbDZjVrc5JjK0V8QFWciBhQsP9Kc3unEhHbOQ&s"},
			{"Hũ thủy tinh nắp thiếc 100ml", "SP04", 8.0, 300, catChaiLo, "Hũ thủy tinh nhỏ đựng yến sào, thực phẩm khô", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcR3qNDXzdOxz1aokDazi9vzi4oAFndwXPX2Skr4TkE7zw&s=10"},
			{"Chai xịt phun sương 100ml", "SP05", 7.0, 150, catChaiLo, "Chai xịt phun sương đựng nước hoa hồng, cồn sát khuẩn", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRetXbPg2aq7zTg67Q2J8vFu0glWWhBQehi3e7Mb_LhjA&s=10"},
			{"Hộp nhựa tròn 1000ml", "SP06", 15.0, 150, catHopNhua, "Hộp nhựa đựng thực phẩm chịu nhiệt dùng được lò vi sóng", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQQujFe3ujP2CBxaKVutJ5Hpi47M3MRCAryitrNn8GYiO0AN_j8qFeIhwM&s=10"},
			{"Hộp nhựa chữ nhật 750ml", "SP07", 12.0, 180, catHopNhua, "Hộp nhựa chữ nhật chia ngăn tiện lợi đựng cơm văn phòng", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRsaGRspoZOrZ3UYSN0dF9_pi50v9uY0NU8i23z6WC9bw&s=10"},
			{"Hộp nhựa vuông 500ml", "SP08", 10.0, 250, catHopNhua, "Hộp nhựa vuông đựng mứt, bánh kẹo, thực phẩm khô", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRqivoIZDa7Pa27dRe9kZM7qw4wJxZfqVpD-MTunjvg7w&s=10"},
			{"Khay nhựa trong đựng quả", "SP09", 6.0, 400, catHopNhua, "Khay nhựa trong suốt có lỗ thoáng khí đựng trái cây, dâu tây", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcThS61leY1eTWbKQ8IPNfUfRBNbDdhgfpwTnX27gZl0ow&s=10"},
			{"Hũ nhựa nắp nhôm 750ml", "SP10", 11.0, 220, catHopNhua, "Hũ nhựa PET dáng cao nắp nhôm xé đựng hạt khô, khô gà", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRSJV0SLEpYnAuuo2JQ86kBA21jmTermdGkWDIQ_jZ4ug&s=10"},
			{"Túi giấy Kraft quai xoắn", "SP11", 3.0, 500, catBaoBi, "Túi giấy Kraft bảo vệ môi trường, dai và chịu lực tốt", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQoDn8dwwUmVL0TYuk7_u99QMA__C8oISVMK3sqMJw51w&s=10"},
			{"Túi ni lông sinh học", "SP12", 2.0, 1000, catBaoBi, "Túi ni lông tự phân hủy sinh học, thân thiện môi trường", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQgEUzQSjCG-_wNMWg0iploL-_qFijACdw78uWSGpsbSg&s=10"},
			{"Hộp carton đóng hàng", "SP13", 4.0, 800, catBaoBi, "Hộp giấy carton 3 lớp đóng hàng COD bảo vệ sản phẩm tốt", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSTBd9DhHEgKPfJK8_Juq9WmF8SHlRyim1Xu4Qy1f4F-ecq47cNf7BOZwA&s=10"},
			{"Ly giấy đựng cafe 12oz", "SP14", 2.5, 600, catBaoBi, "Ly giấy tráng PE 2 lớp đựng đồ uống nóng lạnh tiện lợi", "https://hunufa.vn/wp-content/uploads/2024/06/ly-giay-12oz-360ml-coc-giay-12oz-360ml-3-1.webp"},
			{"Màng co PE bọc thực phẩm", "SP15", 25.0, 80, catBaoBi, "Cuộn màng co PE khổ lớn bọc thực phẩm bảo quản tủ lạnh", "https://thanhbinh.com.vn/upload/images/2022/07/1658742509-single_product1-UD.jpg"},
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
