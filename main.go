package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type healthCheckResult struct {
	IsHealthy bool
}

func main() {
	initConfig()

	run()
}

func run() {
	showOptions()
	runHealthCheck(getOption())
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	path, _ := os.Getwd()
	viper.AddConfigPath(path)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}
}

func showOptions() {
	fmt.Println("Make your choice:")
	fmt.Println("1 - Production")
	fmt.Println("2 - Sandbox")
	fmt.Println("Anything else - You're done!")
}

func getOption() int {
	var choice int
	fmt.Println()
	fmt.Print("Choice: ")
	fmt.Scanf("%d", &choice)
	return choice
}

func runHealthCheck(choice int) {
	var urls []string

	if choice == 2 {
		urls = viper.GetStringSlice("sandboxUrls")
	} else if choice == 1 {
		urls = viper.GetStringSlice("productionUrls")
	} else {
		os.Exit(0)
	}

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)

		checkURL(url, &wg)
	}

	wg.Wait()

	fmt.Println()
	fmt.Println("---------------------")
	fmt.Println()
	run()
}

func checkURL(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := resty.New()

	r, err := client.R().SetResult(&healthCheckResult{}).Get(url)

	if err != nil {
		fmt.Println("Check failed on ", url)
	}

	if r.StatusCode() != 200 {
		fmt.Println("Failed with HTTP Status Code ", r.StatusCode(), "on", url)
	}

	result := r.Result().(*healthCheckResult)

	if result.IsHealthy {
		fmt.Println(url, "is healthy!")
	} else {

		fmt.Println(url, "is not healthy!")
	}
}
