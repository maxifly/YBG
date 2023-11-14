package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

type BackupFileInfo struct {
	Name       string
	CreateDate string
	Size       string
	IsLocal    bool
	IsRemote   bool
}

type BackupResponse struct {
	BFiles []BackupFileInfo
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/index.html",
		"./ui/html/base.html",
	}

	// Используем функцию template.ParseFiles() для чтения файлов шаблона.
	// Если возникла ошибка, мы запишем детальное сообщение ошибки и
	// используя функцию http.Error() мы отправим пользователю
	// ответ: 500 Internal Server Error (Внутренняя ошибка на сервере)
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Затем мы используем метод Execute() для записи содержимого
	// шаблона в тело HTTP ответа. Последний параметр в Execute() предоставляет
	// возможность отправки динамических данных в шаблон.

	bFiles := make([]BackupFileInfo, 0, 0)

	bFiles = append(bFiles, BackupFileInfo{Name: "test1", CreateDate: "01.01", Size: "125", IsLocal: true, IsRemote: true})
	bFiles = append(bFiles, BackupFileInfo{Name: "test2", CreateDate: "02.02", Size: "125", IsLocal: true, IsRemote: true})

	data := BackupResponse{BFiles: bFiles}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", indexHandler)

	log.Printf("Запуск веб-сервера на http://127.0.0.1:%s", port)
	err := http.ListenAndServe(":"+port, mux)
	log.Fatal(err)
}
