package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Mysql   DatabaseConfig           `yaml:"mysql"`
	LibUser LibUserApplicationConfig `yaml:"lib-user"`
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

type GblConfigTable struct {
	Id          int    `db:"id"`
	ConfigKey   string `db:"key"`
	ConfigValue string `db:"value"`
}

type UserCheckJson struct {
	AdminRight string `json:"admin"`
	Message    string `json:"message"`
	Success    bool   `json:"success"`
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "go book libary api.config module here",
		})
	})

	//Get config value by config key
	r.GET("/config/get/", func(c *gin.Context) {
		db := connectMysql()
		applicationConfigKey := c.Query("key")
		sqlString := "SELECT id,`key`,value FROM gbl_config WHERE `key` = ?"

		var resultMessage string
		var success bool
		var applicationConfigs []GblConfigTable

		if applicationConfigKey == "" {
			resultMessage = "no key"
			success = false
		} else {
			err := db.Select(&applicationConfigs, sqlString, applicationConfigKey)
			defer db.Close()
			if err != nil {
				resultMessage = "fail"
				success = false
			} else {
				resultMessage = applicationConfigs[0].ConfigValue
				success = true
			}

		}

		c.JSON(200, gin.H{
			"message": resultMessage,
			"success": success,
		})
	})

	//Add config value by config key
	//Administrator only
	r.POST("/config/add/", func(c *gin.Context) {
		token, err := c.Cookie("token")
		config := readConfig()

		var resultMessage string
		var success bool

		//Read cookie here
		if err != nil {
			resultMessage = "fail in reading cookie"
			success = false
		} else {
			//Get the user information
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
					if responeInfomation.AdminRight == "true" {
						db := connectMysql()
						key := c.Query("key")
						value := c.Query("value")
						sqlString := "SELECT id,`key`,value FROM gbl_config WHERE `key` = ?"

						var applicationConfigs []GblConfigTable

						err := db.Select(&applicationConfigs, sqlString, key)
						if err != nil {
							resultMessage = "fail in check config key"
							success = false
						} else {
							if applicationConfigs == nil {
								//If not duplicate naming
								sqlString = "INSERT INTO gbl_config(id,`key`,value) VALUES(DEFAULT,?,?)"
								_, err = db.Exec(sqlString, key, value)
								defer db.Close()
								if err != nil {
									//If add fail
									resultMessage = "fail in add new config"
									success = false
								} else {
									resultMessage = "add success.key is `" + key + "`.value is `" + value + "`."
									success = true
								}
							} else {
								defer db.Close()
								resultMessage = "duplicate naming"
								success = false
							}
						}
					} else {
						//If User in not adminstrator
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

	//Edit config value by config key
	//Administrator only
	r.POST("/config/edit/", func(c *gin.Context) {
		token, err := c.Cookie("token")
		config := readConfig()

		var resultMessage string
		var success bool

		//Read cookie here
		if err != nil {
			resultMessage = "fail in read cookie"
			success = false
		} else {
			//Get the user information
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
					if responeInfomation.AdminRight == "true" {
						db := connectMysql()
						key := c.Query("key")
						value := c.Query("value")
						sqlStrig := "SELECT id,`key`,value FROM gbl_config WHERE `key` = ?"

						var applicationConfigs []GblConfigTable

						err := db.Select(&applicationConfigs, sqlStrig, key)
						if err != nil {
							resultMessage = "fail in check config key"
							success = false
						} else {
							if applicationConfigs == nil {
								resultMessage = "did not find config called`" + key + "`"
								success = false
							} else {
								sqlStrig = "UPDATE gbl_config SET value = ? WHERE `key` = ?"
								_, err = db.Exec(sqlStrig, value, key)
								if err != nil {
									resultMessage = "fail in editing"
									success = false
								} else {
									resultMessage = "edit successfully.the value of config `" + key + "` is `" + value + "` now."
									success = true
								}
							}
						}
					} else {
						//If User in not adminstrator
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

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "404",
		})
	})

	r.Run(":8081")
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

// func httpPostRequest(url string, contentType string, argument string) string {
// 	payload := strings.NewReader(argument)
// 	response, err := http.Post(url, contentType, payload)
// 	if err != nil {
// 		return "none"
// 	}
// 	defer response.Body.Close()
// 	body, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		return "fail"
// 	}
// 	return string(body)
// }
