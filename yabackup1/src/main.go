package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	debugLog *log.Logger
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

var CId = "859eee7fe42742f485c982b813646431"
var Cs = "5534d58aa29141529ddd9fd6193caf4e"

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

	data := GetTokenResponse{CheckCodeUrl: GetCheckCodeUrl(CId), AlertMessages: alertMessages}
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
		token, err := CreateToken(CId, Cs, r.PostFormValue("check_code"))
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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8099"
	}

	// TODO Надо как-то узнать включен дебаг или нет
	debugLog := log.New(NewNullWriter(), "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile)
	debugLog = log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Test log messages
	debugLog.Println("hello")
	infoLog.Println("hello")
	errorLog.Println("hello")

	// Инициализируем новую структуру с зависимостями приложения.
	app := &Application{
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

	infoLog.Printf("Запуск веб-сервера на http://127.0.0.1:%s", port)
	err := http.ListenAndServe(":"+port, router)
	log.Fatal(err)
}
