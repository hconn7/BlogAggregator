package main

import (
	"fmt"

	"github.com.hconn7/BlogAggregator/internal/config"
)

func main() {
	config := &config.Config{
		DbURL: "postgres://example",
		User:  ""}
	fmt.Println(config)
	fmt.Print(config)

	config.SetUser("Henry")

	fmt.Print(config)
}
