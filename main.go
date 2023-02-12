package main

import (
	"fmt"
	"github.com/Gretamass/kys-backend/db"
	"github.com/Gretamass/kys-backend/user"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type server struct {
	db *db.DB
}

func main() {
	dbc, err := db.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	srv := &server{
		db: dbc,
	}

	r := gin.Default()
	r.SetTrustedProxies([]string{"192.168.68.102"})

	router := r.Group("/user")
	{
		router.GET("/", srv.getUsers)
		router.POST("/", srv.createUser)
		router.PATCH("/:id", srv.updateUser)
		router.DELETE("/:id", srv.deleteUser)
	}

	r.Run()
}

func (s *server) getUsers(c *gin.Context) {
	users, err := s.db.GetUsers()

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if users == nil {
		c.JSON(404, gin.H{"error": "No Users Found"})
		return
	} else {
		c.JSON(200, gin.H{"data": users})
	}
}

func (s *server) createUser(c *gin.Context) {
	var newUser user.User

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad JSON"})
		return
	}

	err := s.db.AddUser(newUser)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}

	c.JSON(200, gin.H{"success": "User added to the database"})
}

func (s *server) updateUser(c *gin.Context) {
	var request user.User

	userId := c.Params.ByName("id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad JSON"})
		return
	}

	err := s.db.UpdateUser(userId, request)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}

	c.JSON(200, gin.H{"message": "User Updated!"})
}
func (s *server) deleteUser(c *gin.Context) {
	var request user.DeleteRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad JSON"})
		return
	}

	err := s.db.DeleteUser(request.Id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}

	c.JSON(200, gin.H{"message": "User Deleted!"})
}
