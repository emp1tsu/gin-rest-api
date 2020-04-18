package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model
	Name     string    `json:"name"`
	Age      int       `json:"age"`
	Birthday time.Time `json:"birthday"`
}

func NewUser() User {
	return User{}
}

func gormConnect() *gorm.DB {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=username dbname=gin_rest_api password=pass sslmode=disable")

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("db connected: ", &db)
	return db
}

func setRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Create
	r.POST("/user", func(c *gin.Context) {
		data := NewUser()
		now := time.Now()
		data.CreatedAt = now
		data.UpdatedAt = now

		if err := c.BindJSON(&data); err != nil {
			c.String(http.StatusBadRequest, "Request is failed: "+err.Error())
		}
		db.NewRecord(data)
		db.Create(&data)
		if db.NewRecord(data) == false {
			c.JSON(http.StatusOK, data)
		}
	})

	// Read
	// 全レコード
	r.GET("/users", func(c *gin.Context) {
		users := []User{}
		db.Find(&users)
		c.JSON(http.StatusOK, users)
	})

	// 1レコード
	r.GET("/user/:id", func(c *gin.Context) {
		user := NewUser()
		id := c.Param("id")

		db.Where("ID = ?", id).First(&user)
		c.JSON(http.StatusOK, user)
	})

	// Update
	r.PUT("/user/:id", func(c *gin.Context) {
		user := NewUser()
		id := c.Param("id")

		data := NewUser()
		if err := c.BindJSON(&data); err != nil {
			c.String(http.StatusBadRequest, "Request is failed: "+err.Error())
		}

		db.Where("ID = ?", id).First(&user).Updates(&data)
	})

	// Delete
	r.DELETE("/user/:id", func(c *gin.Context) {
		user := NewUser()
		id := c.Param("id")

		db.Where("ID = ?", id).Delete(&user)
	})

	return r
}

func main() {
	db := gormConnect()
	r := setRouter(db)

	defer db.Close()

	db.LogMode(true)
	db.Set("gorm:table_options", "ENGINE=InnoDB")
	db.AutoMigrate(&User{})

	r.Run(":8080")
}
