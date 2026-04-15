package main

import (
	"time"
	"fmt"
	"log"
	"os"
	"fullstack-backend/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
var DB *gorm.DB

func initDatabase() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB.AutoMigrate(&models.Item{}, &models.Invoice{}, &models.InvoiceDetail{})

	var count int64
	DB.Model(&models.Item{}).Count(&count)
	if count == 0 {
		items := []models.Item{
			{Code: "BRG-001", Name: "Ban Truk Fuso", Price: 500000000},
			{Code: "BRG-002", Name: "Ban Truk Canter", Price: 300000000},
			{Code: "BRG-003", Name: "Ban Truk Colt Diesel", Price: 200000000},
		}
		DB.Create(&items)
		fmt.Println("Database Seeded scuccessfully")
	}
}

func main() {
	initDatabase()
	app := fiber.New()
	app.Use(cors.New())
	app.Post("api/login", login)
	app.Get("/api/items", func(c *fiber.Ctx) error {
		code := c.Query("code")
		var item models.Item
		if err := DB.Where("code = ?", code).First(&item).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"message" : "Item tidak ditemukan"})
		}
		return c.JSON(item)
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Fullstack Backend is Running!")
	})

	log.Fatal(app.Listen(":8080"))
}

func login(c *fiber.Ctx) error {
	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	var role string
	var userID uint
	
	if input.Username == "admin" && input.Password == "admin123" {
		role = "Admin"
		userID = 1
	} else if input.Username == "user" && input.Password == "kirana123" {
		role = "user"
		userID = 2
	} else {
		return c.Status(401).JSON(fiber.Map{"message": "Username atau password salah"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t, "role" : role})

}