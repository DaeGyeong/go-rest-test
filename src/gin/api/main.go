package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Naver_api_url struct {
	Api_seq          int
	Api_name         string `json:"api_name"`
	Api_type         string `json:"api_type"`
	Api_url          string `json:"api_url"`
	Api_sub_url      string `json:"api_sub_url"`
	Api_content_type string `json:"api_content_type"`
	Api_method       string `json:"api_method"`
	Api_status       string `json:"api_status"`
	Api_query1       string `json:"api_query1"`
	Api_query2       string `json:"api_query2"`
	Api_query3       string `json:"api_query3"`
	Api_query4       string `json:"api_query4"`
	Api_query5       string `json:"api_query5"`
}

func (api *Naver_api_url) ValueCheck(i int, value string) {
	switch i {
	case 0:
		api.Api_seq, _ = strconv.Atoi(value)
	case 1:
		api.Api_name = value
	case 2:
		api.Api_type = value
	case 3:
		api.Api_url = value
	case 4:
		api.Api_sub_url = value
	case 5:
		api.Api_content_type = value
	case 6:
		api.Api_method = value
	case 7:
		api.Api_status = value
	case 8:
		api.Api_query1 = value
	case 9:
		api.Api_query2 = value
	case 10:
		api.Api_query3 = value
	case 11:
		api.Api_query4 = value
	case 12:
		api.Api_query5 = value
	}
}

var db *sql.DB
var err error

func dbConnect() string {
	id := "root"
	password := "password"
	host := "10.113.253.152:13306"
	dbname := "platform_admintool"
	dbcon := id + ":" + password + "@tcp(" + host + ")/" + dbname
	log.Println("DB Connect : ", dbcon)
	return dbcon
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		// c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func ApiGet(c *gin.Context) {
	var (
		api    Naver_api_url
		result gin.H
	)
	id := c.Param("api_seq")
	row, err := db.Query("select * from naver_api_url where api_seq = ?;", id)
	if err != nil {
		log.Println(err.Error())
	}
	coloumns, err := row.Columns()
	if err != nil {
		log.Println(err.Error())
	}
	values := make([]sql.RawBytes, len(coloumns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	err = row.Scan(scanArgs...)
	for row.Next() {
		err = row.Scan(scanArgs...)
		if err != nil {
			log.Println(err.Error())
		}
		var value string
		for i, col := range values {

			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			api.ValueCheck(i, value)
		}
		if err != nil {
			log.Println(err.Error())
		}
	}
	defer row.Close()
	if err != nil {
		log.Println(err.Error())
		result = gin.H{
			"message": "Data not found",
			"result":  nil,
			"count":   0,
		}
	} else {
		result = gin.H{
			"result": api,
			"count":  1,
		}
	}
	c.JSON(http.StatusOK, result)
}

func ApiGetall(c *gin.Context) {
	var (
		api  Naver_api_url
		apis []Naver_api_url
	)
	rows, err := db.Query("select * from naver_api_url;")
	if err != nil {
		log.Println(err.Error())
	}
	coloumns, err := rows.Columns()
	if err != nil {
		log.Println(err.Error())
	}
	values := make([]sql.RawBytes, len(coloumns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Println(err.Error())
		}
		var value string
		for i, col := range values {

			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			api.ValueCheck(i, value)
		}
		apis = append(apis, api)
		if err != nil {
			log.Println(err.Error())
		}
	}
	defer rows.Close()
	c.JSON(http.StatusOK, gin.H{
		"result": apis,
		"count":  len(apis),
	})
}

func ApiPost(c *gin.Context) {
	var api Naver_api_url
	c.BindJSON(&api)
	api_name := api.Api_name
	api_type := api.Api_type
	api_url := api.Api_url
	api_sub_url := api.Api_sub_url
	api_content_type := api.Api_content_type
	api_method := api.Api_method
	api_status := api.Api_status
	api_query1 := api.Api_query1
	api_query2 := api.Api_query2
	api_query3 := api.Api_query3
	api_query4 := api.Api_query4
	api_query5 := api.Api_query5

	log.Println("post value : ", &api)

	row := db.QueryRow("insert into naver_api_url (api_name, api_type, api_url, api_sub_url, api_content_type, "+
		"api_method, api_status, api_query1, api_query2,api_query3, api_query4, api_query5) values(?,?,?,?,?,?,?,?,?,?,?,?);",
		api_name, api_type, api_url, api_sub_url, api_content_type, api_method,
		api_status, api_query1, api_query2, api_query3, api_query4, api_query5)
	err = row.Scan(&api.Api_seq, &api.Api_name, &api.Api_type, &api.Api_url)
	if err != nil {
		log.Println(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Instert",
		"result":  api,
	})
}

func ApiPut(c *gin.Context) {
	var api Naver_api_url
	c.BindJSON(&api)
	id := c.Param("api_seq")

	api_name := api.Api_name
	api_type := api.Api_type
	api_url := api.Api_url
	api_sub_url := api.Api_sub_url
	api_content_type := api.Api_content_type
	api_method := api.Api_method
	api_status := api.Api_status
	api_query1 := api.Api_query1
	api_query2 := api.Api_query2
	api_query3 := api.Api_query3
	api_query4 := api.Api_query4
	api_query5 := api.Api_query5

	log.Println("put value : ", &api)
	stmt, err := db.Prepare("update naver_api_url set api_name=?, api_type=?, api_url=?, api_sub_url=?" +
		", api_content_type=?, api_method=?, api_status=?, api_query1=?, api_query2=?, api_query3=?, api_query4=?" +
		", api_query5=? where api_seq=?;")
	if err != nil {
		log.Println(err.Error())
	}
	_, err = stmt.Exec(api_name, api_type, api_url, api_sub_url, api_content_type, api_method,
		api_status, api_query1, api_query2, api_query3, api_query4, api_query5, id)
	if err != nil {
		log.Println(err.Error())
	}

	defer stmt.Close()
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Updated",
		"result":  api,
	})
}

func ApiDelete(c *gin.Context) {
	// id := c.Query("id")
	id := c.Param("api_seq")
	var api Naver_api_url
	row := db.QueryRow("select * from naver_api_url where api_seq = ?;", id)
	err = row.Scan(&api.Api_seq, &api.Api_name, &api.Api_type, &api.Api_url)
	log.Println("DELETE = ", id, api)
	stmt, err := db.Prepare("delete from naver_api_url where api_seq= ?;")
	if err != nil {
		log.Println(err.Error())
	}
	_, err = stmt.Exec(id)
	if err != nil {
		log.Println(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Deleted",
		"result":  api,
		// fmt.Sprintf("Successfully deleted %s", naver_api_url)
	})
}

func main() {
	db, err = sql.Open("mysql", dbConnect())
	if err != nil {
		log.Println(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err.Error())
	}

	router := gin.Default()
	router.Use(CORSMiddleware())

	router.GET("/api/:api_seq", ApiGet)
	router.GET("/apis", ApiGetall)
	router.POST("/api", ApiPost)
	router.PUT("/api/:api_seq", ApiPut)
	router.DELETE("/api/:api_seq", ApiDelete)

	router.Run(":7777")
}
