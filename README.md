# Hongkue: Helper for Standard Library net/http
A small go package that helps in managing routes for http router standard library. Needs Go version 1.22 or above

### StackMiddlewares()
A function to create stack of middleware for http.hanlder. It's actually a very common midleware stack maker.

Example
```go
func main() {
	router := http.NewServeMux()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world!!"))
	})

	mwList := hongkue.StackMiddlewares(myMiddleware1, myMiddleware2, ...)
    routerWithMw  := mwList(router)
}
```

### GroupRoutes()
A function to group some routes of http.hanlder. The handlers is passed in form of `map[string]func(http.Handler) http.Handler`

Example
```go
func main() {
	r := http.NewServeMux()

	r.HandleFunc("/", home.Home)
	r.HandleFunc("/hello-world", hw.HelloWorld)

	// without middleware
	rc.Group("/article", r, map[string]func(http.Handler) http.Handler{
		"GET /": someHttpHandlerFunc,
		"GET /{id}": someHttpHandlerFunc,
	})

	// With middlewares. The count of the middlewares is unilimited, just make sure the last parameter is for the handlers
	rc.Group("/article", r,myMiddleware1, myMiddleware2, map[string]func(http.Handler) http.Handler{
		"GET /": someHttpHandlerFunc,
		"GET /{id}": someHttpHandlerFunc,
	})
}
```


### GroupUsingHandler()
A function to group some routes of http.hanlder. This uses built-in function http.Handle(), combined with http.StripPrefix(), therefore passing the handlers can be done in a conventional way. 
The used method has a downside: it adds trailing slash for the main subroute path.

Example
```go
func main() {
	r := http.NewServeMux()

	r.HandleFunc("/", home.Home)
	r.HandleFunc("/hello-world", hw.HelloWorld)

	// without middleware
	rc.Group("/comment", r, func(r *http.ServeMux) {
		r.HandleFunc("GET /", someHttpHandlerFunc)
		r.HandleFunc("GET /{id}", someHttpHandlerFunc)
	})

	// With middlewares. The count of the middlewares is unilimited, just make sure the last parameter is for the handlers
	rc.Group("/comment", r, myMiddleware1, myMiddleware2, func(r *http.ServeMux) {
		r.HandleFunc("GET /", someHttpHandlerFunc)
		r.HandleFunc("GET /{id}", someHttpHandlerFunc)
	})
}
```
