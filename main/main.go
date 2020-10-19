package main

import (
	"log"
	"net/http"
	"github.com/yashjjw/apiMeetings/main/models"
	"github.com/yashjjw/apiMeetings/main/routes"
)

func main() {

	client := models.ConnectDatabase("mongodb://localhost:27017")
	router := routes.NewRouteHandler(client.Database("Test").Collection("meetings"))

	http.HandleFunc("/", router.MeetingRoutes)

	log.Println("Listening on PORT 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
