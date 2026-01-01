# Hongkue CRUD Helper
Extra library to help in creating CRUD APIs

## RegCrud()
A function to register CRUD APIs. It's a wrapper of `hongkue.GroupRoutes()` with some extra features.

## Types
A type `CrudBase` is defined to be used as the last parameter of `RegCrud()`. It's a interface that defines the CRUD APIs.

```go
type CrudBase interface {
	Create(w http.ResponseWriter, r *http.Request)
	ReadList(w http.ResponseWriter, r *http.Request)
	ReadDetail(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
```

## Example
```go
import (
	"net/http"
	hc "github.com/karincake/hongkue/crud-helper"
)	

func main() {
	r := http.NewServeMux()

	r.HandleFunc("/", home.Home)
	r.HandleFunc("/hello-world", hw.HelloWorld)

	hc.RegCrud("/article-mw", r,
		h.MapHandlerFunc{
			"GET /{id}/extra": Extra,
		},
		O) // the last parameter is the instance that implements the CrudBase interface
}

// the crud
type myBase struct{}

var O myBase

func (obj myBase) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You are accessing path for creating comment"))
}

func (obj myBase) ReadList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You are accessing path for comment list"))
}

func (obj myBase) ReadDetail(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You accessing path for comment detail with id: " + r.PathValue("id")))
}

func (obj myBase) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You accessing path for comment updating with id: " + r.PathValue("id")))
}

func (obj myBase) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You accessing path for comment deletion with id: " + r.PathValue("id")))
}

func Extra(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You accessing path for comment extra with id: " + r.PathValue("id")))
}


```