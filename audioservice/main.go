package main

import (
	// service "github.com/zenathark/astream/audioservice/server"
	// "os"
	"fmt"
	"github.com/mewkiz/flac"
	"io"
)

func main() {

	// port := os.Getenv("PORT")
	// if len(port) == 0 {
	//	port = "3000"
	// }

	// server := service.NewServer()
	// server.Run(":" + port)
	// song, err := flac.Open("../database/m1.flac")
	// if err != nil {
	//	panic(fmt.Sprintf("Error Opening file"))
	// }

	// defer song.Close()

	// fmt.Printf("Song info %v\n", song.Info)

	// for {
	//	frame, err := song.ParseNext()
	//	if err != nil {
	//		if err == io.EOF {
	//			break
	//		}
	//		panic(fmt.Sprintln(err))
	//	}

	//	for i, subframe := range frame.Subframes {
	//		fmt.Printf("subframe %d sample size %d\n", i, subframe.NSamples)
	//	}
	// }

}
