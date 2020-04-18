package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/tomasen/realip"
)

type Requestor struct {
	Referer   int
	UserAgent int
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/ip", func(c *gin.Context) {
		clientIP := realip.FromRequest(c.Request)

		data := map[string]string{
			"ip-address": clientIP,
		}

		r := c.Query("referer")
		if r != "" {
			referer := c.Request.Referer()
			data["referer"] = referer
		}

		ua := c.Query("user-agent")
		if ua != "" {
			userAgent := c.Request.UserAgent()
			data["user-agent"] = userAgent
		}

		c.JSON(http.StatusOK, data)
	})

	router.GET("/referer", func(c *gin.Context) {
		referer := c.Request.Referer()
		c.JSON(http.StatusOK, map[string]string{
			"referer": referer,
		})
	})

	router.GET("/user-agent", func(c *gin.Context) {
		userAgent := c.Request.UserAgent()

		c.JSON(http.StatusOK, map[string]string{
			"user-agent:": userAgent,
		})
	})

	router.Run(":" + port)
}
