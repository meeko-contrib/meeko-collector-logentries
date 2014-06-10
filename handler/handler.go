/*
   Copyright (C) 2013-2014  The meeko-collector-logentries AUTHORS

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program. If not, see {http://www.gnu.org/licenses/}.
*/

package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/meeko/go-meeko/meeko/services/logging"
)

const EventTypePrefix = "logentries"

const (
	statusUnprocessableEntity = 422
	maxBodySize               = int64(10 << 20)
)

type WebhookHandler struct {
	Logger  *logging.Service
	Forward func(eventType string, eventObject interface{}) error
}

func (handler *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read the request body, up to 10 MB.
	bodyReader := http.MaxBytesReader(w, r.Body, maxBodySize)
	defer bodyReader.Close()

	body, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		http.Error(w, "Request Payload Too Large", http.StatusRequestEntityTooLarge)
		return
	}

	// Unmarshal the event object.
	var eventObject map[string]interface{}
	if err := json.Unmarshal(body, &eventObject); err != nil {
		http.Error(w, "Unexpected Payload: Not Json", http.StatusBadRequest)
		return
	}

	// Publish the event.
	eventTypeValue, ok := eventObject["event"]
	if !ok {
		http.Error(w, "Unexpected Payload: Event Type Not Set", statusUnprocessableEntity)
		return
	}
	eventType, ok := eventTypeValue.(string)
	if !ok {
		http.Error(w, "Unexpected Payload: Invalid Event Type Field Type", statusUnprocessableEntity)
		return
	}
	if err := handler.Forward(EventTypePrefix+"."+eventType, eventObject); err != nil {
		http.Error(w, "Event Not Published", http.StatusInternalServerError)
		handler.Logger.Critical(err)
		// This is a critical error, panic.
		panic(err)
	}

	w.WriteHeader(http.StatusAccepted)
}
