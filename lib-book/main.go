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
	LibStorage LibStorageApplicationConfig `yaml:"lib-config"`
	LibConfiig LibConfigApplicationConfig  `yaml:"lib-storage"`
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

type LibConfigApplicationConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
type LibStorageApplicationConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type UserCheckJson struct {
	AdminRight string `json:"admin"`
	Message    string `json:"message"`
	Success    bool   `json:"success"`
}

type BookList struct {
	Id       int    `json:"id" db:"id"`
	BookName string `json:"book_name" db:"book_name"`
	Type     string `json:"book_type" db:"type"`
	Vision   string `json:"book_vision" db:"vision"`
	BookHash string `json:"book_hash" db:"hash"`
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "go book libary api.book module here",
		})
	})

	r.POST("/book/add/", func(c *gin.Context) {
		var resultMessage string
		var success bool
		var config Config

		bookHash := "none"
		token, err := c.Cookie("token")
		if err != nil {
			resultMessage = "fail in get cookie"
			success = false
		} else {
			//Get the user information
			config = readConfig()
			url := "http://" + config.LibUser.Host + ":" + strconv.Itoa(config.LibUser.Port) + "/user/check/" + "?token=" + token
			respone := httpGetRequest(url)
			//Json unmarshal
			var responeInfomation UserCheckJson
			err := json.Unmarshal([]byte(respone), &responeInfomation)
			//If unmarshal fail
			if err != nil {
				resultMessage = "fail in json unmarshal"
				success = false
			} else {
				if responeInfomation.Success {
					bookName := c.Query("book_name")
					bookType := c.Query("book_type")
					var bookVision string
					if responeInfomation.AdminRight == "true" {
						bookVision = "true"
					} else {
						bookVision = "false"
					}
					db := connectMysql()
					sqlString := "INSERT INTO gbl_book(id,book_name,type,vision,hash) VALUES(DEFAULT,?,?,?,?)"
					bookHash = hashSha256(bookName + bookType + strconv.Itoa(int(time.Now().Unix())))
					_, err := db.Exec(sqlString, bookName, bookType, bookVision, bookHash)
					defer db.Close()
					if err != nil {
						resultMessage = "fail in inserting book information"
						success = false
					} else {
						resultMessage = "success.the vision of the book is `" + bookVision + "`"
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
			"bookhash": bookHash,
			"message":  resultMessage,
			"success":  success,
		})
	})

	r.POST("/book/edit/", func(c *gin.Context) {
		var resultMessage string
		var success bool
		var config Config
		token, err := c.Cookie("token")
		if err != nil {
			resultMessage = "fail in get cookie"
			success = false
		} else {
			//Get the user information
			config = readConfig()
			url := "http://" + config.LibUser.Host + ":" + strconv.Itoa(config.LibUser.Port) + "/user/check/" + "?token=" + token
			respone := httpGetRequest(url)
			//Json unmarshal
			var responeInfomation UserCheckJson
			err := json.Unmarshal([]byte(respone), &responeInfomation)
			//If unmarshal fail
			if err != nil {
				resultMessage = "fail in json unmarshal"
				success = false
			} else {
				if responeInfomation.Success {
					bookName := c.Query("book_name")
					bookType := c.Query("book_type")
					bookVision := c.Query("book_vision")
					if responeInfomation.AdminRight == "true" {
						db := connectMysql()
						sqlString := "UPDATE gbl_book SET book_name=?,type=?,vision=?"
						_, err := db.Exec(sqlString, bookName, bookType, bookVision)
						defer db.Close()
						if err != nil {
							resultMessage = "fail in editing book information"
							success = false
						} else {
							resultMessage = "success.and now the name,type and vision of the book are `" + bookName + "`, `" + bookType + "` and `" + bookVision + "`"
							success = true
						}
					} else {
						resultMessage = "not administrator.beacuse of no right"
						success = false
					}
				} else {
					//If user applaction fail to get user information
					resultMessage = "fail in user application get user information"
					success = false
				}
			}
		}
		c.JSON(200, gin.H{
			"message": resultMessage,
			"success": success,
		})
	})

	r.POST("/book/del/", func(c *gin.Context) {
		var resultMessage string
		var success bool
		var config Config
		token, err := c.Cookie("token")
		if err != nil {
			resultMessage = "fail in get cookie"
			success = false
		} else {
			//Get the user information
			config = readConfig()
			url := "http://" + config.LibUser.Host + ":" + strconv.Itoa(config.LibUser.Port) + "/user/check/" + "?token=" + token
			respone := httpGetRequest(url)
			//Json unmarshal
			var responeInfomation UserCheckJson
			err := json.Unmarshal([]byte(respone), &responeInfomation)
			//If unmarshal fail
			if err != nil {
				resultMessage = "fail in json unmarshal"
				success = false
			} else {
				if responeInfomation.Success {
					bookHash := c.Query("book_hash")
					if responeInfomation.AdminRight == "true" {
						db := connectMysql()
						sqlString := "DELETE FROM gbl_book WHERE hash=?"
						_, err := db.Exec(sqlString, bookHash)
						defer db.Close()
						if err != nil {
							resultMessage = "fail in deleting book information"
							success = false
						} else {
							resultMessage = "success.or no a book has the hash called `" + bookHash + "`"
							success = true
						}
					} else {
						resultMessage = "not administrator.beacuse of no right"
						success = false
					}
				} else {
					//If user applaction fail to get user information
					resultMessage = "fail in user application get user information"
					success = false
				}
			}
		}
		c.JSON(200, gin.H{
			"message": resultMessage,
			"success": success,
		})
	})

	r.GET("/book/search", func(c *gin.Context) {
		var books []BookList
		var resultMessage string
		var success bool
		// var jsonResult []byte

		keyWord := c.Query("key")
		if keyWord == "" {
			db := connectMysql()
			sqlString := "SELECT id,book_name,type,vision,hash FROM gbl_book"
			err := db.Select(&books, sqlString)
			defer db.Close()
			if err != nil {
				resultMessage = "fail in selecting"
				success = false
			} else {
				resultMessage = "success"
				success = true

			}
		} else {
			db := connectMysql()
			sqlString := "SELECT id,book_name,type,vision,hash FROM gbl_book WHERE book_name LIKE ?"
			err := db.Select(&books, sqlString, "%"+keyWord+"%")
			defer db.Close()
			if err != nil {

				fmt.Println(err)

				resultMessage = "fail in selecting"
				success = false
			} else {
				resultMessage = "success"
				success = true

			}
		}
		c.JSON(200, gin.H{
			"message":       resultMessage,
			"success":       success,
			"search_result": books,
		})

	})

	r.Run(":8083")
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
