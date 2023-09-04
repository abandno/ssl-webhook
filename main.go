package main

import (
	"embed"
	"fmt"
	_ "github.com/CodyGuo/godaemon"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	"log"
	"os"
	"ssl-webhook/src"
)

//go:embed view/*
var viewEmbed embed.FS

func main() {
	//indexHtmlContent, err := viewEmbed.ReadFile("view/index.html")
	//if err != nil {
	//	// 处理错误
	//	fmt.Println(err.Error())
	//}
	view := packr.NewBox("./view")
	indexHtmlContent, err := view.FindString("index.html")
	log.Println(indexHtmlContent[0:30])

	file, err := os.OpenFile("log.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)
	log.Println("==== ssl-webhook ====")

	r := gin.Default()
	//r.LoadHTMLGlob("view/*")
	r.GET("/", func(c *gin.Context) {
		//c.HTML(200, "index.html", nil)
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, indexHtmlContent)
	})

	r.GET(src.CONTEXT_PATH+"/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	src.Initialize(r)
	fmt.Println("http://localhost:10010/")
	fmt.Printf("http://localhost:10010%s/ping\n", src.CONTEXT_PATH)
	r.Run(":10010")
}
