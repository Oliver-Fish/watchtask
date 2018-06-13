package main

import "fmt"
import "net/http"

func main() {
	fmt.Println("Starting Webapp on port 8080")
	http.HandleFunc("/", getAll)
	http.ListenAndServe(":8080", nil)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request Made")
	w.Write([]byte(fmt.Sprint("Some Page Text")))

}
