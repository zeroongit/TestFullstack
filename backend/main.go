package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fullstack-backend/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var jwtSecret string

// 1. DATABASE SETUP
func initDatabase() {
	// Load .env dulu sebelum mengambil value
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found, using system environment variables")
	}

	jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "TesTingFullStack321" // Fallback
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrasi Tabel
	DB.AutoMigrate(&models.Item{}, &models.Invoice{}, &models.InvoiceDetail{})

	// Seeding Master Data jika kosong
	var count int64
	DB.Model(&models.Item{}).Count(&count)
	if count == 0 {
		items := []models.Item{
			{Code: "BRG-001", Name: "Ban Truk Fuso", Price: 5000000},
			{Code: "BRG-002", Name: "Ban Truk Canter", Price: 3000000},
			{Code: "BRG-003", Name: "Ban Truk Colt Diesel", Price: 2000000},
		}
		DB.Create(&items)
		fmt.Println("✅ Database Seeded Successfully")
	}
}

// 2. MAIN FUNCTION (Routes & Middleware)
func main() {
	initDatabase()

	app := fiber.New()

	// Standard Middlewares
	app.Use(recover.New()) // Mencegah app mati jika ada panic
	app.Use(logger.New())  // Log setiap request di terminal
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://127.0.0.1:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Public Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("🚀 Fullstack Logistics Backend is Running!")
	})
	app.Post("/api/login", login)
	app.Get("/api/items", getItemByCode)

	// JWT Middleware
	auth := jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtSecret),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token tidak valid atau kadaluwarsa",
			})
		},
	})

	// Protected Routes
	app.Post("/api/invoices", auth, createInvoice)
	app.Get("/api/invoices", auth, getInvoices)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}

// 3. HANDLERS
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
		role = "Admin"; userID = 1
	} else if input.Username == "kerani" && input.Password == "kerani123" {
		role = "Kerani"; userID = 2
	} else {
		return c.Status(401).JSON(fiber.Map{"message": "Username atau password salah"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userID),
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal membuat token"})
	}

	return c.JSON(fiber.Map{"token": t, "role": role})
}

func getItemByCode(c *fiber.Ctx) error {
	code := c.Query("code")
	var item models.Item
	if err := DB.Where("code = ?", code).First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Item tidak ditemukan"})
	}
	return c.JSON(item)
}

func createInvoice(c *fiber.Ctx) error {
	// Safe User Context Retrieval
	rawUser := c.Locals("user")
	if rawUser == nil {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
	}
	userToken := rawUser.(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	var input struct {
		SenderName      string `json:"senderName"`
		SenderAddress   string `json:"senderAddress"`
		ReceiverName    string `json:"receiverName"`
		ReceiverAddress string `json:"receiverAddress"`
		Items           []struct {
			ID       uint `json:"id"`
			Quantity int  `json:"quantity"`
		} `json:"items"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Payload tidak valid"})
	}

	newInvoiceID := uuid.New()

	err := DB.Transaction(func(tx *gorm.DB) error {
		var grandTotal float64 = 0
		invNumber := "INV-" + time.Now().Format("20060102150405")

		invoice := models.Invoice{
			ID:              newInvoiceID,
			InvoiceNumber:   invNumber,
			SenderName:      input.SenderName,
			SenderAddress:   input.SenderAddress,
			ReceiverName:    input.ReceiverName,
			ReceiverAddress: input.ReceiverAddress,
			TotalAmount:     0,
			CreatedBy:       userID,
		}

		if err := tx.Create(&invoice).Error; err != nil {
			return err
		}

		for _, itemInput := range input.Items {
			var masterItem models.Item
			if err := tx.First(&masterItem, itemInput.ID).Error; err != nil {
				return err // Item tidak ditemukan di DB
			}

			subtotal := masterItem.Price * float64(itemInput.Quantity)
			grandTotal += subtotal

			detail := models.InvoiceDetail{
				InvoiceID: invoice.ID,
				ItemID:    masterItem.ID,
				ItemName:  masterItem.Name,
				Quantity:  itemInput.Quantity,
				Price:     masterItem.Price,
				Subtotal:  subtotal,
			}
			if err := tx.Create(&detail).Error; err != nil {
				return err
			}
		}
		return tx.Model(&invoice).Update("total_amount", grandTotal).Error
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal simpan invoice: " + err.Error()})
	}

	go sendWebhook(newInvoiceID.String())

	return c.Status(201).JSON(fiber.Map{
		"message": "Invoice berhasil dibuat!",
		"id":      newInvoiceID,
	})
}

func getInvoices(c *fiber.Ctx) error {
	var invoices []models.Invoice
	if err := DB.Preload("Details").Order("created_at desc").Find(&invoices).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data history"})
	}
	return c.JSON(invoices)
}

func sendWebhook(invoiceID string) {
	url := "https://webhook.site/dummy-endpoint" // Sesuaikan untuk testing
	payload, _ := json.Marshal(map[string]string{
		"event": "invoice.created", "invoice_id": invoiceID,
	})
	http.Post(url, "application/json", bytes.NewBuffer(payload))
}