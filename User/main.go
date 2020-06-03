package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

// GetUser gets the details of a user
func GetUser(c echo.Context) (err error) {
	u1 := c.Get("user").(*jwt.Token)
	claims := u1.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
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
	client.Disconnect(ctx)
	return c.JSON(http.StatusOK, result)
}

// Name returns the name of the microservice
func Name(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, "User-Microservice")
}

func main() {

	e := echo.New()
	r := e.Group("/user")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("/profile", GetUser)
	e.GET("/service/name", Name)
	e.Logger.Fatal(e.Start(":1325"))
}
