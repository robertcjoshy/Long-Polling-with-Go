package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robert/long-polling/handler"
	"github.com/robert/long-polling/pkg"
)

type updates struct {
	Createdat int64
	Message   string `json:"message"`
}

func main() {
	queue := handler.Newcappedqueue[updates](10)
	pub := pkg.Newpubsub()
	router := gin.Default()
	router.GET("/updates/:time", func(ctx *gin.Context) {
		created := ctx.Param("time")
		fmt.Println("time=", created)
		lastupdate, err := strconv.ParseInt(created, 10, 64)

		if err != nil {
			ctx.JSON(400, gin.H{"error": "no timestamp"})
			return
		}

		//var update []updates

		getupdates := func() []updates {
			return pkg.Filter(queue.Copy(), func(update updates) bool {
				return update.Createdat > lastupdate
			})
		}

		// show it user if any update...........
		if update := getupdates(); len(update) > 0 {
			ctx.JSON(200, update)
			return
		}
		ch, close := pub.Subscribe()
		defer close()
		select {
		case <-ctx.Done():
			ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "request timeut"})
			return
		case <-ch:
			ctx.JSON(200, getupdates())
			return
		}
	})

	router.POST("/send", func(ctx *gin.Context) {
		var message updates
		if err := ctx.Bind(&message); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "bad request"})
			return
		}
		fmt.Println(time.Now().Unix())
		queue.Append(updates{
			Createdat: time.Now().Unix(),
			Message:   message.Message,
		})
		pub.Publish()
		ctx.JSON(200, gin.H{"message": "request sent"})
	})
	router.Run()
}
