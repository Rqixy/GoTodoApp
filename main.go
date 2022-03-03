package main

import (
	"fmt"
	
	"sampleapp/app/controllers"
	"sampleapp/app/models"
)

func main() {
	fmt.Println(models.Db)

	controllers.StartMainServer()


}