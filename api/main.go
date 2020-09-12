package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Article represents an article that the api will serve
type Article struct {
	ID      string `json:"ID"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

// Articles represents a list of articles that we save
var Articles []Article

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	// Loop over all of our Articles
	// find the one with the correct id
	// return that article

	for _, article := range Articles {
		if article.ID == key {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {

	// read the request body
	// unmarshall it into an Article object
	// save it to the running list of Articles

	reqBody, _ := ioutil.ReadAll(r.Body)

	var article Article

	json.Unmarshal(reqBody, &article)

	Articles = append(Articles, article)

	json.NewEncoder(w).Encode(article)

}

func deleteArticle(w http.ResponseWriter, r *http.Request) {

	// parse the id path param
	vars := mux.Vars(r)

	id := vars["id"]

	// find the cooresponding article and remove it
	for index, article := range Articles {
		if article.ID == id {
			Articles = append(Articles[:index], Articles[index+1:]...)
		}
	}

}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)

	var updatedArticle Article

	json.Unmarshal(reqBody, &updatedArticle)

	for index, article := range Articles {
		if article.ID == id {
			Articles[index] = updatedArticle
		}
	}
}

func handleRequest() {
	// how we map our endpoints to functions using net/http
	// http.HandleFunc("/", homePage)
	// http.HandleFunc("/articles", returnAllArticles)

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllArticles)
	// order matters here (i think)
	// because READ and DELETE have the same path but one is GET and one is a DELETE
	// need to put DELETE first so that it maps correctly
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", updateArticle).Methods("PUT")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	Articles = []Article{
		Article{ID: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		Article{ID: "2", Title: "Hello 2", Desc: "Description", Content: "Content"},
	}

	handleRequest()
}
