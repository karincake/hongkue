package hongkue

import (
	"net/http"
	"reflect"
)

func checkMiddlewares(mwCandidates []any) []HandlerMw {
	// must have at least router
	wantedType1 := "func(http.Handler) http.Handler"
	wantedType2 := "hongkue.HandlerMw"
	mwList := []HandlerMw{}
	for i := range mwCandidates {
		myType := reflect.TypeOf(mwCandidates[i])
		myTypeStr := myType.String()
		if myTypeStr != wantedType1 && myTypeStr != wantedType2 {
			// check once more if it's array
			if myType.Kind() != reflect.Slice {
				panic("a non middleware is passed - first stage")
			} else {
				subMwCandidates := reflect.ValueOf(mwCandidates[i])
				for j := 0; j < subMwCandidates.Len(); j++ {
					mySubType := reflect.ValueOf(subMwCandidates.Index(j).Interface()).Type()
					mySubTypeStr := mySubType.String()
					if mySubTypeStr != wantedType1 && mySubTypeStr != wantedType2 {
						panic("a non middleware is passed - first stage")
					} else {
						mwList = append(mwList, subMwCandidates.Index(j).Interface().(HandlerMw))
					}
				}
			}
		} else {
			mwList = append(mwList, mwCandidates[i].(func(http.Handler) http.Handler))
		}
	}

	return mwList
}
