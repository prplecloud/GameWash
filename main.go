package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"
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
	Images       struct {
		Platform   string `json:"platform"`
		Background string `json:"background"`
		Studio     string `json:"studio"`
		Gameplay   string `json:"gameplay"`
	} `json:"images"`
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
		article, err := LoadArticles()
		if err != nil {
			fmt.Println("erreur")
			return
		}
		temp.ExecuteTemplate(w, "compet", article)
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

	http.HandleFunc("/form/treatment", FormSubmission)

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
	// Chemin du fichier JSON
	nomFichier := "./data.json"

	// Récupérer les données du formulaire de la requête HTTP
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erreur lors de l'analyse du formulaire", http.StatusInternalServerError)
		return
	}

	fmt.Println(r.Form.Get("categorie"))
	fmt.Println(r.Form.Get("auteur"))

	// Créer une nouvelle instance de Form à partir des données du formulaire
	form := Form{
		Categorie:    r.Form.Get("categorie"),
		Auteur:       r.Form.Get("auteur"),
		Introduction: r.Form.Get("introduction"),
		Texte:        r.Form.Get("texte"),
	}

	fmt.Println(form)
	// Ajouter la date actuelle si elle n'est pas fournie dans le formulaire
	if form.Date == "" {
		form.Date = time.Now().Format("2006-01-02")
	}

	// Ouvrir le fichier en mode lecture/écriture ou le créer s'il n'existe pas
	fichier, err := os.OpenFile(nomFichier, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'ouverture du fichier : %v", err), http.StatusInternalServerError)
		return
	}
	defer fichier.Close()

	// Charger le contenu actuel du fichier
	var forms []Form
	if err := json.NewDecoder(fichier).Decode(&forms); err != nil && err.Error() != "EOF" {
		http.Error(w, fmt.Sprintf("Erreur lors de la lecture du fichier JSON : %v", err), http.StatusInternalServerError)
		return
	}

	// Ajouter la nouvelle forme à la liste
	forms = append(forms, form)

	// Réécrire le fichier avec la nouvelle liste sans tronquer
	fichier.Seek(0, 0)
	if err := fichier.Truncate(0); err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la troncature du fichier : %v", err), http.StatusInternalServerError)
		return
	}
	if _, err := fichier.Seek(0, 0); err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors du positionnement du curseur au début du fichier : %v", err), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(fichier).Encode(forms); err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'écriture du fichier JSON : %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("Ajouté avec succès")
	http.Redirect(w, r, "http://localhost:8080/home", http.StatusSeeOther)
}

func ShowArticles(w http.ResponseWriter, r *http.Request) {
	// Charger les données des articles depuis le fichier JSON
	articles, err := LoadArticles()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors du chargement des articles : %v", err), http.StatusInternalServerError)
		return
	}

	// Charger le modèle HTML
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors du chargement du modèle : %v", err), http.StatusInternalServerError)
		return
	}

	// Exécuter le modèle avec les données des articles
	err = tmpl.Execute(w, articles)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'exécution du modèle : %v", err), http.StatusInternalServerError)
		return
	}
}

func LoadArticles() ([]Form, error) {
	// Charger les articles depuis le fichier JSON
	jsonFilePath := "data.json"
	jsonData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, err
	}

	var articles []Form
	err = json.Unmarshal(jsonData, &articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}
