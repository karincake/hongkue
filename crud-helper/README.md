# Hongkue CRUD Helper
Extra library to help in creating CRUD APIs

## RegCrud()
A function to register CRUD APIs where it is actually a wrapper of `hongkue.GroupRoutes()` with CRUD feature.

## Types
A type `CrudBase` is an interface that defines the CRUD APIs

```go
type CrudBase interface {
	Create(w http.ResponseWriter, r *http.Request)
	ReadList(w http.ResponseWriter, r *http.Request)
	ReadSingle(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
```

It is being used as the last parameter OR the second last parameter of `RegCrud()`:.

If it is used as the second last parameter, then the last parameter will be the type of `MapHandlerFunc`,
therefore more routings can be added by using the `MapHandlerFunc`.

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

 	// the last parameter is the instance that implements the CrudBase interface
	hc.RegCrud("/article", r, O)
	
	// the secondlast parameter is the instance that implements the CrudBase interface
	// and the last parameter is the instance that implements the MapHandlerFunc
	hc.RegCrud("/another-article", r, O,
		h.MapHandlerFunc{
			"GET /{id}/extra": Extra,
		},
	) 
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

func (obj myBase) ReadSingle(w http.ResponseWriter, r *http.Request) {
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