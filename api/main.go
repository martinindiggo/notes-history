package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/appengine/datastore"
)

// Note is a datastore entity that represents a single Note.
type Note struct {
	UID         *datastore.Key `json:"uid" datastore:"-"`
	MeetingID   string         `json:"meetingId"`
	Author      string         `json:"author"`
	Content     string         `json:"content"`
	CreatedDate time.Time      `json:"createdDate"`
}

func init() {
	router := mux.NewRouter().StrictSlash(true)

	router.Methods("GET").Path("/api/notes").Handler(handler(listAll))
	router.Methods("GET").Path("/api/notes/meeting/{id}").Handler(handler(listByMeeting))
	router.Methods("GET").Path("/api/notes/meeting/{id}/latest").Handler(handler(listLatestByMeeting))

	router.Methods("POST").Path("/api/notes").Handler(handler(addNote))

	router.Methods("DELETE").Path("/api/notes/{id}").Handler(handler(deleteNote))

	http.Handle("/", router)
}
