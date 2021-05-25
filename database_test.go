package notion_test

import (
	"encoding/json"
	"testing"

	"github.com/dstotijn/go-notion"
)

func TestShouldBeAbleToCreateQueryWithEmptyFilter(t *testing.T) {
	query := notion.DatabaseQuery{
		PageSize: 10,
	}
	b, _ := json.Marshal(query)

	result := make(map[string]interface{}, 0)
	json.Unmarshal(b, &result)
	if _, ok := result["filter"]; ok {
		t.Errorf("Should be able to marshal DatabaseQuery to a JSON that does not contain filter field\n Got %s", string(b))
	}
}
