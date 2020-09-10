package urlshort_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"urlShortener/urlshort"
)

type pathToUrl struct {
	path string
	url string
	body string
	statusCode int
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}
func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Hello, world!")
}

func TestMapHandler(t *testing.T) {

	pathToUrlMap := map[string]string{
		"/urlshort" :  "https://github.com/gophercises/urlshort",
		"/urlshort-final" : "https://github.com/gophercises/urlshort/tree/solution",
	}

	handler := urlshort.MapHandler(pathToUrlMap,defaultMux())


	testCases := []pathToUrl{
		{   path: "/urlshort",
		    url: "https://github.com/gophercises/urlshort",
		    body: fmt.Sprintf("<a href=\"%s\">Permanent Redirect</a>.","https://github.com/gophercises/urlshort"),
           statusCode: 308,
		},
		{   path: "/urlshort-final",
			url: "https://github.com/gophercises/urlshort/tree/solution",
			body: fmt.Sprintf("<a href=\"%s\">Permanent Redirect</a>.","https://github.com/gophercises/urlshort/tree/solution"),
			statusCode: 308,
		},
		{   path: "/",
			url: "https://localhost",
			body: "Hello, world!",
			statusCode: 200,
		},
		{   path: "/abc",
			url: "https://localhost",
			body: "Hello, world!",
			statusCode: 200,
		},
	}

	for _,val :=range testCases{
		req, err := http.NewRequest("GET", val.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if !(val.statusCode == rr.Code && val.body == strings.TrimSpace(rr.Body.String())) {
			t.Errorf("status code or body are not eqaul, failed to redirect")
		}
	}
}


func TestYAMLHandler(t *testing.T) {

	pathToUrlYaml := `- path: /urlshort
url: https://github.com/gophercises/urlshort
- path: /urlshort-final`

	handler, err := urlshort.YAMLHandler([]byte(pathToUrlYaml),defaultMux())
    if err == nil {
    	t.Error("failed ")
	}

	pathToUrlYaml = `- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution`

	handler,err = urlshort.YAMLHandler([]byte(pathToUrlYaml),defaultMux())
	if err != nil{
		t.Error("failed ")
	}


	testCases := []pathToUrl{
		{   path: "/urlshort",
			url: "https://github.com/gophercises/urlshort",
			body: fmt.Sprintf("<a href=\"%s\">Permanent Redirect</a>.","https://github.com/gophercises/urlshort"),
			statusCode: 308,
		},
		{   path: "/urlshort-final",
			url: "https://github.com/gophercises/urlshort/tree/solution",
			body: fmt.Sprintf("<a href=\"%s\">Permanent Redirect</a>.","https://github.com/gophercises/urlshort/tree/solution"),
			statusCode: 308,
		},
		{   path: "/",
			url: "https://localhost",
			body: "Hello, world!",
			statusCode: 200,
		},
		{   path: "/abc",
			url: "https://localhost",
			body: "Hello, world!",
			statusCode: 200,
		},
	}

	for _,val :=range testCases{
		req, err := http.NewRequest("GET", val.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if !(val.statusCode == rr.Code && val.body == strings.TrimSpace(rr.Body.String())) {
			t.Errorf("status code or body are not eqaul, failed to redirect")
		}
	}

}

func TestJSONHandler(t *testing.T) {

	pathToUrlJson :=`[
	{
		"path": "/urlshort",
		"url": "https://github.com/gophercises/urlshort"
	},
	{
		"path": "/urlshort-final",
	}
]`

	handler, err := urlshort.JSONHandler([]byte(pathToUrlJson),defaultMux())
	if err == nil {
		t.Error("failed ")
	}

	pathToUrlJson =`[
	{
		"path": "/urlshort",
		"url": "https://github.com/gophercises/urlshort"
	},
	{
		"path": "/urlshort-final",
		"url": "https://github.com/gophercises/urlshort/tree/solution"
	}
]`

	handler, err = urlshort.JSONHandler([]byte(pathToUrlJson),defaultMux())

	if err != nil{
		t.Error("failed ")
	}

	testCases := []pathToUrl{
		{   path: "/urlshort",
			url: "https://github.com/gophercises/urlshort",
			body: fmt.Sprintf("<a href=\"%s\">Permanent Redirect</a>.","https://github.com/gophercises/urlshort"),
			statusCode: 308,
		},
		{   path: "/urlshort-final",
			url: "https://github.com/gophercises/urlshort/tree/solution",
			body: fmt.Sprintf("<a href=\"%s\">Permanent Redirect</a>.","https://github.com/gophercises/urlshort/tree/solution"),
			statusCode: 308,
		},
		{   path: "/",
			url: "https://localhost",
			body: "Hello, world!",
			statusCode: 200,
		},
		{   path: "/abc",
			url: "https://localhost",
			body: "Hello, world!",
			statusCode: 200,
		},
	}

	for _,val :=range testCases{
		req, err := http.NewRequest("GET", val.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if !(val.statusCode == rr.Code && val.body == strings.TrimSpace(rr.Body.String())) {
			t.Errorf("status code or body are not eqaul, failed to redirect")
		}
	}

}


