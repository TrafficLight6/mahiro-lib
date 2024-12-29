package main

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ProxyUrlList []string `yaml:"proxy-url-list"`
}

func main() {
	r := gin.Default()
	//Get config
	proxyUrlList := readConfig().ProxyUrlList
	for _, proxyUrl := range proxyUrlList {
		targetURL, err := url.Parse(proxyUrl)
		if err != nil {
			panic(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		r.Any(targetURL.Path+"/*proxyPath", func(c *gin.Context) {

			// Modify the request's host header
			c.Request.Host = targetURL.Host
			c.Request.URL.Host = targetURL.Host
			c.Request.URL.Scheme = targetURL.Scheme
			c.Request.URL.Path = strings.Replace(c.Request.URL.Path, targetURL.Path, "", 1)

			fmt.Println(c.Request.Host)
			fmt.Println(targetURL.Path)
			fmt.Println(c.Request.URL.Host)
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":    "Mahiro-lib Proxy here.There are the Proxy List",
			"proxy_list": proxyUrlList,
		})
	})
	// Start the server
	r.Run(":7622")
}

func readConfig() Config {
	// Open yaml
	var config Config
	file, err := os.Open("config.yaml")
	if err != nil {
		return config
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	// Decode yaml
	err = decoder.Decode(&config)
	if err != nil {
		return config
	}
	return config
}
