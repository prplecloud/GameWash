package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Article struct {
	Categorie string `json:"categorie"`
	Titre     string `json:"titre"`
	Auteur    string `json:"auteur"`
	Contenu   string `json:"contenu"`
	Images    string `json:"images"`
}

type Form struct {
	Categorie    string `json:"categorie"`
	Auteur       string `json:"auteur"`
	Date         string `json:"date"`
	Introduction string `json:"introduction"`
	Texte        string `json:"texte"`
	Images       string `json:"images"`
}

//var logs Login

func main() {

	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Printf(fmt.Sprintf("ERREUR => %s", err.Error()))
		return
	}

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {

		articles, err := LoadArticles()
		if err != nil {
			fmt.Println("erreur load articles")
			return
		}
		articlesChoices := getRandomArticles(articles, 10)
		temp.ExecuteTemplate(w, "home", articlesChoices)
	})

	http.HandleFunc("/compet", func(w http.ResponseWriter, r *http.Request) {

		article, err := LoadArticlesByCategory("Compétitif")
		if err != nil {
			fmt.Println("erreur load articles")
			return
		}
		temp.ExecuteTemplate(w, "compet", article)
	})

	http.HandleFunc("/vrac", func(w http.ResponseWriter, r *http.Request) {
		article, err := LoadArticlesByCategory("Jeux en Vrac")
		if err != nil {
			fmt.Println("erreur load articles")
			return
		}
		temp.ExecuteTemplate(w, "vrac", article)
	})

	http.HandleFunc("/coeur", func(w http.ResponseWriter, r *http.Request) {
		article, err := LoadArticlesByCategory("Coup de Coeur")
		if err != nil {
			fmt.Println("erreur load articles")
			return
		}
		temp.ExecuteTemplate(w, "coeur", article)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "login", nil)
	})

	http.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		recherche := r.URL.Query().Get("content")

		articles, _ := rechercheTitre("data.json", recherche)

		temp.ExecuteTemplate(w, "result", articles)
	})

	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "contact", nil)
	})

	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "error", nil)
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "about", nil)
	})

	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "admin", nil)
	})

	http.HandleFunc("/form/treatment", FormSubmission)

	rootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(rootDoc + "/asset"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))
	fmt.Println("Serveur démarré sur le port 8080...")
	http.ListenAndServe("localhost:8080", nil)
}












