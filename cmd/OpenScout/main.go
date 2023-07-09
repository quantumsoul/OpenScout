package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nsf/termbox-go"
)

type results struct {
	TotalCount        int      `json:"total_count"`
	IncompleteResults bool     `json:"incomplete_results"`
	Items             []result `json:"items"`
}

type result struct {
	HTMLURL    string `json:"html_url"`
	Repository string `json:"repository_url"`
	Title      string `json:"title"`
}

func main() {
	// name := flag.String("name", "Guest", "Specify a name")
	// flag.Parse()
	// fmt.Printf("Hello, %s!\n", *name)
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	currentPage := 1
	results := fetchResults(currentPage)

	renderResults(results)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				return
			case termbox.KeyArrowRight, termbox.KeyCtrlN:
				currentPage = currentPage + 1
				results = fetchResults(currentPage)
				renderResults(results)
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				if currentPage > 1 {
					currentPage = currentPage - 1
					results = fetchResults(currentPage)
					renderResults(results)
				}
			}
		}
	}
}

func fetchResults(page int) []result {
	// Fetch the results from the API for the given page
	// Return the results as a slice or data structure
	url := `https://api.github.com/search/issues?q=label:%22good%20first%20issue%22+state:open+language:go`

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	q := req.URL.Query()
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("per_page", fmt.Sprintf("%d", 30))
	req.URL.RawQuery = q.Encode()

	// Set the User-Agent header
	req.Header.Set("User-Agent", "OpenScout")
	req.Header.Set("Accept", "*/*")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")

	// Send the HTTP request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}

	// Parse the response into  the IssueResponse struct
	var results results
	err = json.Unmarshal(body, &results)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	// fmt.Print(resp.StatusCode)
	return results.Items
}

func renderResults(results []result) {
	// Clear the terminal screen
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	border := strings.Repeat("*", 160)

	fmt.Println(border)

	// Render the results in the terminal UI
	// Display the results in the desired format
	for i, result := range results {
		// Render each result item
		renderResultItem(i, result)
	}

	termbox.Flush()
	fmt.Println("\n")
}

func renderResultItem(index int, result result) {
	// Calculate the row position of the result item based on the index
	row := index + 1

	// Define the column positions for different parts of the result item
	columnRepoName := 1
	columnIssueTitle := 30
	columnIssueLink := 130

	str := result.Repository
	splitStr := strings.Split(str, "/")

	// Render the repository name
	renderText(row, columnRepoName, splitStr[5])

	// Render the issue title
	renderText(row, columnIssueTitle, result.Title)

	// Render the issue link
	renderLink(row, columnIssueLink, result.HTMLURL)

}

func renderText(row, column int, text string) {
	for i, char := range text {
		termbox.SetCell(column+i, row, char, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func renderLink(row, column int, link string) {
	// You can render the link as a highlighted text or underlined text
	// Use termbox functions like SetCell to render each character of the link
	// Apply appropriate styling to make it visually distinct from regular text
	for i, char := range link {
		termbox.SetCell(column+i, row, char, termbox.ColorDefault, termbox.ColorDefault)
	}
}
