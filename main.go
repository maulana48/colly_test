package main

import (
	"fmt"
	"main/controllers"
	"os"
)

func main() {
	if err := os.Setenv("TZ", "Asia/Jakarta"); err != nil {
		fmt.Println(err.Error())
	}

	controller := controllers.NewInstance()
	controller.GetArticleData()
}
