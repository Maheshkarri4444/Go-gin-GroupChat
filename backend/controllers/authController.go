package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"githun.com/Maheshkarri4444/group-chat/auth"
	"githun.com/Maheshkarri4444/group-chat/database"
	"githun.com/Maheshkarri4444/group-chat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	if err != nil {
		return false, "password is incorrect"
	}
	return true, ""
}

func SignUp(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	fmt.Println("user at signup: ", user.Email)
	fmt.Println("user at signup: ", user.Name)

	// if user.Password == nil || user.Name == nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Name and password are required"})
	// 	return
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while checking email"})
		return
	}

	if emailCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		return
	}

	hashedPassword := HashPassword(*user.Password)
	user.Password = &hashedPassword
	user.ID = primitive.NewObjectID()
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User could not be created"})
		return
	}

	userID := user.ID.Hex()
	username := ""
	if user.Name != nil {
		username = *user.Name
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username is missing"})
		return
	}

	token, err, expirationTime := auth.GenerateJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.SetCookie("userid", userID, int(expirationTime.Sub(time.Now()).Seconds()), "/", "", false, true)
	c.SetCookie("username", username, int(expirationTime.Sub(time.Now()).Seconds()), "/", "", false, true)
	c.SetCookie("token", token, int(expirationTime.Sub(time.Now()).Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message":  "User signed up successfully",
		"userid":   userID,
		"username": username,
		"token":    token,
	})
}

func Login(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if user.Email == nil || user.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var foundUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": *user.Email}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Ensure that the password field is not nil before dereferencing
	if foundUser.Password == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password not found in database"})
		return
	}

	passwordIsValid, msg := VerifyPassword(*foundUser.Password, *user.Password)
	if !passwordIsValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
		return
	}

	// Ensure that the username is not nil before dereferencing
	if foundUser.Name == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username not found in database"})
		return
	}

	userID := foundUser.ID.Hex()
	username := *foundUser.Name

	shouldRefresh, err, expirationTime := auth.RefreshToken(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Refresh token error"})
		return
	}

	token := ""
	if shouldRefresh {
		token, err, expirationTime = auth.GenerateJWT(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}
	}

	// Set cookies
	c.SetCookie("userid", userID, int(expirationTime.Sub(time.Now()).Seconds()), "/", "", false, true)
	c.SetCookie("username", username, int(expirationTime.Sub(time.Now()).Seconds()), "/", "", false, true)
	if shouldRefresh {
		c.SetCookie("token", token, int(expirationTime.Sub(time.Now()).Seconds()), "/", "", false, true)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "User logged in successfully",
		"userid":   userID,
		"username": username,
		"token":    token,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("userid", "", -1, "/", "", false, true)
	c.SetCookie("username", "", -1, "/", "", false, true)
	c.SetCookie("token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})

}

func CheckAuth(c *gin.Context) {
	userID, err := c.Cookie("userid")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	username, err := c.Cookie("username")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "User is authenticated",
		"userid":   userID,
		"username": username,
	})
}
