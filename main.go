package main

import (
    "log"
    "net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    log.Println(path)
    if path != "/" {
        http.NotFound(w, r)
        return
    }
    w.Write([]byte("Hello :)"))
}

func viewSnippet(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("you're viewing a snippet"))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.Header().Set("Allow", http.MethodPost)
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    w.Write([]byte("you're creating a new snippet"))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", home)
    mux.HandleFunc("/snippet/view", viewSnippet)
    mux.HandleFunc("/snippet/create", createSnippet)

    log.Println("starting server on port 4000")
    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}
