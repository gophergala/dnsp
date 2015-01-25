package dnsp

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type httpServer struct {
	server *Server
}

func (h *httpServer) index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, err := Asset("web-ui/index.html")
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	} else {
		w.Write(data)
	}
}

func (h *httpServer) logo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, err := Asset("web-ui/logo.png")
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	} else {
		w.Write(data)
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
	fmt.Fprintf(w, "{\"count\":%d}\n", 0) // TODO: h.server.publicListCount())
}

func (h *httpServer) add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := ps.ByName("url")

	h.server.addHostEntry(url)

	// TODO: different response?
	fmt.Fprintf(w, "{added:%q}\n", url)
}

func (h *httpServer) remove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := ps.ByName("url")

	h.server.removeHostEntry(url)

	// TODO: different response?
	fmt.Fprintf(w, "{removed:%q}\n", url)
}

func RunHTTPServer(host string, s *Server) {
	h := httpServer{server: s}

	router := httprouter.New()

	router.GET("/", h.index)
	router.GET("/logo.png", h.logo)

	router.GET("/mode", h.mode)

	// Gets the count for the public blacklist
	router.GET("/blacklist/public", h.publicListCount)

	// Adds a new URL to the list
	router.PUT("/list/:url", h.add)
	// Removes a URL from the list
	router.DELETE("/list/:url", h.remove)

	log.Fatal(http.ListenAndServe(host, router))
}
