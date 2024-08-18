package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"urlshort"
)

var file = flag.String("f", "", "pass in the name of a file with yaml data")

func main() {
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/scrib":          "https://youtube.com",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// 	// Build the YAMLHandler using the mapHandler as the
	// 	// fallback
	yaml := `
/urlshort: https://github.com/gophercises/urlshort
/urlshort-final: https://github.com/gophercises/urlshort/tree/solution
/scribble: https://reddit.com
`
	// If a file argument has been passed we will read the data
	path := *file
	if len(path) > 0 {
		data, err := openYAML(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "main: could not read from file. Using internal data.\n%v\n", err)
		} else {
			yaml = yaml + string(data)
		}
	}
	fmt.Printf("Here is the yaml data:\n%v\n", string(yaml))

	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	jd := `{
"/key1": "https://www.xkcd.com",
    "/key2": "https://www.microsoft.com",
    "/key3": "https://www.isitchristmas.com"
}`
	jsonHandler, err := urlshort.JSONHandler([]byte(jd), yamlHandler)
	if err != nil {
		fmt.Printf("main: encountered error when making json handler: %v\n", err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func openYAML(path string) (data []byte, err error) {
	fp, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "openYAML: could not open file:\n%v\n\nWorking with internal data only...\n", err)
		return data, err
	}
	info, err := fp.Stat()
	if err != nil {
		s := fmt.Sprintf("openYAML: cannot obtain file info: %v\n", err)
		panic(s)
	}
	if info.Size() < 150 {
		fmt.Fprintln(os.Stderr, "The file is small enough to pull in as a whole")
		data, err = io.ReadAll(fp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "openYAML: could not read in data: %v\n", err)
			return []byte{}, err
		} else {
			return data, nil
		}
	} else {
		fmt.Fprintln(os.Stderr, "Using pointless buffer...")
		buf := make([]byte, 50)
		for {
			n, err := fp.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "openYAML: error while trying to buffer file data: %v\n", err)
				return []byte{}, err
			}
			data = append(data, buf[:n]...)
			if err == io.EOF {
				return data, nil
			}
		}
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
