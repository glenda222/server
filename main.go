package main
import "github.com/gin-gonic/gin"
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"server/router"
)

func main() {
	r := router.Router()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"Foo": "Bar" ,
		})
	})
	
	/*
	r := router.Router()

	//Azure app service sets the port in env, can be random
	port := os.Getenv("HTTP_PLATFORM_PORT")

	if port == "" {
		port = "8888"
	}
	fmt.Println("Starting server on the port " + port)
	log.Fatal(http.ListenAndServe("localhost:" + port, r)) // change later

	//log.Fatal(http.ListenAndServe(":8888", r))
 */
	
}

