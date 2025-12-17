package main

import (
	"fmt"

	"github.com/TelitsynNikita/go_core_rest_api"
)

func main() {
	apiService := go_core_rest_api.NewApiService()

	apiService.InitApiService()

	if err := apiService.RunService(); err != nil {
		fmt.Println(err)
	}
}
