package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
		"net/url"
		"strings"
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

func getDirectoryName(urlStr string) (string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
			return "", err
	}
	
	// Extract the path part of the URL
	path := parsedURL.Path
	
	// Split the path into segments
	segments := strings.Split(path, "/")
	
	// Check if there are enough segments to extract the directory name
	if len(segments) < 3 {
			return "", fmt.Errorf("not enough segments in path")
	}
	
	// Return the segment before the last one
	return segments[len(segments)-2], nil
}

func main() {
    http.HandleFunc("/food", func(w http.ResponseWriter, r *http.Request) {
        imageResponse, err := fetchImage()
        if err != nil {
            http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
            return
        }

				// Extract directory name from the image URL
        directoryName, err := getDirectoryName(imageResponse.ImageURL)
        if err != nil {
            http.Error(w, "Failed to process URL", http.StatusInternalServerError)
            return
        }

				// Create a Google search URL
        searchURL := fmt.Sprintf("https://www.google.com/search?q=%s", url.QueryEscape(directoryName))

				// Print the extracted directory name for debugging
        log.Println("Directory Name:", directoryName)

        // Set the content type to HTML
        w.Header().Set("Content-Type", "text/html")

				// Write the HTML content
				fmt.Fprintf(w, `
        <html>
        <head>
            <title>Food Randomizer</title>
            <style>
                body {
                    background-color: #121212;
                    color: #e0e0e0;
                    font-family: Arial, sans-serif;
                    margin: 0;
                    padding: 0;
                    height: 100vh; 
                    display: flex;
                    flex-direction: column; 
                }
                header {
                    position: fixed;
                    top: 0;
                    width: 100%%;
                    background-color: #333;
                    color: #fff;
                    padding: 10px;
                    text-align: center;
                    z-index: 1000; 
                }
                .main-content {
                    flex: 1;
                    display: flex;
                    justify-content: center; 
                    margin-top: 60px;        
                    padding: 20px;
                }
                .content-wrapper {
                    text-align: center; 
                }
                h1 {
                    color: #ffffff;
                }
                a {
                    color: #bb86fc;
                    text-decoration: none;
                }
                a:hover {
                    text-decoration: underline;
                }
                img {
                    max-width: 50%%;
                    height: auto;   
                    border-radius: 8px;
                }
                /* Dark mode media query */
                @media (prefers-color-scheme: dark) {
                    body {
                        background-color: #121212;
                        color: #e0e0e0;
                    }
                }
            </style>
        </head>
        <body>
            <header>
                <h1>Food Randomizer</h1> 
            </header>
            <div class="main-content">
                <div class="content-wrapper">
                    <h1><a href="%s" target="_blank">%s</a></h1>
                    <img src="%s" alt="Image">
                </div>
            </div>
        </body>
        </html>
        `, searchURL, directoryName, imageResponse.ImageURL)
    })


    log.Println("Server started on :3000...")
    if err := http.ListenAndServe(":3000", nil); err != nil {
        log.Fatal(err)
    }
}