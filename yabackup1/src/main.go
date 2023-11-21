package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
)

const FILE_PATH_OPTIONS = "/data/options.json"
const FILE_PATH_TOKEN = "/data/token.json"
const BACKUP_PATH = "/backup"

type ApplOptions struct {
	ClientId                   string `json:"client_id"`
	ClientSecret               string `json:"client_secret"`
	RemotePath                 string `json:"remote_path"`
	RemoteMaximumFilesQuantity int    `json:"remote_maximum_files_quantity"`
	Schedule                   string `json:"schedule"`
	LogLevel                   string `json:"log_level"`
}

type Application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	debugLog *log.Logger
	options  ApplOptions
}

type AlertMessage struct {
	Message string
}
type BackupResponse struct {
	AlertMessages []AlertMessage
	BFiles        []BackupFileInfo
}

type GetTokenResponse struct {
	AlertMessages []AlertMessage
	CheckCodeUrl  string
}

func (app *Application) indexHandler(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("indexHandler")
	files := []string{
		"./ui/html/index.html",
		"./ui/html/base.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	filesInfo, err := GetFilesInfo(app)

	data := BackupResponse{BFiles: filesInfo}

	err = ts.Execute(w, data)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *Application) getTokenForm(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("getTokenForm")
	app.renderTokenForm(w, r, "")
}
func (app *Application) renderTokenForm(w http.ResponseWriter, r *http.Request, errorMessage string) {
	app.infoLog.Println("getTokenForm")
	files := []string{
		"./ui/html/get_token.html",
		"./ui/html/base.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	alertMessages := make([]AlertMessage, 0)
	if errorMessage != "" {
		alertMessages = append(alertMessages, AlertMessage{Message: errorMessage})
	}

	data := GetTokenResponse{CheckCodeUrl: GetCheckCodeUrl(app.options.ClientId), AlertMessages: alertMessages}
	err = ts.Execute(w, data)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *Application) getToken(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("getToken")
	checkCode := r.PostFormValue("check_code")
	if checkCode == "" {
		app.renderTokenForm(w, r, "Check code is required!")
	} else {
		token, err := CreateToken(
			app.options.ClientId,
			app.options.ClientSecret,
			r.PostFormValue("check_code"))
		if err != nil {
			app.errorLog.Println(err.Error())
			http.Error(w, "Create Token Error", 500)
		}
		app.infoLog.Printf("Token %v", token)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//func startUpload(w http.ResponseWriter, r *http.Request) {
//	log.Println("startUpload")
//
//}

func readOptions() (ApplOptions, error) {
	plan, _ := os.ReadFile(FILE_PATH_OPTIONS)
	var data ApplOptions
	err := json.Unmarshal(plan, &data)
	return data, err
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8099"
	}

	debugLog := log.New(NewNullWriter(), "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog := log.New(NewNullWriter(), "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Test log messages

	options, err := readOptions()
	if err != nil {
		panic(fmt.Sprintf("Can not read options: %v", err))
	}

	if options.LogLevel == "DEBUG" {
		debugLog = log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile)
		infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	}
	if options.LogLevel == "INFO" {
		infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	}

	debugLog.Println("hello")
	infoLog.Println("hello")
	errorLog.Println("hello")

	// Инициализируем новую структуру с зависимостями приложения.
	app := &Application{
		options:  options,
		errorLog: errorLog,
		infoLog:  infoLog,
		debugLog: debugLog,
	}

	router := mux.NewRouter()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))
	//router.Handle("/static/", http.StripPrefix("/static", fileServer))

	router.HandleFunc("/", app.indexHandler).Methods("GET")
	router.HandleFunc("/index", app.indexHandler).Methods("GET")
	router.HandleFunc("/get_token", app.getTokenForm).Methods("GET")
	router.HandleFunc("/get_token", app.getToken).Methods("POST")
	//mux.HandleFunc("/start_upload", startUpload)

	errorLog.Printf("(It is not error!!!) Run WEB-Server on http://127.0.0.1:%s", port)
	err = http.ListenAndServe(":"+port, router)
	log.Fatal(err)
}
