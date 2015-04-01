package main

import (
	"bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/thoas/stats"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var Stats = stats.New()

func main() {
	r := gin.New()

	r.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			beginning := time.Now()

			c.Next()

			Stats.End(beginning, c.Writer)
		}
	}())

	r.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, Stats.Data())
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})

	r.Run(":3001")
}
