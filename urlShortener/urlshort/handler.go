package urlshort

import (
	"database/sql"
	"encoding/json"
	"github.com/go-yaml/yaml"
	"net/http"
)

type PathToUrl struct {
	Path string
	Url string
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(res http.ResponseWriter, req  *http.Request) {

		if url, ok := pathsToUrls[req.URL.Path] ; ok {
			http.Redirect(res, req, url, http.StatusPermanentRedirect)
		} else{
			  fallback.ServeHTTP(res,req)
		}
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap,fallback), nil

}

func parseYAML(yml []byte) ([]PathToUrl ,error){

	out := []PathToUrl{}
	if err := yaml.Unmarshal([]byte(yml), &out); err != nil {
		return []PathToUrl{}, err
	}

	return out,nil
}

func buildMap(parsedPathToUrlString []PathToUrl) map[string]string {

	pathMap := make(map[string]string)

	for _, value := range parsedPathToUrlString{
		pathMap[value.Path] = value.Url
	}
	return pathMap
}


func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJSON(json)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJson)
	return MapHandler(pathMap,fallback), nil
}

func parseJSON(jsn []byte) ([]PathToUrl , error){
	out := []PathToUrl{}
	err := json.Unmarshal(jsn, &out)
	if err != nil {
		return []PathToUrl{}, err
	}

	return out,nil
}

func DBHandler(db *sql.DB , tableName string, fallback http.Handler ) (http.HandlerFunc, error){
	parsedDB, err := parseDB(db,tableName)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedDB)
	return MapHandler(pathMap,fallback), nil

}

func parseDB(db *sql.DB, tableName string) ([]PathToUrl, error){
	out := []PathToUrl{}
	var (
		path string
		url string
	)
	rows, err := db.Query("select path, url from " + tableName)

	if err != nil {
		return []PathToUrl{}, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&path, &url)
		if err != nil {
			return []PathToUrl{}, err
		}
		out = append(out,PathToUrl{Path: path, Url: url})
	}

	err = rows.Err()
	if err != nil {
		return []PathToUrl{}, err
	}
	return out, nil
}