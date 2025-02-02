package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// 初始化数据库
func Init() error {
	// 确保数据目录存在
	dbDir := "./data"
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return err
	}

	// 打开数据库连接
	dbPath := filepath.Join(dbDir, "users.db")
	log.Printf("数据库路径: %s", dbPath)

	// 使用 SQLite3 驱动
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		log.Printf("打开数据库失败: %v", err)
		return err
	}

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		log.Printf("数据库连接测试失败: %v", err)
		return err
	}

	log.Printf("数据库连接成功")

	// 创建用户注册表
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS user_registrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		linux_do_username TEXT NOT NULL UNIQUE,
		linux_do_trust_level INTEGER NOT NULL,
		last_name TEXT NOT NULL,
		first_name TEXT NOT NULL,
		school_email TEXT NOT NULL UNIQUE,
		backup_email TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(createTableSQL); err != nil {
		log.Printf("创建表失败: %v", err)
		return err
	}

	log.Printf("数据库表创建成功")

	DB = db
	return nil
}

// 检查用户是否已注册
func CheckUserRegistration(username string) (bool, error) {
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM user_registrations WHERE linux_do_username = ?)", username).Scan(&exists)
	return exists, err
}

// 检查邮箱是否已被注册
func CheckEmailRegistration(email string) (bool, error) {
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM user_registrations WHERE school_email = ?)", email).Scan(&exists)
	return exists, err
}

// 创建新的用户注册
func CreateUserRegistration(registration *UserRegistration) error {
	_, err := DB.Exec(`
		INSERT INTO user_registrations (
			linux_do_username,
			linux_do_trust_level,
			last_name,
			first_name,
			school_email,
			backup_email
		) VALUES (?, ?, ?, ?, ?, ?)`,
		registration.LinuxDoUsername,
		registration.LinuxDoTrustLevel,
		registration.LastName,
		registration.FirstName,
		registration.SchoolEmail,
		registration.BackupEmail,
	)
	return err
}

// UserRegistration 结构体
type UserRegistration struct {
	ID                int64  `json:"id"`
	LinuxDoUsername   string `json:"linux_do_username"`
	LinuxDoTrustLevel int    `json:"linux_do_trust_level"`
	LastName          string `json:"last_name"`
	FirstName         string `json:"first_name"`
	SchoolEmail       string `json:"school_email"`
	BackupEmail       string `json:"backup_email"`
}
