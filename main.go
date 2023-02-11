package main

import (
	"fmt"
	"github.com/Gretamass/kys-backend/models"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	err := models.ConnectDatabase()
	checkErr(err)

	fmt.Println(models.ConnectDatabase())

	r := gin.Default()
	r.SetTrustedProxies([]string{"192.168.68.102"})

	router := r.Group("/users")
	{
		router.POST("/create", createUser)
		router.GET("/", getUsers)
		router.POST("/update/:id", updateUser)
		router.GET("/delete/:id", deleteUser)
	}

	r.Run()
}

func createUser(c *gin.Context) {
	c.JSON(200, gin.H{"message": "A new User Created!"})
}

func getUsers(c *gin.Context) {
	users, err := models.GetUsers()
	checkErr(err)

	if users == nil {
		c.JSON(404, gin.H{"error": "No Users Found"})
		return
	} else {
		c.JSON(200, gin.H{"data": users})
	}
}

func updateUser(c *gin.Context) {
	c.JSON(200, gin.H{"message": "User Updated!"})
}
func deleteUser(c *gin.Context) {
	c.JSON(200, gin.H{"message": "User Deleted!"})
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
