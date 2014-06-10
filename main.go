/*
  Copyright (C) 2014  The meeko-collector-logentries AUTHORS

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

package main

import (
	"github.com/meeko-contrib/meeko-collector-logentries/handler"

	"github.com/meeko-contrib/go-meeko-webhook-receiver/receiver"
	"github.com/meeko/go-meeko/agent"
)

func main() {
	var (
		logger    = agent.Logging()
		publisher = agent.PubSub()
	)
	receiver.ListenAndServe(&handler.WebhookHandler{
		logger,
		func(eventType string, eventObject interface{}) error {
			logger.Infof("Forwarding %s", eventType)
			return publisher.Publish(eventType, eventObject)
		},
	})
}
