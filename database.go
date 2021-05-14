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

// Database property types.
type (
	TitleDatabaseProperty          struct{}
	RichTextDatabaseProperty       struct{}
	DateDatabaseProperty           struct{}
	PeopleDatabaseProperty         struct{}
	FileDatabaseProperty           struct{}
	CheckboxDatabaseProperty       struct{}
	URLDatabaseProperty            struct{}
	EmailDatabaseProperty          struct{}
	PhoneNumberDatabaseProperty    struct{}
	CreatedTimeDatabaseProperty    struct{}
	CreatedByDatabaseProperty      struct{}
	LastEditedTimeDatabaseProperty struct{}
	LastEditedByDatabaseProperty   struct{}

	NumberDatabaseProperty struct {
		Format string `json:"format"`
	}
	SelectDatabaseProperty struct {
		Options []SelectOptions `json:"options"`
	}
	FormulaDatabaseProperty struct {
		Expression string `json:"expression"`
	}
	RelationDatabaseProperty struct {
		DatabaseID     string  `json:"database_id"`
		SyncedPropName *string `json:"synced_property_name"`
		SyncedPropID   *string `json:"synced_property_id"`
	}
	RollupDatabaseProperty struct {
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

	Title          *TitleDatabaseProperty          `json:"title"`
	RichText       *RichTextDatabaseProperty       `json:"rich_text"`
	Number         *NumberDatabaseProperty         `json:"number"`
	Select         *SelectDatabaseProperty         `json:"select"`
	MultiSelect    *SelectDatabaseProperty         `json:"multi_select"`
	Date           *DateDatabaseProperty           `json:"date"`
	People         *PeopleDatabaseProperty         `json:"people"`
	File           *FileDatabaseProperty           `json:"file"`
	Checkbox       *CheckboxDatabaseProperty       `json:"checkbox"`
	URL            *URLDatabaseProperty            `json:"url"`
	Email          *EmailDatabaseProperty          `json:"email"`
	PhoneNumber    *PhoneNumberDatabaseProperty    `json:"phone_number"`
	Formula        *FormulaDatabaseProperty        `json:"formula"`
	Relation       *RelationDatabaseProperty       `json:"relation"`
	Rollup         *RollupDatabaseProperty         `json:"rollup"`
	CreatedTime    *CreatedTimeDatabaseProperty    `json:"created_time"`
	CreatedBy      *CreatedByDatabaseProperty      `json:"created_by"`
	LastEditedTime *LastEditedTimeDatabaseProperty `json:"last_edited_time"`
	LastEditedBy   *LastEditedByDatabaseProperty   `json:"last_edited_by"`
}

func (c *Client) FindDatabaseByID(id string) (db Database, err error) {
	req, err := c.newRequest("GET", "/databases/"+id, nil)
	if err != nil {
		return Database{}, fmt.Errorf("invalid URL: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Database{}, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return Database{}, fmt.Errorf("notion: failed to get database: %w", parseErrorResponse(res))
	}

	err = json.NewDecoder(res.Body).Decode(&db)
	if err != nil {
		return Database{}, fmt.Errorf("failed to parse HTTP response: %w", err)
	}

	return db, nil
}
