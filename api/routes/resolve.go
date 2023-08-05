package routes

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/rohitchauraisa1997/url-shortener/database"
)

func ResolveURL(c *fiber.Ctx) error {
	rdb := database.CreateClient(0)
	defer rdb.Close()

	url := c.Params("url")
	// important in order to not let the browser cache the url on its end..
	// otherwise if the same user visits the same endpoint again using the shortened URL
	// the ResolveURL function wont be hit and hence we cant maintain the counter value correctly!!
	c.Set("Cache-Control", "no-store, max-age=0")

	value, err := rdb.Get(database.Ctx, url).Result()
	if err != nil {
		if err == redis.Nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found in database"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server err = " + err.Error(),
		})
	}

	rdb1 := database.CreateClient(1)
	defer rdb1.Close()

	hgetAllResult := rdb1.HGetAll(context.Background(), url)
	if hgetAllResult.Err() != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server err = " + err.Error(),
		})
	}

	hgetAllResponse, err := hgetAllResult.Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server err = " + err.Error(),
		})
	}

	var urlHitsCount uint64
	if hgetAllResponse["urlHits"] == "" {
		urlHitsCount = 1
	} else {
		urlHitsCount, err = strconv.ParseUint(hgetAllResponse["urlHits"], 10, 64)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "server err = " + err.Error(),
			})
		}
		urlHitsCount += 1
	}

	err = rdb1.HSet(context.Background(), url, map[string]interface{}{"urlHits": urlHitsCount}).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server err = " + err.Error(),
		})
	}

	return c.Redirect(value, 301)
}
