package routes

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rohitchauraisa1997/url-shortener/database"
	"github.com/rohitchauraisa1997/url-shortener/helpers"
)

type request struct {
	URL    string        `json:"url"`
	Expiry time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset string        `json:"rate_limit_reset_in"`
}

func ShortenURL(c *fiber.Ctx) error {
	// db0 is just used to store the url, where key is the custom shortenedUrl and value is the actual url/link.
	// db1 is used for rate limiting logic for each user and analytics regarding url hits and ttl of each url/link.
	rdb1 := database.CreateClient(1)
	defer rdb1.Close()

	rdb0 := database.CreateClient(0)
	defer rdb0.Close()

	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Could not parse request"})
	}

	// check if the input is an actual URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// check for domain error
	if helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "URL cant be accessed :-)"})
	}

	// enforce http
	body.URL = helpers.EnforceHTTP(body.URL)

	// rate limiting logic for each user.
	val, err := rdb1.Get(database.Ctx, c.IP()).Result()

	if err != nil {
		if err == redis.Nil {
			err = rdb1.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*time.Minute).Err()
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "server err = " + err.Error(),
			})
		}
	} else {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			ttl, _ := rdb1.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":               "Rate limit exceeded",
				"rate_limit_reset_in": fmt.Sprintf("%d mins", (ttl / time.Nanosecond / time.Minute)),
			})
		}
	}

	id, err := generateNewUrlId()
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "error while generating custom shortened url",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 1440
	}

	err = rdb0.Set(database.Ctx, id, body.URL, body.Expiry*time.Minute).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server err = " + err.Error(),
		})
	}

	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: "30 mins",
	}

	// decrement on the basis of IP to prevent DDOS on the system.
	rdb1.Decr(database.Ctx, c.IP())

	err = rdb1.HSet(context.Background(), id, map[string]interface{}{"createdAt": time.Now(), "urlHits": 0, "lifespanInMins": body.Expiry}).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server err = " + err.Error(),
		})
	}

	val, _ = rdb1.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := rdb1.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = fmt.Sprintf("%d mins", (ttl / time.Nanosecond / time.Minute))

	resp.CustomShort = "http://" + os.Getenv("DOMAIN") + "/" + id
	return c.Status(fiber.StatusOK).JSON(resp)
}

func generateNewUrlId() (string, error) {
	rdb1 := database.CreateClient(1)
	defer rdb1.Close()
	var id string

	allKeys, err := rdb1.Keys(database.Ctx, "*").Result()
	if err != nil {
		return "", err
	}

	for {
		id = uuid.New().String()[:6]
		// Check if the generated id already exists in allKeys
		found := false
		for _, redisKey := range allKeys {
			if id == redisKey {
				found = true
				break
			}
		}

		if !found {
			break // Found a unique id, break the loop
		}
	}

	return id, nil
}
