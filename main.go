package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CbIPOKGIT/lift/configs"
	"github.com/CbIPOKGIT/lift/internal/mainboard"
)

func main() {
	// web.StartServer()
	mb, err := mainboard.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Do search")
	res := mb.P485.Search()
	log.Println(res)

	commands := []string{"ATS?", "ATO?", "ATQUERY=255,\"ATSEARCH\""}

	for _, command := range commands {
		log.Printf("Executing command '%s'\n", command)
		if data, err := mb.GetData(command); err == nil {
			log.Println("Success")
			log.Println(data)
		} else {
			log.Println("Error")
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}

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
