package main

import (
	"fmt"
	"github.com/Gretamass/kys-backend/db"
	"github.com/Gretamass/kys-backend/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
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

	loginRouter := r.Group("/login")
	{
		loginRouter.POST("/", srv.loginUser)
	}

	sneakerRouter := r.Group("/sneaker")
	{
		sneakerRouter.GET("/", srv.getSneakers)
		sneakerRouter.GET("/info", srv.getSneakersInfo)
		sneakerRouter.GET("/:id", srv.getSneakerInfo)
		sneakerRouter.GET("/availability", srv.getSneakersAvailability)
		sneakerRouter.GET("/:id/scrapper", srv.getSneakerScrapper)
		//TODO: add missing routers
		//sneakerRouter.GET("/:id", srv.getSneakerById)
		//sneakerRouter.POST("/", srv.createSneaker)
		//sneakerRouter.PATCH("/:id", srv.updateSneaker)
		//sneakerRouter.DELETE("/:id", srv.deleteSneaker)
	}

	providerRouter := r.Group("/provider")
	{
		providerRouter.GET("/", srv.getProviders)
		providerRouter.GET("/:id", srv.getProviderById)
		//TODO: add missing routers
		//providerRouter.POST("/", srv.createProvider)
		//providerRouter.PATCH("/:id", srv.updateProvider)
		//providerRouter.DELETE("/:id", srv.deleteProvider)
	}

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

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
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect ID"})
		return
	}

	if err := s.db.DeleteUser(id); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"message": "User Deleted!"})
}

func (s *server) loginUser(c *gin.Context) {
	var user user.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad JSON"})
		return
	}

	userExists, err := s.db.LoginUser(user)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if userExists != true {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect email or password"})
		return
	}

	// generate JWT with user ID as claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
	})

	// sign the token with a secret key
	signedToken, err := token.SignedString([]byte("mysecretkey"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": signedToken})
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

	if err := s.db.DeleteAdmin(id); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"message": "Admin Deleted!"})
}

// SNEAKER handlers
func (s *server) getSneakers(c *gin.Context) {
	sneakers, err := s.db.GetSneakers()

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if sneakers == nil {
		c.JSON(404, gin.H{"error": "No Sneakers Found"})
		return
	} else {
		c.JSON(200, gin.H{"data": sneakers})
	}
}

func (s *server) getSneakersInfo(c *gin.Context) {
	sneakers, err := s.db.GetSneakersInfo()

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if sneakers == nil {
		c.JSON(404, gin.H{"error": "No Sneakers Found"})
		return
	} else {
		c.JSON(200, gin.H{"data": sneakers})
	}
}

func (s *server) getSneakerInfo(c *gin.Context) {
	idStr := c.Params.ByName("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sneaker ID is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect ID"})
		return
	}

	sneakerInfo, err := s.db.GetSneakerInfo(id)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"data": sneakerInfo})
}

func (s *server) getSneakersAvailability(c *gin.Context) {
	sneakers, err := s.db.GetSneakersAvailability()

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if sneakers == nil {
		c.JSON(404, gin.H{"error": "No Sneakers Found"})
		return
	} else {
		c.JSON(200, gin.H{"data": sneakers})
	}
}

func (s *server) getSneakerScrapper(c *gin.Context) {
	idStr := c.Params.ByName("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sneaker ID is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect ID"})
		return
	}

	sneakerInfo, err := s.db.GetSneakerScrapper(id)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"data": sneakerInfo})
}

// PROVIDER handlers
func (s *server) getProviders(c *gin.Context) {
	sneakers, err := s.db.GetProviders()

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if sneakers == nil {
		c.JSON(404, gin.H{"error": "No Providers Found"})
		return
	} else {
		c.JSON(200, gin.H{"data": sneakers})
	}
}

func (s *server) getProviderById(c *gin.Context) {
	idStr := c.Params.ByName("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider ID is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect ID"})
		return
	}

	providerInfo, err := s.db.GetProviderById(id)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"data": providerInfo})
}
