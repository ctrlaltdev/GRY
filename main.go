package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ctrlaltdev/go-redir-yourself/utils"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	PORT   int
	FOLDER string = ".redir-yourself"
)

func LogMiddleware(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

func GetURL(slug string) (url string, err error) {
	userDir, userDirErr := os.UserHomeDir()
	if userDirErr != nil {
		utils.LogErr(userDirErr)
		return "", userDirErr
	}

	content, fileErr := ioutil.ReadFile(filepath.Join(userDir, FOLDER, slug))
	if fileErr != nil {
		utils.LogErr(fileErr)
		return "", fileErr
	}

	url = string(content)

	return url, nil
}

func RedirYourself(w http.ResponseWriter, r *http.Request) {
	var (
		slug string
		url  string
		err  error
	)
	vars := mux.Vars(r)

	if vars["slug"] != "" {
		slug = vars["slug"]

		url, err = GetURL(slug)
		utils.HTTPCheckErr(w, err)
	} else {
		url = "https://ctrlalt.dev/"
	}

	if err == nil {
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
	}
}

func main() {
	portStr, portSet := os.LookupEnv("PORT")
	if portSet {
		port, err := strconv.ParseInt(portStr, 10, 64)
		utils.CheckErr(err)
		PORT = int(port)
	} else {
		PORT = 3000
	}

	folderStr, folderSet := os.LookupEnv("GO_REDIR_YOURSELF_FOLDER")
	if folderSet {
		FOLDER = folderStr
	}

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.Use(LogMiddleware)

	r.HandleFunc("/", RedirYourself).Methods("GET", "HEAD")
	r.HandleFunc("/{slug:(?:[a-zA-Z0-9_-]+)}", RedirYourself).Methods("GET", "HEAD")

	srv := &http.Server{
		Handler:      handlers.CompressHandler(r),
		Addr:         fmt.Sprintf(":%d", PORT),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("\n")
	log.Printf("starting server on port %d\n", PORT)
	fmt.Printf("\n")

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("\n")
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	srv.Shutdown(context.Background())

	fmt.Printf("\n")
	log.Println("stopping server")
	fmt.Printf("\n")
	os.Exit(0)
}
