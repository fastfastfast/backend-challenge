package main

import (
	"backend-challenge/service/authenticate"
	"backend-challenge/service/deletebyid"
	"backend-challenge/service/fetchuserbyid"
	"backend-challenge/service/insert"
	"backend-challenge/service/listalluser"
	"backend-challenge/service/updatebyid"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email" validate:"required"`
	Password  string             `json:"password" bson:"password"`
	CreatedAt string             `json:"createdAt" bson:"createdAt"`
}

var userCollection = "users"

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	mongoClient, err := connectMongo()
	if err != nil {
		e.Logger.Fatal(err)
	}
	testdb := mongoClient.Database("testdb")
	userDB := testdb.Collection(userCollection)
	err = createUniqueEmailIndex(userDB)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go startUserCounter(&wg, userDB)
	wg.Wait()

	handlerInsert := &insert.InsertHandler{
		Service: insert.NewInsertService(),
	}
	r := e.Group("/")
	r.Use(echojwt.JWT([]byte("secret")))
	r.POST("register", func(c echo.Context) error {
		return handlerInsert.HandlerInsert(c, userDB)
	})

	handlerAuthenticate := &authenticate.AuthenticateHandler{
		Service: authenticate.NewAuthenticateService(),
	}
	e.POST("/Authenticate", func(c echo.Context) error {
		return handlerAuthenticate.HandlerAuthenticate(c, userDB)
	})

	handlerUpdatebyId := &updatebyid.UpdateUserByIDHandler{
		Service: updatebyid.NewUpdateUserByIDService(),
	}
	e.POST("/update", func(c echo.Context) error {
		return handlerUpdatebyId.HandlerUpdateUserByID(c, userDB)
	})

	handlerListAllUser := &listalluser.ListAllUserHandler{
		Service: listalluser.NewListAllUserService(),
	}
	e.GET("/users", func(c echo.Context) error {
		return handlerListAllUser.HandlerListAllUser(c, userDB)
	})

	handlerFetchUserByID := &fetchuserbyid.UserHandler{
		Service: fetchuserbyid.NewUserService(),
	}
	e.GET("/user/:id", func(c echo.Context) error {
		return handlerFetchUserByID.HandlerFetchUserByID(c, userDB)
	})

	handlerDeleteByID := &deletebyid.DeleteByIDHandler{
		Service: deletebyid.NewUpdateUserByIDService(),
	}
	e.DELETE("/user/:id", func(c echo.Context) error {
		return handlerDeleteByID.HandlerDeleteByID(c, userDB)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func connectMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("✅ Connected to MongoDB!")
	return client, nil
}

func createUniqueEmailIndex(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}}, // 1 คือ ascending
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

func startUserCounter(wg *sync.WaitGroup, collection *mongo.Collection) {
	defer wg.Done()

	for {
		time.Sleep(10 * time.Second)

		count, err := countUsers(collection)
		if err != nil {
			fmt.Printf("❌ Error counting users: %v\n", err)
			continue
		}

		fmt.Printf("📢 Number of users: %d\n", count)
	}
}

func countUsers(collection *mongo.Collection) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}
