package dnsp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

type httpServer struct {
	server *Server
}

func (h *httpServer) index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	index, err := os.Open("/Users/leavengood/go/src/github.com/gophergala/dnsp/web-ui/index.html")
	if err == nil {
		io.Copy(w, index)
	} else {
		fmt.Fprintf(w, "Error: %v", err)
	}
}

func (h *httpServer) logo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	index, err := os.Open("/Users/leavengood/go/src/github.com/gophergala/dnsp/web-ui/logo.png")
	if err == nil {
		io.Copy(w, index)
	} else {
		fmt.Fprintf(w, "Error: %v", err)
	}
}

func (h *httpServer) mode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	mode := "black"
	if h.server.white {
		mode = "white"
	}
	fmt.Fprintf(w, "{\"mode\":%q}\n", mode)
}

func (h *httpServer) publicListCount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.server.publicEntriesCount())
}

func (h *httpServer) list(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.server.privateHostEntries())
}

func (h *httpServer) add(white bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		if h.server.white && !white {
			w.WriteHeader(422)
			w.Write([]byte(`{"error":"unprocessable_entity","message":"server is running in whitelist mode"`))
			return
		}
		if !h.server.white && white {
			w.WriteHeader(422)
			w.Write([]byte(`{"error":"unprocessable_entity","message":"server is running in blacklist mode"`))
			return
		}

		h.server.addPrivateHostEntry(ps.ByName("url"))
		w.WriteHeader(http.StatusAccepted)
	}
}

func (h *httpServer) remove(white bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		if h.server.white && !white {
			w.WriteHeader(422)
			w.Write([]byte(`{"error":"unprocessable_entity","message":"server is running in whitelist mode"`))
			return
		}
		if !h.server.white && white {
			w.WriteHeader(422)
			w.Write([]byte(`{"error":"unprocessable_entity","message":"server is running in blacklist mode"`))
			return
		}

		h.server.removePrivateHostEntry(ps.ByName("url"))
		w.WriteHeader(http.StatusAccepted)
	}
}

func RunHTTPServer(host string, s *Server) {
	h := httpServer{server: s}

	router := httprouter.New()

	router.GET("/", h.index)
	router.GET("/logo.png", h.logo)

	router.GET("/mode", h.mode)

	// Gets the count for the public blacklist
	router.GET("/blacklist/public", h.publicListCount)

	// Gets the personal blacklist
	router.GET("/blacklist", h.list)
	// Adds a new URL to the blacklist
	router.PUT("/blacklist/:url", h.add(false))
	// Removes a URL from the blacklist
	router.DELETE("/blacklist/:url", h.remove(false))

	// Gets the personal whitelist
	router.GET("/whitelist", h.list)
	// Adds a new URL to the whitelist
	router.PUT("/whitelist/:url", h.add(true))
	// Removes a URL from the whitelist
	router.DELETE("/whitelist/:url", h.remove(true))

	log.Fatal(http.ListenAndServe(host, router))
}
