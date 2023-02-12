package main

import (
	"fmt"
	"github.com/Gretamass/kys-backend/db"
	"github.com/Gretamass/kys-backend/user"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
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

	userRouter := r.Group("/user")
	{
		userRouter.GET("/", srv.getUsers)
		userRouter.GET("/:id", srv.getUserById)
		userRouter.POST("/", srv.createUser)
		userRouter.PATCH("/:id", srv.updateUser)
		userRouter.DELETE("/:id", srv.deleteUser)
	}

	adminRouter := r.Group("/admin")
	{
		adminRouter.GET("/", srv.getAdmins)
		adminRouter.GET("/:id", srv.getAdminById)
		adminRouter.POST("/", srv.createAdmin)
		adminRouter.PATCH("/:id", srv.updateAdmin)
		adminRouter.DELETE("/:id", srv.deleteAdmin)
	}

	r.Run()
}

// USER handlers
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

func (s *server) getUserById(c *gin.Context) {
	idStr := c.Params.ByName("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect ID"})
		return
	}

	user, err := s.db.GetUserById(id)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"data": user})
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
		return
	}

	c.JSON(200, gin.H{"success": "User added to the database"})
}

func (s *server) updateUser(c *gin.Context) {
	var request user.User

	idStr := c.Params.ByName("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect ID"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad JSON"})
		return
	}

	if err := s.db.UpdateUser(id, request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"message": "User Updated!"})
}

func (s *server) deleteUser(c *gin.Context) {
	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	if err := s.db.DeleteUser(id); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"message": "User Deleted!"})
}

// ADMIN handlers
func (s *server) getAdmins(c *gin.Context) {
	admins, err := s.db.GetAdmins()

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if admins == nil || len(admins) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No Admins Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": admins})
	}
}

func (s *server) getAdminById(c *gin.Context) {
	idStr := c.Params.ByName("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "admin ID is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect ID"})
		return
	}

	admin, err := s.db.GetAdminById(id)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"data": admin})
}

func (s *server) createAdmin(c *gin.Context) {
	var newAdmin user.Admin

	if err := c.BindJSON(&newAdmin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad JSON"})
		return
	}

	err := s.db.AddAdmin(newAdmin)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"success": "Admin added to the database"})
}

func (s *server) updateAdmin(c *gin.Context) {
	var request user.Admin

	idStr := c.Params.ByName("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "admin ID is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect ID"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad JSON"})
		return
	}

	if err := s.db.UpdateAdmin(id, request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"message": "Admin Updated!"})
}

func (s *server) deleteAdmin(c *gin.Context) {

}
