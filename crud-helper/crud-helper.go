package crudhelper

import (
	"maps"
	"net/http"
	"reflect"

	hk "github.com/karincake/hongkue"
)

func RegCrud(path string, r *http.ServeMux, mwAndRouter ...any) {
	sLength := len(mwAndRouter)
	mwCandidatesStart := 1
	extraHandlers := hk.MapHandlerFunc{}

	if sLength > 1 {
		if _, ok := mwAndRouter[sLength-2].(hk.MapHandlerFunc); ok {
			mwCandidatesStart = 2
			extraHandlers = mwAndRouter[sLength-2].(hk.MapHandlerFunc)
		}
	}

	mwCandidates := mwAndRouter[:sLength-mwCandidatesStart]
	mwList := []hk.HandlerMw{}
	for i := range mwCandidates {
		// have to do it manually, since casting directly results unexpected result
		myType := reflect.TypeOf(mwCandidates[i])
		if myType.String() != "func(http.Handler) http.Handler" {
			panic("non middleware included as middleware")
		}
		mwList = append(mwList, mwCandidates[i].(func(http.Handler) http.Handler))

		// if g, okHandler := mwCandidates[i].(func(http.Handler) http.Handler); !okHandler {
		// 	panic("non middleware included")
		// } else {
		// 	mwList = append(mwList, g)
		// }
	}

	c, ok := mwAndRouter[sLength-1].(CrudBase)
	if !ok {
		panic("non CrudBase used in the last paramter")
	}

	crudHandlers := hk.MapHandlerFunc{
		"POST /":       c.Create,
		"GET /":        c.ReadList,
		"GET /{id}":    c.ReadDetail,
		"PATCH /{id}":  c.Update,
		"DELETE /{id}": c.Delete,
	}
	maps.Copy(crudHandlers, extraHandlers)

	hk.GroupRoutes(path, r, mwList, crudHandlers)
}
