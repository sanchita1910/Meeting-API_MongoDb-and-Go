package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client // making a client

type Participants1 struct { //Particpant attributes
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Rsvp  string `json:"rsvp,omitempty" bson:"rsvp,omitempty"`
}

//type starttime1 struct {
//  then := time.Date(
//       2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
//   p(then)
//}

type Meeting struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Start             time.Time          `json:"start,omitempty" bson:"start,omitempty"`
	End               time.Time          `json:"end,omitempty" bson:"end,omitempty"`
	Participants      []Participants1    `json:"participants,omitempty" bson:"participants,omitempty"`
	Creationtimestamp time.Time          `json:"creationtimestamp,omitempty" bson:"creationtimestamp,omitempty"`
}

type Meetings2 []Meeting

type urlhandler struct { //making a handler which uses mongodb client
	sync.Mutex
	db *mongo.Client
}

func (ph *urlhandler) ServeHTTP(response http.ResponseWriter, request *http.Request) { //servehttp for routing
	parts := strings.Split(request.URL.Path, "/")
	if len(parts) < 2 {
		// return 404
	}
	if request.Method == "GET" { //seeing if it's get or post and nesting meeting url condition under it
		if parts[1] == "meeting" && len(parts) == 3 { //get single meeting
			ph.getmeeting(response, request)
		} else if parts[1] == "meetings" {
			ph.getlistofmeetings(response, request)
		} else {
			// return 404
		}
	} else if request.Method == "POST" && parts[1] == "meetings" { //post a meeting
		ph.makemeeting(response, request)
	}
}

func (ph *urlhandler) makemeeting(response http.ResponseWriter, request *http.Request) { // A Function For the schedule a meeting endpoint
	response.Header().Set("content-type", "application/json")
	var person Meeting
	_ = json.NewDecoder(request.Body).Decode(&person)
	collection := ph.db.Database("appointydb").Collection("meetingsdb")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

func (ph *urlhandler) getmeeting(response http.ResponseWriter, request *http.Request) { // A Function For get a meeting using ID Endpoint
	response.Header().Set("content-type", "application/json")
	parts := strings.Split(request.URL.Path, "/")
	if len(parts) != 3 {
		// handle error
	}
	id, err := primitive.ObjectIDFromHex(parts[len(parts)-1])
	if err != nil {
		// handle error
	}
	var person Meeting
	collection := ph.db.Database("appointydb").Collection("meetingsdb")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err1 := collection.FindOne(ctx, Meeting{ID: id}).Decode(&person)
	if err1 != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(person)
}

func (ph *urlhandler) getlistofmeetings(response http.ResponseWriter, request *http.Request) { // A function to return all meeting which took place within the given time
	response.Header().Set("content-type", "application/json")
	var meetingsdb []Meeting
	collection := ph.db.Database("appointydb").Collection("meetingsdb")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	start1, ok := request.URL.Query()["start"]
	if !ok || len(start1) == 0 {
		// handle missing starttime
		return
	}
	startTimeStamp, err := time.Parse(time.RFC3339, start1[0])
	if err != nil {
		// invalid timestamp
		return
	}
	end1, ok := request.URL.Query()["end"]
	if !ok || len(end1) == 0 {
		// handle missing endtime
		return
	}
	endTimeStamp, err := time.Parse(time.RFC3339, end1[0])
	if err != nil {
		// invalid timestamp
		return
	}
	filter := bson.M{
		"starttime": bson.M{"$gte": startTimeStamp},
		"endtime":   bson.M{"$lte": endTimeStamp},
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Meeting
		cursor.Decode(&person)
		meetingsdb = append(meetingsdb, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meetingsdb)
}

func newurlhandler(client *mongo.Client) *urlhandler { //handler uses mongoDB client
	return &urlhandler{
		db: client,
	}
}

func (ph *urlhandler) getparticipantmeetingsemailid(response http.ResponseWriter, request *http.Request) { // A function to get all meeting of the participant using his/her email id
	response.Header().Set("content-type", "application/json") // This Function is currently not fully functional
	var meetingsdb []Meeting
	collection := ph.db.Database("appointydb").Collection("meetingsdb")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	emailfinder, ok := request.URL.Query()["email"]
	//emailfinderfinal, err :=

	//filter :=bson.M{
	///		"email":bson.M
	//}

	//cursor,err :=collection.Find(ctx, )
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Meeting
		cursor.Decode(&person)
		meetingsdb = append(meetingsdb, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meetingsdb)
}
func main() {
	port := ":8083"
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	ph := newurlhandler(client)
	http.Handle("/meetings", ph)
	http.Handle("/meeting", ph)
	http.Handle("/meeting/", ph)

	http.ListenAndServe(port, nil)
}
