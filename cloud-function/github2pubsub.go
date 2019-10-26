package github2pubsub

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/300481/mq"
	"gopkg.in/go-playground/webhooks.v5/github"
)

type pubSubMessage struct {
	GithubEventType github.Event
	GithubEvent     interface{}
}

// NewGCP creates new GCP PubSub struct
func newGCP() *mq.GCP {
	log.Println("Create GCP PubSub message queue config.")
	return &mq.GCP{
		TopicName:   os.Getenv("GCP_TOPIC_NAME"),
		CreateTopic: os.Getenv("GCP_CREATE_TOPIC") == "TRUE",
		ProjectID:   os.Getenv("GCP_PROJECT_ID"),
	}
}

func (p *pubSubMessage) ToJson() (b []byte, err error) {
	b, err = json.Marshal(p)
	return b, err
}

// Send sends a notification of Github Webhook to the topic
func Send(w http.ResponseWriter, r *http.Request) {
	// get the Secret from environment
	secret := os.Getenv("GITHUB_SECRET")

	// get the wanted event types from environment
	events := strings.Split(os.Getenv("GITHUB_EVENTS"), "/")

	// Create a new hook config (with secret)
	hook, err := github.New(github.Options.Secret(secret))
	if err != nil {
		handleError(err, w)
		return
	}

	// parse the hooks payload
	gitHubEvents := make([]github.Event, len(events))
	for i := range events {
		gitHubEvents[i] = github.Event(events[i])
	}
	payload, err := hook.Parse(r, gitHubEvents...)
	if err != nil {
		handleError(err, w)
		return
	}

	eventType := r.Header.Get("X-GitHub-Event")

	message, err := json.Marshal(&pubSubMessage{
		GithubEventType: github.Event(eventType),
		GithubEvent:     payload,
	})
	if err != nil {
		handleError(err, w)
		return
	}

	// publish the notification event to GCP PubSub
	id, err := newGCP().Publish(message)
	if err != nil {
		handleError(err, w)
		return
	}

	log.Printf("id: %s, notification: %s", id, string(message))

	// respond OK
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "{}")
}

func handleError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	log.Printf("error: %s", err.Error())
}
