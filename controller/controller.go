package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dibyajyotid/mongoapi/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb://localhost:27017/"
const dbName = "netflix"
const colName = "watchlist"

// MOST IMPORTANT
var collection *mongo.Collection

// connect with mongoDB
func init() {
	//client options
	clientOption := options.Client().ApplyURI(connectionString)

	//connect to mongoDB
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("mongoDB connection successful")

	collection = client.Database(dbName).Collection(colName)

	//if collection instance is ready
	fmt.Println("collection instance is ready")
}

//MONGO helpers - file

// insert one record
func insertOneMovie(movie model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)
	checkNilErr(err)

	fmt.Println(inserted)

	fmt.Println("Inserted one movie with db id: ", inserted.InsertedID)
}

// update One record - file
func updateOneMovie(movieID string) {
	id, err := primitive.ObjectIDFromHex(movieID) //converts string into the objectID which mongoDB understand
	checkNilErr(err)

	filter := bson.M{"_id": id}
	update := bson.M{"set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	checkNilErr(err)

	fmt.Println("modified count: ", result.ModifiedCount)
}

// delete one record - file
func deleteOneMovie(movieID string) {
	id, err := primitive.ObjectIDFromHex(movieID)
	checkNilErr(err)

	filter := bson.M{"_id": id} // we can also use this directly instead of filter variable
	deleteCount, err := collection.DeleteOne(context.Background(), filter)
	checkNilErr(err)

	fmt.Println("Movie got deleted with deletecount: ", deleteCount)
}

// delete all records from mongo db - file
func deleteAllMovie() {
	deleteResult, err := collection.DeleteMany(context.Background(), bson.D{}, nil)
	checkNilErr(err)

	fmt.Println("Number of Movies Deleted: ", deleteResult.DeletedCount)
}

// get all movies from database - file
func getAllMovie() []primitive.M {
	cursor, err := collection.Find(context.Background(), bson.D{})
	checkNilErr(err)
	defer cursor.Close(context.Background())

	var movies []primitive.M

	//loop through inside it
	for cursor.Next(context.Background()) {
		var movie bson.M
		err := cursor.Decode(&movie)
		checkNilErr(err)

		movies = append(movies, movie)
	}
	return movies
}

func checkNilErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Actual controller - file
func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	allMovies := getAllMovie()

	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)

	json.NewEncoder(w).Encode(movie)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	updateOneMovie(params["id"])

	json.NewEncoder(w).Encode(params["id"])
}

func DeleteOneMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)

	deleteOneMovie(params["id"])

	json.NewEncoder(w).Encode(params["id"])
}

func DelteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	deleteAllMovie()
	json.NewEncoder(w).Encode("Succes!")
}
