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
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"config": client.cfg,
		})
	})

	r.POST("/add-names", func(c *gin.Context) {
		err := client.AddNames()
		if err == nil {
			c.JSON(http.StatusCreated, gin.H{
				"status": "ok",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  err,
			})
		}
	})

	r.GET("/get-names", func(c *gin.Context) {
		names, err := client.GetNames()
		if err == nil {
			c.JSON(http.StatusOK, names)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  err,
			})
		}
	})
	return r
}

func SetupClient() (*Client, error) {
	client := &Client{}

	err := client.Init(ctx, &Config{
		Instance: fmt.Sprintf("%s-instance", ProjectName),
		Project:  "piotrostr-resources",
		Database: fmt.Sprintf("%s-db", ProjectName),
		Table:    fmt.Sprintf("%s-table", ProjectName),
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
