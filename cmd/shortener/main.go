package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/lovelydaemon/url-shortener/internal/rnd"
	"github.com/lovelydaemon/url-shortener/internal/validation"
)

var (
  storeURLToToken = make(map[string]string)
  storeTokenToURL = make(map[string]string)
)

func urlHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodPost {
    if r.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
      http.Error(w, "Bad content type", http.StatusBadRequest)
      return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    if validation.IsValidUrl(string(body)) != nil {
      http.Error(w, "Bad url link", http.StatusBadRequest)
      return
    }

    shortURL, ok := storeURLToToken[string(body)]
    if ok {
      w.Write([]byte(shortURL))
      return
    } else {
      token := rnd.NewRandomString(9)
      shortURL = fmt.Sprintf("http://%s/%s", r.Host, token)

      storeURLToToken[string(body)] = shortURL
      storeTokenToURL[token] = string(body) 
      w.WriteHeader(http.StatusCreated)
      w.Write([]byte(shortURL))
      return
    }

  } else if r.Method == http.MethodGet {
    token := strings.TrimLeft(r.URL.Path, "/")
    u, ok := storeTokenToURL[token]
    if !ok {
      http.Error(w, "Not found", http.StatusBadRequest)
      return
    }
    
    w.Header().Set("Location", u)
    w.WriteHeader(http.StatusTemporaryRedirect)
    

  } else {
    http.Error(w, "Bad Request", http.StatusBadRequest)
    return
  }
}



func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", urlHandler)

  if err := http.ListenAndServe(":3000", mux); err != nil {
    log.Fatal(err)
  }
}
