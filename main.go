package main

import (
	"fmt"
	_ "github.com/CodyGuo/godaemon"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"ssl-webhook/src"
)

func main() {
	file, err := os.OpenFile("log.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)
	log.Println("==== ssl-webhook ====")

	r := gin.Default()
	r.GET(src.CONTEXT_PATH+"/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	src.Initialize(r)

	fmt.Printf("http://localhost:10010%s/ping\n", src.CONTEXT_PATH)
	r.Run(":10010")
}
