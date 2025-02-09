package controller

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"githun.com/Maheshkarri4444/group-chat/database"
	"githun.com/Maheshkarri4444/group-chat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var messageCollection *mongo.Collection = database.OpenCollection(database.Client, "messages")
var groupChatCollection *mongo.Collection = database.OpenCollection(database.Client, "groupchat")

func SendMessage(c *gin.Context) {
	userid, _ := c.Get("userid")
	var message models.Message
	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message.ID = primitive.NewObjectID()
	message.UserID, _ = primitive.ObjectIDFromHex(userid.(string))
	message.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	message.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	var groupChat models.GroupChat
	err := groupChatCollection.FindOne(context.TODO(), bson.M{}).Decode(&groupChat)
	if err == mongo.ErrNoDocuments {
		groupChat = models.GroupChat{
			ID:           primitive.NewObjectID(),
			Participants: []primitive.ObjectID{message.UserID},
			Messages:     []primitive.ObjectID{message.ID},
			CreatedAt:    primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt:    primitive.NewDateTimeFromTime(time.Now()),
		}
		groupChatCollection.InsertOne(context.TODO(), groupChat)
	} else {
		update := bson.M{
			"$addToSet": bson.M{"messages": message.ID, "participants": message.UserID},
			"$set":      bson.M{"updated_at": primitive.NewDateTimeFromTime(time.Now())},
		}
		groupChatCollection.UpdateOne(context.TODO(), bson.M{"_id": groupChat.ID}, update)
	}
	messageCollection.InsertOne(context.TODO(), message)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": message})

}

func GetMessages(c *gin.Context) {
	var groupChat models.GroupChat
	err := groupChatCollection.FindOne(context.TODO(), bson.M{}).Decode(&groupChat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "group chat not found"})
		return
	}

	if len(groupChat.Messages) == 0 {
		c.JSON(http.StatusOK, gin.H{"messages": []models.Message{}})
		return
	}

	var messages []models.Message
	cursor, err := messageCollection.Find(
		context.TODO(),
		bson.M{"_id": bson.M{"$in": groupChat.Messages}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retreving messages"})
		return
	}

	if err = cursor.All(context.TODO(), &messages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error decoding messages"})
		return
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.Time().Before(messages[j].CreatedAt.Time())
	})

	c.JSON(http.StatusOK, messages)

}
