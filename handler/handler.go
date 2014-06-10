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
	"net/http"
	"strings"

	"github.com/meeko/go-meeko/meeko/services/logging"
)

const EventTypePrefix = "logentries"

const statusUnprocessableEntity = 422

type WebhookHandler struct {
	Logger  *logging.Service
	Forward func(eventType string, eventObject interface{}) error
}

func (handler *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Unmarshal the event object.
	payload := r.FormValue("payload")
	if payload == "" {
		http.Error(w, "Payload Form Field Missing", statusUnprocessableEntity)
		return
	}

	var eventObject map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(payload))
	if err := decoder.Decode(&eventObject); err != nil {
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
