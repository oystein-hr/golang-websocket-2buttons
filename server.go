package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrader med default settings
var upgrader = websocket.Upgrader{}

// HTTP Server definisjon
var httpServer = &http.Server{
	Addr: ":8080",
}

var valid_origin = map[string]bool{
	"http://127.0.0.1:8080": true,
	"http://localhost:8080": true,
}

// Sette opp websocket funksjoner, mappet til /ws
func ws(w http.ResponseWriter, r *http.Request) {
	// Definerer tillate origins

	// Sjekke om det er tillate origins som kobler til
	upgrader.CheckOrigin = func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		log.Println("Origin: ", origin)
		return valid_origin[origin]
	}

	// Oppgraderer HTTP til Websocket, og sjekker for errors
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	// Beskjed som blir sent når klient kobler til
	msgType := 1 // Setter meldingtype til Text
	msg := []byte("Trykk på en av knappene.")
	conn.WriteMessage(msgType, msg)

	// Loop for å lytte til beskjeder fra klient, og sende svar
	for {
		rMsgType, rMsg, rErr := conn.ReadMessage()
		if rErr != nil {
			return
		}

		switch string(rMsg) {
		case "1":
			response := []byte("Knapp nr.1")
			conn.WriteMessage(rMsgType, response)
		case "2":
			response := []byte("Knapp nr.2")
			conn.WriteMessage(rMsgType, response)
		}
	}
}

// Sette opp filen som skal bli servert på home (/)
func httpHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client.html")
}

func main() {
	log.Println("Starting services. Press CTRL-C to exit.")

	// Koble funksjoner til URL path
	http.HandleFunc("/", httpHome)
	http.HandleFunc("/ws", ws)

	// Starte HTTP server
	httpServer.ListenAndServe()
}
