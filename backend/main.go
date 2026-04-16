package main

import (
	"fmt"
	"fullstack-backend/models"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// setup database
func initDatabase() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}
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

// function main
func main() {
	initDatabase()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://127.0.0.1:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))
	app.Post("api/login", login)
	app.Get("/api/items", func(c *fiber.Ctx) error {
		code := c.Query("code")
		var item models.Item
		if err := DB.Where("code = ?", code).First(&item).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"message": "Item tidak ditemukan"})
		}
		return c.JSON(item)
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Fullstack Backend is Running!")
	})
	app.Post("/api/invoices", createInvoice)

	log.Fatal(app.Listen(":8080"))
}

// function login
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
	} else if input.Username == "kerani" && input.Password == "kerani123" {
		role = "Kerani"
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

	return c.JSON(fiber.Map{"token": t, "role": role})

}

func createInvoice(c *fiber.Ctx) error {
	var input struct {
		SenderName      string  `json:"senderName"`
		SenderAddress   string  `json:"senderAddress"`
		ReceiverName    string  `json:"receiverName"`
		ReceiverAddress string  `json:"receiverAddress"`
		TotalAmount     float64 `json:"total_amount"`
		Items           []struct {
			ID       uint    `json:"id"`
			Name     string  `json:"name"`
			Price    float64 `json:"price"`
			Quantity int     `json:"quantity"`
			Subtotal float64 `json:"subtotal"`
		} `json:"items"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Format data tidak valid"})
	}

	err := DB.Transaction(func(tx *gorm.DB) error {
		// 1. Generate UUID dan Nomor Invoice Sederhana
		newID := uuid.New()
		invNumber := "INV-" + time.Now().Format("20060102150405")

		invoice := models.Invoice{
			ID:              newID,
			InvoiceNumber:   invNumber,
			SenderName:      input.SenderName,
			SenderAddress:   input.SenderAddress,
			ReceiverName:    input.ReceiverName,
			ReceiverAddress: input.ReceiverAddress,
			TotalAmount:     input.TotalAmount,
		}

		if err := tx.Create(&invoice).Error; err != nil {
			return err
		}

		// 2. Simpan Detail (Gunakan Nama Field sesuai Model kamu: Name)
		for _, item := range input.Items {
			detail := models.InvoiceDetail{
				InvoiceID: invoice.ID,
				ItemID:    item.ID,
				Name:      item.Name, // Sesuai model: Name string `json:"name"`
				Quantity:  item.Quantity,
				Price:     item.Price,
				Subtotal:  item.Subtotal,
			}
			if err := tx.Create(&detail).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menyimpan: " + err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Invoice berhasil dibuat!"})
}
