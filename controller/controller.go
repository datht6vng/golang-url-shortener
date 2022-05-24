package controller

import (
	"database/sql"
	"fmt"
	"net/url"
	"server_go/cache"
	"server_go/model"
	"server_go/util"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

var BaseUrl string
var Model = model.Model{}
var Cache = cache.Cache{}
var baseTime = time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC).UnixNano()

func InitController(baseUrl string) {
	if baseUrl == "" {
		BaseUrl = "http://localhost:8080/"
	} else {
		BaseUrl = baseUrl
	}
	Model.Connect()
	Cache.Connect()
}
func ValidateController(ctx *fiber.Ctx) error {
	requestData := model.Url{}
	requestData.Url = ctx.Params("url")
	if len(requestData.Url) < 6 {
		return ctx.JSON(&fiber.Map{"error": "Invalid Short Url"})
	}
	urlPart := requestData.Url[:len(requestData.Url)-5]
	userSignature := requestData.Url[len(requestData.Url)-5:]
	systemSignature := util.SignUrl(urlPart)
	if userSignature != systemSignature {
		return ctx.JSON(&fiber.Map{"error": "Invalid Short Url"})
	}
	return ctx.Next()
}
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
		return ctx.JSON(&fiber.Map{"url": BaseUrl + shortUrl, "error": nil})
	}
	// err that is not "not found"
	if err != redis.Nil {
		return ctx.JSON(&fiber.Map{"url": nil, "error": err.Error()})
	}
	// Find in database
	urlRecord, err := Model.FindLongUrl(requestData.Url)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("Fail here")
		return ctx.JSON(&fiber.Map{"url": nil, "error": err.Error()})
	}
	if urlRecord.ShortUrl != "" && urlRecord.ExpireTime.Before(time.Now().UTC()) {
		// Expire
		err = Cache.Set(urlRecord.ShortUrl, urlRecord.LongUrl, 24)
		err = Cache.Set(urlRecord.LongUrl, urlRecord.ShortUrl, 24)
		if err != nil {
			return ctx.JSON(&fiber.Map{"url": nil, "error": err.Error()})
		}
		return ctx.JSON(&fiber.Map{"url": BaseUrl + urlRecord.ShortUrl, "error": nil})
	}
	// insert DB
	newShortUrl := ""
	channelModel := make(chan struct{})
	channelCache := make(chan struct{})
	newID := util.TrimTimeStamp(time.Now().UnixNano(), 10)
	var errModel, errCache error
	go func() {
		newShortUrl = util.GenerateShortLink(newID)
		errModel = Model.InsertUrl(newID, newShortUrl, requestData.Url, time.Now().AddDate(0, 0, 3).UTC(), 0)
		channelModel <- struct{}{}
	}()
	go func() {
		errCache = Cache.Set(newShortUrl, requestData.Url, 24)
		errCache = Cache.Set(requestData.Url, newShortUrl, 24)
		channelCache <- struct{}{}
	}()
	<-channelCache
	<-channelModel
	if errModel != nil {
		return ctx.JSON(&fiber.Map{"url": nil, "error": errModel})
	}
	if errCache != nil {
		return ctx.JSON(&fiber.Map{"url": nil, "error": errCache})
	}
	return ctx.JSON(&fiber.Map{"url": BaseUrl + newShortUrl, "error": nil})
}
func GetUrlController(ctx *fiber.Ctx) error {
	requestData := model.Url{}
	requestData.Url = ctx.Params("url")
	longUrl, err := Cache.Get(requestData.Url)
	if err == nil {
		return ctx.JSON(&fiber.Map{"url": longUrl, "error": nil})
	}
	if err != redis.Nil {
		return ctx.JSON(&fiber.Map{"url": nil, "error": err.Error()})
	}
	urlRecord, err := Model.FindShortUrl(requestData.Url)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(&fiber.Map{"url": urlRecord.LongUrl, "error": "Url Not Found!"})
		}
		return ctx.JSON(&fiber.Map{"url": urlRecord.LongUrl, "error": err.Error()})
	}
	if urlRecord.ExpireTime.Before(time.Now().UTC()) {
		return ctx.JSON(&fiber.Map{"url": "", "error": "Url Is Expired!"})
	}
	err = Cache.Set(urlRecord.ShortUrl, urlRecord.LongUrl, 24)
	err = Cache.Set(urlRecord.LongUrl, urlRecord.ShortUrl, 24)
	return ctx.JSON(&fiber.Map{"url": urlRecord.LongUrl, "error": nil})
}

func GetResetCache(ctx *fiber.Ctx) error {
	err := Cache.Flush()
	if err != nil {
		return ctx.JSON(&fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(&fiber.Map{"error": nil})
}
func GetResetDB(ctx *fiber.Ctx) error {
	err := Model.DeleteUrl("", "")
	if err != nil {
		return ctx.JSON(&fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(&fiber.Map{"error": nil})
}
