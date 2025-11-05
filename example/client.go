package main

import (
	"context"
	"log"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	v1 "github.com/moweilong/art-design-pro-go/pkg/api/apiserver/v1"
)

func main() {
	callHTTP()

}

func callHTTP() {
	conn, err := transhttp.NewClient(
		context.Background(),
		transhttp.WithMiddleware(
			recovery.Recovery(),
		),
		transhttp.WithEndpoint("localhost:5555"),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := v1.NewUserCenterHTTPClient(conn)

	// 创建用户
	createUserReply, err := client.CreateUser(context.Background(), &v1.CreateUserRequest{
		Username: "art",
		Nickname: "art",
		Password: "art(#)888",
		Email:    "art@example.com",
		Phone:    "13800000000",
	})
	if err != nil {
		if !strings.Contains(err.Error(), "Duplicate entry") {
			log.Fatal(err)
		}
		log.Println("[http] Duplicate create user")
	} else {
		log.Printf("[http] CreateUser %s\n", createUserReply.UserID)
	}

	reply, err := client.Login(context.Background(), &v1.LoginRequest{Username: "art", Password: "art(#)888"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] Login %s\n", reply.RefreshToken)

	// returns error
	_, err = client.Login(context.Background(), &v1.LoginRequest{Username: "art", Password: "badpassword"})
	if err != nil {
		log.Printf("[http] Login error: %v\n", err)
	}
	if errors.IsBadRequest(err) {
		log.Printf("[http] Login error is invalid argument: %v\n", err)
	}
}
