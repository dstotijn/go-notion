package notion_test

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/google/go-cmp/cmp"
)

type mockRoundtripper struct {
	fn func(*http.Request) (*http.Response, error)
}

func (m *mockRoundtripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return m.fn(r)
}

func mustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("calls option funcs", func(t *testing.T) {
		t.Parallel()

		funcCalled := false
		opt := func(c *notion.Client) {
			funcCalled = true
		}

		_ = notion.NewClient("secret-api-key", opt)

		exp := true
		got := funcCalled

		if exp != got {
			t.Error("expected option func to be called.")
		}
	})

	t.Run("calls option func with client", func(t *testing.T) {
		t.Parallel()

		var clientArg *notion.Client
		opt := func(c *notion.Client) {
			clientArg = c
		}

		exp := notion.NewClient("secret-api-key", opt)
		got := clientArg

		if exp != got {
			t.Errorf("option func called with incorrect *Client value (expected: %+v, got: %+v)", exp, got)
		}
	})
}

func TestFindDatabaseByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expDatabase    notion.Database
		expError       error
	}{
		{
			name: "successful response",
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "database",
						"id": "668d797c-76fa-4934-9b05-ad288df2d136",
						"created_time": "2020-03-17T19:10:04.968Z",
						"last_edited_time": "2020-03-17T21:49:37.913Z",
						"title": [
							{
								"type": "text",
								"text": {
									"content": "Grocery List",
									"link": null
								},
								"annotations": {
									"bold": false,
									"italic": false,
									"strikethrough": false,
									"underline": false,
									"code": false,
									"color": "default"
								},
								"plain_text": "Grocery List",
								"href": null
							}
						],
						"properties": {
							"Name": {
								"id": "title",
								"type": "title",
								"title": {}
							},
							"Description": {
								"id": "J@cS",
								"type": "rich_text",
								"text": {}
							},
							"In stock": {
								"id": "{xYx",
								"type": "checkbox",
								"checkbox": {}
							},
							"Food group": {
								"id": "TJmr",
								"type": "select",
								"select": {
									"options": [
										{
											"id": "96eb622f-4b88-4283-919d-ece2fbed3841",
											"name": "🥦Vegetable",
											"color": "green"
										},
										{
											"id": "bb443819-81dc-46fb-882d-ebee6e22c432",
											"name": "🍎Fruit",
											"color": "red"
										},
										{
											"id": "7da9d1b9-8685-472e-9da3-3af57bdb221e",
											"name": "💪Protein",
											"color": "yellow"
										}
									]
								}
							},
							"Price": {
								"id": "cU^N",
								"type": "number",
								"number": {
									"format": "dollar"
								}
							},
							"Cost of next trip": {
								"id": "p:sC",
								"type": "formula",
								"formula": {
									"expression": "if(prop(\"In stock\"), 0, prop(\"Price\"))"
								}
							},
							"Last ordered": {
								"id": "]\\R[",
								"type": "date",
								"date": {}
							},
							"Meals": {
								"id": "lV]M",
								"type": "relation",
								"relation": {
									"database_id": "668d797c-76fa-4934-9b05-ad288df2d136",
									"synced_property_name": "Related to Test database (Relation Test)",
									"synced_property_id": "IJi<"
								}
							},
							"Number of meals": {
								"id": "Z\\Eh",
								"type": "rollup",
								"rollup": {
									"rollup_property_name": "Name",
									"relation_property_name": "Meals",
									"rollup_property_id": "title",
									"relation_property_id": "mxp^",
									"function": "count"
								}
							},
							"Store availability": {
								"id": "=_>D",
								"type": "multi_select",
								"multi_select": {
									"options": [
										{
											"id": "d209b920-212c-4040-9d4a-bdf349dd8b2a",
											"name": "Duc Loi Market",
											"color": "blue"
										},
										{
											"id": "6c3867c5-d542-4f84-b6e9-a420c43094e7",
											"name": "Gus's Community Market",
											"color": "yellow"
										}
									]
								}
							},
							"+1": {
								"id": "aGut",
								"type": "people",
								"people": {}
							},
							"Photo": {
								"id": "aTIT",
								"type": "files",
								"files": {}
							}
						}
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expDatabase: notion.Database{
				ID:             "668d797c-76fa-4934-9b05-ad288df2d136",
				CreatedTime:    mustParseTime(time.RFC3339, "2020-03-17T19:10:04.968Z"),
				LastEditedTime: mustParseTime(time.RFC3339, "2020-03-17T21:49:37.913Z"),
				Title: []notion.RichText{
					{
						Type: notion.RichTextTypeText,
						Text: &notion.Text{
							Content: "Grocery List",
						},
						Annotations: &notion.Annotations{
							Color: notion.ColorDefault,
						},
						PlainText: "Grocery List",
					},
				},
				Properties: notion.DatabaseProperties{
					"Name": notion.DatabaseProperty{
						ID:   "title",
						Type: notion.DBPropTypeTitle,
					},
					"Description": notion.DatabaseProperty{
						ID:   "J@cS",
						Type: notion.DBPropTypeRichText,
					},
					"In stock": notion.DatabaseProperty{
						ID:   "{xYx",
						Type: notion.DBPropTypeCheckbox,
					},
					"Food group": notion.DatabaseProperty{
						ID:   "TJmr",
						Type: notion.DBPropTypeSelect,
						Select: &notion.SelectMetadata{
							Options: []notion.SelectOptions{
								{
									ID:    "96eb622f-4b88-4283-919d-ece2fbed3841",
									Name:  "🥦Vegetable",
									Color: notion.ColorGreen,
								},
								{
									ID:    "bb443819-81dc-46fb-882d-ebee6e22c432",
									Name:  "🍎Fruit",
									Color: notion.ColorRed,
								},
								{
									ID:    "7da9d1b9-8685-472e-9da3-3af57bdb221e",
									Name:  "💪Protein",
									Color: notion.ColorYellow,
								},
							},
						},
					},
					"Price": notion.DatabaseProperty{
						ID:   "cU^N",
						Type: notion.DBPropTypeNumber,
						Number: &notion.NumberMetadata{
							Format: notion.NumberFormatDollar,
						},
					},
					"Cost of next trip": {
						ID:   "p:sC",
						Type: notion.DBPropTypeFormula,
						Formula: &notion.FormulaMetadata{
							Expression: `if(prop("In stock"), 0, prop("Price"))`,
						},
					},
					"Last ordered": notion.DatabaseProperty{
						ID:   "]\\R[",
						Type: notion.DBPropTypeDate,
					},
					"Meals": notion.DatabaseProperty{
						ID:   "lV]M",
						Type: notion.DBPropTypeRelation,
						Relation: &notion.RelationMetadata{
							DatabaseID:     "668d797c-76fa-4934-9b05-ad288df2d136",
							SyncedPropName: "Related to Test database (Relation Test)",
							SyncedPropID:   "IJi<",
						},
					},
					"Number of meals": notion.DatabaseProperty{
						ID:   "Z\\Eh",
						Type: notion.DBPropTypeRollup,
						Rollup: &notion.RollupMetadata{
							RollupPropName:   "Name",
							RelationPropName: "Meals",
							RollupPropID:     "title",
							RelationPropID:   "mxp^",
							Function:         "count",
						},
					},
					"Store availability": notion.DatabaseProperty{
						ID:   "=_>D",
						Type: notion.DBPropTypeMultiSelect,
						MultiSelect: &notion.SelectMetadata{
							Options: []notion.SelectOptions{
								{
									ID:    "d209b920-212c-4040-9d4a-bdf349dd8b2a",
									Name:  "Duc Loi Market",
									Color: notion.ColorBlue,
								},
								{
									ID:    "6c3867c5-d542-4f84-b6e9-a420c43094e7",
									Name:  "Gus's Community Market",
									Color: notion.ColorYellow,
								},
							},
						},
					},
					"+1": notion.DatabaseProperty{
						ID:   "aGut",
						Type: notion.DBPropTypePeople,
					},
					"Photo": {
						ID:   "aTIT",
						Type: "files",
					},
				},
			},
			expError: nil,
		},
		{
			name: "error response",
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "error",
						"status": 400,
						"code": "validation_error",
						"message": "foobar"
					}`,
				)
			},
			respStatusCode: http.StatusBadRequest,
			expDatabase:    notion.Database{},
			expError:       errors.New("notion: failed to find database: foobar (code: validation_error, status: 400)"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			httpClient := &http.Client{
				Transport: &mockRoundtripper{fn: func(r *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: tt.respStatusCode,
						Status:     http.StatusText(tt.respStatusCode),
						Body:       ioutil.NopCloser(tt.respBody(r)),
					}, nil
				}},
			}
			client := notion.NewClient("secret-api-key", notion.WithHTTPClient(httpClient))
			db, err := client.FindDatabaseByID(context.Background(), "00000000-0000-0000-0000-000000000000")

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expDatabase, db); diff != "" {
				t.Fatalf("database not equal (-exp, +got):\n%v", diff)
			}
		})
	}
}
