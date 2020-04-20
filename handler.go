package debugger

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	cdp "github.com/chromedp/cdproto"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type version struct {
	Browser         string
	ProtocolVersion string `json:"Protocol-Version"`
}

// Handler implements a HTTP handler that implements the standard debugger protocol.
func Handler(d Driver) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request", r.Method, r.URL.Path)

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		targets, err := d.List()
		if err != nil {
			writeJSON(w, nil, err)
			return
		}

		switch r.URL.Path {
		case "/json/version":
			writeJSON(w, version{d.Name(), "1.1"}, nil)
			return
		case "/json", "/json/", "/json/list", "/json/list/":
			writeJSON(w, targets, err)
			return
		}
		for _, target := range targets {
			if r.URL.Path != "/"+target.ID {
				continue
			}

			log.Println("New debug session for", target.ID)
			conn, err := websocket.Accept(w, r, nil)
			if err != nil {
				log.Println(err)
			} else {
				handleConnection(conn, d)
			}
			return
		}
	})
}

func handleConnection(conn *websocket.Conn, d Driver) {
	defer conn.Close(websocket.StatusInternalError, "fin fin")
	for {
		var msg cdp.Message
		if err := wsjson.Read(context.Background(), conn, &msg); err != nil {
			log.Println("Unexpected error", err)
			return
		}

		log.Println("Got message", msg.ID, msg.Method, string(msg.Params))
		d.Dispatch(d, &msg)
		msg.Params = nil
		msg.Method = ""
		if err := wsjson.Write(context.Background(), conn, &msg); err != nil {
			log.Println("Unexpected message write error", err)
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, data interface{}, err error) {
	if err == nil {
		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(data)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
