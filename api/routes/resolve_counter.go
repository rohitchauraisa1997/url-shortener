package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/gofiber/fiber/v2"
	"github.com/rohitchauraisa1997/url-shortener/database"
	"github.com/rohitchauraisa1997/url-shortener/helpers"
	"github.com/rohitchauraisa1997/url-shortener/models"
)

func GetUrlResolutionAnalytics(c *fiber.Ctx) error {
	// db0 is just used to store the url, where key is the custom shortenedUrl and value is the actual url/link.
	// db1 is used for rate limiting logic for each user and analytics regarding url hits and ttl of each url/link.
	rdb0 := database.CreateClient(0)
	defer rdb0.Close()
	rdb1 := database.CreateClient(1)
	defer rdb1.Close()
	resp := map[string]models.UrlAnalyticDetails{}

	keys, err := rdb0.Keys(database.Ctx, "*").Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server err = " + err.Error(),
		})
	}
	if len(keys) == 0 {
		return c.Status(fiber.StatusOK).JSON(resp)
	}

	vals, err := rdb0.MGet(database.Ctx, keys...).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server err = " + err.Error(),
		})
	}

	mappedShortenedUrlAndUrl := make(map[string]string)
	for ctr, key := range keys {
		mappedShortenedUrlAndUrl[key] = vals[ctr].(string)
	}

	// using pipelines for faster and more efficient calls.
	pipe := rdb1.Pipeline()
	commandMapper := map[string]*redis.StringCmd{}
	for _, key := range keys {
		commandMapper[key] = pipe.HGet(database.Ctx, key, "createdAt")
	}
	_, err = pipe.Exec(database.Ctx)
	if err != nil && err != redis.Nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server err = " + err.Error(),
		})
	}
	var mappedUrlsAndCreatedAt = make(map[string]time.Time)
	// iterate through the commands and their responses from the pipeline execution.
	for _, v := range commandMapper {
		if v.Err() == redis.Nil {
			continue
		}
		args := v.Args()
		redisKey := args[1].(string)
		createdAtTimeAsString := v.Val()
		// 2024-01-01T20:01:05.251021+09:00
		layout := "2006-01-02T15:04:05.99999999-07:00"

		parsedTime, err := time.Parse(layout, createdAtTimeAsString)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "server err = " + err.Error(),
			})
		}
		mappedUrlsAndCreatedAt[redisKey] = parsedTime
	}

	commandMapper = map[string]*redis.StringCmd{}
	for _, key := range keys {
		commandMapper[key] = pipe.HGet(database.Ctx, key, "urlHits")
	}
	_, err = pipe.Exec(database.Ctx)
	if err != nil && err != redis.Nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server err = " + err.Error(),
		})
	}
	var mappedUrlAndHits = make(map[string]uint64)
	// iterate through the commands and their responses from the pipeline execution.
	for _, v := range commandMapper {
		if v.Err() == redis.Nil {
			continue
		}
		args := v.Args()
		redisKey := args[1].(string)
		hits, err := strconv.ParseUint(v.Val(), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "server err = " + err.Error(),
			})
		}
		mappedUrlAndHits[redisKey] = hits
	}

	// for getting lifespans for all keys in one go.
	commandMapper = map[string]*redis.StringCmd{}
	for _, key := range keys {
		commandMapper[key] = pipe.HGet(database.Ctx, key, "lifespanInMins")
	}
	_, err = pipe.Exec(database.Ctx)
	if err != nil && err != redis.Nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server err = " + err.Error(),
		})
	}

	var mappedUrlAndLifeSpans = make(map[string]uint64)
	// iterate through the commands and their responses from the pipeline execution.
	for _, v := range commandMapper {
		if v.Err() == redis.Nil {
			continue
		}
		args := v.Args()
		redisKey := args[1].(string)
		hits, err := strconv.ParseUint(v.Val(), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "server err = " + err.Error(),
			})
		}
		mappedUrlAndLifeSpans[redisKey] = hits
	}

	for _, key := range keys {
		var routeResolutionDetails models.UrlAnalyticDetails
		val0, ok1 := mappedUrlAndLifeSpans[key]
		val1, ok2 := mappedUrlsAndCreatedAt[key]
		if ok1 && ok2 {
			timeSinceCreatedAt := uint64(time.Since(val1).Minutes())
			routeResolutionDetails.TTL = val0 - (timeSinceCreatedAt)
		}
		val2, ok := mappedShortenedUrlAndUrl[key]
		if ok {
			routeResolutionDetails.URL = val2
		}
		val3, ok := mappedUrlAndHits[key]
		if ok {
			routeResolutionDetails.Hits = val3
		}
		shortenedUrlKey := "http://" + os.Getenv("DOMAIN") + "/" + key
		resp[shortenedUrlKey] = routeResolutionDetails
	}

	sortedResponse := helpers.SortResponse(resp)

	return c.Status(fiber.StatusOK).JSON(sortedResponse)
}
