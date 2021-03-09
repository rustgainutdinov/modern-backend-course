package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	filePath := os.Args[1]
	paths := getPathsFromJson(filePath)
	connectRedirects(paths)
	err := http.ListenAndServe(":8091", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func getPathsFromJson(filePath string) map[string]interface{} {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		fmt.Println(err)
	}
	return result["paths"].(map[string]interface{})
}

func connectRedirects(paths map[string]interface{}) {
	for short, long := range paths {
		http.HandleFunc(short, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, long.(string), 301)
		})
	}
	http.HandleFunc("/", handler404)
}

func handler404(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "ERRROR PANIC AAAAA")
}
