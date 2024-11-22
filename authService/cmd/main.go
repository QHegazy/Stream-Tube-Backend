package main

import (
	"authService/auth_grpc"
	"authService/auth_server"

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
		auth_server.Auth()
	}()
	wg.Wait()
}
