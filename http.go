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
	fmt.Fprintf(w, "{\"count\":%d}\n", 1234) //h.server.publicListCount())
}

func (h *httpServer) list(white bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// var urls []string
		urls := []string{"1", "2", "3"}
		if white {
			// urls = {"1", "2"} //h.server.whitelist()
		} else {
			// urls = {"3", "4"} //h.server.blacklist()
		}

		encoder := json.NewEncoder(w)
		encoder.Encode(urls)
	}
}

func (h *httpServer) add(white bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		url := ps.ByName("url")
		if white {
			// h.server.addToWhitelist(url)
		} else {
			// h.server.addToBlacklist(url)
		}

		// TODO: response?
		fmt.Fprintf(w, "{add:%q}\n", url)
	}
}

func (h *httpServer) remove(white bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		url := ps.ByName("url")
		if white {
			// h.server.removeFromWhitelist(url)
		} else {
			// h.server.removeFromBlacklist(url)
		}

		// TODO: response?
		fmt.Fprintf(w, "{remove:%q}\n", url)
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
	router.GET("/blacklist", h.list(false))
	// Adds a new URL to the blacklist
	router.PUT("/blacklist/:url", h.add(false))
	// Removes a URL from the blacklist
	router.DELETE("/blacklist/:url", h.remove(false))

	// Gets the personal whitelist
	router.GET("/whitelist", h.list(true))
	// Adds a new URL to the whitelist
	router.PUT("/whitelist/:url", h.add(true))
	// Removes a URL from the whitelist
	router.DELETE("/whitelist/:url", h.remove(true))

	log.Fatal(http.ListenAndServe(host, router))
}
