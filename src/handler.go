package src

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func Initialize(engine *gin.Engine) {
	config := GetConfig()
	engine.GET(config.ContextPath+"/hello", func(c *gin.Context) {
		log.Println("/hello")
		c.JSON(http.StatusOK, gin.H{
			"message": "world",
		})
	})

	ohttps := engine.Group(config.ContextPath + "/ohttps")
	{
		ohttps.GET("/hello", func(c *gin.Context) {
			log.Println("/ohttps/hello")
			c.JSON(http.StatusOK, gin.H{
				"success": "hello ohttps",
			})
		})
		// [OHTTPS - 免费HTTPS证书、自动更新、自动部署](https://ohttps.com/docs/cloud/webhook/webhook)
		ohttps.POST("/deploy", func(c *gin.Context) {
			log.Println("/ohttps/deploy")
			//request := make(map[string]interface{})
			var request OhttpsSslDeployRequest
			fmt.Println(c.Query("a")) // url 参数
			err := c.ShouldBind(&request)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   err.Error(),
				})
				log.Panic("[ERROR] ", err)
				return
			}

			log.Println(request.Payload.CertificateName, request.Payload.CertificateDomains)
			// 验签
			md5 := md5.New()
			io.WriteString(md5, strconv.FormatInt(request.Timestamp, 10)+":"+config.CallbackToken)
			md5sum := hex.EncodeToString(md5.Sum(nil))
			if md5sum != request.Sign {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   "签名验证不通过",
				})
				log.Panic("[ERROR] 签名验证不通过, request.Sign: ", request.Sign, ", request.Timestamp: ", request.Timestamp, ", md5sum: ", md5sum)
				return
			}

			// 按关联的域名, 生成对应的目录 先带 .tmp 后缀,
			// 然后, 正在使用的目录更名带 日期后缀 备份, 新的去掉 tmp 后缀
			for _, domain := range request.Payload.CertificateDomains {
				domainCertPath := domain
				if strings.HasPrefix(domain, "*") {
					// 泛域名
					domainCertPath = domain[2:]
				}
				tmpCertPath := config.NginxCertBasePath + "/" + domainCertPath + ".tmp"
				os.MkdirAll(tmpCertPath, os.ModePerm)

				certKeyFile, _ := os.OpenFile(tmpCertPath+"/cert.key", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
				defer certKeyFile.Close()
				certKeyFile.WriteString(request.Payload.CertificateCertKey)

				fullchainFile, _ := os.OpenFile(tmpCertPath+"/fullchain.cer", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
				defer fullchainFile.Close()
				fullchainFile.WriteString(request.Payload.CertificateFullchainCerts)

				// 原目录备份
				os.Rename(config.NginxCertBasePath+"/"+domainCertPath, config.NginxCertBasePath+"/"+domainCertPath+"."+time.Now().Format("20060102150405"))
				// 启用新的
				os.Rename(tmpCertPath, tmpCertPath[:len(tmpCertPath)-4])
				log.Printf("部署 %s 成功, 路径: %s\n", domain, tmpCertPath[:len(tmpCertPath)-4])
			}

			// nginx reload 新证书才能生效
			msg, err := execNginxReload()

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"msg":     msg,
			})

		})
	}

	nginx := engine.Group(config.ContextPath + "/nginx")
	{
		//nginx -s reload
		nginx.GET("/reload", func(c *gin.Context) {
			msg, err := execNginxReload()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"msg":     err.Error(),
				})
				log.Printf("`nginx -s reload` err: %s\n", err)
			} else {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"msg":     msg,
				})
			}
		})
	}

}

func execNginxReload() (string, error) {
	cmd := exec.Command("nginx", "-s", "reload")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("`nginx -s reload` err: %s\n", err)
		return err.Error(), err
	}
	log.Printf("`nginx -s reload`: %s\n", string(out))
	return string(out), nil
}

type OhttpsSslDeployRequest struct {
	Timestamp int64 `form:"timestamp"` // 请求时间戳
	Payload   struct {
		CertificateName           string   `form:"certificateName"`           // 证书ID
		CertificateDomains        []string `form:"certificateDomains"`        // 证书关联域名
		CertificateCertKey        string   `form:"certificateCertKey"`        // 证书私钥(PEM格式)
		CertificateFullchainCerts string   `form:"certificateFullchainCerts"` // 证书(包含证书和中间证书)(PEM格式)
		CertificateExpireAt       int64    `form:"certificateExpireAt"`       // 证书过期时间
	} `form:"payload"`
	Sign string `form:"sign"` // 请求签名，`${timestamp}:${回调令牌}`的32位小写md5值
}
