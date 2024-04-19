package controllers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt-project/database"
	"github.com/golang-jwt-project/helpers"
	"github.com/golang-jwt-project/models"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func HashPassword(password string) string {
	byte, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(byte)
}

func Singup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(user)
		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		var count int64
		err := database.DB.Where("email = ?", user.Email).Count(&count)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking for the emai"})
		}

		password := HashPassword(user.Password)
		user.Password = password

		err = database.DB.Where("phone = ?", user.Phone).Count(&count)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking for the phone number"})
		}

		if count > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "email or phone number is already exists"})
		}

		user.Created_at = time.Now()
		user.Updated_at = time.Now()
		user.ID = rand.Int()
		user.User_id = strconv.Itoa(user.ID)
		token, refreshToken, _ := helpers.GenerateAllTokens(user.Email, user.First_Name, user.Last_Name, user.User_type, user.User_id)
		user.Token = token
		user.Refresh_token = refreshToken

		if err := database.DB.Create(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User items was not created"})
		}
		ctx.JSON(http.StatusOK, gin.H{"success": "Data inserted successfully"})
	}

}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User
		var foundUser models.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		if err := database.DB.Where("email = ?", user.Email).Find(&foundUser).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
		}

		passIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
		if !passIsValid {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		if foundUser.Email == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}
		token, refreshToken, _ := helpers.GenerateAllTokens(foundUser.Email, foundUser.First_Name, foundUser.Last_Name, foundUser.User_type, foundUser.User_id)
		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		if err := database.DB.Where("user_id = ?", foundUser.User_id).Find(&foundUser).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("user_id")
		if err := helpers.MatchUserTypeToUid(ctx, userId); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		database.DB.Where("user_id = ?", userId).Find(&user)
		ctx.JSON(http.StatusOK, user)
	}
}
