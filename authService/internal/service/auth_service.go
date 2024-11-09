package service

import (
	"authService/internal/models"
	"authService/internal/repository"
	"context"
	"fmt"
	"log"
)


 

func CreateUser(){
	userRepo := repository.NewUserRepository()
	ctx := context.Background()
	user := models.User{
		Username:    "testuser",
		AuthMethod:  "local",
		Email:       "test@example.com",
		Status:      "active",
	}
	
	result := userRepo.Insert(ctx, user)
	res := <-result
	if res.Error != nil {
		log.Fatal(res.Error)
	}
	fmt.Println(res.Data)


}