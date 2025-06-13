package main

// import (
// 	"github.com/dopp1e/webingo/backend/internal/db"
// 	"github.com/dopp1e/webingo/backend/internal/service"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// func main() {
// 	conn, err := gorm.Open(postgres.Open("host=localhost user=admin password=adminpassword dbname=webingo port=5432 sslmode=disable TimeZone=UTC"), &gorm.Config{
// 		// DisableForeignKeyConstraintWhenMigrating: true,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	service := service.NewService(conn)

// 	db.Seed(*service)
// }
