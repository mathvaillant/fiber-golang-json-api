package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BookHandlers(route fiber.Router, db *gorm.DB) {
	route.Get("/", func(c *fiber.Ctx) error {
		// Queries
		title := c.Query("title")
		status := c.Query("status")
		author := c.Query("author")
		year := c.QueryInt("year")

		userId := int(c.Locals("userId").(float64))
		books := new([]Book)

		query := db.Where("user_id = ?", userId)

		if title != "" {
			query = query.Where("title LIKE ?", "%"+title+"%")
		}
		if status != "" {
			query = query.Where("status = ?", status)
		}
		if author != "" {
			query = query.Where("author = ?", author)
		}
		if year != 0 {
			query = query.Where("year = ?", year)
		}

		if err := query.Find(&books).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Books not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(books)
	})

	route.Get("/:id", func(c *fiber.Ctx) error {
		bookId, _ := c.ParamsInt("id")
		userId := int(c.Locals("userId").(float64))
		book := new(Book)

		if err := db.Where("id = ? AND user_id = ?", bookId, userId).First(&book).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(book)
	})

	route.Post("/", func(c *fiber.Ctx) error {
		book := new(Book)
		book.UserID = int(c.Locals("userId").(float64))

		if err := c.BodyParser(book); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := db.Create(&book).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(book)
	})

	route.Put("/:id", func(c *fiber.Ctx) error {
		bookId, _ := c.ParamsInt("id")
		userId := int(c.Locals("userId").(float64))
		book := new(Book)

		if err := db.Where("id = ? AND user_id = ?", bookId, userId).First(&book).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}

		if err := c.BodyParser(book); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := db.Save(&book).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(book)
	})

	route.Delete("/:id", func(c *fiber.Ctx) error {
		bookId, _ := c.ParamsInt("id")
		userId := int(c.Locals("userId").(float64))
		book := new(Book)

		if err := db.Where("id = ? AND user_id = ?", bookId, userId).First(&book).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}

		if err := db.Delete(&book).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}
