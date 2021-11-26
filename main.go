package main

import (
	"context"
	"fmt"
	"log"
	"mybot/binance"
	"mybot/bot"
	"mybot/configs"
	"mybot/repository/users"
	usersService "mybot/service/users"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	initTimeZone()
	config := configs.GetConfig()
	// mongo
	ctx := context.TODO()
	mgClient := connectMongo(config.Mongo.Username, config.Mongo.Password)
	db := mgClient.Database(config.Mongo.DbName)
	emaCollection := db.Collection(config.Mongo.EmaCollection)
	userCollection := db.Collection(config.Mongo.UserCollection)

	_ = ctx
	_ = emaCollection
	_ = userCollection

	usrRepo := users.NewUserRepository(userCollection, ctx)
	userService := usersService.NewUserService(usrRepo)
	u := users.User{
		Id:     "",
		Name:   "HadesGod",
		Active: []string{},
		Equity: 2500,
	}
	userService.Create(u)

	/** binance */
	binanceClient := binance.NewBinanceClient("")
	/* bot processes */
	go bot.EMACrossProcess(binanceClient)

	e := echo.New()
	e.Logger.Fatal(e.Start(fmt.Sprintf("localhost:%d", config.EchoApp.Port)))
}

func connectMongo(username string, password string) *mongo.Client {
	clientOptions := options.Client().
		ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@hdgcluster.xmgsx.mongodb.net/myFirstDatabase?retryWrites=true&w=majority", username, password))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func initTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	time.Local = ict
}
