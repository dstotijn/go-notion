package notion

import (
	"encoding/json"
	"fmt"
	"time"
)

// Database is a resource on the Notion platform.
// See: https://developers.notion.com/reference/database
type Database struct {
	ID             string             `json:"id"`
	CreatedTime    time.Time          `json:"created_time"`
	LastEditedTime time.Time          `json:"last_edited_time"`
	Title          []RichText         `json:"title"`
	Properties     DatabaseProperties `json:"properties"`
}

type DatabaseProperties map[string]DatabaseProperty

// Database property metadata types.
type (
	NumberMetadata struct {
		Format string `json:"format"`
	}
	SelectMetadata struct {
		Options []SelectOptions `json:"options"`
	}
	FormulaMetadata struct {
		Expression string `json:"expression"`
	}
	RelationMetadata struct {
		DatabaseID     string  `json:"database_id"`
		SyncedPropName *string `json:"synced_property_name"`
		SyncedPropID   *string `json:"synced_property_id"`
	}
	RollupMetadata struct {
		RelationPropName string `json:"relation_property_name"`
		RelationPropID   string `json:"relation_property_id"`
		RollupPropName   string `json:"rollup_property_name"`
		RollupPropID     string `json:"rollup_property_id"`
		Function         string `json:"function"`
	}
)

type SelectOptions struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type DatabaseProperty struct {
	ID   string `json:"id"`
	Type string `json:"type"`

	Number      *NumberMetadata   `json:"number"`
	Select      *SelectMetadata   `json:"select"`
	MultiSelect *SelectMetadata   `json:"multi_select"`
	Formula     *FormulaMetadata  `json:"formula"`
	Relation    *RelationMetadata `json:"relation"`
	Rollup      *RollupMetadata   `json:"rollup"`
}

// Metadata returns the underlying property metadata, based on its `type` field.
// When type is unknown/unmapped or doesn't have additional properies, `nil` is returned.
func (prop DatabaseProperty) Metadata() interface{} {
	switch prop.Type {
	case "number":
		return prop.Number
	case "select":
		return prop.Select
	case "multi_select":
		return prop.MultiSelect
	case "formula":
		return prop.Formula
	case "relation":
		return prop.Relation
	case "rollup":
		return prop.Rollup
	default:
		return nil
	}
}

func (c *Client) FindDatabaseByID(id string) (db Database, err error) {
	req, err := c.newRequest("GET", "/databases/"+id, nil)
	if err != nil {
		return Database{}, fmt.Errorf("notion: invalid URL: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Database{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return Database{}, fmt.Errorf("notion: failed to find database: %w", parseErrorResponse(res))
	}

	err = json.NewDecoder(res.Body).Decode(&db)
	if err != nil {
		return Database{}, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return db, nil
}
