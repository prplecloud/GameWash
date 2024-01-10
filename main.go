package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

/*type Image struct {
	Platform   string `json:"platform"`
	Background string `json:"background"`
	Studio     string `json:"studio"`
	Gameplay   string `json:"gameplay"`
}*/

/* type Article struct {
	Categorie string  `json:"categorie"`
	Titre     string  `json:"titre"`
	Auteur    string  `json:"auteur"`
	Contenu   string  `json:"contenu"`
	Images    []Image `json:"images"`
	URL       string  `json:"url"`
}

type Glaoui struct {
	Articles []Article `json:"articles"`
}*/

type Article struct {
	Categorie string `json:"categorie"`
	Titre     string `json:"titre"`
	Auteur    string `json:"auteur"`
	Contenu   string `json:"contenu"`
	Images    struct {
		Platform   string `json:"platform"`
		Background string `json:"background"`
		Studio     string `json:"studio"`
		Gameplay   string `json:"gameplay"`
	} `json:"images"`
}

type Form struct {
	Categorie    string `json:"categorie"`
	Auteur       string `json:"auteur"`
	Date         string `json:"date"`
	Introduction string `json:"introduction"`
	Texte        string `json:"texte"`
}

//var logs Login <-- login var for login.html

func main() {

	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Printf(fmt.Sprintf("ERREUR => %s", err.Error()))
		return
	}

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {

		temp.ExecuteTemplate(w, "home", nil)
	})

	http.HandleFunc("/compet", func(w http.ResponseWriter, r *http.Request) {

		temp.ExecuteTemplate(w, "compet", nil)
	})

	http.HandleFunc("/vrac", func(w http.ResponseWriter, r *http.Request) {

		temp.ExecuteTemplate(w, "vrac", nil)
	})

	http.HandleFunc("/coeur", func(w http.ResponseWriter, r *http.Request) {

		temp.ExecuteTemplate(w, "coeur", nil)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		temp.ExecuteTemplate(w, "login", nil)
	})

	//http.HandleFunc("/login/treatment", func(w http.ResponseWriter, r *http.Request) {
	//logs = Login{r.FormValue("Username"), r.FormValue("Passwd")}

	//temp.ExecuteTemplate(w, "login", nil)
	//})

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

	rootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(rootDoc + "/asset"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))
	fmt.Println("Serveur démarré sur le port 8080...")
	Json()
	http.ListenAndServe("localhost:8080", nil)

}

func Json() {

	jsonFilePath := "base.json"

	jsonData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier JSON :", err)
		return
	}

	var glaouiData []Article

	err = json.Unmarshal(jsonData, &glaouiData)
	if err != nil {
		fmt.Println("Erreur lors du marshal de la struct en JSON :", err)
		return
	}

	fmt.Println(glaouiData)
}

/*func handleFormSubmission(w http.ResponseWriter, r *http.Request) {

	// Decode form values
	var form Form

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		fmt.Println("Erreur lors du décodage du fichier JSON :", err)
		return
	}

	// Read existing json data
	var articles []Article

	// Unmarshal json

	err = json.Unmarshal(jsonData, &articles)
	if err != nil {
		fmt.Println("Erreur lors du marshal de la struct en JSON :", err)
		return
	}

	// Append new article
	articles = append(articles, Article{
		Categorie: form.Categorie,
		// map other fields
	})

	// Marshall back to json
	data, err := json.Marshal(articles)

	// Write updated json to file

	http.Redirect(w, r, "/", http.StatusSeeOther)

}*/
