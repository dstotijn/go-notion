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

	params := notion.CreatePageParams{
		ParentType: notion.ParentTypePage,
		ParentID:   parentPageID,
		Title: []notion.RichText{
			{
				Text: &notion.Text{
					Content: "Create Page Example",
				},
			},
		},
		Children: []notion.Block{
			notion.Heading1Block{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Heading 1",
						},
					},
				},
			},
			notion.Heading2Block{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Heading 2",
						},
					},
				},
			},
			notion.Heading3Block{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Heading 3",
						},
					},
				},
			},
			notion.ParagraphBlock{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "This is a paragraph.",
						},
					},
				},
			},
			notion.CalloutBlock{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "This is a callout.",
						},
					},
				},
				Icon: &notion.Icon{
					Type:  notion.IconTypeEmoji,
					Emoji: notion.StringPtr("💁"),
				},
			},
			notion.QuoteBlock{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: `"Assumption is the mother of all fuck-ups."`,
						},
					},
				},
			},
			notion.BulletedListItemBlock{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Bullet list item",
						},
					},
				},
			},
			notion.NumberedListItemBlock{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Numbered list item",
						},
					},
				},
			},
			notion.ToDoBlock{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Include to do item",
						},
					},
				},
				Checked: notion.BoolPtr(true),
			},
			notion.ToggleBlock{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Toggle",
						},
					},
				},
				Children: []notion.Block{
					notion.ParagraphBlock{
						RichText: []notion.RichText{
							{
								Text: &notion.Text{
									Content: "Toggled content.",
								},
							},
						},
					},
				},
			},
			notion.CodeBlock{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: `fmt.Println("Hello, world!)`,
						},
					},
				},
				Language: notion.StringPtr("go"),
				Caption: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Print `Hello, world!` to standard output.",
						},
					},
				},
			},
			notion.EmbedBlock{
				URL: "https://www.youtube.com/watch?v=8BETOsW4Y8g",
			},
			notion.ImageBlock{
				Type: notion.FileTypeExternal,
				External: &notion.FileExternal{
					URL: "https://picsum.photos/600/200.jpg",
				},
			},
			notion.VideoBlock{
				Type: notion.FileTypeExternal,
				External: &notion.FileExternal{
					URL: "https://download.samplelib.com/mp4/sample-5s.mp4",
				},
			},
			notion.FileBlock{
				Type: notion.FileTypeExternal,
				External: &notion.FileExternal{
					URL: "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
				},
				Caption: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Example file.",
						},
					},
				},
			},
			notion.PDFBlock{
				Type: notion.FileTypeExternal,
				External: &notion.FileExternal{
					URL: "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
				},
				Caption: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Example PDF file.",
						},
					},
				},
			},
			notion.BookmarkBlock{
				URL: "https://v0x.nl",
				Caption: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Homepage of David Stotijn.",
						},
					},
				},
			},
			notion.EquationBlock{
				Expression: "e=mc^2",
			},
			notion.DividerBlock{},
			notion.TableOfContentsBlock{},
			notion.BreadcrumbBlock{},
			notion.ColumnListBlock{
				Children: []notion.ColumnBlock{
					{
						Children: []notion.Block{
							notion.ParagraphBlock{
								RichText: []notion.RichText{
									{
										Text: &notion.Text{
											Content: "Column One",
										},
									},
								},
							},
						},
					},
					{
						Children: []notion.Block{
							notion.ParagraphBlock{
								RichText: []notion.RichText{
									{
										Text: &notion.Text{
											Content: "Column One",
										},
									},
								},
							},
						},
					},
				},
			},
			notion.TemplateBlock{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Create callout template.",
						},
					},
				},
				Children: []notion.Block{
					notion.CalloutBlock{
						RichText: []notion.RichText{
							{
								Text: &notion.Text{
									Content: "Placeholder callout text.",
								},
							},
						},
					},
				},
			},
			notion.SyncedBlock{
				SyncedFrom: nil,
				Children: []notion.Block{
					notion.CalloutBlock{
						RichText: []notion.RichText{
							{
								Text: &notion.Text{
									Content: "Callout in original synced block.",
								},
							},
						},
					},
				},
			},
			notion.TableBlock{
				TableWidth:      1,
				HasColumnHeader: true,
				Children: []notion.Block{
					notion.TableRowBlock{
						Cells: [][]notion.RichText{
							{
								{
									Text: &notion.Text{
										Content: "Column 1",
									},
								},
							},
						},
					},
					notion.TableRowBlock{
						Cells: [][]notion.RichText{
							{
								{
									Text: &notion.Text{
										Content: "Column 1 content.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	page, err := client.CreatePage(ctx, params)
	if err != nil {
		log.Fatalf("Failed to create page: %v", err)
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

	// Pretty print parsed `notion.Page` value.
	litter.Dump(page)
}
