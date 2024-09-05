package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

const apiURL = "https://foodish-api.com/api"

type ImageResponse struct {
    ImageURL string `json:"image"`
}

func fetchImage() (ImageResponse, error) {
    var response ImageResponse
    resp, err := http.Get(apiURL)
    if err != nil {
        return response, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return response, fmt.Errorf("API request failed with status: %s", resp.Status)
    }

    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return response, err
    }

    return response, nil
}

func main() {
    http.HandleFunc("/images", func(w http.ResponseWriter, r *http.Request) {
        imageResponse, err := fetchImage()
        if err != nil {
            http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
            return
        }

        // Set the content type to HTML
        w.Header().Set("Content-Type", "text/html")

        // Write the HTML header
        fmt.Fprint(w, "<html><body>")

        // Display the image
        fmt.Fprintf(w, `<img src="%s" alt="Image" style="max-width: 100%%; height: auto;"><br>`, imageResponse.ImageURL)

        // Write the HTML footer
        fmt.Fprint(w, "</body></html>")
    })

    log.Println("Server started on :3000...")
    if err := http.ListenAndServe(":3000", nil); err != nil {
        log.Fatal(err)
    }
}