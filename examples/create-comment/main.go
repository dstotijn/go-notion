package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/sanity-io/litter"
)

type httpTransport struct {
	w io.Writer
}

// RoundTrip implements http.RoundTripper. It multiplexes the read HTTP response
// data to an io.Writer for debugging.
func (t *httpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	res.Body = io.NopCloser(io.TeeReader(res.Body, t.w))

	return res, nil
}

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("NOTION_API_KEY")
	buf := &bytes.Buffer{}
	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &httpTransport{w: buf},
	}
	client := notion.NewClient(apiKey, notion.WithHTTPClient(httpClient))

	var parentPageID string

	flag.StringVar(&parentPageID, "parentPageId", "", "Parent page ID.")
	flag.Parse()

	params := notion.CreateCommentParams{
		ParentPageID: parentPageID,
		RichText: []notion.RichText{
			{
				Text: &notion.Text{
					Content: "This is an example comment.",
				},
			},
		},
	}
	comment, err := client.CreateComment(ctx, params)
	if err != nil {
		log.Fatalf("Failed to create comment: %v", err)
	}

	decoded := map[string]interface{}{}
	if err := json.NewDecoder(buf).Decode(&decoded); err != nil {
		log.Fatal(err)
	}

	// Pretty print JSON reponse.
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(decoded); err != nil {
		log.Fatal(err)
	}

	// Pretty print parsed `notion.Comment` value.
	litter.Dump(comment)
}
