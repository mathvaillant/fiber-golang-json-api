package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Download(route fiber.Router, db *gorm.DB) {
	route.Get("/", func(c *fiber.Ctx) error {
		format := c.Query("format", "json")
		userId := int(c.Locals("userId").(float64))
		books := new([]Book)

		// Fetch books for the specific user
		if err := db.Where("user_id = ?", userId).Find(&books).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch books",
			})
		}

		var fileName string

		switch format {
		case "json":
			fileName = "books.json"

			// Create a JSON file
			file, err := os.Create(fileName)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create JSON file",
				})
			}
			defer file.Close()

			// Write book data to the JSON file
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(books); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to write JSON file",
				})
			}

		case "csv":
			fileName = "books.csv"

			// Create a CSV file
			file, err := os.Create(fileName)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create CSV file",
				})
			}
			defer file.Close()

			// Write book data to the CSV file
			writer := csv.NewWriter(file)

			// Write CSV headers
			writer.Write([]string{"ID", "Title", "Status", "Author", "Year"})

			// Write book data
			for _, book := range *books {
				writer.Write([]string{
					fmt.Sprintf("%d", book.ID),
					book.Title,
					string(book.Status),
					book.Author,
					fmt.Sprintf("%d", book.Year),
				})
			}

			// Flush and check for errors
			writer.Flush()
			if err := writer.Error(); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to write CSV file",
				})
			}

		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid format. Use 'json' or 'csv'",
			})
		}

		// Send the file for download and delete it afterward
		defer os.Remove(fileName)
		return c.Download(fileName)
	})
}
