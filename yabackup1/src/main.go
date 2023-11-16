package main

import (
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

type BackupResponse struct {
	BFiles []BackupFileInfo
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

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.indexHandler)
	mux.HandleFunc("/index", app.indexHandler)
	//mux.HandleFunc("/start_upload", startUpload)

	infoLog.Printf("Запуск веб-сервера на http://127.0.0.1:%s", port)
	err := http.ListenAndServe(":"+port, mux)
	log.Fatal(err)
}
