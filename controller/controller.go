package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"server_go/cache"
	"server_go/model"
	"server_go/util"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	model   *model.Model
	cache   *cache.Cache
	baseUrl string
}

func (this *Controller) Init() {
	this.baseUrl = os.Getenv("DOMAIN")
	if this.baseUrl == "" {
		this.baseUrl = "http://localhost:8080/"
	}
	this.model = new(model.Model)
	this.cache = new(cache.Cache)
	this.model.Connect()
	this.cache.Connect()
	// Routine to delete expired record and reset ID
	go func() {
		for {
			// err := this.model.DeleteExpiredRecord()
			// if err != nil {
			// 	log.Println(err.Error())
			// }
			currentID, err := this.model.GetMaxID()
			if err != nil {
				log.Println(err.Error())
			}
			this.cache.Set("CurrentID", currentID, -1)
			log.Println("Delete expired record and reset ID!")
			time.Sleep(24 * time.Hour)
		}
	}()
}
func (this *Controller) Close() {
	this.model.Close()
	this.cache.Flush()
}

func (this *Controller) GetNextID() string {
	return fmt.Sprint(this.cache.Increase("CurrentID"))
}

func (this *Controller) ErrorController(ctx *fiber.Ctx, err error) error {
	// Default 500 statuscode
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		// Override status code if fiber.Error type
		code = e.Code
	}
	// Set Content-Type: text/plain; charset=utf-8
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	log.Println(err.Error())
	// Return statuscode with error message
	if code == 404 {
		return ctx.Status(404).Render("404", nil)
	}
	return ctx.Status(500).Render("500", nil)
}

func (this *Controller) ValidateController(ctx *fiber.Ctx) error {
	requestData := new(model.Url)
	requestData.Url = ctx.Params("url")
	if len(requestData.Url) < 5 {
		return ctx.Status(fiber.StatusNotFound).Render("404", nil)
	}
	urlPart := requestData.Url[:len(requestData.Url)-4]
	userSignature := requestData.Url[len(requestData.Url)-4:]
	systemSignature := util.SignUrl(urlPart)
	if userSignature != systemSignature {
		return ctx.Status(fiber.StatusNotFound).Render("404", nil)
	}
	return ctx.Next()
}
func (this *Controller) GetIndexController(ctx *fiber.Ctx) error {
	return ctx.Render("index", &fiber.Map{})
}
func (this *Controller) PostGenUrlController(ctx *fiber.Ctx) error {
	requestData := new(model.Url)
	if err := ctx.BodyParser(requestData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"url": nil, "error": err.Error()})

	}

	// check is a valid url
	_, err := url.ParseRequestURI(requestData.Url)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"url": nil, "error": err.Error()})
	}
	// find in database
	shortUrl, err := this.cache.Get(requestData.Url)

	// found in cache
	if err == nil {
		return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{"url": this.baseUrl + shortUrl, "error": nil})
	}
	// err that is not "not found"
	if err != redis.Nil {
		return err
	}
	// Find in database
	urlRecord, err := this.model.FindLongUrl(requestData.Url)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if urlRecord.ShortUrl != "" && urlRecord.ExpireTime.Before(time.Now().UTC()) {
		// Expire
		err = this.cache.Set(urlRecord.ShortUrl, urlRecord.LongUrl, 24)
		err = this.cache.Set(urlRecord.LongUrl, urlRecord.ShortUrl, 24)
		if err != nil {
			return err
		}
		return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{"url": this.baseUrl + urlRecord.ShortUrl, "error": nil})
	}
	// insert DB
	channelModel := make(chan struct{})
	channelCache := make(chan struct{})
	newID := this.GetNextID()
	var errModel, errCache error
	newShortUrl := util.GenerateShortLink(newID)
	go func() {
		errModel = this.model.InsertUrl(newID, newShortUrl, requestData.Url, time.Now().AddDate(0, 0, 3))
		channelModel <- struct{}{}
	}()
	go func() {
		errCache = this.cache.Set(newShortUrl, requestData.Url, 24)
		if errCache != nil {
			channelCache <- struct{}{}
			return
		}
		errCache = this.cache.Set(requestData.Url, newShortUrl, 24)
		channelCache <- struct{}{}
	}()
	<-channelCache
	<-channelModel
	if errModel != nil {
		return err
	}
	if errCache != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{"url": this.baseUrl + newShortUrl, "error": nil})
}
func (this *Controller) GetUrlController(ctx *fiber.Ctx) error {
	requestData := new(model.Url)
	requestData.Url = ctx.Params("url")
	longUrl, err := this.cache.Get(requestData.Url)
	if err == nil {
		return ctx.Redirect(longUrl)
	}
	if err != redis.Nil {
		return err
	}
	urlRecord, err := this.model.FindShortUrl(requestData.Url)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).Render("404", nil)
		}
		return err
	}
	// Url Expire
	if urlRecord.ExpireTime.Before(time.Now()) {
		return ctx.Status(fiber.StatusGone).Render("410", nil)
	}
	err = this.cache.Set(urlRecord.ShortUrl, urlRecord.LongUrl, 24)
	err = this.cache.Set(urlRecord.LongUrl, urlRecord.ShortUrl, 24)
	if err != nil {
		return err
	}
	return ctx.Redirect(longUrl)
}

func (this *Controller) GetResetCache(ctx *fiber.Ctx) error {
	err := this.cache.Flush()
	if err != nil {
		return err
	}
	log.Println("Flush cache!")
	return ctx.JSON(&fiber.Map{"error": nil})
}
func (this *Controller) GetResetDB(ctx *fiber.Ctx) error {
	err := this.model.DeleteUrl("", "")
	if err != nil {
		return err
	}
	currentID, _ := this.model.GetMaxID()
	this.cache.Set("CurrentID", currentID, -1)
	log.Println("Reset max valid url ID!")
	return ctx.JSON(&fiber.Map{"error": nil})
}

// use this or run routine
func (this *Controller) GetResetID(ctx *fiber.Ctx) error {
	currentID, err := this.model.GetMaxID()
	this.cache.Set("CurrentID", currentID, -1)
	log.Println("Reset max valid url ID!")
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{"error": err.Error()})
}
