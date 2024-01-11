package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type Image struct {
	Platform   string `json:"platform"`
	Background string `json:"background"`
	Studio     string `json:"studio"`
	Gameplay   string `json:"gameplay"`
}

/* type Article struct {
	Categorie string  `json:"categorie"`
	Titre     string  `json:"titre"`
	Auteur    string  `json:"auteur"`
	Contenu   string  `json:"contenu"`
	Images    []Image `json:"images"`
	URL       string  `json:"url"`
}
*/

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

//var logs Login

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

	jsonFilePath := "./base.json"

	jsonData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier JSON :", err)
		return
	}

	var ArticleData []Article

	err = json.Unmarshal(jsonData, &ArticleData)
	if err != nil {
		fmt.Println("Erreur lors du marshal de la struct en JSON :", err)
		return
	}
	fmt.Println(ArticleData)
}

func FormSubmission(w http.ResponseWriter, r *http.Request) {

	var form Form

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		fmt.Println("Erreur de décodage", err)
		return
	}

	// Lecture JSON
	jsonFile, err := os.Open("data.json")
	if err != nil {
		fmt.Println("Erreur de lecture", err)
		return
	}
	defer jsonFile.Close()

	var articles []Form

	// Unmarshal JSON
	err = json.NewDecoder(jsonFile).Decode(&articles)
	if err != nil {
		fmt.Println("Erreur d'unMarshal", err)
		return
	}

	// Append article
	articles = append(articles, Form{
		Categorie:    form.Categorie,
		Auteur:       form.Auteur,
		Date:         form.Date,
		Introduction: form.Introduction,
		Texte:        form.Texte,
	})

	// Marshal nouvelles data
	Data, err := json.MarshalIndent(articles, "", "  ")
	if err != nil {
		fmt.Println("Erreur de Marshal", err)
		return
	}

	// Ecriture dans le JSON
	err = os.WriteFile("data.json", Data, 0644)
	if err != nil {
		fmt.Println("Erreur d'écriture", err)
		return
	}

	http.Redirect(w, r, "http://localhost:8080/form/treatment", http.StatusSeeOther)

}
