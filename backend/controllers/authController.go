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
	fmt.Println("Stored Hashed Password:", userPassword)
	fmt.Println("Provided Password:", providedPassword)

	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	if err != nil {
		return false, "password is incorrect"
	}
	return true, ""
}

func SignUp(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	defer cancel()

	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
	}

	if emailCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user with this email already exists"})
		return
	}

	password := HashPassword(*user.Password)
	user.Password = &password
	user.ID = primitive.NewObjectID()
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		msg := fmt.Sprintf("user item was not %s", "created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
	}

	userID := user.ID.Hex()
	username := *user.Name

	token, err, expirationTime := auth.GenerateJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while generating token"})
		return
	}
	//c.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var foundUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "email is not valid"})
		return
	}

	passwordIsValid, msg := VerifyPassword(*foundUser.Password, *user.Password)
	if !passwordIsValid {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	if foundUser.Email == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	userID := foundUser.ID.Hex()
	username := *foundUser.Name

	shouldRefresh, err, expirationTime := auth.RefreshToken(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "refresh token error", "err": err.Error()})
		return
	}
	token := ""
	if shouldRefresh {
		token, err, expirationTime = auth.GenerateJWT(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while generating token"})
			return
		}
	}

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
