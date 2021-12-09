package notion_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
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
						},
						    "parent": {
								"type": "page_id",
								"page_id": "b8595b75-abd1-4cad-8dfe-f935a8ef57cb"
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
						ID:    "title",
						Type:  notion.DBPropTypeTitle,
						Title: &notion.EmptyMetadata{},
					},
					"Description": notion.DatabaseProperty{
						ID:   "J@cS",
						Type: notion.DBPropTypeRichText,
					},
					"In stock": notion.DatabaseProperty{
						ID:       "{xYx",
						Type:     notion.DBPropTypeCheckbox,
						Checkbox: &notion.EmptyMetadata{},
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
						Date: &notion.EmptyMetadata{},
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
						ID:     "aGut",
						Type:   notion.DBPropTypePeople,
						People: &notion.EmptyMetadata{},
					},
					"Photo": {
						ID:    "aTIT",
						Type:  "files",
						Files: &notion.EmptyMetadata{},
					},
				},
				Parent: notion.Parent{
					Type:   notion.ParentTypePage,
					PageID: "b8595b75-abd1-4cad-8dfe-f935a8ef57cb",
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

func TestQueryDatabase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		query          *notion.DatabaseQuery
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expPostBody    map[string]interface{}
		expResponse    notion.DatabaseQueryResponse
		expError       error
	}{
		{
			name: "with query, successful response",
			query: &notion.DatabaseQuery{
				Filter: &notion.DatabaseQueryFilter{
					Property: "Name",
					Text: &notion.TextDatabaseQueryFilter{
						Contains: "foobar",
					},
				},
				Sorts: []notion.DatabaseQuerySort{
					{
						Property:  "Name",
						Timestamp: notion.SortTimeStampCreatedTime,
						Direction: notion.SortDirAsc,
					},
					{
						Property:  "Date",
						Timestamp: notion.SortTimeStampLastEditedTime,
						Direction: notion.SortDirDesc,
					},
				},
			},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "list",
						"results": [
							{
								"object": "page",
								"id": "7c6b1c95-de50-45ca-94e6-af1d9fd295ab",
								"created_time": "2021-05-18T17:50:22.371Z",
								"last_edited_time": "2021-05-18T17:50:22.371Z",
								"parent": {
									"type": "database_id",
									"database_id": "39ddfc9d-33c9-404c-89cf-79f01c42dd0c"
								},
								"archived": false,
								"url": "https://www.notion.so/Avocado-251d2b5f268c4de2afe9c71ff92ca95c",
								"properties": {
									"Date": {
										"id": "Q]uT",
										"type": "date",
										"date": {
											"start": "2021-05-18T12:49:00.000-05:00",
											"end": null
										}
									},
									"Name": {
										"id": "title",
										"type": "title",
										"title": [
											{
												"type": "text",
												"text": {
													"content": "Foobar",
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
												"plain_text": "Foobar",
												"href": null
											}
										]
									},
									"Age": {
										"id": "$9nb",
										"type": "number",
										"number": 42
									},
									"People": {
										"id": "1#nc",
										"type": "people",
										"people": [
											{
												"id": "be32e790-8292-46df-a248-b784fdf483cf",
												"name": "Jane Doe",
												"avatar_url": "https://example.com/image.png",
												"type": "person",
												"person": {
													"email": "jane@example.com"
												}
											}
										]
									},
									"Files": {
										"id": "!$9x",
										"type": "files",
										"files": [
											{
												"name": "foobar.pdf"
											}
										]
									},
									"Checkbox": {
										"id": "49S@",
										"type": "checkbox",
										"checkbox": true
									},
									"URL": {
										"id": "93$$",
										"type": "url",
										"url": "https://example.com"
									},
									"Email": {
										"id": "xb3Q",
										"type": "email",
										"email": "jane@example.com"
									},
									"PhoneNumber": {
										"id": "c2#Q",
										"type": "phone_number",
										"phone_number": "867-5309"
									},
									"CreatedTime": {
										"id": "s#0s",
										"type": "created_time",
										"created_time": "2021-05-24T15:44:09.123Z"
									},
									"CreatedBy": {
										"id": "49S@",
										"type": "created_by",
										"created_by": {
											"id": "be32e790-8292-46df-a248-b784fdf483cf",
											"name": "Jane Doe",
											"avatar_url": "https://example.com/image.png",
											"type": "person",
											"person": {
												"email": "jane@example.com"
											}
										}
									},
									"LastEditedTime": {
										"id": "x#0s",
										"type": "last_edited_time",
										"last_edited_time": "2021-05-24T15:44:09.123Z"
									},
									"LastEditedBy": {
										"id": "x9S@",
										"type": "last_edited_by",
										"last_edited_by": {
											"id": "be32e790-8292-46df-a248-b784fdf483cf",
											"name": "Jane Doe",
											"avatar_url": "https://example.com/image.png",
											"type": "person",
											"person": {
												"email": "jane@example.com"
											}
										}
									},
									"Calculation": {
										"id": "s(4f",
										"type": "formula",
										"formula": {
											"type": "number",
											"number": 42
										}
									},
									"Relation": {
										"id": "Cxl[",
										"type": "relation",
										"relation": [
											{
												"id": "2be9597f-693f-4b87-baf9-efc545d38ebe"
											}
										]
									},
									"Rollup": {
										"id": "xyA}",
										"type": "rollup",
										"rollup": {
											"type": "array",
											"array": [
												{
													"type": "number",
													"number": 42
												}
											]
										}
									}
								}
							}
						],
						"next_cursor": "A^hd",
						"has_more": true
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody: map[string]interface{}{
				"filter": map[string]interface{}{
					"property": "Name",
					"text": map[string]interface{}{
						"contains": "foobar",
					},
				},
				"sorts": []interface{}{
					map[string]interface{}{
						"property":  "Name",
						"timestamp": "created_time",
						"direction": "ascending",
					},
					map[string]interface{}{
						"property":  "Date",
						"timestamp": "last_edited_time",
						"direction": "descending",
					},
				},
			},
			expResponse: notion.DatabaseQueryResponse{
				Results: []notion.Page{
					{
						ID:             "7c6b1c95-de50-45ca-94e6-af1d9fd295ab",
						CreatedTime:    mustParseTime(time.RFC3339Nano, "2021-05-18T17:50:22.371Z"),
						LastEditedTime: mustParseTime(time.RFC3339Nano, "2021-05-18T17:50:22.371Z"),
						URL:            "https://www.notion.so/Avocado-251d2b5f268c4de2afe9c71ff92ca95c",
						Parent: notion.Parent{
							Type:       notion.ParentTypeDatabase,
							DatabaseID: "39ddfc9d-33c9-404c-89cf-79f01c42dd0c",
						},
						Archived: false,
						Properties: notion.DatabasePageProperties{
							"Date": notion.DatabasePageProperty{
								ID:   "Q]uT",
								Type: notion.DBPropTypeDate,
								Date: &notion.Date{
									Start: mustParseDateTime("2021-05-18T12:49:00.000-05:00"),
								},
							},
							"Name": notion.DatabasePageProperty{
								ID:   "title",
								Type: notion.DBPropTypeTitle,
								Title: []notion.RichText{
									{
										Type: notion.RichTextTypeText,
										Text: &notion.Text{
											Content: "Foobar",
										},
										PlainText: "Foobar",
										Annotations: &notion.Annotations{
											Color: notion.ColorDefault,
										},
									},
								},
							},
							"Age": notion.DatabasePageProperty{
								ID:     "$9nb",
								Type:   notion.DBPropTypeNumber,
								Number: notion.Float64Ptr(42),
							},
							"People": notion.DatabasePageProperty{
								ID:   "1#nc",
								Type: notion.DBPropTypePeople,
								People: []notion.User{
									{
										ID:        "be32e790-8292-46df-a248-b784fdf483cf",
										Name:      "Jane Doe",
										AvatarURL: notion.StringPtr("https://example.com/image.png"),
										Type:      "person",
										Person: &notion.Person{
											Email: "jane@example.com",
										},
									},
								},
							},
							"Files": notion.DatabasePageProperty{
								ID:   "!$9x",
								Type: notion.DBPropTypeFiles,
								Files: []notion.File{
									{
										Name: "foobar.pdf",
									},
								},
							},
							"Checkbox": notion.DatabasePageProperty{
								ID:       "49S@",
								Type:     notion.DBPropTypeCheckbox,
								Checkbox: notion.BoolPtr(true),
							},
							"Calculation": notion.DatabasePageProperty{
								ID:   "s(4f",
								Type: notion.DBPropTypeFormula,
								Formula: &notion.FormulaResult{
									Type:   notion.FormulaResultTypeNumber,
									Number: notion.Float64Ptr(float64(42)),
								},
							},
							"URL": notion.DatabasePageProperty{
								ID:   "93$$",
								Type: notion.DBPropTypeURL,
								URL:  notion.StringPtr("https://example.com"),
							},
							"Email": notion.DatabasePageProperty{
								ID:    "xb3Q",
								Type:  notion.DBPropTypeEmail,
								Email: notion.StringPtr("jane@example.com"),
							},
							"PhoneNumber": notion.DatabasePageProperty{
								ID:          "c2#Q",
								Type:        notion.DBPropTypePhoneNumber,
								PhoneNumber: notion.StringPtr("867-5309"),
							},
							"CreatedTime": notion.DatabasePageProperty{
								ID:          "s#0s",
								Type:        notion.DBPropTypeCreatedTime,
								CreatedTime: notion.TimePtr(mustParseTime(time.RFC3339Nano, "2021-05-24T15:44:09.123Z")),
							},
							"CreatedBy": notion.DatabasePageProperty{
								ID:   "49S@",
								Type: notion.DBPropTypeCreatedBy,
								CreatedBy: &notion.User{
									ID:        "be32e790-8292-46df-a248-b784fdf483cf",
									Name:      "Jane Doe",
									AvatarURL: notion.StringPtr("https://example.com/image.png"),
									Type:      "person",
									Person: &notion.Person{
										Email: "jane@example.com",
									},
								},
							},
							"LastEditedTime": notion.DatabasePageProperty{
								ID:             "x#0s",
								Type:           notion.DBPropTypeLastEditedTime,
								LastEditedTime: notion.TimePtr(mustParseTime(time.RFC3339Nano, "2021-05-24T15:44:09.123Z")),
							},
							"LastEditedBy": notion.DatabasePageProperty{
								ID:   "x9S@",
								Type: notion.DBPropTypeLastEditedBy,
								LastEditedBy: &notion.User{
									ID:        "be32e790-8292-46df-a248-b784fdf483cf",
									Name:      "Jane Doe",
									AvatarURL: notion.StringPtr("https://example.com/image.png"),
									Type:      "person",
									Person: &notion.Person{
										Email: "jane@example.com",
									},
								},
							},
							"Relation": notion.DatabasePageProperty{
								ID:   "Cxl[",
								Type: notion.DBPropTypeRelation,
								Relation: []notion.Relation{
									{
										ID: "2be9597f-693f-4b87-baf9-efc545d38ebe",
									},
								},
							},
							"Rollup": notion.DatabasePageProperty{
								ID:   "xyA}",
								Type: notion.DBPropTypeRollup,
								Rollup: &notion.RollupResult{
									Type: notion.RollupResultTypeArray,
									Array: []notion.DatabasePageProperty{
										{
											Type:   notion.DBPropTypeNumber,
											Number: notion.Float64Ptr(42),
										},
									},
								},
							},
						},
					},
				},
				HasMore:    true,
				NextCursor: notion.StringPtr("A^hd"),
			},
			expError: nil,
		},
		{
			name:  "without query, doesn't send POST body",
			query: nil,
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "list",
						"results": [],
						"next_cursor": null,
						"has_more": false
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody:    nil,
			expResponse: notion.DatabaseQueryResponse{
				Results:    []notion.Page{},
				HasMore:    false,
				NextCursor: nil,
			},
			expError: nil,
		},
		{
			name:  "with non nil query, but without fields, omits all fields from POST body",
			query: &notion.DatabaseQuery{},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "list",
						"results": [],
						"next_cursor": null,
						"has_more": false
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody:    map[string]interface{}{},
			expResponse: notion.DatabaseQueryResponse{
				Results:    []notion.Page{},
				HasMore:    false,
				NextCursor: nil,
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
			expResponse:    notion.DatabaseQueryResponse{},
			expError:       errors.New("notion: failed to query database: foobar (code: validation_error, status: 400)"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			httpClient := &http.Client{
				Transport: &mockRoundtripper{fn: func(r *http.Request) (*http.Response, error) {
					postBody := make(map[string]interface{})

					err := json.NewDecoder(r.Body).Decode(&postBody)
					if err != nil && err != io.EOF {
						t.Fatal(err)
					}

					if len(tt.expPostBody) == 0 && len(postBody) != 0 {
						t.Errorf("unexpected post body: %+v", postBody)
					}

					if len(tt.expPostBody) != 0 && len(postBody) == 0 {
						t.Errorf("post body not equal (expected %+v, got: nil)", tt.expPostBody)
					}

					if len(tt.expPostBody) != 0 && len(postBody) != 0 {
						if diff := cmp.Diff(tt.expPostBody, postBody); diff != "" {
							t.Errorf("post body not equal (-exp, +got):\n%v", diff)
						}
					}

					return &http.Response{
						StatusCode: tt.respStatusCode,
						Status:     http.StatusText(tt.respStatusCode),
						Body:       ioutil.NopCloser(tt.respBody(r)),
					}, nil
				}},
			}
			client := notion.NewClient("secret-api-key", notion.WithHTTPClient(httpClient))
			resp, err := client.QueryDatabase(context.Background(), "00000000-0000-0000-0000-000000000000", tt.query)

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expResponse, resp); diff != "" {
				t.Fatalf("response not equal (-exp, +got):\n%v", diff)
			}
		})
	}
}
func TestCreateDatabase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		params         notion.CreateDatabaseParams
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expPostBody    map[string]interface{}
		expResponse    notion.Database
		expError       error
	}{
		{
			name: "successful response",
			params: notion.CreateDatabaseParams{
				ParentPageID: "b0668f48-8d66-4733-9bdb-2f82215707f7",
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Foobar",
						},
					},
				},
				Properties: notion.DatabaseProperties{
					"Title": notion.DatabaseProperty{
						Type:  notion.DBPropTypeTitle,
						Title: &notion.EmptyMetadata{},
					},
				},
				Icon: &notion.Icon{
					Type:  notion.IconTypeEmoji,
					Emoji: notion.StringPtr("✌️"),
				},
			},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "database",
						"id": "b89664e3-30b4-474a-9cce-c72a4827d1e4",
						"created_time": "2021-07-20T20:09:00.000Z",
						"last_edited_time": "2021-07-20T20:09:00.000Z",
						"title": [
							{
								"type": "text",
								"text": {
									"content": "Foobar",
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
								"plain_text": "Foobar",
								"href": null
							}
						],
						"properties": {
							"Title": {
								"id": "title",
								"type": "title",
								"title": {}
							}
						},
						"parent": {
							"type": "page_id",
							"page_id": "b0668f48-8d66-4733-9bdb-2f82215707f7"
						},
						"icon": {
							"type": "emoji",
							"emoji": "✌️"
						}
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody: map[string]interface{}{
				"parent": map[string]interface{}{
					"type":    "page_id",
					"page_id": "b0668f48-8d66-4733-9bdb-2f82215707f7",
				},
				"title": []interface{}{
					map[string]interface{}{
						"text": map[string]interface{}{
							"content": "Foobar",
						},
					},
				},
				"properties": map[string]interface{}{
					"Title": map[string]interface{}{
						"type":  "title",
						"title": map[string]interface{}{},
					},
				},
				"icon": map[string]interface{}{
					"type":  "emoji",
					"emoji": "✌️",
				},
			},
			expResponse: notion.Database{
				ID:             "b89664e3-30b4-474a-9cce-c72a4827d1e4",
				CreatedTime:    mustParseTime(time.RFC3339Nano, "2021-07-20T20:09:00Z"),
				LastEditedTime: mustParseTime(time.RFC3339Nano, "2021-07-20T20:09:00Z"),
				Parent: notion.Parent{
					Type:   notion.ParentTypePage,
					PageID: "b0668f48-8d66-4733-9bdb-2f82215707f7",
				},
				Title: []notion.RichText{
					{
						Type: notion.RichTextTypeText,
						Text: &notion.Text{
							Content: "Foobar",
						},
						Annotations: &notion.Annotations{
							Color: notion.ColorDefault,
						},
						PlainText: "Foobar",
					},
				},
				Properties: notion.DatabaseProperties{
					"Title": notion.DatabaseProperty{
						ID:    "title",
						Type:  notion.DBPropTypeTitle,
						Title: &notion.EmptyMetadata{},
					},
				},
				Icon: notion.Icon{
					Type:  notion.IconTypeEmoji,
					Emoji: notion.StringPtr("✌️"),
				},
			},
			expError: nil,
		},
		{
			name: "error response",
			params: notion.CreateDatabaseParams{
				ParentPageID: "b0668f48-8d66-4733-9bdb-2f82215707f7",
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Foobar",
						},
					},
				},
				Properties: notion.DatabaseProperties{
					"Title": notion.DatabaseProperty{
						Type:  notion.DBPropTypeTitle,
						Title: &notion.EmptyMetadata{},
					},
				},
			},
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
			expPostBody: map[string]interface{}{
				"parent": map[string]interface{}{
					"type":    "page_id",
					"page_id": "b0668f48-8d66-4733-9bdb-2f82215707f7",
				},
				"title": []interface{}{
					map[string]interface{}{
						"text": map[string]interface{}{
							"content": "Foobar",
						},
					},
				},
				"properties": map[string]interface{}{
					"Title": map[string]interface{}{
						"type":  "title",
						"title": map[string]interface{}{},
					},
				},
			},
			expResponse: notion.Database{},
			expError:    errors.New("notion: failed to create database: foobar (code: validation_error, status: 400)"),
		},
		{
			name: "parent id required error",
			params: notion.CreateDatabaseParams{
				Properties: notion.DatabaseProperties{},
			},
			expResponse: notion.Database{},
			expError:    errors.New("notion: invalid database params: parent page ID is required"),
		},
		{
			name: "database properties required error",
			params: notion.CreateDatabaseParams{
				ParentPageID: "b0668f48-8d66-4733-9bdb-2f82215707f7",
			},
			expResponse: notion.Database{},
			expError:    errors.New("notion: invalid database params: database properties are required"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			httpClient := &http.Client{
				Transport: &mockRoundtripper{fn: func(r *http.Request) (*http.Response, error) {
					postBody := make(map[string]interface{})

					err := json.NewDecoder(r.Body).Decode(&postBody)
					if err != nil && err != io.EOF {
						t.Fatal(err)
					}

					if len(tt.expPostBody) == 0 && len(postBody) != 0 {
						t.Errorf("unexpected post body: %#v", postBody)
					}

					if len(tt.expPostBody) != 0 && len(postBody) == 0 {
						t.Errorf("post body not equal (expected %+v, got: nil)", tt.expPostBody)
					}

					if len(tt.expPostBody) != 0 && len(postBody) != 0 {
						if diff := cmp.Diff(tt.expPostBody, postBody); diff != "" {
							t.Errorf("post body not equal (-exp, +got):\n%v", diff)
						}
					}

					return &http.Response{
						StatusCode: tt.respStatusCode,
						Status:     http.StatusText(tt.respStatusCode),
						Body:       ioutil.NopCloser(tt.respBody(r)),
					}, nil
				}},
			}
			client := notion.NewClient("secret-api-key", notion.WithHTTPClient(httpClient))
			page, err := client.CreateDatabase(context.Background(), tt.params)

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expResponse, page); diff != "" {
				t.Fatalf("response not equal (-exp, +got):\n%v", diff)
			}
		})
	}
}

func TestFindPageByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expPage        notion.Page
		expError       error
	}{
		{
			name: "successful response",
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "page",
						"id": "606ed832-7d79-46de-bbed-5b4896e7bc02",
						"created_time": "2021-05-19T18:34:00.000Z",
						"last_edited_time": "2021-05-19T18:34:00.000Z",
						"parent": {
							"type": "page_id",
							"page_id": "b0668f48-8d66-4733-9bdb-2f82215707f7"
						},
						"archived": false,
						"url": "https://www.notion.so/Avocado-251d2b5f268c4de2afe9c71ff92ca95c",
						"properties": {
							"title": {
								"id": "title",
								"type": "title",
								"title": [
									{
										"type": "text",
										"text": {
											"content": "Lorem ipsum",
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
										"plain_text": "Lorem ipsum",
										"href": null
									}
								]
							}
						}
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPage: notion.Page{
				ID:             "606ed832-7d79-46de-bbed-5b4896e7bc02",
				CreatedTime:    mustParseTime(time.RFC3339Nano, "2021-05-19T18:34:00.000Z"),
				LastEditedTime: mustParseTime(time.RFC3339Nano, "2021-05-19T18:34:00.000Z"),
				URL:            "https://www.notion.so/Avocado-251d2b5f268c4de2afe9c71ff92ca95c",
				Parent: notion.Parent{
					Type:   notion.ParentTypePage,
					PageID: "b0668f48-8d66-4733-9bdb-2f82215707f7",
				},
				Properties: notion.PageProperties{
					Title: notion.PageTitle{
						Title: []notion.RichText{
							{
								Type: notion.RichTextTypeText,
								Text: &notion.Text{
									Content: "Lorem ipsum",
								},
								Annotations: &notion.Annotations{
									Color: notion.ColorDefault,
								},
								PlainText: "Lorem ipsum",
							},
						},
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
						"status": 404,
						"code": "object_not_found",
						"message": "foobar"
					}`,
				)
			},
			respStatusCode: http.StatusNotFound,
			expPage:        notion.Page{},
			expError:       errors.New("notion: failed to find page: foobar (code: object_not_found, status: 404)"),
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
			page, err := client.FindPageByID(context.Background(), "00000000-0000-0000-0000-000000000000")

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expPage, page); diff != "" {
				t.Fatalf("page not equal (-exp, +got):\n%v", diff)
			}
		})
	}
}

func TestCreatePage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		params         notion.CreatePageParams
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expPostBody    map[string]interface{}
		expResponse    notion.Page
		expError       error
	}{
		{
			name: "successful response",
			params: notion.CreatePageParams{
				ParentType: notion.ParentTypePage,
				ParentID:   "b0668f48-8d66-4733-9bdb-2f82215707f7",
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Foobar",
						},
					},
				},
				Children: []notion.Block{
					{
						Object: "block",
						Type:   notion.BlockTypeParagraph,
						Paragraph: &notion.RichTextBlock{
							Text: []notion.RichText{
								{
									Text: &notion.Text{
										Content: "Lorem ipsum dolor sit amet.",
									},
								},
							},
						},
					},
				},
			},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "page",
						"id": "276ee233-e426-4ed0-9986-6b22af8550df",
						"created_time": "2021-05-19T19:34:05.068Z",
						"last_edited_time": "2021-05-19T19:34:05.069Z",
						"parent": {
							"type": "page_id",
							"page_id": "b0668f48-8d66-4733-9bdb-2f82215707f7"
						},
						"archived": false,
						"url": "https://www.notion.so/Avocado-251d2b5f268c4de2afe9c71ff92ca95c",
						"properties": {
							"title": {
								"id": "title",
								"type": "title",
								"title": [
									{
										"type": "text",
										"text": {
											"content": "Foobar",
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
										"plain_text": "Foobar",
										"href": null
									}
								]
							}
						}
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody: map[string]interface{}{
				"parent": map[string]interface{}{
					"page_id": "b0668f48-8d66-4733-9bdb-2f82215707f7",
				},
				"properties": map[string]interface{}{
					"title": []interface{}{
						map[string]interface{}{
							"text": map[string]interface{}{
								"content": "Foobar",
							},
						},
					},
				},
				"children": []interface{}{
					map[string]interface{}{
						"object": "block",
						"type":   "paragraph",
						"paragraph": map[string]interface{}{
							"text": []interface{}{
								map[string]interface{}{
									"text": map[string]interface{}{
										"content": "Lorem ipsum dolor sit amet.",
									},
								},
							},
						},
					},
				},
			},
			expResponse: notion.Page{
				ID:             "276ee233-e426-4ed0-9986-6b22af8550df",
				CreatedTime:    mustParseTime(time.RFC3339Nano, "2021-05-19T19:34:05.068Z"),
				LastEditedTime: mustParseTime(time.RFC3339Nano, "2021-05-19T19:34:05.069Z"),
				URL:            "https://www.notion.so/Avocado-251d2b5f268c4de2afe9c71ff92ca95c",
				Parent: notion.Parent{
					Type:   notion.ParentTypePage,
					PageID: "b0668f48-8d66-4733-9bdb-2f82215707f7",
				},
				Properties: notion.PageProperties{
					Title: notion.PageTitle{
						Title: []notion.RichText{
							{
								Type: notion.RichTextTypeText,
								Text: &notion.Text{
									Content: "Foobar",
								},
								Annotations: &notion.Annotations{
									Color: notion.ColorDefault,
								},
								PlainText: "Foobar",
							},
						},
					},
				},
			},
			expError: nil,
		},
		{
			name: "database parent, successful response",
			params: notion.CreatePageParams{
				ParentType: notion.ParentTypeDatabase,
				ParentID:   "b0668f48-8d66-4733-9bdb-2f82215707f7",
				DatabasePageProperties: &notion.DatabasePageProperties{
					"title": notion.DatabasePageProperty{
						Title: []notion.RichText{
							{
								Text: &notion.Text{
									Content: "Foobar",
								},
							},
						},
					},
				},
				Children: []notion.Block{
					{
						Object: "block",
						Type:   notion.BlockTypeParagraph,
						Paragraph: &notion.RichTextBlock{
							Text: []notion.RichText{
								{
									Text: &notion.Text{
										Content: "Lorem ipsum dolor sit amet.",
									},
								},
							},
						},
					},
				},
			},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "page",
						"id": "276ee233-e426-4ed0-9986-6b22af8550df",
						"created_time": "2021-05-19T19:34:05.068Z",
						"last_edited_time": "2021-05-19T19:34:05.069Z",
						"parent": {
							"type": "database_id",
							"database_id": "b0668f48-8d66-4733-9bdb-2f82215707f7"
						},
						"archived": false,
						"properties": {
							"title": {
								"id": "title",
								"title": [
									{
										"text": {
											"content": "Foobar",
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
										"href": null
									}
								]
							}
						}
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody: map[string]interface{}{
				"parent": map[string]interface{}{
					"database_id": "b0668f48-8d66-4733-9bdb-2f82215707f7",
				},
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"title": []interface{}{
							map[string]interface{}{
								"text": map[string]interface{}{
									"content": "Foobar",
								},
							},
						},
					},
				},
				"children": []interface{}{
					map[string]interface{}{
						"object": "block",
						"type":   "paragraph",
						"paragraph": map[string]interface{}{
							"text": []interface{}{
								map[string]interface{}{
									"text": map[string]interface{}{
										"content": "Lorem ipsum dolor sit amet.",
									},
								},
							},
						},
					},
				},
			},
			expResponse: notion.Page{
				ID:             "276ee233-e426-4ed0-9986-6b22af8550df",
				CreatedTime:    mustParseTime(time.RFC3339Nano, "2021-05-19T19:34:05.068Z"),
				LastEditedTime: mustParseTime(time.RFC3339Nano, "2021-05-19T19:34:05.069Z"),
				Parent: notion.Parent{
					Type:       notion.ParentTypeDatabase,
					DatabaseID: "b0668f48-8d66-4733-9bdb-2f82215707f7",
				},
				Properties: notion.DatabasePageProperties{
					"title": notion.DatabasePageProperty{
						ID: "title",
						Title: []notion.RichText{
							{
								Text: &notion.Text{
									Content: "Foobar",
								},
								Annotations: &notion.Annotations{
									Color: notion.ColorDefault,
								},
							},
						},
					},
				},
			},
			expError: nil,
		},
		{
			name: "error response",
			params: notion.CreatePageParams{
				ParentType: notion.ParentTypePage,
				ParentID:   "b0668f48-8d66-4733-9bdb-2f82215707f7",
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Foobar",
						},
					},
				},
			},
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
			expPostBody: map[string]interface{}{
				"parent": map[string]interface{}{
					"page_id": "b0668f48-8d66-4733-9bdb-2f82215707f7",
				},
				"properties": map[string]interface{}{
					"title": []interface{}{
						map[string]interface{}{
							"text": map[string]interface{}{
								"content": "Foobar",
							},
						},
					},
				},
			},
			expResponse: notion.Page{},
			expError:    errors.New("notion: failed to create page: foobar (code: validation_error, status: 400)"),
		},
		{
			name: "parent type required error",
			params: notion.CreatePageParams{
				ParentID: "b0668f48-8d66-4733-9bdb-2f82215707f7",
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Foobar",
						},
					},
				},
			},
			expResponse: notion.Page{},
			expError:    errors.New("notion: invalid page params: parent type is required"),
		},
		{
			name: "parent id required error",
			params: notion.CreatePageParams{
				ParentType: notion.ParentTypePage,
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Foobar",
						},
					},
				},
			},
			expResponse: notion.Page{},
			expError:    errors.New("notion: invalid page params: parent ID is required"),
		},
		{
			name: "page title required error",
			params: notion.CreatePageParams{
				ParentType: notion.ParentTypePage,
				ParentID:   "b0668f48-8d66-4733-9bdb-2f82215707f7",
			},
			expResponse: notion.Page{},
			expError:    errors.New("notion: invalid page params: title is required when parent type is page"),
		},
		{
			name: "database properties required error",
			params: notion.CreatePageParams{
				ParentType: notion.ParentTypeDatabase,
				ParentID:   "b0668f48-8d66-4733-9bdb-2f82215707f7",
			},
			expResponse: notion.Page{},
			expError:    errors.New("notion: invalid page params: database page properties is required when parent type is database"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			httpClient := &http.Client{
				Transport: &mockRoundtripper{fn: func(r *http.Request) (*http.Response, error) {
					postBody := make(map[string]interface{})

					err := json.NewDecoder(r.Body).Decode(&postBody)
					if err != nil && err != io.EOF {
						t.Fatal(err)
					}

					if len(tt.expPostBody) == 0 && len(postBody) != 0 {
						t.Errorf("unexpected post body: %#v", postBody)
					}

					if len(tt.expPostBody) != 0 && len(postBody) == 0 {
						t.Errorf("post body not equal (expected %+v, got: nil)", tt.expPostBody)
					}

					if len(tt.expPostBody) != 0 && len(postBody) != 0 {
						if diff := cmp.Diff(tt.expPostBody, postBody); diff != "" {
							t.Errorf("post body not equal (-exp, +got):\n%v", diff)
						}
					}

					return &http.Response{
						StatusCode: tt.respStatusCode,
						Status:     http.StatusText(tt.respStatusCode),
						Body:       ioutil.NopCloser(tt.respBody(r)),
					}, nil
				}},
			}
			client := notion.NewClient("secret-api-key", notion.WithHTTPClient(httpClient))
			page, err := client.CreatePage(context.Background(), tt.params)

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expResponse, page); diff != "" {
				t.Fatalf("response not equal (-exp, +got):\n%v", diff)
			}
		})
	}
}

func TestUpdatePageProps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		params         notion.UpdatePageParams
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expPostBody    map[string]interface{}
		expResponse    notion.Page
		expError       error
	}{
		{
			name: "page props, successful response",
			params: notion.UpdatePageParams{
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Foobar",
						},
					},
				},
			},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "page",
						"id": "cb261dc5-6c85-4767-8585-3852382fb466",
						"created_time": "2021-05-14T09:15:46.796Z",
						"last_edited_time": "2021-05-22T15:54:31.116Z",
						"parent": {
							"type": "page_id",
							"page_id": "b0668f48-8d66-4733-9bdb-2f82215707f7"
						},
						"archived": false,
						"url": "https://www.notion.so/Avocado-251d2b5f268c4de2afe9c71ff92ca95c",
						"properties": {
							"title": {
								"id": "title",
								"type": "title",
								"title": [
									{
										"type": "text",
										"text": {
											"content": "Lorem ipsum",
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
										"plain_text": "Lorem ipsum",
										"href": null
									}
								]
							}
						}
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody: map[string]interface{}{
				"properties": map[string]interface{}{
					"title": []interface{}{
						map[string]interface{}{
							"text": map[string]interface{}{
								"content": "Foobar",
							},
						},
					},
				},
			},
			expResponse: notion.Page{
				ID:             "cb261dc5-6c85-4767-8585-3852382fb466",
				CreatedTime:    mustParseTime(time.RFC3339Nano, "2021-05-14T09:15:46.796Z"),
				LastEditedTime: mustParseTime(time.RFC3339Nano, "2021-05-22T15:54:31.116Z"),
				URL:            "https://www.notion.so/Avocado-251d2b5f268c4de2afe9c71ff92ca95c",
				Parent: notion.Parent{
					Type:   notion.ParentTypePage,
					PageID: "b0668f48-8d66-4733-9bdb-2f82215707f7",
				},
				Properties: notion.PageProperties{
					Title: notion.PageTitle{
						Title: []notion.RichText{
							{
								Type: notion.RichTextTypeText,
								Text: &notion.Text{
									Content: "Lorem ipsum",
								},
								Annotations: &notion.Annotations{
									Color: notion.ColorDefault,
								},
								PlainText: "Lorem ipsum",
							},
						},
					},
				},
			},
			expError: nil,
		},
		{
			name: "page icon, successful response",
			params: notion.UpdatePageParams{
				Icon: &notion.Icon{
					Type: notion.IconTypeExternal,
					External: &notion.IconExternal{
						URL: "https://www.notion.so/front-static/pages/pricing/pro.png",
					},
				},
			},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "page",
						"id": "cb261dc5-6c85-4767-8585-3852382fb466",
						"created_time": "2021-05-14T09:15:46.796Z",
						"last_edited_time": "2021-05-22T15:54:31.116Z",
						"parent": {
							"type": "page_id",
							"page_id": "b0668f48-8d66-4733-9bdb-2f82215707f7"
						},
						"icon": {
							"type": "external",
							"external": {
								"url": "https://www.notion.so/front-static/pages/pricing/pro.png"
							}
						},
						"archived": false,
						"url": "https://www.notion.so/Avocado-251d2b5f268c4de2afe9c71ff92ca95c",
						"properties": {
							"title": {
								"id": "title",
								"type": "title",
								"title": [
									{
										"type": "text",
										"text": {
											"content": "Lorem ipsum",
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
										"plain_text": "Lorem ipsum",
										"href": null
									}
								]
							}
						}
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody: map[string]interface{}{
				"icon": map[string]interface{}{
					"type": "external",
					"external": map[string]interface{}{
						"url": "https://www.notion.so/front-static/pages/pricing/pro.png",
					},
				},
			},
			expResponse: notion.Page{
				ID:             "cb261dc5-6c85-4767-8585-3852382fb466",
				CreatedTime:    mustParseTime(time.RFC3339Nano, "2021-05-14T09:15:46.796Z"),
				LastEditedTime: mustParseTime(time.RFC3339Nano, "2021-05-22T15:54:31.116Z"),
				URL:            "https://www.notion.so/Avocado-251d2b5f268c4de2afe9c71ff92ca95c",
				Parent: notion.Parent{
					Type:   notion.ParentTypePage,
					PageID: "b0668f48-8d66-4733-9bdb-2f82215707f7",
				},
				Icon: &notion.Icon{
					Type: notion.IconTypeExternal,
					External: &notion.IconExternal{
						URL: "https://www.notion.so/front-static/pages/pricing/pro.png",
					},
				},
				Properties: notion.PageProperties{
					Title: notion.PageTitle{
						Title: []notion.RichText{
							{
								Type: notion.RichTextTypeText,
								Text: &notion.Text{
									Content: "Lorem ipsum",
								},
								Annotations: &notion.Annotations{
									Color: notion.ColorDefault,
								},
								PlainText: "Lorem ipsum",
							},
						},
					},
				},
			},
			expError: nil,
		},
		{
			name: "database page props, successful response",
			params: notion.UpdatePageParams{
				DatabasePageProperties: &notion.DatabasePageProperties{
					"Name": notion.DatabasePageProperty{
						Title: []notion.RichText{
							{
								Text: &notion.Text{
									Content: "Lorem ipsum",
								},
							},
						},
					},
				},
			},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "page",
						"id": "e4f419a7-f01f-4d5b-af58-ff4786a429fe",
						"created_time": "2021-05-17T17:56:00.000Z",
						"last_edited_time": "2021-05-22T16:24:23.007Z",
						"parent": {
							"type": "database_id",
							"database_id": "4cb17949-f08d-4d5c-ab50-fe6ba689d2c8"
						},
						"archived": false,
						"properties": {
							"Name": {
								"id": "title",
								"type": "title",
								"title": [
									{
										"type": "text",
										"text": {
											"content": "Lorem ipsum",
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
										"plain_text": "Lorem ipsum",
										"href": null
									}
								]
							}
						}
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody: map[string]interface{}{
				"properties": map[string]interface{}{
					"Name": map[string]interface{}{
						"title": []interface{}{
							map[string]interface{}{
								"text": map[string]interface{}{
									"content": "Lorem ipsum",
								},
							},
						},
					},
				},
			},
			expResponse: notion.Page{
				ID:             "e4f419a7-f01f-4d5b-af58-ff4786a429fe",
				CreatedTime:    mustParseTime(time.RFC3339Nano, "2021-05-17T17:56:00.000Z"),
				LastEditedTime: mustParseTime(time.RFC3339Nano, "2021-05-22T16:24:23.007Z"),
				Parent: notion.Parent{
					Type:       notion.ParentTypeDatabase,
					DatabaseID: "4cb17949-f08d-4d5c-ab50-fe6ba689d2c8",
				},
				Properties: notion.DatabasePageProperties{
					"Name": notion.DatabasePageProperty{
						ID:   "title",
						Type: notion.DBPropTypeTitle,
						Title: []notion.RichText{
							{
								Type: notion.RichTextTypeText,
								Text: &notion.Text{
									Content: "Lorem ipsum",
								},
								Annotations: &notion.Annotations{
									Color: notion.ColorDefault,
								},
								PlainText: "Lorem ipsum",
							},
						},
					},
				},
			},
			expError: nil,
		},
		{
			name: "error response",
			params: notion.UpdatePageParams{
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Foobar",
						},
					},
				},
			},
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
			expPostBody: map[string]interface{}{
				"properties": map[string]interface{}{
					"title": []interface{}{
						map[string]interface{}{
							"text": map[string]interface{}{
								"content": "Foobar",
							},
						},
					},
				},
			},
			expResponse: notion.Page{},
			expError:    errors.New("notion: failed to update page properties: foobar (code: validation_error, status: 400)"),
		},
		{
			name:        "missing any params",
			params:      notion.UpdatePageParams{},
			expResponse: notion.Page{},
			expError:    errors.New("notion: invalid page params: at least one of database page properties, title or icon is required"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			httpClient := &http.Client{
				Transport: &mockRoundtripper{fn: func(r *http.Request) (*http.Response, error) {
					postBody := make(map[string]interface{})

					err := json.NewDecoder(r.Body).Decode(&postBody)
					if err != nil && err != io.EOF {
						t.Fatal(err)
					}

					if len(tt.expPostBody) == 0 && len(postBody) != 0 {
						t.Errorf("unexpected post body: %#v", postBody)
					}

					if len(tt.expPostBody) != 0 && len(postBody) == 0 {
						t.Errorf("post body not equal (expected %+v, got: nil)", tt.expPostBody)
					}

					if len(tt.expPostBody) != 0 && len(postBody) != 0 {
						if diff := cmp.Diff(tt.expPostBody, postBody); diff != "" {
							t.Errorf("post body not equal (-exp, +got):\n%v", diff)
						}
					}

					return &http.Response{
						StatusCode: tt.respStatusCode,
						Status:     http.StatusText(tt.respStatusCode),
						Body:       ioutil.NopCloser(tt.respBody(r)),
					}, nil
				}},
			}
			client := notion.NewClient("secret-api-key", notion.WithHTTPClient(httpClient))
			page, err := client.UpdatePageProps(context.Background(), "00000000-0000-0000-0000-000000000000", tt.params)

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expResponse, page); diff != "" {
				t.Fatalf("response not equal (-exp, +got):\n%v", diff)
			}
		})
	}
}

func TestFindBlockChildrenById(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		query          *notion.PaginationQuery
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expQueryParams url.Values
		expResponse    notion.BlockChildrenResponse
		expError       error
	}{
		{
			name: "with query, successful response",
			query: &notion.PaginationQuery{
				StartCursor: "7c6b1c95-de50-45ca-94e6-af1d9fd295ab",
				PageSize:    42,
			},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "list",
						"results": [
							{
								"object": "block",
								"id": "ae9c9a31-1c1e-4ae2-a5ee-c539a2d43113",
								"created_time": "2021-05-14T09:15:00.000Z",
								"last_edited_time": "2021-05-14T09:15:00.000Z",
								"has_children": false,
								"type": "paragraph",
								"paragraph": {
									"text": [
										{
											"type": "text",
											"text": {
												"content": "Lorem ipsum dolor sit amet.",
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
											"plain_text": "Lorem ipsum dolor sit amet.",
											"href": null
										}
									]
								}
							}
						],
						"next_cursor": "A^hd",
						"has_more": true
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expQueryParams: url.Values{
				"start_cursor": []string{"7c6b1c95-de50-45ca-94e6-af1d9fd295ab"},
				"page_size":    []string{"42"},
			},
			expResponse: notion.BlockChildrenResponse{
				Results: []notion.Block{
					{
						Object:         "block",
						ID:             "ae9c9a31-1c1e-4ae2-a5ee-c539a2d43113",
						CreatedTime:    notion.TimePtr(mustParseTime(time.RFC3339Nano, "2021-05-14T09:15:00.000Z")),
						LastEditedTime: notion.TimePtr(mustParseTime(time.RFC3339Nano, "2021-05-14T09:15:00.000Z")),
						Type:           notion.BlockTypeParagraph,
						Paragraph: &notion.RichTextBlock{
							Text: []notion.RichText{
								{
									Type: notion.RichTextTypeText,
									Text: &notion.Text{
										Content: "Lorem ipsum dolor sit amet.",
									},
									Annotations: &notion.Annotations{
										Color: notion.ColorDefault,
									},
									PlainText: "Lorem ipsum dolor sit amet.",
								},
							},
						},
					},
				},
				HasMore:    true,
				NextCursor: notion.StringPtr("A^hd"),
			},
			expError: nil,
		},
		{
			name:  "without query, successful response",
			query: nil,
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "list",
						"results": [],
						"next_cursor": null,
						"has_more": false
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expQueryParams: nil,
			expResponse: notion.BlockChildrenResponse{
				Results:    []notion.Block{},
				HasMore:    false,
				NextCursor: nil,
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
			expResponse:    notion.BlockChildrenResponse{},
			expError:       errors.New("notion: failed to find block children: foobar (code: validation_error, status: 400)"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			httpClient := &http.Client{
				Transport: &mockRoundtripper{fn: func(r *http.Request) (*http.Response, error) {
					q := r.URL.Query()

					if len(tt.expQueryParams) == 0 && len(q) != 0 {
						t.Errorf("unexpected query params: %+v", q)
					}

					if len(tt.expQueryParams) != 0 && len(q) == 0 {
						t.Errorf("query params not equal (expected %+v, got: nil)", tt.expQueryParams)
					}

					if len(tt.expQueryParams) != 0 && len(q) != 0 {
						if diff := cmp.Diff(tt.expQueryParams, q); diff != "" {
							t.Errorf("query params not equal (-exp, +got):\n%v", diff)
						}
					}

					return &http.Response{
						StatusCode: tt.respStatusCode,
						Status:     http.StatusText(tt.respStatusCode),
						Body:       ioutil.NopCloser(tt.respBody(r)),
					}, nil
				}},
			}
			client := notion.NewClient("secret-api-key", notion.WithHTTPClient(httpClient))
			resp, err := client.FindBlockChildrenByID(context.Background(), "00000000-0000-0000-0000-000000000000", tt.query)

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expResponse, resp); diff != "" {
				t.Fatalf("response not equal (-exp, +got):\n%v", diff)
			}
		})
	}
}

func TestAppendBlockChildren(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		children       []notion.Block
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expPostBody    map[string]interface{}
		expBlock       notion.Block
		expError       error
	}{
		{
			name: "successful response",
			children: []notion.Block{
				{
					Type: notion.BlockTypeParagraph,
					Paragraph: &notion.RichTextBlock{
						Text: []notion.RichText{
							{
								Text: &notion.Text{
									Content: "Lorem ipsum dolor sit amet.",
								},
							},
						},
					},
				},
			},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "block",
						"id": "cb261dc5-6c85-4767-8585-3852382fb466",
						"created_time": "2021-05-14T09:15:46.796Z",
						"last_edited_time": "2021-05-22T20:31:43.231Z",
						"has_children": true,
						"type": "child_page",
						"child_page": {
							"title": "Sub page"
						}
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody: map[string]interface{}{
				"children": []interface{}{
					map[string]interface{}{
						"object": "block",
						"type":   "paragraph",
						"paragraph": map[string]interface{}{
							"text": []interface{}{
								map[string]interface{}{
									"text": map[string]interface{}{
										"content": "Lorem ipsum dolor sit amet.",
									},
								},
							},
						},
					},
				},
			},
			expBlock: notion.Block{
				Object:         "block",
				ID:             "cb261dc5-6c85-4767-8585-3852382fb466",
				CreatedTime:    notion.TimePtr(mustParseTime(time.RFC3339Nano, "2021-05-14T09:15:46.796Z")),
				LastEditedTime: notion.TimePtr(mustParseTime(time.RFC3339Nano, "2021-05-22T20:31:43.231Z")),
				HasChildren:    true,
				Type:           notion.BlockTypeChildPage,
				ChildPage: &notion.ChildPage{
					Title: "Sub page",
				},
			},
			expError: nil,
		},
		{
			name: "error response",
			children: []notion.Block{
				{
					Type: notion.BlockTypeParagraph,
					Paragraph: &notion.RichTextBlock{
						Text: []notion.RichText{
							{
								Text: &notion.Text{
									Content: "Lorem ipsum dolor sit amet.",
								},
							},
						},
					},
				},
			},
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
			expPostBody: map[string]interface{}{
				"children": []interface{}{
					map[string]interface{}{
						"object": "block",
						"type":   "paragraph",
						"paragraph": map[string]interface{}{
							"text": []interface{}{
								map[string]interface{}{
									"text": map[string]interface{}{
										"content": "Lorem ipsum dolor sit amet.",
									},
								},
							},
						},
					},
				},
			},
			expBlock: notion.Block{},
			expError: errors.New("notion: failed to append block children: foobar (code: validation_error, status: 400)"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			httpClient := &http.Client{
				Transport: &mockRoundtripper{fn: func(r *http.Request) (*http.Response, error) {
					postBody := make(map[string]interface{})

					err := json.NewDecoder(r.Body).Decode(&postBody)
					if err != nil && err != io.EOF {
						t.Fatal(err)
					}

					if len(tt.expPostBody) == 0 && len(postBody) != 0 {
						t.Errorf("unexpected post body: %#v", postBody)
					}

					if len(tt.expPostBody) != 0 && len(postBody) == 0 {
						t.Errorf("post body not equal (expected %+v, got: nil)", tt.expPostBody)
					}

					if len(tt.expPostBody) != 0 && len(postBody) != 0 {
						if diff := cmp.Diff(tt.expPostBody, postBody); diff != "" {
							t.Errorf("post body not equal (-exp, +got):\n%v", diff)
						}
					}

					return &http.Response{
						StatusCode: tt.respStatusCode,
						Status:     http.StatusText(tt.respStatusCode),
						Body:       ioutil.NopCloser(tt.respBody(r)),
					}, nil
				}},
			}
			client := notion.NewClient("secret-api-key", notion.WithHTTPClient(httpClient))
			block, err := client.AppendBlockChildren(context.Background(), "00000000-0000-0000-0000-000000000000", tt.children)

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expBlock, block); diff != "" {
				t.Fatalf("response not equal (-exp, +got):\n%v", diff)
			}
		})
	}
}

func TestFindUserByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expUser        notion.User
		expError       error
	}{
		{
			name: "successful response",
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "user",
						"id": "be32e790-8292-46df-a248-b784fdf483cf",
						"name": "Jane Doe",
						"avatar_url": "https://example.com/avatar.png",
						"type": "person",
						"person": {
							"email": "jane@example.com"
						}
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expUser: notion.User{
				ID:        "be32e790-8292-46df-a248-b784fdf483cf",
				Name:      "Jane Doe",
				AvatarURL: notion.StringPtr("https://example.com/avatar.png"),
				Type:      "person",
				Person: &notion.Person{
					Email: "jane@example.com",
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
						"status": 404,
						"code": "object_not_found",
						"message": "foobar"
					}`,
				)
			},
			respStatusCode: http.StatusNotFound,
			expUser:        notion.User{},
			expError:       errors.New("notion: failed to find user: foobar (code: object_not_found, status: 404)"),
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
			user, err := client.FindUserByID(context.Background(), "00000000-0000-0000-0000-000000000000")

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expUser, user); diff != "" {
				t.Fatalf("user not equal (-exp, +got):\n%v", diff)
			}
		})
	}
}

func TestListUsers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		query          *notion.PaginationQuery
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expQueryParams url.Values
		expResponse    notion.ListUsersResponse
		expError       error
	}{
		{
			name: "with query, successful response",
			query: &notion.PaginationQuery{
				StartCursor: "7c6b1c95-de50-45ca-94e6-af1d9fd295ab",
				PageSize:    42,
			},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "list",
						"results": [
							{
								"object": "user",
								"id": "be32e790-8292-46df-a248-b784fdf483cf",
								"name": "Jane Doe",
								"avatar_url": "https://example.com/avatar.png",
								"type": "person",
								"person": {
									"email": "jane@example.com"
								}
							},
							{
								"object": "user",
								"id": "25c9cc08-1afd-4d22-b9e6-31b0f6e7b44f",
								"name": "Johnny 5",
								"avatar_url": null,
								"type": "bot",
								"bot": {}
							}
						],
						"next_cursor": "A^hd",
						"has_more": true
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expQueryParams: url.Values{
				"start_cursor": []string{"7c6b1c95-de50-45ca-94e6-af1d9fd295ab"},
				"page_size":    []string{"42"},
			},
			expResponse: notion.ListUsersResponse{
				Results: []notion.User{
					{
						ID:        "be32e790-8292-46df-a248-b784fdf483cf",
						Name:      "Jane Doe",
						AvatarURL: notion.StringPtr("https://example.com/avatar.png"),
						Type:      "person",
						Person: &notion.Person{
							Email: "jane@example.com",
						},
					},
					{
						ID:   "25c9cc08-1afd-4d22-b9e6-31b0f6e7b44f",
						Name: "Johnny 5",
						Type: "bot",
						Bot:  &notion.Bot{},
					},
				},
				HasMore:    true,
				NextCursor: notion.StringPtr("A^hd"),
			},
			expError: nil,
		},
		{
			name:  "without query, successful response",
			query: nil,
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "list",
						"results": [],
						"next_cursor": null,
						"has_more": false
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expQueryParams: nil,
			expResponse: notion.ListUsersResponse{
				Results:    []notion.User{},
				HasMore:    false,
				NextCursor: nil,
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
			expResponse:    notion.ListUsersResponse{},
			expError:       errors.New("notion: failed to list users: foobar (code: validation_error, status: 400)"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			httpClient := &http.Client{
				Transport: &mockRoundtripper{fn: func(r *http.Request) (*http.Response, error) {
					q := r.URL.Query()

					if len(tt.expQueryParams) == 0 && len(q) != 0 {
						t.Errorf("unexpected query params: %+v", q)
					}

					if len(tt.expQueryParams) != 0 && len(q) == 0 {
						t.Errorf("query params not equal (expected %+v, got: nil)", tt.expQueryParams)
					}

					if len(tt.expQueryParams) != 0 && len(q) != 0 {
						if diff := cmp.Diff(tt.expQueryParams, q); diff != "" {
							t.Errorf("query params not equal (-exp, +got):\n%v", diff)
						}
					}

					return &http.Response{
						StatusCode: tt.respStatusCode,
						Status:     http.StatusText(tt.respStatusCode),
						Body:       ioutil.NopCloser(tt.respBody(r)),
					}, nil
				}},
			}
			client := notion.NewClient("secret-api-key", notion.WithHTTPClient(httpClient))
			resp, err := client.ListUsers(context.Background(), tt.query)

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expResponse, resp); diff != "" {
				t.Fatalf("response not equal (-exp, +got):\n%v", diff)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		opts           *notion.SearchOpts
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expPostBody    map[string]interface{}
		expResponse    notion.SearchResponse
		expError       error
	}{
		{
			name: "with query, successful response",
			opts: &notion.SearchOpts{
				Query: "foobar",
				Filter: &notion.SearchFilter{
					Property: "object",
					Value:    "database",
				},
				Sort: &notion.SearchSort{
					Direction: notion.SortDirAsc,
					Timestamp: notion.SearchSortTimestampLastEditedTime,
				},
				StartCursor: "39ddfc9d-33c9-404c-89cf-79f01c42dd0c",
				PageSize:    42,
			},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "list",
						"results": [
							{
								"object": "database",
								"id": "668d797c-76fa-4934-9b05-ad288df2d136",
								"created_time": "2020-03-17T19:10:04.968Z",
								"last_edited_time": "2020-03-17T21:49:37.913Z",
								"title": [
									{
										"type": "text",
										"text": {
											"content": "Foobar",
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
										"plain_text": "Foobar",
										"href": null
									}
								],
								"properties": {
									"Name": {
										"id": "title",
										"type": "title",
										"title": {}
									}
								}
							},
							{
								"object": "page",
								"id": "276ee233-e426-4ed0-9986-6b22af8550df",
								"created_time": "2021-05-19T19:34:05.068Z",
								"last_edited_time": "2021-05-19T19:34:05.069Z",
								"parent": {
									"type": "page_id",
									"page_id": "b0668f48-8d66-4733-9bdb-2f82215707f7"
								},
								"archived": false,
								"properties": {
									"title": {
										"id": "title",
										"type": "title",
										"title": [
											{
												"type": "text",
												"text": {
													"content": "Foobar",
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
												"plain_text": "Foobar",
												"href": null
											}
										]
									}
								}
							}
						],
						"next_cursor": "A^hd",
						"has_more": true
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody: map[string]interface{}{
				"query": "foobar",
				"filter": map[string]interface{}{
					"property": "object",
					"value":    "database",
				},
				"sort": map[string]interface{}{
					"direction": "ascending",
					"timestamp": "last_edited_time",
				},
				"start_cursor": "39ddfc9d-33c9-404c-89cf-79f01c42dd0c",
				"page_size":    float64(42),
			},
			expResponse: notion.SearchResponse{
				Results: notion.SearchResults{
					notion.Database{
						ID:             "668d797c-76fa-4934-9b05-ad288df2d136",
						CreatedTime:    mustParseTime(time.RFC3339, "2020-03-17T19:10:04.968Z"),
						LastEditedTime: mustParseTime(time.RFC3339, "2020-03-17T21:49:37.913Z"),
						Title: []notion.RichText{
							{
								Type: notion.RichTextTypeText,
								Text: &notion.Text{
									Content: "Foobar",
								},
								Annotations: &notion.Annotations{
									Color: notion.ColorDefault,
								},
								PlainText: "Foobar",
							},
						},
						Properties: notion.DatabaseProperties{
							"Name": notion.DatabaseProperty{
								ID:    "title",
								Type:  notion.DBPropTypeTitle,
								Title: &notion.EmptyMetadata{},
							},
						},
					},
					notion.Page{
						ID:             "276ee233-e426-4ed0-9986-6b22af8550df",
						CreatedTime:    mustParseTime(time.RFC3339Nano, "2021-05-19T19:34:05.068Z"),
						LastEditedTime: mustParseTime(time.RFC3339Nano, "2021-05-19T19:34:05.069Z"),
						Parent: notion.Parent{
							Type:   notion.ParentTypePage,
							PageID: "b0668f48-8d66-4733-9bdb-2f82215707f7",
						},
						Properties: notion.PageProperties{
							Title: notion.PageTitle{
								Title: []notion.RichText{
									{
										Type: notion.RichTextTypeText,
										Text: &notion.Text{
											Content: "Foobar",
										},
										Annotations: &notion.Annotations{
											Color: notion.ColorDefault,
										},
										PlainText: "Foobar",
									},
								},
							},
						},
					},
				},
				HasMore:    true,
				NextCursor: notion.StringPtr("A^hd"),
			},
			expError: nil,
		},
		{
			name: "without query, doesn't send POST body",
			opts: nil,
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "list",
						"results": [],
						"next_cursor": null,
						"has_more": false
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody:    nil,
			expResponse: notion.SearchResponse{
				Results:    notion.SearchResults{},
				HasMore:    false,
				NextCursor: nil,
			},
			expError: nil,
		},
		{
			name: "with non nil query, but without fields, omits all fields from POST body",
			opts: &notion.SearchOpts{},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "list",
						"results": [],
						"next_cursor": null,
						"has_more": false
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody:    map[string]interface{}{},
			expResponse: notion.SearchResponse{
				Results:    notion.SearchResults{},
				HasMore:    false,
				NextCursor: nil,
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
			expResponse:    notion.SearchResponse{},
			expError:       errors.New("notion: failed to search: foobar (code: validation_error, status: 400)"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			httpClient := &http.Client{
				Transport: &mockRoundtripper{fn: func(r *http.Request) (*http.Response, error) {
					postBody := make(map[string]interface{})

					err := json.NewDecoder(r.Body).Decode(&postBody)
					if err != nil && err != io.EOF {
						t.Fatal(err)
					}

					if len(tt.expPostBody) == 0 && len(postBody) != 0 {
						t.Errorf("unexpected post body: %+v", postBody)
					}

					if len(tt.expPostBody) != 0 && len(postBody) == 0 {
						t.Errorf("post body not equal (expected %+v, got: nil)", tt.expPostBody)
					}

					if len(tt.expPostBody) != 0 && len(postBody) != 0 {
						if diff := cmp.Diff(tt.expPostBody, postBody); diff != "" {
							t.Errorf("post body not equal (-exp, +got):\n%v", diff)
						}
					}

					return &http.Response{
						StatusCode: tt.respStatusCode,
						Status:     http.StatusText(tt.respStatusCode),
						Body:       ioutil.NopCloser(tt.respBody(r)),
					}, nil
				}},
			}
			client := notion.NewClient("secret-api-key", notion.WithHTTPClient(httpClient))
			resp, err := client.Search(context.Background(), tt.opts)

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expResponse, resp); diff != "" {
				t.Fatalf("response not equal (-exp, +got):\n%v", diff)
			}
		})
	}
}
