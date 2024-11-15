package main

import (
	"authService/auth_grpc"
	"authService/oauth_server"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done() 
		auth_grpc.AuthGrpcServer()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done() 
		oauth_server.OAuth();
	}()
	wg.Wait()
}
