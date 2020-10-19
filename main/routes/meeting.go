package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"strings"
	"math/rand"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/yashjjw/apiMeetings/main/models"
)

type RouteHandler struct {
	Meeting *mongo.Collection
}

func NewRouteHandler(Meeting *mongo.Collection) *RouteHandler {
	return &RouteHandler{
		Meeting: Meeting,
	}
}

//******************generate random string for ID *****************

var seededRand *rand.Rand = rand.New(
  rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}

//*****************function to get meeting details of a participant or meetings within a time frame **************

func (router *RouteHandler) getMeeting(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	startArr, noneStart := query["start"]
	endArr, noneEnd := query["end"]
	if  (!noneStart || !noneEnd || len(startArr) <=0 || len(endArr) <=0){

		participantArr, noneParticipant := query["participant"]
		if  (!noneParticipant || len(participantArr) <=0) {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}
		participant := participantArr[0]

		participantUnwind := bson.D{{
			"$unwind", "$participants",
		}}

		matching := bson.D{{
			"$match", bson.M{
				"participants.email": participant,
			},
		}}

		meetingsEx, err := router.Meeting.Aggregate(context.TODO(), mongo.Pipeline{participantUnwind, matching})

		var meetings []bson.M
		if err = meetingsEx.All(context.TODO(), &meetings); err != nil {
			log.Fatal(err.Error())
			http.Error(w, "Server Error", http.StatusInternalServerError)
		}

		if len(meetings) != 0 {
			json.NewEncoder(w).Encode(meetings)
			return
		}

	}

	start, _ := strconv.Atoi(startArr[0])
	end, _ := strconv.Atoi(endArr[0])

	meetingsEx, err := router.Meeting.Find(context.TODO(), bson.D{
		{
		"startTime",
		bson.D{{
			"$gte", start,
		}}}, 
		{
		"endTime",
		bson.D{{
			"$lte", end,
		}},
	}})

	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	var meetings []bson.M

	if err = meetingsEx.All(context.TODO(), &meetings); err != nil {
		log.Fatal(err.Error())
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(meetings)
	return
}

//*****************function to schedule meetings******************************

func (router *RouteHandler) scheduleMeeting(w http.ResponseWriter, r *http.Request) {
	var meet models.Meeting

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&meet)
	meet.CreationTime = time.Now().Unix()
	date= meet.startTime+" MST"  //"02 Jan 06 15:04"
	parse_time, _ := time.Parse(time.RFC822, date)
	meet.startTime=parse_time.Unix()

	const charset = "abcdefghijklmnopqrstuvwxyz" + "0123456789"
	meet.ID=StringWithCharset(10, charset)

	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if meet.StartTime > meet.EndTime {
		http.Error(w, "startTime can not be lesser than endTime", http.StatusBadRequest)
	}

	emails := make(bson.A, len(meet.Participants))
	for i := 0; i < len(meet.Participants); i++ {
		emails[i] = meet.Participants[i].Email
	}

	participantUnwind := bson.D{{
		"$unwind", "$participants",
	}}
	//To check overlapping meetings **match if start or end time of new meeting lies between start and end of exising meeting
	timeOverlap := bson.D{{
		"$match", bson.M{
			"$or": bson.A{
				bson.M{
					"startTime": bson.M{"$lte": meet.StartTime},
					"endTime":   bson.M{"$gte": meet.StartTime},
				},
				bson.M{
					"startTime": bson.M{"$lte": meet.EndTime},
					"endTime":   bson.M{"$gte": meet.EndTime},
				},

			},
		},
	}}
	// To check participants with rspv yes in other meetings
	matchingParticipant := bson.D{{
		"$match", bson.M{
			"participants.rsvp": bson.M{
				"$in": bson.A{"Yes"},
			},
			"participants.email": bson.M{
				"$in": emails,
			},
		},
	}}

	meetingsEx, err := router.Meeting.Aggregate(context.TODO(), mongo.Pipeline{participantUnwind, timeOverlap, matchingParticipant})

	var meetings []bson.M
	if err = meetingsEx.All(context.TODO(), &meetings); err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}

	if len(meetings) != 0 {
		json.NewEncoder(w).Encode("The meeting timing is overlapping with other meetings of participants")
		return
	}

	_ , err = router.Meeting.InsertOne(context.TODO(), meet)

	if err != nil {
		log.Fatal(err)
		return
	}

	json.NewEncoder(w).Encode(meet)
}

//*****************function to get meeting details by ID *****************************

func (router *RouteHandler) getMeetingByID(w http.ResponseWriter, r *http.Request) {
	
	p := strings.Split(r.URL.Path, "/")
	id := p[len(p)-1]

	var meeting bson.M

	err := router.Meeting.FindOne(context.TODO(), bson.D{{
		"id", id,
	}}).Decode(&meeting)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(meeting)
}



func (router *RouteHandler) MeetingRoutes(w http.ResponseWriter, r *http.Request) {
	p:=r.URL.Path
	if(p=="/meetings"){
		switch r.Method{
		case "GET":
			router.getMeeting(w,r)
		case "POST":
			router.scheduleMeeting(w,r)

		default:
			http.Error(w, "Method Not Found", http.StatusMethodNotAllowed)
		}
	} else{
		switch r.Method {
		case "GET":
			router.getMeetingByID(w, r)
	
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}




