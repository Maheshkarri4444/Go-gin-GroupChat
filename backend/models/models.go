package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      *string            `json:"username" bson:"username"`
	Email     *string            `json:"email" bson:"email"`
	Password  *string            `json:"password" bson:"password"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt primitive.DateTime `json:"updated_at" bson:"updated_at"`
}

type Message struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Username  string             `json:"username" bson:"username"`
	Content   string             `json:"content" bson:"content"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt primitive.DateTime `json:"updated_at" bson:"updated_at"`
}

type GroupChat struct {
	ID           primitive.ObjectID   `bson:"_id"`
	Participants []primitive.ObjectID `json:"participants" bson:"participants"`
	Messages     []primitive.ObjectID `json:"messages,omitempty" bson:"messages,omitempty"`
	CreatedAt    primitive.DateTime   `json:"created_at" bson:"created_at"`
	UpdatedAt    primitive.DateTime   `json:"updated_at" bson:"updated_at"`
}
