package crudhelper

import "net/http"

type CrudBase interface {
	Create(w http.ResponseWriter, r *http.Request)
	ReadList(w http.ResponseWriter, r *http.Request)
	ReadSingle(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
