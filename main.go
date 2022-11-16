package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const ProjectName = "global_spanner_api"

var ctx = context.Background()

func SetupRouter(client *Client) *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.GET("/", func(c *gin.Context) {
		if err := client.CheckHealth(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"config": client.cfg,
		})
	})

	r.POST("/add-names", func(c *gin.Context) {
		err := client.AddNames()
		if err == nil {
			c.JSON(http.StatusCreated, nil)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
		}
	})

	r.GET("/get-names", func(c *gin.Context) {
		names, err := client.GetNames()
		if err == nil {
			c.JSON(http.StatusOK, names)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}
	})

	return r
}

func SetupClient() (*Client, error) {
	client := &Client{}

	err := client.Init(ctx, &Config{
		Instance: fmt.Sprintf("%s_instance", ProjectName),
		Project:  "piotrostr-resources",
		Database: fmt.Sprintf("%s_db", ProjectName),
		Table:    fmt.Sprintf("%s_table", ProjectName),
	})
	if err != nil {
		log.Fatal(err)
	}

	return client, err
}

func main() {
	client, err := SetupClient()
	if err != nil {
		log.Fatal(err)
	}

	router := SetupRouter(client)

	if err = router.Run(":8080"); err != nil {
		client.Teardown()
		log.Fatal(err)
	}
}
