package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the token either from the Auth header or from the Cookie
		cookieToken := c.Cookies("jwt")
		var tokenString string

		if cookieToken != "" {
			// If the token is in the Cookies, use it
			log.Warn("token from Cookies, using it...")
			tokenString = cookieToken
		} else {
			// If the token is not in the Cookies, try to get it from the Authorization header
			log.Warn("empty token from Cookies, trying to get it from the Authorization header...")

			// Get the Authorization header
			authHeader := c.Get("Authorization")

			// If the Authorization header is empty, return an error
			if authHeader == "" {
				log.Warn("empty authorization header")

				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"status":  "fail",
					"message": "Unauthorized",
				})
			}

			// Split the Authorization header to get the token
			tokenParts := strings.Split(authHeader, " ")

			// If the token is not in the format
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				log.Warn("invalid token parts")

				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"status":  "fail",
					"message": "Unauthorized",
				})
			}

			// Get the token from the token parts
			tokenString = tokenParts[1]
		}

		secret := []byte("super-secret-key")

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		// If the token is invalid, return an error and clear the cookie
		if err != nil || !token.Valid {
			log.Warn("invalid token")
			c.ClearCookie("jwt")

			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "Unauthorized",
			})
		}

		// Get the userId from the token
		userId := token.Claims.(jwt.MapClaims)["userId"]

		// Check if the user exists in the database, if not, return an error and clear the cookie
		if err := db.Model(&User{}).Where("id = ?", userId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("user not found in the db")
			c.ClearCookie("jwt")

			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "Unauthorized",
			})
		}

		// Set the userId in the Locals
		c.Locals("userId", userId)

		// Continue to the next middleware
		return c.Next()
	}
}
