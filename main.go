package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dibyajyotid/mongoapi/router"
)

func main() {
	fmt.Println("mongoDB API")

	fmt.Println("Serve is getting started...")

	//getting th router
	r := router.Router()

	log.Fatal(http.ListenAndServe(":4000", r))
	fmt.Println("Listening at port 4000")
}
