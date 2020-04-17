package main

import "strings"

type hub struct {
	channels       map[string]*channel
	clients        map[string]*client
	commands       chan command
	deregistration chan *client
	registration   chan *client
}

func newHub() *hub {
	return &hub{
		registration:   make(chan *client),
		deregistration: make(chan *client),
		clients:        make(map[string]*client),
		channels:       make(map[string]*channel),
		commands:       make(chan command),
	}
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.registration:
			h.register(client)
		case client := <-h.deregistration:
			h.unregister(client)
		case cmd := <-h.commands:
			switch cmd.id {
			case JOIN:
				h.joinChannel(cmd.sender, cmd.recipient)
			case LEAVE:
				h.leaveChannel(cmd.sender, cmd.recipient)
			case MSG:
				h.message(cmd.sender, cmd.recipient, cmd.body)
			case USRS:
				h.listUsers(cmd.sender)
			case CHNS:
				h.listChannels(cmd.sender)
			default:
				panic("Incorrect command")
			}
		}
	}
}

func (h *hub) register(c *client) {
	if _, exists := h.clients[c.username]; exists {
		c.username = ""
		c.conn.Write([]byte("Error : username already taken \n"))
	} else {
		h.clients[c.username] = c
		c.conn.Write([]byte("OK\n"))
	}
}

func (h *hub) unregister(c *client) {
	if _, exists := h.clients[c.username]; exists {
		delete(h.clients, c.username)
		for _, channel := range h.channels {
			delete(channel.clients, c)
		}
	}
}

func (h *hub) joinChannel(u string, c string) {
	if client, ok := h.clients[u]; ok {
		if channel, ok := h.channels[c]; ok {
			// Channel exists, Join
			channel.clients[client] = true
		} else {
			// Channel doesn't exists, create and join
			ch := newChannel(c)
			ch.clients[client] = true
			h.channels[c] = ch
		}
		client.conn.Write([]byte("OK\n"))
	}
}

func (h *hub) leaveChannel(u string, c string) {
	if client, ok := h.clients[u]; ok {
		if channel, ok := h.channels[c]; ok {
			delete(channel.clients, client)
		}
	}
}

func (h *hub) message(u, r string, m []byte) {
	if sender, ok := h.clients[u]; ok {
		switch r[0] {
		case '#':
			if channel, ok := h.channels[r]; ok {
				if _, ok := channel.clients[sender]; ok {
					channel.broadcast(sender.username, m)
				} else {
					sender.conn.Write([]byte("Error : No such channel exists"))
				}

			}
		case '@':
			if user, ok := h.clients[r]; ok {
				msg := append([]byte(user.username+":"), m...)
				msg = append(msg, '\n')
				user.conn.Write(msg)
			} else {
				sender.conn.Write([]byte("Error : No such user is present"))
			}
		default:
			sender.conn.Write([]byte("Error : Invalid MSG command"))
		}

	}
}
func (h *hub) listChannels(u string) {
	if client, ok := h.clients[u]; ok {
		var names []string
		if len(h.channels) == 0 {
			client.conn.Write([]byte("Error : No such channel is found"))
		}
		for c := range h.channels {
			names = append(names, "#"+c+" ")
		}
		resp := strings.Join(names, ", ")
		client.conn.Write([]byte(resp + "\n"))
	}
}

func (h *hub) listUsers(u string) {
	if client, ok := h.clients[u]; ok {
		var names []string
		for c, _ := range h.clients {
			names = append(names, "@"+c+" ")
		}
		resp := strings.Join(names, ", ")
		client.conn.Write([]byte(resp + "\n"))
	}
}
