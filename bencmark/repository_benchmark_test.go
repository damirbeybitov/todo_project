package repository_test

import (
	"database/sql"
	"testing"

	"github.com/damirbeybitov/todo_project/internal/user/repository"
	_ "github.com/go-sql-driver/mysql"
)

func BenchmarkCheckUserInDB(b *testing.B) {
	db := setupTestDB(b)
	repo := repository.NewRepository(db)
	tx, _ := db.Begin()

	// Run the benchmark function b.N times
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.AddUserToDB(tx, "username", "email", "password")
	}
	for i := 0; i < b.N; i++ {
		_ = repo.CheckUserInDB(tx, "username", "email")
	}
	b.ReportAllocs()
    b.ReportMetric(float64(b.N), "iterations")
}

func BenchmarkAddUserToDB(b *testing.B) {
	db := setupTestDB(b)
	repo := repository.NewRepository(db)
	tx, _ := db.Begin()

	// Run the benchmark function b.N times
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.AddUserToDB(tx, "username", "email", "password")
	}
	b.ReportAllocs()
    b.ReportMetric(float64(b.N), "iterations")
}

func BenchmarkCheckPassword(b *testing.B) {
	db := setupTestDB(b)
	repo := repository.NewRepository(db)

	// Run the benchmark function b.N times
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = repo.CheckPassword("username", "password")
	}
	b.ReportAllocs()
    b.ReportMetric(float64(b.N), "iterations")
}

func BenchmarkDeleteUserFromDB(b *testing.B) {
	db := setupTestDB(b)
	repo := repository.NewRepository(db)
	tx, _ := db.Begin()

	// Run the benchmark function b.N times
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = repo.DeleteUserFromDB(tx, "username")
	}
	b.ReportAllocs()
    b.ReportMetric(float64(b.N), "iterations")
}

func setupTestDB(t testing.TB) *sql.DB {
	// Connect to the test database
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/to_do")
	if err != nil {
		t.Fatalf("Error connecting to MySQL: %v", err)
	}
	return db
}
