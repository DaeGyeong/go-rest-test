package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	Id         int
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
}

var db *sql.DB
var err error

func dbConnect() string {
	id := "root"
	password := "enghks1@"
	host := "10.113.253.152:13306"
	dbname := "test"
	dbcon := id + ":" + password + "@tcp(" + host + ")/" + dbname
	// gorm.Open(“mysql”, “user:pass@tcp(127.0.0.1:3306)/database?charset=utf8&parseTime=True&loc=Local”)
	log.Println("DB Connect : ", dbcon)
	return dbcon
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func ApiGet(c *gin.Context) {
	var (
		person Person
		result gin.H
	)
	id := c.Param("id")
	row := db.QueryRow("select * from person where id = ?;", id)
	err = row.Scan(&person.Id, &person.First_Name, &person.Last_Name)
	if err != nil {
		result = gin.H{
			"result": nil,
			"count":  0,
		}
	} else {
		result = gin.H{
			"result": person,
			"count":  1,
		}
	}
	c.JSON(http.StatusOK, result)
}

func ApiGetall(c *gin.Context) {
	var (
		person  Person
		persons []Person
	)
	rows, err := db.Query("select * from person;")
	if err != nil {
		log.Println(err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&person.Id, &person.First_Name, &person.Last_Name)
		persons = append(persons, person)
		if err != nil {
			log.Println(err.Error())
		}
	}
	defer rows.Close()
	c.JSON(http.StatusOK, gin.H{
		"result": persons,
		"count":  len(persons),
	})
}

func ApiPost(c *gin.Context) {
	var person Person

	c.BindJSON(&person)
	first_name := person.First_Name
	last_name := person.Last_Name
	log.Println("post value : ", first_name, last_name)

	row := db.QueryRow("insert into person (first_name, last_name) values(?,?);", first_name, last_name)
	err = row.Scan(&person.Id, &person.First_Name, &person.Last_Name)
	if err != nil {
		log.Println(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprintf(" %s, %s successfully created", first_name, last_name),
	})
}

func ApiPut(c *gin.Context) {
	var buffer bytes.Buffer
	id := c.Param("id")
	var person Person
	c.BindJSON(&person)
	first_name := person.First_Name
	last_name := person.Last_Name
	log.Println("put value : ", first_name, last_name)
	stmt, err := db.Prepare("update person set first_name=?, last_name=? where id=?;")
	if err != nil {
		log.Println(err.Error())
	}
	_, err = stmt.Exec(first_name, last_name, id)
	if err != nil {
		log.Println(err.Error())
	}

	buffer.WriteString(first_name)
	buffer.WriteString(" ")
	buffer.WriteString(last_name)
	defer stmt.Close()
	name := buffer.String()
	c.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprintf("Successfully updated to %s", name),
	})
}

func ApiDelete(c *gin.Context) {
	// id := c.Query("id")
	id := c.Param("id")
	var person Person
	row := db.QueryRow("select * from person where id = ?;", id)
	err = row.Scan(&person.Id, &person.First_Name, &person.Last_Name)
	log.Println("DELETE = ", id, person)
	stmt, err := db.Prepare("delete from person where id= ?;")
	if err != nil {
		log.Println(err.Error())
	}
	_, err = stmt.Exec(id)
	if err != nil {
		log.Println(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted",
		"result":  person,
		// fmt.Sprintf("Successfully deleted %s", person)
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

	router.GET("/person/:id", ApiGet)
	router.GET("/persons", ApiGetall)
	router.POST("/person", ApiPost)
	router.PUT("/person/:id", ApiPut)
	router.DELETE("/person/:id", ApiDelete)

	router.Run(":3000")
}
