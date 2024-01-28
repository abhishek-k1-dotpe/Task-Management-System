package main

import (
	"fmt"
	"os"
	"user-service/db"
	"user-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	routes.InitRoutes(router)
	err := db.InitMySQL()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	err = router.Run(":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
