package service

import (
	"bytes"
	"encoding/binary"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"

	"github.com/mewkiz/flac"
)

// NewServer returns a new microservice server for music
func NewServer() *negroni.Negroni {

	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter)
	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/test", getMHandler(formatter)).Methods("GET")
}

func getMHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK,
			struct{ Test string }{"This is a test"})
	}
}

// GetSubFrame gets the subframe at index frameNo from song (FLAC)
func GetSubFrame(song flac.Stream, frameNo int) {
	seek(frameNo)
	frame = song.ParseNext()

}

func seek(song flac.Stream, offset int) error {
	for i := 0; i < frameNo-1; i++ {
		_, err := song.Next()
		if err != nil {
			return error.Error("Frame Index Out Of Bounds")
		}
	}
	return nil
}
