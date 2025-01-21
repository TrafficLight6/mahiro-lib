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
	Mysql    DatabaseConfig            `yaml:"mysql"`
	LibProxy LibProxyApplicationConfig `yaml:"lib-proxy"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}

type LibProxyApplicationConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type UserCheckJson struct {
	AdminRight string `json:"admin"`
	Message    string `json:"message"`
	Success    bool   `json:"success"`
}

type BookList struct {
	Id        int    `json:"id" db:"id"`
	BookName  string `json:"book_name" db:"book_name"`
	BookCover string `json:"book_cover" db:"book_cover"`
	Type      string `json:"book_type" db:"type"`
	Vision    string `json:"book_vision" db:"vision"`
	BookHash  string `json:"book_hash" db:"hash"`
}

type ChapterList struct {
	Id              int    `json:"id" db:"id"`
	BookId          int    `json:"book_id" db:"book_id"`
	ChapterName     string `json:"chapter_name" db:"name"`
	ChapterHash     string `json:"chapter_hash" db:"hash"`
	ChapterFileList string `json:"chapter_file_list" db:"file_list"`
}

type ResultChapterList struct {
	Id              int      `json:"id" db:"id"`
	BookId          int      `json:"book_id" db:"book_id"`
	ChapterName     string   `json:"chapter_name" db:"name"`
	ChapterHash     string   `json:"chapter_hash" db:"hash"`
	ChapterFileList []string `json:"chapter_file_list" db:"file_list"`
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
			url := "http://" + config.LibProxy.Host + ":" + strconv.Itoa(config.LibProxy.Port) + "/user/check/" + "?token=" + token
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
					bookCover := c.Query("book_cover")
					var bookVision string
					if responeInfomation.AdminRight == "true" {
						bookVision = "true"
					} else {
						bookVision = "false"
					}
					db := connectMysql()
					sqlString := "INSERT INTO gbl_book(id,book_name,book_cover,type,vision,hash) VALUES(DEFAULT,?,?,?,?,?)"
					bookHash = hashSha256(bookName + bookType + strconv.Itoa(int(time.Now().Unix())))
					_, err := db.Exec(sqlString, bookName, bookCover, bookType, bookVision, bookHash)
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
			url := "http://" + config.LibProxy.Host + ":" + strconv.Itoa(config.LibProxy.Port) + "/user/check/" + "?token=" + token
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
					bookCover := c.Query("book_cover")
					bookHash := c.Query("book_hash")
					if bookHash == "" {
						resultMessage = "no book hash"
						success = false
					} else {
						if responeInfomation.AdminRight == "true" {
							db := connectMysql()
							sqlString := "UPDATE gbl_book SET book_name=?,book_cover=?,type=?,vision=? WHERE hash=?"
							_, err := db.Exec(sqlString, bookName, bookCover, bookType, bookVision, bookHash)
							defer db.Close()
							if err != nil {
								resultMessage = "fail in editing book information"
								success = false
							} else {
								resultMessage = "success.and now the name,type and vision of the book are `" + bookName + "`,`" + bookCover + "`, `" + bookType + "` and `" + bookVision + "`"
								success = true
							}
						} else {
							resultMessage = "not administrator.beacuse of no right"
							success = false
						}
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
			url := "http://" + config.LibProxy.Host + ":" + strconv.Itoa(config.LibProxy.Port) + "/user/check/" + "?token=" + token
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

		keyWord := c.Query("key")
		if keyWord == "" {
			db := connectMysql()
			sqlString := "SELECT id,book_name,book_cover,type,vision,hash FROM gbl_book"
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
			sqlString := "SELECT id,book_name,book_cover,type,vision,hash FROM gbl_book WHERE book_name LIKE ?"
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

	r.GET("/book/gethash/", func(c *gin.Context) {
		db := connectMysql()
		var bookId int
		var books []BookList
		var resultMessage string
		var success bool
		var bookHash string

		bookId, err := strconv.Atoi(c.Query("book_id"))
		if err != nil {
			resultMessage = "fail in reading book id"
			success = false
			bookHash = ""
		} else {
			db.Select(&books, "SELECT id,book_name,book_cover,type,vision,hash FROM gbl_book WHERE id = ?", bookId)
			if books == nil {
				resultMessage = "can not selecting a book which id is `" + strconv.Itoa(bookId) + "`"
				success = false
				bookHash = ""
			} else {
				resultMessage = "success"
				success = true
				bookHash = books[0].BookHash
			}
		}
		c.JSON(200, gin.H{
			"message": resultMessage,
			"success": success,
			"hash":    bookHash,
		})
	})

	r.GET("/book/get/", func(c *gin.Context) {
		db := connectMysql()
		bookHash := c.Query("book_hash")

		var resultMessage string
		var success bool
		var result BookList

		if bookHash == "" {
			resultMessage = "Arugument cannot be empty"
			success = false
			result = BookList{}
		} else {
			var books []BookList
			db.Select(&books, "SELECT id,book_name,book_cover,type,vision,hash FROM gbl_book WHERE hash = ?", bookHash)
			if books == nil {
				resultMessage = "Can not selecting a book which hash is `" + bookHash + "`"
				success = false
				result = BookList{}
			} else {
				resultMessage = "success"
				success = true
				result = books[0]
			}
		}
		c.JSON(200, gin.H{
			"message": resultMessage,
			"success": success,
			"book":    result,
		})
	})

	//----------------------------chapter api----------------------------------
	r.GET("/book/chapter/get/", func(c *gin.Context) {
		db := connectMysql()
		chapterHash := c.Query("chapter_hash")
		var chapterList []ChapterList
		var chapterFileList []string
		var success bool
		var resultMessage string
		var reasult ResultChapterList

		db.Select(&chapterList, "SELECT id,book_id,name,hash,file_list FROM gbl_chapter WHERE hash = ?", chapterHash)
		if chapterList == nil {
			success = false
			resultMessage = "Can not selecting a chapter which hash is `" + chapterHash + "`"
			reasult = ResultChapterList{}
		} else {
			jsonData := chapterList[0].ChapterFileList
			err := json.Unmarshal([]byte(jsonData), &chapterFileList)
			if err != nil {
				success = false
				resultMessage = "Fail in json unmarshal"
				reasult = ResultChapterList{}
			} else {
				success = true
				resultMessage = "success"
				reasult = ResultChapterList{
					Id:              chapterList[0].Id,
					BookId:          chapterList[0].BookId,
					ChapterName:     chapterList[0].ChapterName,
					ChapterHash:     chapterList[0].ChapterHash,
					ChapterFileList: chapterFileList,
				}
			}
		}
		c.JSON(200, gin.H{
			"message": resultMessage,
			"success": success,
			"chapter": reasult,
		})
	})

	r.GET("/book/chapter/gethash/", func(c *gin.Context) {
		db := connectMysql()
		var chapterList []ChapterList
		var resultMessage string
		var success bool
		var chapterHash string
		chapterId := c.Query("id")
		if chapterId == "" {
			resultMessage = "Argument cannot be empty"
			success = false
			chapterHash = ""
		} else {
			db.Select(&chapterList, "SELECT id,book_id,name,hash,file_list FROM gbl_chapter WHERE id = ?", chapterId)
			if chapterList == nil {
				resultMessage = "Can not selecting a chapter which id is `" + chapterId + "`"
				success = false
				chapterHash = ""
			} else {
				resultMessage = "success"
				success = true
				chapterHash = chapterList[0].ChapterHash
			}
		}
		c.JSON(200, gin.H{
			"message": resultMessage,
			"success": success,
			"hash":    chapterHash,
		})
	})

	r.POST("/book/chapter/add/", func(c *gin.Context) {
		var resultMessage string
		var success bool
		var config Config

		db := connectMysql()
		chapterHash := "none"
		token, err := c.Cookie("token")
		file_list := c.Query("file_list")
		bookHash := c.Query("book_hash")
		ChapterName := c.Query("chapter_name")
		if err != nil {
			resultMessage = "Fail in get cookie"
			success = false
		} else {
			config = readConfig()
			url := "http://" + config.LibProxy.Host + ":" + strconv.Itoa(config.LibProxy.Port) + "/user/check/" + "?token=" + token
			respone := httpGetRequest(url)
			var responeInfomation UserCheckJson
			err = json.Unmarshal([]byte(respone), &responeInfomation)
			if err != nil {
				resultMessage = "Fail in json unmarshal"
				success = false
			} else {
				if responeInfomation.AdminRight == "false" {
					resultMessage = "You are not Admin"
					success = false
				} else {
					if file_list == "" || bookHash == "" || ChapterName == "" {
						resultMessage = "Arguments are not enough and empty"
						success = false
					} else {
						var books []BookList
						db.Select(&books, "SELECT id,book_name,book_cover,type,vision,hash FROM gbl_book WHERE hash = ?", bookHash)
						if books == nil {
							resultMessage = "Can not adding a chapter in the book which hash is `" + bookHash + "`"
							success = false
						} else {
							chapterHash = hashSha256(bookHash + file_list + strconv.Itoa(int(time.Now().Unix())))
							sqlString := "INSERT INTO gbl_chapter(id,book_id,name,hash,file_list) VALUES(DEFAULT,?,?,?,?)"
							_, err = db.Exec(sqlString, books[0].Id, ChapterName, chapterHash, file_list)
							if err != nil {
								resultMessage = "Fail in inserting"
								success = false
							} else {
								resultMessage = "Success.The hash of new chapter is `" + chapterHash + "` in the book which hash is `" + bookHash + "`"
								success = true
							}
						}
					}

				}
			}

		}
		c.JSON(200, gin.H{
			"message":      resultMessage,
			"success":      success,
			"chapter_hash": chapterHash,
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
