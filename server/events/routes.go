package events

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Stream struct {
	Type string
	Msg  string
}

type Client struct {
	uuid string
	ch   chan Stream
}

type Event struct {
	Uuid string
	Type string
	Msg  string
}

type EventRouter struct {
	register   chan *Client
	unregister chan string
	event      chan *Event
	clients    map[string]chan Stream
}

func SetupRouter(g *gin.RouterGroup) chan *Event {
	h := NewEventRouter()
	g.GET("", h.GetEvents)
	g.PUT("", h.PutEvent)

	go h.Run()

	return h.event
}

func NewEventRouter() *EventRouter {
	return &EventRouter{
		register:   make(chan *Client),
		unregister: make(chan string),
		event:      make(chan *Event),
		clients:    make(map[string]chan Stream),
	}
}

func (h *EventRouter) Run() {
	for {
		select {
		case client := <-h.register:
			h.AddClient(client.uuid, client.ch)
		case uuid := <-h.unregister:
			h.RemoveClient(uuid)
		case event := <-h.event:
			if ch, ok := h.clients[event.Uuid]; ok {
				ch <- Stream{Type: event.Type, Msg: event.Msg}
			}
		}
	}
}

func (h *EventRouter) AddClient(uuid string, c chan Stream) {
	h.clients[uuid] = c
}

func (h *EventRouter) RemoveClient(uuid string) {
	delete(h.clients, uuid)
}

func (h *EventRouter) PutEvent(c *gin.Context) {
	var event Event
	if err := c.BindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.event <- &event
}

func (h *EventRouter) GetEvents(c *gin.Context) {
	uuid := uuid.New().String()[:8]
	ch := make(chan Stream)
	h.register <- &Client{uuid: uuid, ch: ch}
	defer func() {
		h.unregister <- uuid
	}()

	fmt.Println("Connected", uuid)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Status(200)
	c.Writer.Flush()

	c.SSEvent("message", "connected")
	c.SSEvent("uuid", uuid)
	c.Writer.Flush()

	for {
		select {
		case <-c.Writer.CloseNotify():
			return
		case stream := <-ch:
			c.SSEvent(stream.Type, stream.Msg)
			c.Writer.Flush()
		}
	}
}
