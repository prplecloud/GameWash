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
		fmt.Println(article)
		temp.ExecuteTemplate(w, "vrac", article)
	})

	http.HandleFunc("/coeur", func(w http.ResponseWriter, r *http.Request) {
		article, err := LoadArticlesByCategory("Coup de Coeur")
		if err != nil {
			fmt.Println("erreur load articles")
			return
		}
		fmt.Println(article)
		temp.ExecuteTemplate(w, "coeur", article)
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
		Images:       headerFile.Filename,
	}

	fmt.Println(form)
	// Ajouter la date actuelle si elle n'est pas fournie dans le formulaire
	if form.Date == "" {
		form.Date = time.Now().Format("2006-01-02")
	}

	dataForms, errForms := LoadArticles()
	if errForms != nil {
		fmt.Println("1")
		http.Error(w, fmt.Sprintf("Erreur lors de l'ouverture du fichier : %v", errForms), http.StatusInternalServerError)
		return
	}

	// Ajouter la nouvelle forme à la liste
	dataForms = append(dataForms, form)
	fmt.Println(dataForms)

	dataWrite, errWrite := json.Marshal(dataForms)
	if errWrite != nil {
		fmt.Println("2")
		http.Error(w, fmt.Sprintf("Erreur lors de l'ouverture du fichier : %v", errWrite), http.StatusInternalServerError)
		return
	}

	errWriteFile := os.WriteFile(nomFichier, dataWrite, fs.FileMode(0644))
	if errWriteFile != nil {
		fmt.Println("3")
		http.Error(w, fmt.Sprintf("Erreur lors de l'ouverture du fichier : %v", errWriteFile), http.StatusInternalServerError)
		return
	}

	fmt.Println("Ajouté avec succès")
	http.Redirect(w, r, "http://localhost:8080/home", http.StatusSeeOther)
}

func ShowArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := LoadArticles()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors du chargement des articles : %v", err), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors du chargement du modèle : %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, articles)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'exécution du modèle : %v", err), http.StatusInternalServerError)
		return
	}
}

func LoadArticles() ([]Form, error) {
	fileData, err := ioutil.ReadFile("data.json")
	if err != nil {
		return nil, err
	}

	var forms []Form

	err = json.Unmarshal(fileData, &forms)
	if err != nil {
		return nil, err
	}

	return forms, nil
}

func getRandomArticles(liste []Form, nombreElements int) []Form {
	rand.Seed(time.Now().UnixNano())

	if len(liste) <= nombreElements {
		return liste
	}

	resultat := make([]Form, nombreElements)

	copieListe := append([]Form{}, liste...)

	for i := 0; i < nombreElements; i++ {
		indiceAleatoire := rand.Intn(len(copieListe))

		resultat[i] = copieListe[indiceAleatoire]

		copieListe = append(copieListe[:indiceAleatoire], copieListe[indiceAleatoire+1:]...)
	}

	return resultat
}

func LoadArticlesByCategory(category string) ([]Form, error) {
	fileData, err := ioutil.ReadFile("data.json")
	if err != nil {
		return nil, err
	}

	var allForms []Form

	err = json.Unmarshal(fileData, &allForms)
	if err != nil {
		return nil, err
	}

	var specifics []Form
	for _, form := range allForms {
		if form.Categorie == category {
			specifics = append(specifics, form)
		}
	}

	return specifics, nil
}
