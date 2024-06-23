package server

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"io/fs"
	"net/http"
)

//go:generate cp -r ../../ui/dist/helm-live/browser ./static
//go:embed static/*
var staticAssets embed.FS

// Event for broadcasting
type Event struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan string

	// New client connections
	NewClients chan chan string

	// Closed client connections
	ClosedClients chan chan string

	// Total client connections
	TotalClients map[chan string]bool
}

// New event messages are broadcast to all registered client connection channels
type ClientChan chan string

// RunServer run the http server
func RunServer(stream *Event, registersFnc []func(engine *gin.Engine), staticPath string, port int) {

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	for _, fnc := range registersFnc {
		fnc(router)
	}

	// Add event-streaming render endpoint
	router.GET("/stream", StreamHeadersMiddleware(), stream.serveHTTP(), func(c *gin.Context) {
		v, ok := c.Get("clientChan")
		if !ok {
			return
		}
		clientChan, ok := v.(ClientChan)
		if !ok {
			return
		}
		c.Stream(func(w io.Writer) bool {
			// Stream message to client from message channel
			if msg, ok := <-clientChan; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	})

	// If statics file are provided, use it otherwise load embedded statics
	var staticsFs http.FileSystem
	if staticPath != "" {
		staticsFs = http.Dir(staticPath)
	} else {
		subFS, err := fs.Sub(staticAssets, "static")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get statics file system")
		}
		staticsFs = http.FS(subFS)
	}

	// Add http client side statics
	router.StaticFS("/ui", staticsFs)

	log.Info().Msg("Watch mode enabled. Watching for file changes...")
	log.Info().Msgf("Local:  http://localhost:%d/ui", port)

	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}

}

// NewServer Initialize event and Start processing requests
func NewServer() (event *Event) {
	event = &Event{
		Message:       make(chan string),
		NewClients:    make(chan chan string),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[chan string]bool),
	}

	go event.listen()

	return
}

// It Listens all incoming requests from clients.
// Handles addition and removal of clients and broadcast messages to clients.
func (stream *Event) listen() {
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			log.Debug().Msgf("Client added. %d registered clients", len(stream.TotalClients))

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			close(client)
			log.Debug().Msgf("Removed client. %d registered clients", len(stream.TotalClients))

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.TotalClients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

func (stream *Event) serveHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize client channel
		clientChan := make(ClientChan)

		// Send new connection to event server
		stream.NewClients <- clientChan

		defer func() {
			// Send closed connection to event server
			stream.ClosedClients <- clientChan
		}()

		c.Set("clientChan", clientChan)

		c.Next()
	}
}

func StreamHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}
