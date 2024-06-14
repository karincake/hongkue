package hongkue

import (
	"fmt"
	"net/http"
	"strings"
)

// type HandlerMw func(http.Handler) http.Handler
type MapHandlerFunc map[string]func(http.ResponseWriter, *http.Request)

// Create stack of middleware to be used at once
// Note that this is very common practice
func StackMiddlewares(mws ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			next = mws[i](next)
		}
		return next
	}
}

// Create subroute of a path
// WARNING: this is work around to avoid redirect behavior regarding
// trailing slash from http.StripPrefix
func GroupRoutes(path string, router *http.ServeMux, mwAndRouter ...any) {
	// must have at least router
	sLength := len(mwAndRouter)
	if sLength == 0 {
		panic("requires third parameter for the subrouter")
	}
	routes, okMux := mwAndRouter[sLength-1].(MapHandlerFunc)
	if !okMux {
		panic("last parameter should satistifes `map[string]func(http.ResponseWriter, *http.Request)`")
	}

	// check middleware availability
	mwCandidates := mwAndRouter[:sLength-1]
	mwList := []func(http.Handler) http.Handler{}
	for i := range mwCandidates {
		if g, okHandler := mwCandidates[i].(func(http.Handler) http.Handler); !okHandler {
			panic("non middleware included")
		} else {
			mwList = append(mwList, g)
		}
	}
	mStack := StackMiddlewares(mwList...)

	//
	subRouterStatus := false
	subRouter := http.NewServeMux()
	for key := range routes {
		keys := strings.Fields(key)
		if len(keys) == 1 {
			if key == "/" {
				tempMux := http.NewServeMux()
				tempMux.HandleFunc("/", routes[key])
				router.Handle(path, mStack(tempMux))
			} else {
				subRouterStatus = true
				subRouter.HandleFunc(key, routes[key])
			}
		} else {
			if keys[1] == "/" {
				tempMux := http.NewServeMux()
				tempMux.HandleFunc(keys[0]+" /", routes[key])
				router.Handle(path, mStack(tempMux))
			} else {
				subRouterStatus = true
				subRouter.HandleFunc(keys[0]+" "+keys[1], routes[key])
			}
		}
	}
	if subRouterStatus {
		fmt.Println("new", path+"/", path)
		router.Handle(path+"/", http.StripPrefix(path, mStack(subRouter)))
	}
}

// Create subroute of a path
// The mwAndRouter, according to its name, acts as middlewares and router.
// The middlewares are optional, the router is mandataory for it will be
// the one to create the subroutes.
func GroupRoutesUsingHandler(path string, router *http.ServeMux, mwAndRouter ...any) {
	// must have at least router
	sLength := len(mwAndRouter)
	if sLength == 0 {
		panic("requires third parameter for the subrouter")
	}
	subRouterCall, okMux := mwAndRouter[sLength-1].(func(*http.ServeMux))
	if !okMux {
		panic("last parameter should satistifes `func(*http.ServeMux)`")
	}

	// check middleware
	mwCandidates := mwAndRouter[:sLength-1]
	mwList := []func(http.Handler) http.Handler{}
	for i := range mwCandidates {
		if g, okHandler := mwCandidates[i].(func(http.Handler) http.Handler); !okHandler {
			panic("non middleware included")
		} else {
			mwList = append(mwList, g)
		}
	}
	mStack := StackMiddlewares(mwList...)

	// execute sub rouiter then use the middleware
	subRouter := http.NewServeMux()
	subRouterCall(subRouter)
	fmt.Println("convent", path+"/", path)
	router.Handle(path+"/", http.StripPrefix(path, mStack(subRouter)))
}
