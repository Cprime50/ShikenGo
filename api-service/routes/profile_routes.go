package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/Cprime50/api-service/client"
	"github.com/Cprime50/api-service/middleware"
	"github.com/gin-gonic/gin"
	// import middleware
	// import client
)

var (
	timeout = time.Second
)

func RegisterProfileRoutes(r *gin.Engine, client *auth.Client) {

	routes := r.Group("/profile")
	routes.Use(middleware.Auth(client))
	{
		routes.POST("/create", CreateProfile)
		routes.PUT("/update", UpdateProfile)
		routes.GET("/:id", GetProfileByID)
		routes.DELETE("/delete/:id", DeleteProfile)
	}
	routes.Use(middleware.Auth(client), middleware.RoleAuth("admin"))
	{
		routes.GET("/profiles", GetProfiles)
	}

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Shiken-Go")
	})

}

func CreateProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	profile, err := client.CreateUpdateProfile(c, ctx, c.Request.Method)
	if err != nil {
		log.Println("Error creating profile:", err)
		c.JSON(http.StatusBadRequest, gin.H{"Failed to create profile ": err})
		return
	}
	c.JSON(http.StatusCreated, profile)
}

func UpdateProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	profile, err := client.CreateUpdateProfile(c, ctx, c.Request.Method)
	if err != nil {
		log.Println("Error updating profile:", err)
		c.JSON(http.StatusBadRequest, gin.H{"Failed to update profile ": err})
		return
	}
	c.JSON(http.StatusCreated, profile)
}

func GetProfileByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	uid, ok := getAuthUserID(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	id := c.Param("id")
	if !isAdmin(c) && uid != id {
		log.Println("Error uid and id don't match, user is unauthorized to access")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	profile, err := client.GetProfile(ctx, id)
	if err != nil {
		log.Println("Error fetching profile:", err)
		c.JSON(http.StatusBadRequest, gin.H{"Failed to fetch profile ": err})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func GetProfiles(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	profiles, err := client.GetAllProfiles(ctx)
	if err != nil {
		log.Println("Error fetching profiles:", err)
		c.JSON(http.StatusBadRequest, gin.H{"Failed to fetch profiles ": err})
		return
	}
	c.JSON(http.StatusOK, profiles)
}

func DeleteProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	uid, ok := getAuthUserID(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	id := c.Param("id")
	if !isAdmin(c) && uid != id {
		log.Println("Error uid and id don't match, user is unauthorized to access")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err := client.DeleteProfile(ctx, id)
	if err != nil {
		log.Println("Error deleting profiles:", err)
		c.JSON(http.StatusBadRequest, gin.H{"Failed to delete profile ": err})
		return
	}
	c.Status(http.StatusOK)
}

func getAuthUserID(ctx *gin.Context) (string, bool) {
	user, exists := ctx.Get("user")
	if !exists {
		return "", false
	}
	userID := user.(*middleware.User).UserID

	return userID, true
}

func isAdmin(ctx *gin.Context) bool {
	user, exists := ctx.Get("user")
	if !exists {
		return false
	}
	role := user.(*middleware.User).Role
	return role == "admin"
}
