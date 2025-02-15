package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pquerna/otp/totp"
)

// Global configuration variables
var (
	PORT         int
	FOLDER       string = ".GRY"
	STORAGE_PATH string
	TOTP_SECRET  string
)

// LogMiddleware wraps handlers with combined logging functionality
func LogMiddleware(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

// HandleError writes an appropriate error response based on the error type
func HandleError(w http.ResponseWriter, err error, defaultStatus int) {
	LogErr(err)
	switch err.Error() {
	case "slug already exists":
		http.Error(w, "Nope! Path already exists.", http.StatusConflict)
	case "slug does not exist":
		http.Error(w, "Nope! Path does not exist.", http.StatusNotFound)
	default:
		http.Error(w, "Nope! Sorry.", defaultStatus)
	}
}

// RedirYourself handles redirect requests, either to the target URL or the default homepage
func RedirYourself(w http.ResponseWriter, r *http.Request) {
	var (
		slug     string
		location string
		err      error
	)
	vars := mux.Vars(r)

	if vars["slug"] != "" {
		slug = vars["slug"]

		location, err = GetURL(slug)
		if err != nil {
			HandleError(w, err, http.StatusNotFound)
			return
		}
	} else {
		location = "https://ctrlalt.dev/GRY/"
	}

	if err == nil {
		w.Header().Set("Location", location)
		w.WriteHeader(http.StatusFound)
	}
}

// HealthCheck responds to health check requests
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// ValidateTOTP validates the TOTP token
func ValidateTOTP(token string) bool {
	return totp.Validate(token, TOTP_SECRET)
}

// CheckAuthorization checks the authorization header
func CheckAuthorization(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	if ValidateTOTP(authHeader) {
		return true
	}

	return false
}

// CreateRedir handles creation of new URL redirects
func CreateRedir(w http.ResponseWriter, r *http.Request) {
	if !CheckAuthorization(r) {
		HandleError(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	var (
		slug     string
		location string
		err      error
	)

	vars := mux.Vars(r)

	if vars["slug"] == "" {
		HandleError(w, errors.New("path required"), http.StatusBadRequest)
		return
	}

	slug = vars["slug"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		HandleError(w, err, http.StatusInternalServerError)
		return
	}

	location = string(body)
	if err := ValidateURL(location); err != nil {
		HandleError(w, err, http.StatusBadRequest)
		return
	}

	err = CreateURL(slug, location)
	if err != nil {
		HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Created"))
}

// UpdateRedir handles updating existing URL redirects
func UpdateRedir(w http.ResponseWriter, r *http.Request) {
	if !CheckAuthorization(r) {
		HandleError(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	var (
		slug     string
		location string
		err      error
	)

	vars := mux.Vars(r)

	if vars["slug"] == "" {
		HandleError(w, errors.New("path required"), http.StatusBadRequest)
		return
	}

	slug = vars["slug"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		HandleError(w, err, http.StatusInternalServerError)
		return
	}

	location = string(body)
	if err := ValidateURL(location); err != nil {
		HandleError(w, err, http.StatusBadRequest)
		return
	}

	if err := UpdateURL(slug, location); err != nil {
		HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteRedir(w http.ResponseWriter, r *http.Request) {
	if !CheckAuthorization(r) {
		HandleError(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	var (
		slug string
		err  error
	)

	vars := mux.Vars(r)

	if vars["slug"] == "" {
		HandleError(w, errors.New("path required"), http.StatusBadRequest)
		return
	}

	slug = vars["slug"]

	err = DeleteURL(slug)
	if err != nil {
		LogErr(err)
		HandleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// init initializes configuration from environment variables
func init() {
	portStr, portSet := os.LookupEnv("GRY_PORT")
	if portSet {
		port, err := strconv.ParseInt(portStr, 10, 64)
		CheckErr(err)
		PORT = int(port)
	} else {
		PORT = 3000
	}

	folderStr, folderSet := os.LookupEnv("GRY_FOLDER")
	if folderSet {
		FOLDER = folderStr
	}

	userDir, userDirErr := os.UserHomeDir()
	if userDirErr != nil {
		LogErr(userDirErr)
	}

	STORAGE_PATH = filepath.Join(userDir, FOLDER)

	totpSecret, totpSecretSet := os.LookupEnv("GRY_TOTP_SECRET")
	if totpSecretSet {
		TOTP_SECRET = totpSecret
	} else {
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "https://ln.0x5f.info",
			AccountName: "GRY",
		})
		if err != nil {
			LogErr(err)
		}
		TOTP_SECRET = key.Secret()
		fmt.Printf("No TOTP_SECRET provided, generated one\n")
		fmt.Printf("TOTP_SECRET: %s\n", TOTP_SECRET)
		fmt.Printf("TOTP_SECRET won't persist, securely save it and define GRY_TOTP_SECRET\n")
		fmt.Printf("In production environment, you should define your own TOTP_SECRET and not use the default generated one\n")
	}
}

// main sets up and runs the HTTP server with graceful shutdown
func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.Use(LogMiddleware)

	r.HandleFunc("/.well-known/health", HealthCheck).Methods("GET", "HEAD")
	r.HandleFunc("/", RedirYourself).Methods("GET", "HEAD")

	r.HandleFunc("/{slug:(?:[a-zA-Z0-9_-]+)}", RedirYourself).Methods("GET", "HEAD")
	r.HandleFunc("/{slug:(?:[a-zA-Z0-9_-]+)}", CreateRedir).Methods("POST", "PUT")
	r.HandleFunc("/{slug:(?:[a-zA-Z0-9_-]+)}", UpdateRedir).Methods("PATCH")
	r.HandleFunc("/{slug:(?:[a-zA-Z0-9_-]+)}", DeleteRedir).Methods("DELETE")

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
