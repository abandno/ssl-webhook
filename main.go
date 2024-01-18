package main

import (
	"embed"
	"fmt"
	_ "github.com/CodyGuo/godaemon"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"os"
	"ssl-webhook/src"
)

//go:embed view
var view embed.FS

func main() {
	configLog()
	config := src.GetConfig()
	//view := packr.NewBox("./view")
	//indexHtmlContent, _ := view.FindString("index.html")
	//log.Println("Load view OK:\n", indexHtmlContent[0:30])

	log.Println("==== ssl-webhook ====")

	r := gin.Default()
	//r.LoadHTMLGlob("view/*")
	viewfs, _ := fs.Sub(view, "view")
	//http.Handle("/", http.FileServer(http.FS(viewfs)))
	r.GET("/", func(c *gin.Context) {
		//c.HTML(200, "index.html", nil)
		c.Header("Content-Type", "text/html; charset=utf-8")
		//c.String(200, indexHtmlContent)
		c.FileFromFS("/", http.FS(viewfs))
	})

	r.GET(config.ContextPath+"/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	src.Initialize(r)
	fmt.Println("http://localhost:10010/")
	fmt.Printf("http://localhost:10010%s/ping\n", config.ContextPath)
	r.Run(":10010")
}

func configLog() {
	file, err := os.OpenFile("log.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
}
