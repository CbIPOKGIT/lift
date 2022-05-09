package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CbIPOKGIT/lift/configs"
	"github.com/CbIPOKGIT/lift/internal/mainboard"
)

func main() {
	mb, err := mainboard.New()
	if err != nil {
		log.Fatal(err)
	}
	mb.Listen()

	http.HandleFunc("/on", func(w http.ResponseWriter, r *http.Request) {
		cmd := "lift_on"
		if data, err := mb.GetData(configs.TranslateCommand(cmd)); err == nil {
			fmt.Fprintf(w, "Success. Lift ON. Response %s", data)
		} else {
			fmt.Fprintf(w, "Error: %s", err.Error())
		}
	})

	http.HandleFunc("/off", func(w http.ResponseWriter, r *http.Request) {
		cmd := "lift_off"
		if data, err := mb.GetData(configs.TranslateCommand(cmd)); err == nil {
			fmt.Fprintf(w, "Success. Lift OFF. Response %s", data)
		} else {
			fmt.Fprintf(w, "Error: %s", err.Error())
		}
	})

	if err := http.ListenAndServe(":9001", nil); err == nil {
		log.Println("Server started")
	} else {
		log.Fatal(err)
	}
}
