package main
//import "github.com/gin-gonic/gin"
import (
	"fmt"
	"log"
	"net/http"

	"server/router"
)

func main() {
	/*
 	router := gin.Default();

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"Foo": "Bar" ,
		})
	})

	
	fmt.Println("Starting server on the port 8888...")
	router.Run("localhost:8888")
  */
	
	r := router.Router()

	fmt.Println("Starting server on the port 8888...")
	log.Fatal(http.ListenAndServe("localhost:8888", r)) // change later

	//log.Fatal(http.ListenAndServe(":8888", r))
	
}

