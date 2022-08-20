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
	"github.com/google/go-cmp/cmp/cmpopts"
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
						"created_by": {
							"object": "user",
							"id": "71e95936-2737-4e11-b03d-f174f6f13087"
						},
						"last_edited_by": {
							"object": "user",
							"id": "5ba97cc9-e5e0-4363-b33a-1d80a635577f"
						},
						"url": "https://www.notion.so/668d797c76fa49349b05ad288df2d136",
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
									"type": "dual_property",
									"dual_property": {
										"synced_property_name": "Related to Test database (Relation Test)",
										"synced_property_id": "IJi<"
									}
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
									"function": "count_all"
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
				CreatedBy: notion.BaseUser{
					ID: "71e95936-2737-4e11-b03d-f174f6f13087",
				},
				LastEditedBy: notion.BaseUser{
					ID: "5ba97cc9-e5e0-4363-b33a-1d80a635577f",
				},
				URL: "https://www.notion.so/668d797c76fa49349b05ad288df2d136",
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
							DatabaseID: "668d797c-76fa-4934-9b05-ad288df2d136",
							Type:       notion.RelationTypeDualProperty,
							DualProperty: &notion.DualPropertyRelation{
								SyncedPropID:   "IJi<",
								SyncedPropName: "Related to Test database (Relation Test)",
							},
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
							Function:         notion.RollupFunctionCountAll,
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
					DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
						RichText: &notion.TextPropertyFilter{
							Contains: "foobar",
						},
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
									"Name": {
										"id": "title"
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
					"rich_text": map[string]interface{}{
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
						Properties: notion.PageProperties{
							"Name": notion.PagePropertyID{
								ID: "title",
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
				Description: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Lorem ipsum dolor sit amet.",
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
				Cover: &notion.Cover{
					Type: notion.FileTypeExternal,
					External: &notion.FileExternal{
						URL: "https://example.com/image.png",
					},
				},
				IsInline: true,
			},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "database",
						"id": "b89664e3-30b4-474a-9cce-c72a4827d1e4",
						"created_time": "2021-07-20T20:09:00.000Z",
						"last_edited_time": "2021-07-20T20:09:00.000Z",
						"url": "https://www.notion.so/b89664e330b4474a9ccec72a4827d1e4",
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
						"description": [
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
						},
						"cover": {
							"type": "external",
							"external": {
								"url": "https://example.com/image.png"
							}
						},
						"is_inline": true
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
				"description": []interface{}{
					map[string]interface{}{
						"text": map[string]interface{}{
							"content": "Lorem ipsum dolor sit amet.",
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
				"cover": map[string]interface{}{
					"type": "external",
					"external": map[string]interface{}{
						"url": "https://example.com/image.png",
					},
				},
				"is_inline": true,
			},
			expResponse: notion.Database{
				ID:             "b89664e3-30b4-474a-9cce-c72a4827d1e4",
				CreatedTime:    mustParseTime(time.RFC3339Nano, "2021-07-20T20:09:00Z"),
				LastEditedTime: mustParseTime(time.RFC3339Nano, "2021-07-20T20:09:00Z"),
				URL:            "https://www.notion.so/b89664e330b4474a9ccec72a4827d1e4",
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
				Description: []notion.RichText{
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
				Properties: notion.DatabaseProperties{
					"Title": notion.DatabaseProperty{
						ID:    "title",
						Type:  notion.DBPropTypeTitle,
						Title: &notion.EmptyMetadata{},
					},
				},
				Icon: &notion.Icon{
					Type:  notion.IconTypeEmoji,
					Emoji: notion.StringPtr("✌️"),
				},
				Cover: &notion.Cover{
					Type: notion.FileTypeExternal,
					External: &notion.FileExternal{
						URL: "https://example.com/image.png",
					},
				},
				IsInline: true,
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

func TestUpdateDatabase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		params         notion.UpdateDatabaseParams
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expPostBody    map[string]interface{}
		expResponse    notion.Database
		expError       error
	}{
		{
			name: "successful response",
			params: notion.UpdateDatabaseParams{
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Updated title",
						},
					},
				},
				Description: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Updated description.",
						},
					},
				},
				Properties: map[string]*notion.DatabaseProperty{
					"New": {
						Type:     notion.DBPropTypeRichText,
						RichText: &notion.EmptyMetadata{},
					},
					"Removed": nil,
				},
				Icon: &notion.Icon{
					Type:  notion.IconTypeEmoji,
					Emoji: notion.StringPtr("✌️"),
				},
				Cover: &notion.Cover{
					Type: notion.FileTypeExternal,
					External: &notion.FileExternal{
						URL: "https://example.com/image.png",
					},
				},
				IsInline: notion.BoolPtr(true),
			},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "database",
						"id": "668d797c-76fa-4934-9b05-ad288df2d136",
						"created_time": "2020-03-17T19:10:04.968Z",
						"last_edited_time": "2020-03-17T21:49:37.913Z",
						"url": "https://www.notion.so/668d797c76fa49349b05ad288df2d136",
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
						"description": [
							{
								"type": "text",
								"text": {
									"content": "Updated description.",
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
								"plain_text": "Updated description.",
								"href": null
							}
						],
						"properties": {
							"Name": {
								"id": "title",
								"type": "title",
								"title": {}
							},
							"New": {
								"id": "J@cS",
								"type": "rich_text",
								"text": {}
							}
						},
						"parent": {
							"type": "page_id",
							"page_id": "b8595b75-abd1-4cad-8dfe-f935a8ef57cb"
						},
						"icon": {
							"type": "emoji",
							"emoji": "✌️"
						},
						"cover": {
							"type": "external",
							"external": {
								"url": "https://example.com/image.png"
							}
						},
						"is_inline": true
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody: map[string]interface{}{
				"title": []interface{}{
					map[string]interface{}{
						"text": map[string]interface{}{
							"content": "Updated title",
						},
					},
				},
				"description": []interface{}{
					map[string]interface{}{
						"text": map[string]interface{}{
							"content": "Updated description.",
						},
					},
				},
				"properties": map[string]interface{}{
					"New": map[string]interface{}{
						"type":      "rich_text",
						"rich_text": map[string]interface{}{},
					},
					"Removed": nil,
				},
				"icon": map[string]interface{}{
					"type":  "emoji",
					"emoji": "✌️",
				},
				"cover": map[string]interface{}{
					"type": "external",
					"external": map[string]interface{}{
						"url": "https://example.com/image.png",
					},
				},
				"is_inline": true,
			},
			expResponse: notion.Database{
				ID:             "668d797c-76fa-4934-9b05-ad288df2d136",
				CreatedTime:    mustParseTime(time.RFC3339, "2020-03-17T19:10:04.968Z"),
				LastEditedTime: mustParseTime(time.RFC3339, "2020-03-17T21:49:37.913Z"),
				URL:            "https://www.notion.so/668d797c76fa49349b05ad288df2d136",
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
				Description: []notion.RichText{
					{
						Type: notion.RichTextTypeText,
						Text: &notion.Text{
							Content: "Updated description.",
						},
						Annotations: &notion.Annotations{
							Color: notion.ColorDefault,
						},
						PlainText: "Updated description.",
					},
				},
				Properties: notion.DatabaseProperties{
					"Name": notion.DatabaseProperty{
						ID:    "title",
						Type:  notion.DBPropTypeTitle,
						Title: &notion.EmptyMetadata{},
					},
					"New": notion.DatabaseProperty{
						ID:   "J@cS",
						Type: notion.DBPropTypeRichText,
					},
				},
				Parent: notion.Parent{
					Type:   notion.ParentTypePage,
					PageID: "b8595b75-abd1-4cad-8dfe-f935a8ef57cb",
				},
				Icon: &notion.Icon{
					Type:  notion.IconTypeEmoji,
					Emoji: notion.StringPtr("✌️"),
				},
				Cover: &notion.Cover{
					Type: notion.FileTypeExternal,
					External: &notion.FileExternal{
						URL: "https://example.com/image.png",
					},
				},
				IsInline: true,
			},
			expError: nil,
		},
		{
			name: "error response",
			params: notion.UpdateDatabaseParams{
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Updated title",
						},
					},
				},
				Properties: map[string]*notion.DatabaseProperty{
					"New": {
						Type:     notion.DBPropTypeRichText,
						RichText: &notion.EmptyMetadata{},
					},
					"Removed": nil,
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
				"title": []interface{}{
					map[string]interface{}{
						"text": map[string]interface{}{
							"content": "Updated title",
						},
					},
				},
				"properties": map[string]interface{}{
					"New": map[string]interface{}{
						"type":      "rich_text",
						"rich_text": map[string]interface{}{},
					},
					"Removed": nil,
				},
			},
			expResponse: notion.Database{},
			expError:    errors.New("notion: failed to update database: foobar (code: validation_error, status: 400)"),
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
			updatedDB, err := client.UpdateDatabase(context.Background(), "00000000-0000-0000-0000-000000000000", tt.params)

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expResponse, updatedDB); diff != "" {
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
						"created_by": {
							"object": "user",
							"id": "71e95936-2737-4e11-b03d-f174f6f13087"
						},
						"last_edited_time": "2021-05-19T18:34:00.000Z",
						"last_edited_by": {
							"object": "user",
							"id": "5ba97cc9-e5e0-4363-b33a-1d80a635577f"
						},
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
				ID:          "606ed832-7d79-46de-bbed-5b4896e7bc02",
				CreatedTime: mustParseTime(time.RFC3339Nano, "2021-05-19T18:34:00.000Z"),
				CreatedBy: &notion.BaseUser{
					ID: "71e95936-2737-4e11-b03d-f174f6f13087",
				},
				LastEditedTime: mustParseTime(time.RFC3339Nano, "2021-05-19T18:34:00.000Z"),
				LastEditedBy: &notion.BaseUser{
					ID: "5ba97cc9-e5e0-4363-b33a-1d80a635577f",
				},
				URL: "https://www.notion.so/Avocado-251d2b5f268c4de2afe9c71ff92ca95c",
				Parent: notion.Parent{
					Type:   notion.ParentTypePage,
					PageID: "b0668f48-8d66-4733-9bdb-2f82215707f7",
				},
				Properties: notion.PageProperties{
					"title": notion.PagePropertyID{
						ID: "title",
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
					&notion.ParagraphBlock{
						RichText: []notion.RichText{
							{
								Text: &notion.Text{
									Content: "Lorem ipsum dolor sit amet.",
								},
							},
						},
					},
				},
				Icon: &notion.Icon{
					Type: notion.IconTypeExternal,
					External: &notion.FileExternal{
						URL: "https://example.com/icon.png",
					},
				},
				Cover: &notion.Cover{
					Type: notion.FileTypeExternal,
					External: &notion.FileExternal{
						URL: "https://example.com/cover.png",
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
								"id": "title"
							}
						},
						"icon": {
							"type": "external",
							"external": {
								"url": "https://example.com/icon.png"
							}
						},
						"cover": {
							"type": "external",
							"external": {
								"url": "https://example.com/cover.png"
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
						"paragraph": map[string]interface{}{
							"rich_text": []interface{}{
								map[string]interface{}{
									"text": map[string]interface{}{
										"content": "Lorem ipsum dolor sit amet.",
									},
								},
							},
						},
					},
				},
				"icon": map[string]interface{}{
					"type": "external",
					"external": map[string]interface{}{
						"url": "https://example.com/icon.png",
					},
				},
				"cover": map[string]interface{}{
					"type": "external",
					"external": map[string]interface{}{
						"url": "https://example.com/cover.png",
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
					"title": notion.PagePropertyID{
						ID: "title",
					},
				},
				Icon: &notion.Icon{
					Type: notion.IconTypeExternal,
					External: &notion.FileExternal{
						URL: "https://example.com/icon.png",
					},
				},
				Cover: &notion.Cover{
					Type: notion.FileTypeExternal,
					External: &notion.FileExternal{
						URL: "https://example.com/cover.png",
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
					&notion.ParagraphBlock{
						RichText: []notion.RichText{
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
						"paragraph": map[string]interface{}{
							"rich_text": []interface{}{
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
				Properties: notion.PageProperties{
					"title": notion.PagePropertyID{
						ID: "title",
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

func TestUpdatePage(t *testing.T) {
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
				DatabasePageProperties: notion.DatabasePageProperties{
					"Name": notion.DatabasePageProperty{
						Title: []notion.RichText{
							{
								Text: &notion.Text{
									Content: "Foobar",
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
					"Name": map[string]interface{}{
						"title": []interface{}{
							map[string]interface{}{
								"text": map[string]interface{}{
									"content": "Foobar",
								},
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
					"title": notion.PagePropertyID{
						ID: "title",
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
					External: &notion.FileExternal{
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
					External: &notion.FileExternal{
						URL: "https://www.notion.so/front-static/pages/pricing/pro.png",
					},
				},
				Properties: notion.PageProperties{
					"title": notion.PagePropertyID{
						ID: "title",
					},
				},
			},
			expError: nil,
		},
		{
			name: "page archived, successful response",
			params: notion.UpdatePageParams{
				Archived: notion.BoolPtr(true),
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
						"cover": {
							"type": "external",
							"external": {
								"url": "https://example.com/image.png"
							}
						},
						"archived": true,
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
				"archived": true,
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
				Archived: true,
				Cover: &notion.Cover{
					Type: notion.FileTypeExternal,
					External: &notion.FileExternal{
						URL: "https://example.com/image.png",
					},
				},
				Properties: notion.PageProperties{
					"title": notion.PagePropertyID{
						ID: "title",
					},
				},
			},
			expError: nil,
		},
		{
			name: "page cover, successful response",
			params: notion.UpdatePageParams{
				Cover: &notion.Cover{
					Type: notion.FileTypeExternal,
					External: &notion.FileExternal{
						URL: "https://example.com/image.png",
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
						"cover": {
							"type": "external",
							"external": {
								"url": "https://example.com/image.png"
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
				"cover": map[string]interface{}{
					"type": "external",
					"external": map[string]interface{}{
						"url": "https://example.com/image.png",
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
				Cover: &notion.Cover{
					Type: notion.FileTypeExternal,
					External: &notion.FileExternal{
						URL: "https://example.com/image.png",
					},
				},
				Properties: notion.PageProperties{
					"title": notion.PagePropertyID{
						ID: "title",
					},
				},
			},
			expError: nil,
		},
		{
			name: "error response",
			params: notion.UpdatePageParams{
				DatabasePageProperties: notion.DatabasePageProperties{
					"Name": notion.DatabasePageProperty{
						Title: []notion.RichText{
							{
								Text: &notion.Text{
									Content: "Foobar",
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
				"properties": map[string]interface{}{
					"Name": map[string]interface{}{
						"title": []interface{}{
							map[string]interface{}{
								"text": map[string]interface{}{
									"content": "Foobar",
								},
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
			expError:    errors.New("notion: invalid page params: at least one of database page properties, archived, icon or cover is required"),
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
			page, err := client.UpdatePage(context.Background(), "00000000-0000-0000-0000-000000000000", tt.params)

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

func TestFindPagePropertyByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		query          *notion.PaginationQuery
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expQueryParams url.Values
		expResponse    notion.PagePropResponse
		expError       error
	}{
		{
			name: "paginated property item, with query, successful response",
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
								"object": "property_item",
								"type": "rich_text",
								"rich_text": {
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
			expResponse: notion.PagePropResponse{
				Results: []notion.PagePropItem{
					{
						Type: notion.DBPropTypeRichText,
						RichText: notion.RichText{
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
				HasMore:    true,
				NextCursor: "A^hd",
			},
			expError: nil,
		},
		{
			name:  "paginated property item, successful response",
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
			expResponse: notion.PagePropResponse{
				Results:    []notion.PagePropItem{},
				HasMore:    false,
				NextCursor: "",
			},
			expError: nil,
		},
		{
			name:  "simple property item, successful response",
			query: nil,
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "property_item",
						"type": "number",
						"number": 42
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expQueryParams: nil,
			expResponse: notion.PagePropResponse{
				PagePropItem: notion.PagePropItem{
					Type:   notion.DBPropTypeNumber,
					Number: 42,
				},
			},
			expError: nil,
		},
		{
			name:  "rollup property item with aggregation, successful response",
			query: nil,
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "list",
						"results": [
							{
								"object": "property_item",
								"type": "relation",
								"relation": {
									"id": "de5d73e8-3748-40fa-9102-f1290fe2444b"
								}
							},
							{
								"object": "property_item",
								"type": "relation",
								"relation": {
									"id": "164325b0-4c9e-416b-ba9c-037b4c9acdfd"
								}
							},
							{
								"object": "property_item",
								"type": "relation",
								"relation": {
									"id": "456baa98-3239-4c1f-b0ea-bdae945aaf33"
								}
							}
						],
						"has_more": true,
						"type": "property_item",
						"property_item": {
							"id": "aBcD123",
							"next_url": "https://api.notion.com/v1/pages/b55c9c91-384d-452b-81db-d1ef79372b75/properties/aBcD123?start_cursor=some-next-cursor-value",
							"type": "rollup",
							"rollup": {
								"type": "date",
								"date": {
									"start": "2021-10-07T14:42:00.000+00:00",
									"end": null
								},
								"function": "latest_date"
							}
						}
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expQueryParams: nil,
			expResponse: notion.PagePropResponse{
				PagePropItem: notion.PagePropItem{
					Type: notion.DBPropTypePropertyItem,
				},
				PropertyItem: notion.PagePropListItem{
					ID:   "aBcD123",
					Type: notion.DBPropTypeRollup,
					Rollup: notion.RollupResult{
						Type: notion.RollupResultTypeDate,
						Date: &notion.Date{
							Start: mustParseDateTime("2021-10-07T14:42:00.000+00:00"),
						},
					},
					NextURL: "https://api.notion.com/v1/pages/b55c9c91-384d-452b-81db-d1ef79372b75/properties/aBcD123?start_cursor=some-next-cursor-value",
				},
				HasMore: true,
				Results: []notion.PagePropItem{
					{
						Type: notion.DBPropTypeRelation,
						Relation: notion.Relation{
							ID: "de5d73e8-3748-40fa-9102-f1290fe2444b",
						},
					},
					{
						Type: notion.DBPropTypeRelation,
						Relation: notion.Relation{
							ID: "164325b0-4c9e-416b-ba9c-037b4c9acdfd",
						},
					},
					{
						Type: notion.DBPropTypeRelation,
						Relation: notion.Relation{
							ID: "456baa98-3239-4c1f-b0ea-bdae945aaf33",
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
						"status": 400,
						"code": "validation_error",
						"message": "foobar"
					}`,
				)
			},
			respStatusCode: http.StatusBadRequest,
			expResponse:    notion.PagePropResponse{},
			expError:       errors.New("notion: failed to find page property: foobar (code: validation_error, status: 400)"),
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
			resp, err := client.FindPagePropertyByID(context.Background(), "page-id", "prop-id", tt.query)

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

func TestFindBlockChildrenById(t *testing.T) {
	t.Parallel()

	type blockFields struct {
		id             string
		createdTime    time.Time
		lastEditedTime time.Time
		hasChildren    bool
		archived       bool
	}

	tests := []struct {
		name           string
		query          *notion.PaginationQuery
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expQueryParams url.Values
		expResponse    notion.BlockChildrenResponse
		expBlockFields []blockFields
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
									"rich_text": [
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
					&notion.ParagraphBlock{
						RichText: []notion.RichText{
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
				HasMore:    true,
				NextCursor: notion.StringPtr("A^hd"),
			},
			expBlockFields: []blockFields{
				{
					id:             "ae9c9a31-1c1e-4ae2-a5ee-c539a2d43113",
					createdTime:    mustParseTime(time.RFC3339, "2021-05-14T09:15:00.000Z"),
					lastEditedTime: mustParseTime(time.RFC3339, "2021-05-14T09:15:00.000Z"),
					hasChildren:    false,
					archived:       false,
				},
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

			if diff := cmp.Diff(tt.expResponse, resp, cmpopts.IgnoreUnexported(notion.ParagraphBlock{})); diff != "" {
				t.Fatalf("response not equal (-exp, +got):\n%v", diff)
			}

			if len(tt.expBlockFields) != len(resp.Results) {
				t.Fatalf("expected %v result(s), got %v", len(tt.expBlockFields), len(resp.Results))
			}

			for i, exp := range tt.expBlockFields {
				if exp.id != resp.Results[i].ID() {
					t.Fatalf("id not equal (expected: %v, got: %v)", exp.id, resp.Results[i].ID())
				}

				if exp.createdTime != resp.Results[i].CreatedTime() {
					t.Fatalf("createdTime not equal (expected: %v, got: %v)", exp.createdTime, resp.Results[i].CreatedTime())
				}

				if exp.lastEditedTime != resp.Results[i].LastEditedTime() {
					t.Fatalf("lastEditedTime not equal (expected: %v, got: %v)", exp.lastEditedTime, resp.Results[i].LastEditedTime())
				}

				if exp.hasChildren != resp.Results[i].HasChildren() {
					t.Fatalf("hasChildren not equal (expected: %v, got: %v)", exp.hasChildren, resp.Results[i].HasChildren())
				}

				if exp.archived != resp.Results[i].Archived() {
					t.Fatalf("archived not equal (expected: %v, got: %v)", exp.archived, resp.Results[i].Archived())
				}
			}
		})
	}
}

func TestAppendBlockChildren(t *testing.T) {
	t.Parallel()

	type blockFields struct {
		id             string
		createdTime    time.Time
		lastEditedTime time.Time
		hasChildren    bool
		archived       bool
	}

	tests := []struct {
		name           string
		children       []notion.Block
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expPostBody    map[string]interface{}
		expResponse    notion.BlockChildrenResponse
		expBlockFields []blockFields
		expError       error
	}{
		{
			name: "successful response",
			children: []notion.Block{
				&notion.ParagraphBlock{
					RichText: []notion.RichText{
						{
							Text: &notion.Text{
								Content: "Lorem ipsum dolor sit amet.",
							},
						},
					},
				},
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
									"rich_text": [
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
			expPostBody: map[string]interface{}{
				"children": []interface{}{
					map[string]interface{}{
						"paragraph": map[string]interface{}{
							"rich_text": []interface{}{
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
			expResponse: notion.BlockChildrenResponse{
				Results: []notion.Block{
					&notion.ParagraphBlock{
						RichText: []notion.RichText{
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
				HasMore:    true,
				NextCursor: notion.StringPtr("A^hd"),
			},
			expBlockFields: []blockFields{
				{
					id:             "ae9c9a31-1c1e-4ae2-a5ee-c539a2d43113",
					createdTime:    mustParseTime(time.RFC3339, "2021-05-14T09:15:00.000Z"),
					lastEditedTime: mustParseTime(time.RFC3339, "2021-05-14T09:15:00.000Z"),
					hasChildren:    false,
					archived:       false,
				},
			},
			expError: nil,
		},
		{
			name: "error response",
			children: []notion.Block{
				&notion.ParagraphBlock{
					RichText: []notion.RichText{
						{
							Text: &notion.Text{
								Content: "Lorem ipsum dolor sit amet.",
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
						"paragraph": map[string]interface{}{
							"rich_text": []interface{}{
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
			expResponse: notion.BlockChildrenResponse{},
			expError:    errors.New("notion: failed to append block children: foobar (code: validation_error, status: 400)"),
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
			resp, err := client.AppendBlockChildren(context.Background(), "00000000-0000-0000-0000-000000000000", tt.children)

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expResponse, resp, cmpopts.IgnoreUnexported(notion.ParagraphBlock{})); diff != "" {
				t.Fatalf("response not equal (-exp, +got):\n%v", diff)
			}

			if len(tt.expBlockFields) != len(resp.Results) {
				t.Fatalf("expected %v result(s), got %v", len(tt.expBlockFields), len(resp.Results))
			}

			for i, exp := range tt.expBlockFields {
				if exp.id != resp.Results[i].ID() {
					t.Fatalf("id not equal (expected: %v, got: %v)", exp.id, resp.Results[i].ID())
				}

				if exp.createdTime != resp.Results[i].CreatedTime() {
					t.Fatalf("createdTime not equal (expected: %v, got: %v)", exp.createdTime, resp.Results[i].CreatedTime())
				}

				if exp.lastEditedTime != resp.Results[i].LastEditedTime() {
					t.Fatalf("lastEditedTime not equal (expected: %v, got: %v)", exp.lastEditedTime, resp.Results[i].LastEditedTime())
				}

				if exp.hasChildren != resp.Results[i].HasChildren() {
					t.Fatalf("hasChildren not equal (expected: %v, got: %v)", exp.hasChildren, resp.Results[i].HasChildren())
				}

				if exp.archived != resp.Results[i].Archived() {
					t.Fatalf("archived not equal (expected: %v, got: %v)", exp.archived, resp.Results[i].Archived())
				}
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
				BaseUser: notion.BaseUser{
					ID: "be32e790-8292-46df-a248-b784fdf483cf",
				},
				Name:      "Jane Doe",
				AvatarURL: "https://example.com/avatar.png",
				Type:      notion.UserTypePerson,
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

func TestFindCurrentUser(t *testing.T) {
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
						"type": "bot",
						"bot": {
							"owner": {
								"type": "user",
								"user": {
									"object": "user",
									"id": "5389a034-eb5c-47b5-8a9e-f79c99ef166c",
									"name": "Jane Doe",
									"avatar_url": "https://example.com/avatar.png",
									"type": "person",
									"person": {
										"email": "jane@example.com"
									}
								}
							}
						}
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expUser: notion.User{
				BaseUser: notion.BaseUser{
					ID: "be32e790-8292-46df-a248-b784fdf483cf",
				},
				Type: notion.UserTypeBot,
				Bot: &notion.Bot{
					Owner: notion.BotOwner{
						Type: notion.BotOwnerTypeUser,
						User: &notion.User{
							BaseUser: notion.BaseUser{
								ID: "5389a034-eb5c-47b5-8a9e-f79c99ef166c",
							},
							Name:      "Jane Doe",
							AvatarURL: "https://example.com/avatar.png",
							Type:      notion.UserTypePerson,
							Person: &notion.Person{
								Email: "jane@example.com",
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
			expUser:        notion.User{},
			expError:       errors.New("notion: failed to find current user: foobar (code: object_not_found, status: 404)"),
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
			user, err := client.FindCurrentUser(context.Background())

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
						BaseUser: notion.BaseUser{
							ID: "be32e790-8292-46df-a248-b784fdf483cf",
						},
						Name:      "Jane Doe",
						AvatarURL: "https://example.com/avatar.png",
						Type:      notion.UserTypePerson,
						Person: &notion.Person{
							Email: "jane@example.com",
						},
					},
					{
						BaseUser: notion.BaseUser{
							ID: "25c9cc08-1afd-4d22-b9e6-31b0f6e7b44f",
						},
						Name: "Johnny 5",
						Type: notion.UserTypeBot,
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
								"url": "https://www.notion.so/668d797c76fa49349b05ad288df2d136",
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
						URL:            "https://www.notion.so/668d797c76fa49349b05ad288df2d136",
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
							"title": notion.PagePropertyID{
								ID: "title",
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

func TestFindBlockByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		blockID           string
		respBody          func(r *http.Request) io.Reader
		respStatusCode    int
		expBlock          notion.Block
		expID             string
		expParent         notion.Parent
		expCreatedTime    time.Time
		expCreatedBy      notion.BaseUser
		expLastEditedTime time.Time
		expLastEditedBy   notion.BaseUser
		expHasChildren    bool
		expArchived       bool
		expError          error
	}{
		{
			name:    "successful response",
			blockID: "test-block-id",
			respBody: func(r *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "block",
						"id": "048e165e-352d-4119-8128-e46c3527d95c",
						"parent": {
							"type": "page_id",
							"page_id": "59833787-2cf9-4fdf-8782-e53db20768a5"
						},
						"created_time": "2021-10-02T06:09:00.000Z",
						"created_by": {
							"object": "user",
							"id": "71e95936-2737-4e11-b03d-f174f6f13087"
						},
						"last_edited_time": "2021-10-02T06:31:00.000Z",
						"last_edited_by": {
							"object": "user",
							"id": "5ba97cc9-e5e0-4363-b33a-1d80a635577f"
						},
						"has_children": true,
						"archived": false,
						"type": "child_page",
						"child_page": {
							"title": "test title"
						}
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expBlock: &notion.ChildPageBlock{
				Title: "test title",
			},
			expID: "048e165e-352d-4119-8128-e46c3527d95c",
			expParent: notion.Parent{
				Type:   notion.ParentTypePage,
				PageID: "59833787-2cf9-4fdf-8782-e53db20768a5",
			},
			expCreatedTime: mustParseTime(time.RFC3339, "2021-10-02T06:09:00Z"),
			expCreatedBy: notion.BaseUser{
				ID: "71e95936-2737-4e11-b03d-f174f6f13087",
			},
			expLastEditedTime: mustParseTime(time.RFC3339, "2021-10-02T06:31:00Z"),
			expLastEditedBy: notion.BaseUser{
				ID: "5ba97cc9-e5e0-4363-b33a-1d80a635577f",
			},
			expHasChildren: true,
			expArchived:    false,
			expError:       nil,
		},
		{
			name: "error response not found",
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "error",
						"status": 404,
						"code": "object_not_found",
						"message": "Could not find block with ID: test id."
					}`,
				)
			},
			respStatusCode: http.StatusNotFound,
			expBlock:       nil,
			expError:       errors.New("notion: failed to find block: Could not find block with ID: test id. (code: object_not_found, status: 404)"),
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
			block, err := client.FindBlockByID(context.Background(), tt.blockID)

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expBlock, block, cmpopts.IgnoreUnexported(notion.ChildPageBlock{})); diff != "" {
				t.Fatalf("user not equal (-exp, +got):\n%v", diff)
			}

			if block != nil {
				if tt.expID != block.ID() {
					t.Fatalf("id not equal (expected: %v, got: %v)", tt.expID, block.ID())
				}

				if tt.expParent != block.Parent() {
					t.Fatalf("parent not equal (expected: %+v, got: %+v)", tt.expParent, block.Parent())
				}

				if tt.expCreatedTime != block.CreatedTime() {
					t.Fatalf("createdTime not equal (expected: %v, got: %v)", tt.expCreatedTime, block.CreatedTime())
				}

				if tt.expCreatedBy != block.CreatedBy() {
					t.Fatalf("createdBy not equal (expected: %v, got: %v)", tt.expCreatedBy, block.CreatedBy())
				}

				if tt.expLastEditedTime != block.LastEditedTime() {
					t.Fatalf("lastEditedTime not equal (expected: %v, got: %v)", tt.expLastEditedTime, block.LastEditedTime())
				}

				if tt.expLastEditedBy != block.LastEditedBy() {
					t.Fatalf("lastEditedBy not equal (expected: %v, got: %v)", tt.expLastEditedBy, block.LastEditedBy())
				}

				if tt.expHasChildren != block.HasChildren() {
					t.Fatalf("hasChildren not equal (expected: %v, got: %v)", tt.expHasChildren, block.HasChildren())
				}

				if tt.expArchived != block.Archived() {
					t.Fatalf("archived not equal (expected: %v, got: %v)", tt.expArchived, block.Archived())
				}
			}
		})
	}
}

func TestUpdateBlock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		block             notion.Block
		respBody          func(r *http.Request) io.Reader
		respStatusCode    int
		expPostBody       map[string]interface{}
		expResponse       notion.Block
		expID             string
		expCreatedTime    time.Time
		expLastEditedTime time.Time
		expHasChildren    bool
		expArchived       bool
		expError          error
	}{
		{
			name: "successful response",
			block: &notion.ParagraphBlock{
				RichText: []notion.RichText{
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
						"object": "block",
						"id": "048e165e-352d-4119-8128-e46c3527d95c",
						"created_time": "2021-10-02T06:09:00.000Z",
						"last_edited_time": "2021-10-02T06:31:00.000Z",
						"has_children": true,
						"archived": false,
						"type": "paragraph",
						"paragraph": {
							"rich_text": [
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
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody: map[string]interface{}{
				"paragraph": map[string]interface{}{
					"rich_text": []interface{}{
						map[string]interface{}{
							"text": map[string]interface{}{
								"content": "Foobar",
							},
						},
					},
				},
			},
			expResponse: &notion.ParagraphBlock{
				RichText: []notion.RichText{
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
			expID:             "048e165e-352d-4119-8128-e46c3527d95c",
			expCreatedTime:    mustParseTime(time.RFC3339, "2021-10-02T06:09:00Z"),
			expLastEditedTime: mustParseTime(time.RFC3339, "2021-10-02T06:31:00Z"),
			expHasChildren:    true,
			expArchived:       false,
			expError:          nil,
		},
		{
			name: "error response",
			block: &notion.ParagraphBlock{
				RichText: []notion.RichText{
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
				"paragraph": map[string]interface{}{
					"rich_text": []interface{}{
						map[string]interface{}{
							"text": map[string]interface{}{
								"content": "Foobar",
							},
						},
					},
				},
			},
			expResponse: nil,
			expError:    errors.New("notion: failed to update block: foobar (code: validation_error, status: 400)"),
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
			updatedBlock, err := client.UpdateBlock(context.Background(), "00000000-0000-0000-0000-000000000000", tt.block)

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expResponse, updatedBlock, cmpopts.IgnoreUnexported(notion.ParagraphBlock{})); diff != "" {
				t.Fatalf("response not equal (-exp, +got):\n%v", diff)
			}

			if updatedBlock != nil {
				if tt.expID != updatedBlock.ID() {
					t.Fatalf("id not equal (expected: %v, got: %v)", tt.expID, updatedBlock.ID())
				}

				if tt.expCreatedTime != updatedBlock.CreatedTime() {
					t.Fatalf("createdTime not equal (expected: %v, got: %v)", tt.expCreatedTime, updatedBlock.CreatedTime())
				}

				if tt.expLastEditedTime != updatedBlock.LastEditedTime() {
					t.Fatalf("lastEditedTime not equal (expected: %v, got: %v)", tt.expLastEditedTime, updatedBlock.LastEditedTime())
				}

				if tt.expHasChildren != updatedBlock.HasChildren() {
					t.Fatalf("hasChildren not equal (expected: %v, got: %v)", tt.expHasChildren, updatedBlock.HasChildren())
				}

				if tt.expArchived != updatedBlock.Archived() {
					t.Fatalf("archived not equal (expected: %v, got: %v)", tt.expArchived, updatedBlock.Archived())
				}
			}
		})
	}
}

func TestDeleteBlock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		respBody          func(r *http.Request) io.Reader
		respStatusCode    int
		expResponse       notion.Block
		expID             string
		expCreatedTime    time.Time
		expLastEditedTime time.Time
		expHasChildren    bool
		expArchived       bool
		expError          error
	}{
		{
			name: "successful response",
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "block",
						"id": "048e165e-352d-4119-8128-e46c3527d95c",
						"created_time": "2021-10-02T06:09:00.000Z",
						"last_edited_time": "2021-10-02T06:31:00.000Z",
						"has_children": true,
						"archived": true,
						"type": "paragraph",
						"paragraph": {
							"rich_text": [
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
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expResponse: &notion.ParagraphBlock{
				RichText: []notion.RichText{
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
			expID:             "048e165e-352d-4119-8128-e46c3527d95c",
			expCreatedTime:    mustParseTime(time.RFC3339, "2021-10-02T06:09:00Z"),
			expLastEditedTime: mustParseTime(time.RFC3339, "2021-10-02T06:31:00Z"),
			expHasChildren:    true,
			expArchived:       true,
			expError:          nil,
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
			expResponse:    nil,
			expError:       errors.New("notion: failed to delete block: foobar (code: validation_error, status: 400)"),
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
			deletedBlock, err := client.DeleteBlock(context.Background(), "00000000-0000-0000-0000-000000000000")

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expResponse, deletedBlock, cmpopts.IgnoreUnexported(notion.ParagraphBlock{})); diff != "" {
				t.Fatalf("response not equal (-exp, +got):\n%v", diff)
			}

			if deletedBlock != nil {
				if tt.expID != deletedBlock.ID() {
					t.Fatalf("id not equal (expected: %v, got: %v)", tt.expID, deletedBlock.ID())
				}

				if tt.expCreatedTime != deletedBlock.CreatedTime() {
					t.Fatalf("createdTime not equal (expected: %v, got: %v)", tt.expCreatedTime, deletedBlock.CreatedTime())
				}

				if tt.expLastEditedTime != deletedBlock.LastEditedTime() {
					t.Fatalf("lastEditedTime not equal (expected: %v, got: %v)", tt.expLastEditedTime, deletedBlock.LastEditedTime())
				}

				if tt.expHasChildren != deletedBlock.HasChildren() {
					t.Fatalf("hasChildren not equal (expected: %v, got: %v)", tt.expHasChildren, deletedBlock.HasChildren())
				}

				if tt.expArchived != deletedBlock.Archived() {
					t.Fatalf("archived not equal (expected: %v, got: %v)", tt.expArchived, deletedBlock.Archived())
				}
			}
		})
	}
}
