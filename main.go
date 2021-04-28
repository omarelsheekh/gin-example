package main

import(
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
    "gorm.io/driver/sqlite"
	"net/http"
	"strconv"
	_ "fmt"
) 

type User struct {
	ID   uint           `gorm:"primaryKey"`
	Name string			`gorm:"unique;not null"`
  }

func main() {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&User{})

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})
	r.POST("/user", func(c *gin.Context) {
		// Create
		name := c.GetHeader("name")
		var u User
		if name == ""{
			c.String(http.StatusBadRequest, "name must be included in headers")
		} else if db.Where("name = ?", name).Find(&u).RowsAffected > 0 {
			c.String(http.StatusBadRequest, "name %s already exists", name)
		} else {
			db.Create(&User{
				Name: name,
			})
			c.String(http.StatusOK, "added user %s", name)
		}
		
	})
	r.GET("/user", func(c *gin.Context) {
		var users []User
		db.Find(&users)
		c.JSON(http.StatusOK, gin.H{"users":users})
	})
	r.GET("/user/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		var user User
		if err != nil || db.First(&user, id).Error != nil{
			// handle error
			c.String(http.StatusNotFound, "not found")
		}else{
			c.JSON(http.StatusOK, gin.H{"user":user})
		}
	})
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

/* References
https://tour.golang.org/list
https://github.com/gin-gonic/gin
https://gorm.io/docs/index.html
*/