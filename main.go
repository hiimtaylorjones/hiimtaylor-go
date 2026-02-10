  package main

  import (
      "fmt"
      "log"
      "net/http"
  )

  func main() {
      http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
          fmt.Fprintf(w, "Hello, World!")
      })

      log.Println("Server starting on http://localhost:3000")
      log.Fatal(http.ListenAndServe(":3000", nil))
  }