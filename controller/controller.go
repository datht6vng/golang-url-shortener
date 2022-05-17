package controller

import (
	"database/sql"
	"net/url"
	"server_go/cache"
	"server_go/model"
	"server_go/util"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

var baseUrl = "http://localhost/"
var Model = model.Model{}
var Cache = cache.Cache{}

func GetIndexController(ctx *fiber.Ctx) error {
	return ctx.Render("index", &fiber.Map{})
}
func PostGenUrlController(ctx *fiber.Ctx) error {
	requestData := model.Url{}
	if err := ctx.BodyParser(&requestData); err != nil {
		return ctx.JSON(&fiber.Map{"url": nil, "error": err.Error()})
	}
	// check is a valid url
	_, err := url.ParseRequestURI(requestData.Url)
	if err != nil {
		return ctx.JSON(&fiber.Map{"url": nil, "error": err.Error()})
	}
	// find in database
	// need to find in cache ??
	shortUrl, err := Cache.Get(requestData.Url)

	// found in cache
	if err == nil {
		return ctx.JSON(&fiber.Map{"url": baseUrl + shortUrl, "error": nil})
	}
	// err that is not "not found"
	if err != redis.Nil {
		return ctx.JSON(&fiber.Map{"url": nil, "error": err.Error()})
	}
	// Find in database
	urlRecord, err := Model.FindLongUrl(requestData.Url)
	if err != nil && err != sql.ErrNoRows {
		return ctx.JSON(&fiber.Map{"url": nil, "error": err.Error()})
	}
	if urlRecord.ShortUrl != "" {
		// Expire

		err = Cache.Set(urlRecord.ShortUrl, urlRecord.LongUrl, 24)
		err = Cache.Set(urlRecord.LongUrl, urlRecord.ShortUrl, 24)
		if err != nil {
			return ctx.JSON(&fiber.Map{"url": nil, "error": err.Error()})
		}
		return ctx.JSON(&fiber.Map{"url": urlRecord.ShortUrl, "error": nil})
	}
	// insert DB
	newShortUrl := util.GenerateShortLink(requestData.Url)
	channelModel := make(chan struct{})
	channelCache := make(chan struct{})
	go func() {
		err = Model.InsertUrl(newShortUrl, requestData.Url, time.Now().UTC().AddDate(0, 0, 3), 0)
		channelModel <- struct{}{}
	}()

	go func() {
		err = Cache.Set(newShortUrl, requestData.Url, 24)
		err = Cache.Set(requestData.Url, newShortUrl, 24)
		channelCache <- struct{}{}
	}()
	<-channelCache
	<-channelModel
	return ctx.JSON(&fiber.Map{"url": baseUrl + newShortUrl, "error": err})
}
func GetUrlController(ctx *fiber.Ctx) error {
	requestData := model.Url{}
	requestData.Url = ctx.Params("url")

	result, err := Model.FindShortUrl(requestData.Url)
	if err != nil {
		return ctx.JSON(&fiber.Map{"url": result.LongUrl, "error": err.Error()})
	}
	return ctx.JSON(&fiber.Map{"url": result.LongUrl, "error": nil})
}
func InitDBController(ctx *fiber.Ctx) error {
	err := Model.CreateModel()
	if err != nil {
		return ctx.JSON(&fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(&fiber.Map{"error": nil})
}
func GetReset(ctx *fiber.Ctx) error {
	err := Cache.Flush()
	if err != nil {
		return ctx.JSON(&fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(&fiber.Map{"error": nil})
}
