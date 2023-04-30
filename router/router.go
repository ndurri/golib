package router

import "net/http"

type Handler func(w http.ResponseWriter, r *http.Request)

type Router struct {
	List   Handler
	Get    Handler
	Post   Handler
	Patch  Handler
	Delete Handler
}

func (rtr Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/" && rtr.List != nil:
		rtr.List(w, r)
	case r.Method == http.MethodGet && rtr.Get != nil:
		rtr.Get(w, r)
	case r.Method == http.MethodPost && rtr.Post != nil:
		rtr.Post(w, r)
	case r.Method == http.MethodDelete && rtr.Delete != nil:
		rtr.Delete(w, r)
	case r.Method == http.MethodPatch && rtr.Patch != nil:
		rtr.Patch(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
