package github2pubsub

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/300481/mq"
	"gopkg.in/go-playground/webhooks.v5/github"
)

// pubSubMessage contains the type of webhook and its payload
type pubSubMessage struct {
	Method string
	Header http.Header
	Body   []byte
}

// newGCP creates new GCP PubSub publisher struct
func newGCP() *mq.GCP {
	return &mq.GCP{
		TopicName:   os.Getenv("GCP_TOPIC_NAME"),
		CreateTopic: os.Getenv("GCP_CREATE_TOPIC") == "TRUE",
		ProjectID:   os.Getenv("GCP_PROJECT_ID"),
	}
}

// cloneHeader copies the request Header into a new string map
func cloneHeader(header http.Header) http.Header {
	var b bytes.Buffer

	r, err := http.NewRequest("POST", "http://example.com", &b)
	if err != nil {
		log.Println("error while cloning Header")
		return nil
	}

	r.Header = header
	return r.Header
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
	_, err = hook.Parse(r, gitHubEvents...)
	if err != nil {
		handleError(err, w)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(err, w)
		return
	}

	message, err := json.Marshal(&pubSubMessage{
		Method: r.Method,
		Header: cloneHeader(r.Header),
		Body:   body,
	})
	if err != nil {
		handleError(err, w)
		return
	}

	// publish the message to GCP PubSub
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

// error handling function
func handleError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	log.Printf("error: %s", err.Error())
}
