package main

import (
	"flag"
	"fmt"
	"github.com/mewkiz/flac"
	"io"
	"os"
)

func main() {
	filepath := flag.String("p", "", "Path of the file to split raw")
	filename := flag.String("f", "", "File name to split raw")

	flag.Parse()

	song, err := flac.Open(fmt.Sprintf("%s%s.flac", *filepath, *filename))
	if err != nil {
		panic(fmt.Sprintf("Error Opening file %s", *filename))
	}
	defer song.Close()
	out := make([]*os.File, 2)
	out[0], err = os.OpenFile(fmt.Sprintf("%s%s_ch1.flac", *filepath, *filename),
		os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Error Opening output file %s%s_ch1.flac", *filepath, *filename))
	}
	out[1], err = os.OpenFile(fmt.Sprintf("%s%s_ch2.flac", *filepath, *filename),
		os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Error Opening output file %s%s_ch1.flac", *filepath, *filename))
	}

	for {
		frame, err := song.ParseNext()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(fmt.Sprintln(err))
		}

		for i, subframe := range frame.Subframes {
			for _, sample := range subframe.Samples {
				fmt.Fprintln(out[i], sample)
			}
		}
	}
}
