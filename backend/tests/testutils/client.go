package testutils

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"g_gen/internal/env"
	"g_gen/internal/infra/db"
	applogger "g_gen/internal/infra/logger"
)

// SetupTestDB テスト用のDBクライアントをセットアップする
func SetupTestDB(t *testing.T) db.Client {
	e, err := env.NewValues()
	if err != nil {
		t.Fatalf("failed to load environment variables: %v", err)
	}
	client, err := db.NewSQLHandler(&db.DatabaseConfig{
		Host:     e.TestDB.TestDatabaseHost,
		Port:     e.TestDB.TestDatabasePort,
		User:     e.TestDB.TestDatabaseUsername,
		Password: e.TestDB.TestDatabasePassword,
		DBName:   e.TestDB.TestDatabaseName,
		SSLMode:  "disable",
		Timezone: "Asia/Tokyo",
	}, applogger.New(applogger.DefaultConfig()))
	if err != nil {
		t.Errorf("failed to connect to test database: %v", err)
	}

	return client
}

// TruncateAllTables テスト用のDBの全テーブルをトランケートする
func TruncateAllTables(t *testing.T, client db.Client) {
	// トランザクションを開始
	tx := client.Conn(context.Background()).Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	// 全テーブルをトランケート
	if err := tx.Exec("TRUNCATE TABLE prefectures, municipalities RESTART IDENTITY CASCADE").Error; err != nil {
		tx.Rollback()
		t.Fatalf("failed to truncate tables: %v", err)
	}

	// トランザクションをコミット
	if err := tx.Commit().Error; err != nil {
		t.Fatalf("failed to commit transaction: %v", err)
	}
}

func NewTestClient(t *testing.T) (db.Client, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %s", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn: sqlDB,
	})

	mockDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm: %s", err)
	}

	return &db.SQLHandler{Driver: mockDB}, mock
}
