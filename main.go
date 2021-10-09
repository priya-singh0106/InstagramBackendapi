package main

import (
	"context"
	"encoding/json"
	"fmt"
	"httprouter"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
}
type Post struct {
	ID               primitive.ObjectID  `json:"id"                  bson:"_id"`
	caption          string              `json:"caption,omitempty"   bson:"caption"`
	imageurl         string              `json:"imageurl,omitempty"  bson:"imageurl"`
	posted_timestamp primitive.Timestamp `json:"posted_timestamp,omitempty" bson:"posted_timestamp"`
}

var client *mongo.Client
func Pagination(request *http.Request, FindOptions *options.FindOptions) (int64, int64) {
	if request.URL.Query().Get("page") != "" && request.URL.Query().Get("limit") != "" {
		page, _ := strconv.ParseInt(request.URL.Query().Get("page"), 10, 32)
		limit, _ := strconv.ParseInt(request.URL.Query().Get("limit"), 10, 32)
		if page == 1 {
			FindOptions.SetSkip(0)
			FindOptions.SetLimit(limit)
			return page, limit
		}

		FindOptions.SetSkip((page - 1) * limit)
		FindOptions.SetLimit(limit)
		return page, limit

	}
	FindOptions.SetSkip(0)
	FindOptions.SetLimit(0)
	return 0, 0
}

func CreateUser(response http.ResponseWriter, request *http.Request) {

	response.Header().Add("content-type", "application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	collection := client.Database("users").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)

}
func GetUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user User
	collection := client.Database("users").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(user)
}
func CreatePost(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var post Post
	json.NewDecoder(request.Body).Decode(&post)
	collection := client.Database("users").Collection("postlist")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, post)
	json.NewEncoder(response).Encode(result)

}
func GetPost(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var post Post
	collection := client.Database("users").Collection("postlist")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Post{ID: id}).Decode(&post)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return

	}
	json.NewEncoder(response).Encode(post)

}
func GetAllPost(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var postlist []Post
	collection := client.Database("users").Collection("postlist")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{Post{ID: id}})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var post Post
		cursor.Decode(&post)
		postlist = append(postlist, post)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return

	}
	json.NewEncoder(response).Encode(postlist)
}

func main() {
	fmt.Println("Starting the applictation...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, "mongodb://localhost:2701")
	router := httprouter.New()
	router.POST("/users,CreateUser")
	router.GET("/users/{id},GetUser")
	router.POST("/posts,CreatePost")
	router.GET("/posts/{id},GetPost")
	router.GET("/posts/users/{id},GetAllPost")
	http.ListenAndServe(":12345", router)

}
