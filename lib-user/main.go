package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Mysql    DatabaseConfig           `yaml:"mysql"`
	LibProxy LibProxyApplcationConfig `yaml:"lib-config"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}

type LibProxyApplcationConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type GblUserTable struct {
	Id         int    `db:"id"`
	Username   string `db:"username"`
	Password   string `db:"password"`
	AdminRight string `db:"admin"`
}

type GblToeknTable struct {
	Id      int    `db:"id"`
	Userid  int    `db:"user_id"`
	Token   string `db:"token"`
	Dietime int    `db:"dietime"`
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "go book libary api.user module here",
		})
	})

	//Log in.
	//To get the cookie.
	//Cookie will save in cilent.If you have "remember" argument.
	//So if your request do not have that argument,please save it in JS of cilent.
	r.POST("/user/login/", func(c *gin.Context) {
		//new sql connect
		db := connectMysql()

		//get query
		username := c.Query("username")
		password := hashSha256(c.Query("password"))
		remember := c.Query("remember")
		var success bool
		var resultMessage string
		var users []GblUserTable
		var cookie string
		var cookieLife int
		sqlString := "SELECT id FROM gbl_user WHERE username = ? AND password = ?"
		err := db.Select(&users, sqlString, username, password)

		if err != nil {
			success = false
			resultMessage = "fail in selecting user information"
			cookie = ""
			cookieLife = -1
		}
		if users == nil {
			success = false
			resultMessage = "do not have a user call `" + username + "` in the database"
			cookie = ""
			cookieLife = -1
		} else {
			success = true
		}

		if success == true {
			cookie = hashSha256(username + password + strconv.Itoa(int(time.Now().Unix())))
			if remember == "true" {
				cookieLife = 3156000
			} else {
				cookieLife = 86400
			}
			sql := "INSERT INTO gbl_token(id,user_id,token,dietime) VALUES(DEFAULT,?,?,?)"

			_, err := db.Exec(sql, users[0].Id, cookie, int(time.Now().Unix())+cookieLife)
			if err != nil {
				resultMessage = "fail in inserting new token"
				success = false
			} else {
				resultMessage = "success"
				success = true
				// if remember == "true" {
				// 	c.SetCookie("token", cookie, cookieLife, "/", "", false, true)
				// } else {
				// 	c.SetCookie("token", cookie, -1, "/", "", false, true)
				// }
			}
		}
		defer db.Close()
		c.JSON(200, gin.H{
			"message":     resultMessage,
			"success":     success,
			"cookie":      cookie,
			"cookie_life": cookieLife,
		})
	})

	//Check user token.
	//Other modules will use it
	r.GET("/user/check/", func(c *gin.Context) {
		db := connectMysql()
		var tokens []GblToeknTable
		var users []GblUserTable
		var success bool
		var resultMessage string
		var admin string
		cookie := c.Query("token")
		sqlString := "SELECT id,user_id,token,dietime FROM gbl_token WHERE token=?"
		err := db.Select(&tokens, sqlString, cookie)
		if err != nil {
			success = false
			resultMessage = "select error"
		}
		if tokens == nil {
			success = false
			resultMessage = "connot get user infomation"
			admin = "fail"
		} else {
			success = true
			resultMessage = "success"
			sqlString := "SELECT id,username,password,admin FROM gbl_user WHERE id=?"
			err := db.Select(&users, sqlString, tokens[0].Userid)
			if err != nil {
				admin = "fail"
			} else {
				admin = users[0].AdminRight
			}
		}
		defer db.Close()
		c.JSON(200, gin.H{
			"success": success,
			"message": resultMessage,
			"admin":   admin,
		})

	})

	//Add a new user.
	//Administrator only.
	r.POST("/user/add/", func(c *gin.Context) {
		db := connectMysql()
		var tokens []GblToeknTable
		var users []GblUserTable
		var success bool
		var resultMessage string
		cookie, err := c.Cookie("token")
		if err != nil {
			success = false
			resultMessage = "fail in get cookie"
		} else {
			sqlString := "SELECT id,user_id,token,dietime FROM gbl_token WHERE token=?"
			err = db.Select(&tokens, sqlString, cookie)
			if err != nil {
				success = false
				resultMessage = "select error"
			}
			if tokens == nil {
				success = false
				resultMessage = "connot get user infomation"
			} else {
				success = true
				resultMessage = "success"
				sqlString := "SELECT id,username,password,admin FROM gbl_user WHERE id=?"
				err := db.Select(&users, sqlString, tokens[0].Userid)
				if err != nil {
					success = false
					resultMessage = "fail in selecting"
				} else {
					if users[0].AdminRight == "true" {
						userName := c.Query("username")
						password := hashSha256(c.Query("password"))
						adminRight := c.Query("admin")
						sqlString = "INSERT INTO gbl_user (id,username,password,admin) VALUES (DEFAULT,?,?,?)"
						_, err := db.Exec(sqlString, userName, password, adminRight)
						if err != nil {
							resultMessage = "fail in inserting date"
							success = false
						} else {
							resultMessage = "success.the name,password hash ,and administration right of new user are `" + userName + "`,`" + password + "` and `" + adminRight + "`"
							success = true
						}
					} else {
						resultMessage = "not administrator.beacuse of no right"
						success = false
					}
				}
			}
			defer db.Close()

		}
		c.JSON(200, gin.H{
			"success": success,
			"message": resultMessage,
		})

	})

	//Edit user information.
	//Also administrator only.
	//Or if you edit your own account,don't need administration right.But you cannot use it to give you administration right
	r.POST("/user/edit/", func(c *gin.Context) {
		db := connectMysql()
		var tokens []GblToeknTable
		var users []GblUserTable
		var success bool
		var resultMessage string
		cookie, err := c.Cookie("token")
		if err != nil {
			success = false
			resultMessage = "fail in get cookie"
		} else {
			sqlString := "SELECT id,user_id,token,dietime FROM gbl_token WHERE token=?"
			err = db.Select(&tokens, sqlString, cookie)
			if err != nil {
				success = false
				resultMessage = "select error"
			}
			if tokens == nil {
				success = false
				resultMessage = "connot get user infomation"
			} else {
				success = true
				resultMessage = "success"
				sqlString := "SELECT id,username,password,admin FROM gbl_user WHERE id=?"
				err := db.Select(&users, sqlString, tokens[0].Userid)
				if err != nil {
					success = false
					resultMessage = "fail in selecting"
				} else {
					if users[0].AdminRight == "true" || users[0].Id == tokens[0].Userid {
						userId := c.Query("user_id")
						username := c.Query("username")
						password := hashSha256(c.Query("password"))
						var adminRight string
						if users[0].AdminRight == "true" {
							adminRight = c.Query("admin")

						} else {
							adminRight = "false"
						}

						if userId == "" || username == "" || password == "" || adminRight == "" {
							success = false
							resultMessage = "arguments not enough"
						} else {
							sqlString = "UPDATE gbl_user SET username=?,password=?,admin=? WHERE id = ?"
							_, err = db.Exec(sqlString, username, password, adminRight, userId)
							if err != nil {
								success = false
								resultMessage = "fail in updating"
							} else {
								resultMessage = "success.the username,hash of password ,and administration right of the user with the id `" + userId + "` are `" + username + "`,`" + password + "` and `" + adminRight + "` now(if some is null,it means no changing)"
								success = true
							}
						}

					} else {
						resultMessage = "not administrator.beacuse of no right"
						success = false
					}
				}
			}
			defer db.Close()

		}
		c.JSON(200, gin.H{
			"success": success,
			"message": resultMessage,
		})

	})

	//Delet user
	//WHY ALSO ADMINISTRATOR ONLY?F**K!
	r.POST("/user/del/", func(c *gin.Context) {
		db := connectMysql()
		var tokens []GblToeknTable
		var users []GblUserTable
		var success bool
		var resultMessage string
		cookie, err := c.Cookie("token")
		if err != nil {
			success = false
			resultMessage = "fail in get cookie"
		} else {
			sqlString := "SELECT id,user_id,token,dietime FROM gbl_token WHERE token=?"
			err = db.Select(&tokens, sqlString, cookie)
			if err != nil {
				success = false
				resultMessage = "select error"
			}
			if tokens == nil {
				success = false
				resultMessage = "connot get user infomation"
			} else {
				success = true
				resultMessage = "success"
				sqlString := "SELECT id,username,password,admin FROM gbl_user WHERE id=?"
				err := db.Select(&users, sqlString, tokens[0].Userid)
				if err != nil {
					success = false
					resultMessage = "fail in selecting"
				} else {
					if users[0].AdminRight == "true" {
						userId, err := strconv.Atoi(c.Query("user_id"))
						if err != nil {
							resultMessage = "fail in reading user id"
							success = false
						} else {
							if userId == users[0].Id {
								resultMessage = "cannot delete youself"
								success = false
							} else {
								sqlString = "DELETE FROM gbl_user WHERE id = ?"
								_, err := db.Exec(sqlString, userId)
								if err != nil {
									resultMessage = "fail in deleting user"
									success = false
								} else {
									resultMessage = "success"
									success = true
								}
							}
						}
					} else {
						resultMessage = "not administrator.beacuse of no right"
						success = false
					}
				}
			}
			defer db.Close()

		}
		c.JSON(200, gin.H{
			"success": success,
			"message": resultMessage,
		})

	})

	//ALSO GET THE F**KING USERS LIST?F**K YOU TrafficLight6!!!
	r.GET("/user/search/", func(c *gin.Context) {
		db := connectMysql()
		var tokens []GblToeknTable
		var users []GblUserTable
		var userResult []GblUserTable
		var success bool
		var resultMessage string
		cookie, err := c.Cookie("token")
		if err != nil {
			success = false
			resultMessage = "fail in get cookie"
			userResult = nil
		} else {
			sqlString := "SELECT id,user_id,token,dietime FROM gbl_token WHERE token=?"
			err = db.Select(&tokens, sqlString, cookie)
			if err != nil {
				success = false
				resultMessage = "select error"
				userResult = nil
			}
			if tokens == nil {
				success = false
				resultMessage = "connot get user infomation"
				userResult = nil
			} else {
				success = true
				resultMessage = "success"
				sqlString := "SELECT id,username,password,admin FROM gbl_user WHERE id=?"
				err := db.Select(&users, sqlString, tokens[0].Userid)
				if err != nil {
					success = false
					resultMessage = "fail in selecting"
					userResult = nil
				} else {
					if users[0].AdminRight == "true" {
						keyWord := c.Query("key")
						sqlString = "SELECT id,username,password,admin FROM gbl_user WHERE username LIKE ?"
						err = db.Select(&userResult, sqlString, "%"+keyWord+"%")
						if err != nil {
							resultMessage = "fail in selecting user"
							success = false
							userResult = nil
						} else {
							resultMessage = "success"
							success = true

						}
					} else {
						resultMessage = "not administrator.beacuse of no right"
						success = false
						userResult = nil
					}
				}
			}
			defer db.Close()

		}
		c.JSON(200, gin.H{
			"success": success,
			"message": resultMessage,
			"result":  userResult,
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "404",
		})
	})

	r.Run(":8080")
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

func hashSha256(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	bytes := hash.Sum(nil)
	sum := hex.EncodeToString(bytes)
	return sum
}
