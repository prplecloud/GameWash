package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"time"
)

/*
	type Image struct {
		Platform   string `json:"platform"`
		Background string `json:"background"`
		Studio     string `json:"studio"`
		Gameplay   string `json:"gameplay"`
	}

*/

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
		temp.ExecuteTemplate(w, "home", nil)
	})

	http.HandleFunc("/compet", func(w http.ResponseWriter, r *http.Request) {

		article, err := LoadArticles()
		if err != nil {
			fmt.Println("erreur load articles")
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
	jsonFilePath := "base.json"
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
	nomFichier := "data.json"

	// Récupérer les données du formulaire de la requête HTTP
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erreur lors de l'analyse du formulaire", http.StatusInternalServerError)
		return
	}

	fmt.Println(r.Form.Get("categorie"))

	fmt.Println(r.Form.Get("auteur"))

	dataFile, headerFile, errFile := r.FormFile("file")
	if errFile != nil {

		fmt.Println("erreur avec le fichier....")
	}
	defer dataFile.Close()

	File, errOpen := os.OpenFile(("./asset/uploads/" + headerFile.Filename), os.O_CREATE, 0644)
	if errOpen != nil {

		fmt.Println("Erreur ouverture")
	}

	defer File.Close()

	_, errCopy := io.Copy(File, dataFile)
	if errCopy != nil {

		fmt.Println("erreur avec la copie du fichier....")
	}

	// Créer une nouvelle instance de Form à partir des données du formulaire
	form := Form{
		Categorie:    r.Form.Get("categorie"),
		Auteur:       r.Form.Get("auteur"),
		Introduction: r.Form.Get("introduction"),
		Texte:        r.Form.Get("texte"),
		Images:       "test",
	}

	fmt.Println(form)
	// Ajouter la date actuelle si elle n'est pas fournie dans le formulaire
	if form.Date == "" {
		form.Date = time.Now().Format("2006-01-02")
	}

	dataForms, errForms := LoadArticles()
	if errForms != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'ouverture du fichier : %v", errForms), http.StatusInternalServerError)
		return
	}

	// Ajouter la nouvelle forme à la liste
	dataForms = append(dataForms, form)
	fmt.Println(dataForms)

	dataWrite, errWrite := json.Marshal(dataForms)
	if errWrite != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'ouverture du fichier : %v", errWrite), http.StatusInternalServerError)
		return
	}

	errWriteFile := os.WriteFile(nomFichier, dataWrite, fs.FileMode(0644))
	if errWriteFile != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'ouverture du fichier : %v", errWriteFile), http.StatusInternalServerError)
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
		fmt.Println("1")
		return nil, err
	}

	var articles []Form

	err = json.Unmarshal(jsonData, &articles)
	if err != nil {
		fmt.Println("2")
		return nil, err
	}
	fmt.Println("3")
	return articles, nil
}

/*
// Fonction de recherche dans les données
func search(query string, articles Article) Article {

	var results Article

	for _, article := range articles {

		if strings.Contains(strings.ToLower(article.Titre), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(article.Contenu), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(article.Auteur), strings.ToLower(query)) {
			results = append(results, article)
		}
	}

	return results
}

func searchBar(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Empty search query", http.StatusBadRequest)
		return
	}

	articles, err := LoadArticles()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading JSON: %s", err), http.StatusInternalServerError)
		return
	}

	results := search(query, articles)

	// Afficher les résultats dans le navigateur
	fmt.Println(w, "Résultats de la recherche pour '%s':\n", query)
	for _, result := range results {

		fmt.Println(w, "Catégorie: %s, Titre: %s, Auteur: %s, Contenu: %s, Images: %s\n",
			result.Categorie, result.Titre, result.Auteur, result.Contenu, result.Images)
	}
}
*/
