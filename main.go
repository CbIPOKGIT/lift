package main

import (
	"log"

	"github.com/CbIPOKGIT/lift/internal/logger"
	"github.com/CbIPOKGIT/lift/internal/mainboard"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading ENV file")
	}

	// if err := database.InitDB(); err != nil {
	// 	log.Fatal(err)
	// }
}

func main() {

	logger := logger.New()

	mb, err := mainboard.New()
	if err != nil {
		log.Fatal(err)
	}

	mb.Listen(logger)

	select {}

	// http.HandleFunc("/on", func(w http.ResponseWriter, r *http.Request) {
	// 	cmd := "lift_on"
	// 	if data, err := mb.GetData(configs.TranslateCommand(cmd)); err == nil {
	// 		fmt.Fprintf(w, "Success. Lift ON. Response %s", data)
	// 	} else {
	// 		fmt.Fprintf(w, "Error: %s", err.Error())
	// 	}
	// })

	// http.HandleFunc("/off", func(w http.ResponseWriter, r *http.Request) {
	// 	// cmd := "lift_off"
	// 	commands := []string{`ATLCDNEWSCREEN="Text"`, "ATLCDLIGHTOFF", "ATRESET"}

	// 	for _, command := range commands {
	// 		if data, err := mb.GetData(command); err == nil {
	// 			log.Printf("Command %s - %s\n", command, data)
	// 		} else {
	// 			fmt.Fprintf(w, "Error: %s", err.Error())
	// 		}
	// 		time.Sleep(time.Second)
	// 	}
	// })

	// if err := http.ListenAndServe(":9001", nil); err == nil {
	// 	log.Println("Server started")
	// } else {
	// 	log.Fatal(err)
	// }
}
