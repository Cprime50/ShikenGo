package routes

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"firebase.google.com/go/v4/auth"
	"github.com/Cprime50/api-service/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(r *gin.Engine, client *auth.Client) {

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.Auth(client), middleware.RoleAuth("admin"))
	{
		adminRoutes.POST("/make", func(ctx *gin.Context) {
			makeAdmin(ctx, client)
		})
		adminRoutes.DELETE("/remove", func(ctx *gin.Context) {
			removeAdmin(ctx, client)
		})
	}

} 

type EmailInput struct {
	Email string `json:"email"`
}

func makeAdmin(ctx *gin.Context, client *auth.Client) {
	var input EmailInput

	if err := ctx.BindJSON(&input); err != nil {
		log.Print("error making admin: invalid json format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	emailOk := ValidateEmailInput(input.Email)
	if !emailOk {
		log.Print("error making admin: Invalid email format")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Email format"})
		return
	}

	err := middleware.MakeAdmin(ctx.Request.Context(), client, input.Email)
	if err != nil {
		log.Print("error making admin:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s is now an admin", input.Email)})
}

func removeAdmin(ctx *gin.Context, client *auth.Client) {
	var input EmailInput

	if err := ctx.BindJSON(&input); err != nil {
		log.Print("error removing admin: invalid json format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	emailOk := ValidateEmailInput(input.Email)
	if !emailOk {
		log.Print("error removing admin: Invalid email format")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Email format"})
		return
	}

	err := middleware.RemoveAdmin(ctx.Request.Context(), client, input.Email)
	if err != nil {
		log.Print("error removing admin:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s admin rights have been revoked", input.Email)})
}

// Regex for email validation
func ValidateEmailInput(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
