package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/Blaze5333/shorten-url/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func ResolveURL(c *fiber.Ctx) error {
	rInr := database.CreateClient(1)
	defer rInr.Close()
	val, _ := rInr.Get(database.Ctx, c.IP()).Result()
	if val == "" {
		_ = rInr.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		valint, _ := strconv.Atoi(val)
		if valint <= 0 {
			limit, _ := rInr.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "rate limit exceeded", "rate_limit_reset": limit.Seconds()})
		}
	}
	url := c.Params("url")
	r := database.CreateClient(0)
	defer r.Close()
	val, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	rInr.Incr(database.Ctx, "counter")
	rInr.Decr(database.Ctx, c.IP())
	return c.Redirect(val, fiber.StatusTemporaryRedirect)

}
