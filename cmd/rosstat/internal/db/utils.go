package db

import (
	"io"
	"log"
	"net/http"
)

func ResponseForbidden(writer http.ResponseWriter, err error) {
	writer.WriteHeader(403)
	_, err = io.WriteString(writer, "Forbidden")
	if err != nil {
		log.Printf("unable to send forbidden error response: %s\n", err)
	}
}

func ResponseBadRequest(writer http.ResponseWriter, err error) {
	writer.WriteHeader(400)
	_, err = io.WriteString(writer, "BadRequest: " + err.Error())
	if err != nil {
		log.Printf("unable to send bad request error response: %s\n", err)
	}
}

func ResponseInternalError(writer http.ResponseWriter, err error) {
	log.Printf("internal error: %s\n", err)

	writer.WriteHeader(500)
	_, err = io.WriteString(writer, "InternalError: " + err.Error())
	if err != nil {
		log.Printf("unable to send internal error response: %s\n", err)
	}
}
