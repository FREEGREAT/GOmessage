package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gomessage.com/users/internal/models"
	"gomessage.com/users/internal/service"
	repository "gomessage.com/users/internal/storage/postgresql"
	"gomessage.com/users/pkg/postgresql"
	"gomessage.com/users/pkg/utils"
)

const connectAttempts = 3

/*
func main() {
	srv := new(server.Server)
	if err := utils.InitConfig(); err != nil {
		log.Fatalf("error init config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading enb var: %s", err.Error())
	}

	postgresqlClient, err := postgresql.NewClient(context.TODO(), connectAttempts, postgresql.StorageConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		log.Fatalf("Failed to init db: %s", err.Error())
	}
	repository := repository.NewRepository(postgresqlClient)
	all, err := repository.FindAll(context.TODO())
	if err != nil {
		return
	}

	for _, usr := range all {
		fmt.Printf("%v", usr)
	}
	if err := srv.Run("8081"); err != nil {
		log.Fatalf("error occure while running http server: %s", err.Error())
	}
}*/

func main() {
	if err := utils.InitConfig(); err != nil {
		log.Fatalf("error init config: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading enb var: %s", err.Error())
	}
	fmt.Println("11")
	postgresqlClient, err := postgresql.NewClient(context.TODO(), connectAttempts, postgresql.StorageConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("Failed to initialize database client: %s", err)
	}
	fmt.Println("22")
	userRepo := repository.NewRepository(postgresqlClient)
	userService := service.CreateNewUserService(userRepo)

	// Використання сервісу для створення користувача
	age := 25
	img := "http://example.com/imagejpg"
	newUser := &models.UserModel{
		Nickname:     "JohnDoe",
		PasswordHash: "hashed_password",
		Email:        "johnddo333e@example.com",
		Age:          &age,
		ImageUrl:     &img,
	}

	if err := userService.CreateUser(context.TODO(), newUser); err != nil {
		log.Fatalf("Failed to create user: %s", err)
	}

	all, err := userService.GetAllUsers(context.TODO())
	if err != nil {
		return
	}
	for _, usr := range all {
		fmt.Printf("%v", usr)
	}

	log.Println("User created successfully")
}
