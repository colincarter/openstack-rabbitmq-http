package handlers

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HandleEvents routes events from OpenStack
func HandleEvents(rabbitEventChan chan []byte) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		body, err := ioutil.ReadAll(req.Body)
		if err == nil {
			rabbitEventChan <- body
			req.Body.Close()
		} else {
			log.Print("Error reading request body")
		}
	}
}

// HandleMeters routes meters from OpenStack
func HandleMeters(rabbitMeterChan chan []byte) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		body, err := ioutil.ReadAll(req.Body)
		if err == nil {
			rabbitMeterChan <- body
			req.Body.Close()
		} else {
			log.Print("Error reading request body")
		}
	}
}

// HandlePing place to ensure service is working
func HandlePing(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	io.WriteString(w, "pong")
}
