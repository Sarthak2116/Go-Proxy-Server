package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type user struct {
	User  string `json:"user" form:"user" query:"user"`
	Email string `json:"email" form:"email" query:"email"`
	PhNo  string `json:"phone" form:"phone" query:"phone"`
	DOB   string `json:"dob" form:"dob" query:"dob"`
	Age   string `json:"age" form:"age" query:"age"`
}

var ctx context.Context

func makeConnectionUser() (m *mongo.Client, err error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://admin:admin@testcluster-2sjwn.mongodb.net/test?retryWrites=true&w=majority"))
	if err != nil {
		return nil, err
	}
	ctx, e := context.WithTimeout(context.Background(), 10*time.Second)
	_ = e
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// AuthUser checks whether a user exists or not
func AuthUser(c echo.Context) (err error) {
	name := c.Request().Header.Get("Username")
	client, err := makeConnectionUser()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, "Not Present!")
	}
	db := client.Database("CaringCompany")
	coll := db.Collection("users")
	filter := bson.M{
		"user": bson.M{
			"$eq": name,
		},
	}
	var result user
	_ = coll.FindOne(context.TODO(), filter).Decode(&result)
	a := result.User
	if a == "" {
		return c.JSON(http.StatusUnauthorized, "User not registered")
	}
	client.Disconnect(ctx)
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = result.User
	claims["admin"] = true

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func main() {

	e := echo.New()
	e.GET("/auth", AuthUser)
	e.Logger.Fatal(e.Start(":1324"))
}
