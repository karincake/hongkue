package main

import (
	"log"
	"net/http"

	"github.com/karincake/hongkue"
)

func main() {
	r := http.NewServeMux()

	r.HandleFunc("/", home)

	hongkue.GroupRoutesUsingHandler("/sub-route-3", r, func(r *http.ServeMux) {
		r.HandleFunc("GET /", subRoute)
		r.HandleFunc("GET /hello-buddy", helloBuddy)
		r.HandleFunc("GET /hello-someone/{name}", helloSomeone)
	})

	hongkue.GroupRoutesUsingHandler("/sub-route-4", r, myMiddleware1, myMiddleware2, func(r *http.ServeMux) {
		r.HandleFunc("GET /", subRoute)
		r.HandleFunc("GET /hello-buddy", helloBuddy)
		r.HandleFunc("GET /hello-someone/{name}", helloSomeone)
	})

	mm := hongkue.StackMiddlewares(myMiddleware0)
	addr := "127.0.0.1:8000"
	srv := &http.Server{
		Addr:    addr,
		Handler: mm(r),
	}

	log.Print("listeing ", addr)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err.Error())
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
		return
	}
	w.Write([]byte("Hello world!!"))
}

func subRoute(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You are accessing a subroute"))
}

func helloBuddy(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello buddy!!"))
}

func helloSomeone(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello " + r.PathValue("name")))
}

func myMiddleware0(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Middleware-0 is called")
		next.ServeHTTP(w, r)
	})
}

func myMiddleware1(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Middleware-1 is called")
		next.ServeHTTP(w, r)
	})
}

func myMiddleware2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Middleware-2 is called")
		next.ServeHTTP(w, r)
	})
}
