package api

import (
	"encoding/json"
	"net/http"
	"time"

	"io/ioutil"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// listAll returns all notes from the datastore.
func listAll(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)

	qry := datastore.NewQuery("Note").Order("MeetingID").Order("-CreatedDate")

	var notes = make([]Note, 0)
	keys, err := qry.GetAll(c, &notes)
	if err != nil {
		return nil, &handlerError{err, err.Error(), http.StatusInternalServerError}
	}

	for i, k := range keys {
		notes[i].UID = k
	}

	return notes, nil
}

// listByMeeting returns all notes for a specified meetingID from the datastore.
func listByMeeting(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)

	meetingID := mux.Vars(r)["id"]

	qry := datastore.NewQuery("Note").Filter("MeetingID =", meetingID).Order("-CreatedDate")

	var notes = make([]Note, 0)
	keys, err := qry.GetAll(c, &notes)
	if err != nil {
		return nil, &handlerError{err, "Could not retrieve notes from datastore.", http.StatusInternalServerError}
	}

	for i, k := range keys {
		notes[i].UID = k
	}

	return notes, nil
}

// listLatestByMeeting returns the latest note for specified meetingID from the datastore.
func listLatestByMeeting(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)

	meetingID := mux.Vars(r)["id"]

	qry := datastore.NewQuery("Note").Filter("MeetingID =", meetingID).Order("-CreatedDate").Limit(1)

	var notes = make([]Note, 0)
	keys, err := qry.GetAll(c, &notes)
	if err != nil {
		return nil, &handlerError{err, "Could not retrieve notes from datastore.", http.StatusInternalServerError}
	}

	for i, k := range keys {
		notes[i].UID = k
	}

	return notes, nil
}

// addNote adds a new note into the datastore.
func addNote(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	var n Note

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, &handlerError{err, "Could not read request.", http.StatusBadRequest}
	}

	if err = json.Unmarshal(payload, &n); err != nil {
		return nil, &handlerError{err, "Could not parse JSON.", 422} // http.StatusUnprocessableEntity
	}

	c := appengine.NewContext(r)
	key := datastore.NewIncompleteKey(c, "Note", nil)

	n.CreatedDate = time.Now()
	key, err = datastore.Put(c, key, &n)
	if err != nil {
		return nil, &handlerError{err, "Could not save to datastore.", http.StatusInternalServerError}
	}

	n.UID = key

	return n, nil
}

// deleteNote removes a note from the datastore.
func deleteNote(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	c := appengine.NewContext(r)

	noteID := mux.Vars(r)["id"]

	id, _ := datastore.DecodeKey(noteID)

	if err := datastore.Delete(c, id); err != nil {
		return nil, &handlerError{err, "Could not delete note from datastore.", http.StatusInternalServerError}
	}

	return true, nil
}
