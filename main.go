package main

import (
	"errors"
	"fruit_shop/config"
	"fruit_shop/router"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg = pflag.StringP("config", "c", "", "fruit_shop config file path.")
)

func main() {
	pflag.Parse()

	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// set gin mode
	gin.SetMode(viper.GetString("runmode"))
	g := gin.New() // Create the gin engine

	middlewares := []gin.HandlerFunc{}

	router.Load(
		g,
		middlewares...,
	)

	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up", err)
		}
		log.Println("The router has been deployed successfully")
	}()

	log.Printf("Start to listen the incoming requests on http address: %s", viper.GetString("addr"))
	log.Printf(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		log.Println(err)
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		log.Println("Waiting for the router, retry in 1 second")
		time.Sleep(time.Second)
	}
	return errors.New("Can not connect to the router")
}
