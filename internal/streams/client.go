package streams

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/jpillora/backoff"
	"github.com/r3labs/sse"
)

const (
	// DefaultURL is the URL of the Wikimedia EventStreams service (sans any stream endpoints)
	DefaultURL string = "https://stream.wikimedia.org/v2/stream"

	// Wikimedia's traffic layer will disconnect clients after 15 minutes (see: https://phabricator.wikimedia.org/T242767).
	// This manifests as an http.http2StreamError (https://golang.org/pkg/net/http/?m=all#http2StreamError) returned by
	// sse.Client#Subscribe (code http2ErrCodeNo), (and this error is NOT handled by sse's ReconnectStrategy).  To work
	// around this, errors that are spaced at least `resetInterval` apart will be tried `retries` times with an
	// exponential back-off.
	resetInterval = time.Minute * 10
	retries       = 3
)

// Client is used to subscribe to the Wikimedia EventStreams service
type Client struct {
	BaseURL       string
	Predicates    map[string]interface{}
	Since         string
	lastTimestamp string
}

// NewClient returns an initialized Client
func NewClient(event string) *Client {
	return &Client{BaseURL: DefaultURL, Predicates: make(map[string]interface{})}
}

// Match adds a new predicate.  Predicates are used to used to establish a match based on the JSON
// attribute name.  Events match only when all predicates do.
func (client *Client) Match(attribute string, value interface{}) *Client {
	client.Predicates[attribute] = value
	return client
}

// LastTimestamp returns the ISO8601 formatted timestamp of the last event received.
func (client *Client) LastTimestamp() string {
	return client.lastTimestamp
}

// RecentChanges subscribes to the recent changes feed. The handler is invoked with a
// RecentChangeEvent once for every matching event received.
func (client *Client) GetStreamData(stream string, handlerFunc interface{}) error {
	var bOff = &backoff.Backoff{}
	for {
		var lastSub time.Time
		err := client.UnmarshalJSONBasedOnStream(stream, handlerFunc, &lastSub)
		if err == nil {
			return err
		}
		if time.Since(lastSub) >= resetInterval {
			bOff.Reset()
		}
		time.Sleep(bOff.Duration())
		if bOff.Attempt() >= retries {
			return err
		}
		client.Since = client.lastTimestamp
	}
}

func (client *Client) UnmarshalJSONBasedOnStream(stream string, handlerFunc interface{}, lastSub *time.Time) error {
	sseClient := sse.NewClient(client.url(stream))

	handler := reflect.ValueOf(handlerFunc)

	if handler.Kind() != reflect.Func {
		return errors.New("handler must be a function")
	}

	return sseClient.Subscribe("message", func(msg *sse.Event) {
		if len(msg.Data) == 0 {
			return
		}

		handlerType := handler.Type()
		if handlerType.NumIn() != 1 {
			log.Println("Handler function must take exactly one parameter")
			return
		}

		eventDataPtr := reflect.New(handlerType.In(0))

		// Navigate through embedded fields to find the actual type
		actualType := eventDataPtr.Type()
		for actualType.Kind() == reflect.Ptr {
			actualType = actualType.Elem()
		}

		if actualType != handlerType.In(0) {
			log.Println("Handler function parameter type doesn't match the event type")
			return
		}

		if err := json.Unmarshal(msg.Data, eventDataPtr.Interface()); err != nil {
			log.Printf("Error deserializing JSON event: %s\n", err)
			return
		}

		// Type assertion to call the specific handler
		handler.Call([]reflect.Value{eventDataPtr.Elem()})
		*lastSub = time.Now()
	})
}

func (client *Client) url(stream string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/%s", client.BaseURL, stream)
	if client.Since != "" {
		fmt.Fprintf(&b, "?since=%s", client.Since)
	}
	return b.String()
}
