package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Mysql      DatabaseConfig              `yaml:"mysql"`
	LibUser    LibUserApplicationConfig    `yaml:"lib-user"`
	LibStorage LibStorageApplicationConfig `yaml:"lib-storage"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}

type LibUserApplicationConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type LibStorageApplicationConfig struct {
	MaxFileSize int `yaml:"maxsize"`
}

type UserCheckJson struct {
	AdminRight string `json:"admin"`
	Message    string `json:"message"`
	Success    bool   `json:"success"`
}

func main() {
	r := gin.Default()
	r.MaxMultipartMemory = int64(readConfig().LibStorage.MaxFileSize)
	r.Static("/storage/get/text", "./storage/text")
	r.Static("/storage/get/picture", "./storage/picture")

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "go book libary api.storage module here",
		})
	})

	//Update .txt file,need token,but not administrator only
	r.POST("/storage/update/text", func(c *gin.Context) {
		var success bool
		var resultMessage string
		var fileList []string

		//Get cookie
		token, err := c.Cookie("token")
		config := readConfig()
		if err != nil {
			success = false
			resultMessage = "fail in getting cookie"
			fileList = nil
		} else {
			//Get user information
			url := "http://" + config.LibUser.Host + ":" + strconv.Itoa(config.LibUser.Port) + "/user/check/" + "?token=" + token
			respone := httpGetRequest(url)

			//Json unmarshal
			var responeInfomation UserCheckJson
			err := json.Unmarshal([]byte(respone), &responeInfomation)
			//If unmarshal fail
			if err != nil {
				resultMessage = "fail in json unmarshal"
				success = false
				fileList = nil
			} else {
				if responeInfomation.Success {
					//Save file
					form, _ := c.MultipartForm()
					files := form.File["file"]
					if files == nil {
						resultMessage = "file form is null"
						success = false
						fileList = nil
					} else {
						db := connectMysql()
						sqlString := "INSERT INTO gbl_storage (id,chapter_id,`type`,`name`) VALUES(DEFAULT,?,?,?)"
						chapterId := c.Query("chapter_id")

						for _, file := range files {
							name := hashSha256(strconv.Itoa(int(time.Now().Unix()))+file.Filename) + "---" + file.Filename
							dst := "./storage/text/" + name
							c.SaveUploadedFile(file, dst)
							fileList = append(fileList, name)
							_, err = db.Exec(sqlString, chapterId, "text", name)
							if err != nil {
								fmt.Println("fail in sql inserting of updating:", err)
							}
						}
						defer db.Close()
						resultMessage = "success"
						success = true

					}

				} else {
					//If user applaction fail to get user information
					resultMessage = "fail in user application get user information"
					success = false
				}
			}
		}

		c.JSON(200, gin.H{
			"message":   resultMessage,
			"success":   success,
			"file_list": fileList,
		})
	})

	//Like the "/storage/update/text".But picture here.
	r.POST("/storage/update/picture", func(c *gin.Context) {
		var success bool
		var resultMessage string
		var fileList []string

		token, err := c.Cookie("token")
		config := readConfig()
		if err != nil {
			success = false
			resultMessage = "fail in getting cookie"
			fileList = nil
		} else {
			url := "http://" + config.LibUser.Host + ":" + strconv.Itoa(config.LibUser.Port) + "/user/check/" + "?token=" + token
			respone := httpGetRequest(url)

			//Json unmarshal
			var responeInfomation UserCheckJson
			err := json.Unmarshal([]byte(respone), &responeInfomation)
			//If unmarshal fail
			if err != nil {
				resultMessage = "fail in json unmarshal"
				success = false
				fileList = nil
			} else {
				if responeInfomation.Success {
					chapterId := c.Query("chapter_id")
					db := connectMysql()
					sqlString := "INSERT INTO gbl_storage (id,chapter_id,`type`,`name`) VALUES(DEFAULT,?,?,?)"
					form, _ := c.MultipartForm()
					files := form.File["file"]
					if files == nil {
						resultMessage = "file form is null"
						success = false
						fileList = nil
					} else {
						for _, file := range files {
							fileName := hashSha256(strconv.Itoa(int(time.Now().Unix()))+file.Filename) + "---" + file.Filename
							dst := "./storage/picture/" + fileName
							c.SaveUploadedFile(file, dst)
							fileList = append(fileList, fileName)
							_, _ = db.Exec(sqlString, chapterId, "picture", fileName)
						}
						defer db.Close()
						resultMessage = "success"
						success = true

					}

				} else {
					//If user applaction fail to get user information
					resultMessage = "fail in user application get user information"
					success = false
				}
			}
		}

		c.JSON(200, gin.H{
			"message":   resultMessage,
			"success":   success,
			"file_list": fileList,
		})
	})

	r.GET("/storage/get/picture", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "picture root diratory here",
			"success": true,
		})
	})

	r.GET("/storage/get/text", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "text root diratory here",
			"success": true,
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "404",
		})
	})

	r.Run(":8082")
}

func readConfig() Config {
	// open yaml
	var config Config
	file, err := os.Open("config.yaml")
	if err != nil {
		return config
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	// decode yaml
	err = decoder.Decode(&config)
	if err != nil {
		return config
	}
	return config
}

func httpGetRequest(url string) string {
	response, err := http.Get(url)
	if err != nil {
		return "fail in request"
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "fail in read respone"
	}
	return string(body)
}

func hashSha256(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	bytes := hash.Sum(nil)
	sum := hex.EncodeToString(bytes)
	return sum
}

func connectMysql() *sqlx.DB {
	var dbObj *sqlx.DB
	config := readConfig()
	database, err := sqlx.Open("mysql", config.Mysql.Username+":"+config.Mysql.Password+"@tcp("+config.Mysql.Host+":"+strconv.Itoa(config.Mysql.Port)+")/"+config.Mysql.Db)
	if err != nil {
		fmt.Println("open mysql failed,", err)
	}
	dbObj = database
	return dbObj
}
