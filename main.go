package main
import (
    "encoding/json"
    "log"
    "github.com/gorilla/mux"
    "net/http"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "fmt"
)

type Person struct{
    ID string `json:"id,bson:"_id,omitempty""`
    Firstname string `json:"firstname,omitempty"`
    Lastname string `json:"lastname,omitempty"` 
    Address *Address `json:"address,omitempty"` 
}

type Address struct{
    City string `json:"city,omitempty"`
    State string `json:"state,omitempty"`
}
var db *mgo.Database
var people []Person


func collection() *mgo.Collection {
	return db.C("persons")
}

func GetPersonEndPoint(w http.ResponseWriter, req *http.Request){
    res := Person{}
    params := mux.Vars(req)
    id := params["id"]
    collection().Find(bson.M{"id": id}).One(&res);
    fmt.Println("ff", res)    
    json.NewEncoder(w).Encode(res)

}
func GetPeopleEndPoint(w http.ResponseWriter, req *http.Request){
    res := []Person{}
    collection().Find(nil).All(&res)
    json.NewEncoder(w).Encode(res)
}
func CreatePersonEndPoint(w http.ResponseWriter, req *http.Request){
    params := mux.Vars(req)
    id := params["id"]
    firstname := req.FormValue("firstname")
    lastname := req.FormValue("lastname")
    city := req.FormValue("city")
    state := req.FormValue("state")
    item := Person{ID: id, Firstname: firstname, Lastname: lastname, Address: &Address{City: city, State:state}}
    fmt.Println("ff", item)    
    collection().Insert(item)
    w.Write([]byte("OK"))
}
func DeletePersonEndPoint(w http.ResponseWriter, req *http.Request){
    params := mux.Vars(req)
    id := params["id"]
    collection().Remove(bson.M{"id": id})
    w.Write([]byte("DELETED"))
    
}
func main() {
    session, err := mgo.Dial("mongodb://HimaniTH1:HimaniTH1@ds161224.mlab.com:61224/vidly")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db = session.DB("vidly")
    router := mux.NewRouter()
    router.HandleFunc("/people", GetPeopleEndPoint).Methods("GET")
    router.HandleFunc("/people/{id}", GetPersonEndPoint).Methods("GET")
    router.HandleFunc("/people/{id}", CreatePersonEndPoint).Methods("POST")
    router.HandleFunc("/people/{id}", DeletePersonEndPoint).Methods("DELETE")
    log.Fatal(http.ListenAndServe(":12345", router))
}