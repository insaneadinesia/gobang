package main

import (
	"context"
	"encoding/json"
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

	reqBody := map[string]any{
		"test_1": "TEST 1",
		"test_2": "TEST 2",
	}

	by, _ := json.Marshal(reqBody)

	resp, status, err := client.Post(ctx, "/post-something", headers, by)
	if err != nil {
		logger.Log.Error(ctx, "error post /post-something", err.Error())
		return
	}

	fmt.Println("Do something with this status: ", status)
	fmt.Println("Do something with this response: ", string(resp))
}
