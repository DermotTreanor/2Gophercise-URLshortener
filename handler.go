package urlshort

import (
	"net/http"
	"gopkg.in/yaml.v3"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	// Make a handler to call a redirect if our path is a key in the URLs map.
	// Otherwise we call the fallback.
	h := func(rw http.ResponseWriter, r *http.Request) {
		p, ok := pathsToUrls[r.URL.Path]
		if !ok {
			fallback.ServeHTTP(rw, r)
		} else {
			http.Redirect(rw, r, p, http.StatusSeeOther)
		}
	}
	return h
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...
	ym := map[string]string{}
	err := yaml.Unmarshal(yml, ym)
	return MapHandler(ym, fallback), err 
}


func JSONHandler(js []byte, fallback http.Handler)(http.HandlerFunc, error){
	mp := map[string]string{}
	return MapHandler(mp, fallback), nil
}