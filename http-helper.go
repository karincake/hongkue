package hongkue

import (
	"net/http"
	"reflect"
	"strings"
)

// A func type of `func(http.Handler) http.Handler`
type HandlerMw func(http.Handler) http.Handler

// A func type of `map[string]func(http.ResponseWriter, *http.Request)`
type MapHandlerFunc map[string]func(http.ResponseWriter, *http.Request)

// Create stack of middleware to be used at once
// Note that this is very common practice
func StackMiddlewares(mws ...HandlerMw) HandlerMw {
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
		panic("last parameter should satisfies `map[string]func(http.ResponseWriter, *http.Request)`")
	}

	// check middleware availability
	mwCandidates := mwAndRouter[:sLength-1]
	mwList := []HandlerMw{}
	for i := range mwCandidates {
		if handler, okHandler := mwCandidates[i].(func(http.Handler) http.Handler); !okHandler {
			// check once more if it's array
			rt := reflect.TypeOf(mwCandidates[i]).Kind()
			switch rt {
			case reflect.Slice:
				s := reflect.ValueOf(mwCandidates[i])
				for i := 0; i < s.Len(); i++ {
					if handlerAgain, okHandler := mwCandidates[i].(func(http.Handler) http.Handler); !okHandler {
						panic("non middleware included")
					} else {
						mwList = append(mwList, handlerAgain)
					}
				}
			default:
				panic("non middleware included")
			}
		} else {
			mwList = append(mwList, handler)
		}
	}
	mStack := StackMiddlewares(mwList...)

	// process routing
	mainPathMethod := map[string]interface{}{} // to record the registered method
	subRouterStatus := false
	subRouter := http.NewServeMux()
	for key := range routes {
		keys := strings.Fields(key)
		if len(keys) == 1 {
			if key == "/" {
				mainPathMethod["ALL"] = nil // not specified method means any method is allowed
				tempMux := http.NewServeMux()
				tempMux.HandleFunc("/", routes[key])
				router.Handle(path, mStack(tempMux))
			} else {
				subRouterStatus = true
				subRouter.HandleFunc(key, routes[key])
			}
		} else {
			if keys[1] == "/" {
				mainPathMethod[keys[0]] = nil // specific method
				tempMux := http.NewServeMux()
				tempMux.HandleFunc(keys[0]+" /", routes[key])
				router.Handle(keys[0]+" "+path, mStack(tempMux))
			} else {
				subRouterStatus = true
				subRouter.HandleFunc(keys[0]+" "+keys[1], routes[key])
			}
		}
	}

	// if no unspecified method registered, means not all method can be processed
	if _, ok := mainPathMethod["ALL"]; !ok {
		router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if _, ok := mainPathMethod[r.Method]; !ok {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("Method not allowed"))
			}
		})
	}

	// if subrouter has other than main path, then process
	if subRouterStatus {
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
	mwList := []HandlerMw{}
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
	router.Handle(path+"/", http.StripPrefix(path, mStack(subRouter)))
}
