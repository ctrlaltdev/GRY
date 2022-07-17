package utils

import (
	"fmt"
	"net/http"
)

func LogErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func HTTPCheckErr(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, "Nope! Sorry.", http.StatusNotFound)
	}
}
