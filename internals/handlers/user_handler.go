package handlers

import (
	"net/http"
	"strings"
	"time"
	"todo_with_auth/internals/config"
	"todo_with_auth/internals/models"
	"todo_with_auth/internals/repositories"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type SignUpUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SigninUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SigninUserResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user_details"`
}

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var registerUser SignUpUser

		err := ctx.BindJSON(&registerUser)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(registerUser.Password) < 6 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters long!"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user := &models.User{
			Email:    registerUser.Email,
			Password: string(hashedPassword),
		}

		createdUser, err := repositories.CreateUser(pool, *user)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, createdUser)
	}
}

func SigninUserHandler(pool *pgxpool.Pool, config *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginRequest SigninUserRequest

		err := ctx.BindJSON(&loginRequest)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repositories.GetUserByEmail(pool, loginRequest.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
			return
		}

		//Upon successful password comparison, create a new JWT Token with suitable expiration time
		claims := jwt.MapClaims{
			"user_id": user.Id,
			"email":   user.Email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(config.JWTSecret))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, SigninUserResponse{Token: tokenString, User: user})
	}
}
