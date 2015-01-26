package dnsp

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type httpServer struct {
	server *Server
}

func (h *httpServer) index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")

	data, err := Asset("web-ui/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *httpServer) logo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "image/png")

	data, err := Asset("web-ui/logo.png")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *httpServer) mode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	mode := "black"
	if h.server.white {
		mode = "white"
	}
	w.Write([]byte(`"` + mode + `"`))
}

func (h *httpServer) publicListCount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	n := 0
	if !h.server.white {
		n = h.server.publicEntriesCount()
	}
	json.NewEncoder(w).Encode(n)
}

func (h *httpServer) list(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.server.privateHostEntries())
}

func (h *httpServer) add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	h.server.addPrivateHostEntry(ps.ByName("url"))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"ok"}`))
}

func (h *httpServer) remove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	h.server.removePrivateHostEntry(ps.ByName("url"))
	w.Write([]byte(`{"status":"ok"}`))
}

func RunHTTPServer(host string, s *Server) {
	h := httpServer{server: s}

	router := httprouter.New()

	router.GET("/", h.index)
	router.GET("/logo.png", h.logo)

	router.GET("/mode", h.mode)

	// Gets the count for the public blacklist
	router.GET("/blacklist/public", h.publicListCount)

	// Gets the current list
	router.GET("/list", h.list)
	// Adds a new URL to the list
	router.PUT("/list/:url", h.add)
	// Removes a URL from the list
	router.DELETE("/list/:url", h.remove)

	log.Fatal(http.ListenAndServe(host, router))
}
