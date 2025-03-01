package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/insaneadinesia/gobang/logger"
	"github.com/insaneadinesia/gobang/rest"
)

func main() {
	ctx := context.Background()

	// initiate logger
	logger.NewLogger(logger.Option{IsEnable: true})

	client := rest.New(rest.Option{
		Address: "https://test-api.free.beeceptor.com",
		Timeout: time.Duration(30 * time.Second),
	})

	headers := http.Header{}
	headers.Add("Content-Type", "application/json")

	resp, status, err := client.Get(ctx, "/get-something", headers)
	if err != nil {
		logger.Log.Error(ctx, "error get /get-something", err.Error())
		return
	}

	fmt.Println("Do something with this status: ", status)
	fmt.Println("Do something with this response: ", string(resp))
}
