package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("[%s] New %s request from %s on %s", time.Now(), c.Method(), c.IP(), c.Path())
		return c.Next()
	}
}

func ErrLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			log.Printf("[%s] Error during %s request from %s on %s: %s", time.Now(), c.Method(), c.IP(), c.Path(), err)
		}

		return err
	}
}
