package entries

import (
	"encoding/json"
	"errors"

	"github.com/xeipuuv/gojsonschema"
)

type AliasEntry struct {
	Type   string `json:"type"`
	Target string `json:"target"`
}

const aliasEntryJSONSchema = `
{
    "$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"properties": {
		"format_version": {
			"const": 1
		},
		"type": {
			"type": "string",
			"const": "alias"
		},
		"target": {
			"type": "string"
		}
	},
	"required": [
		"format_version",
		"type",
		"target"
	],
	"additionalProperties": false
}
`

// NewAliasEntryFromJSON creates a new entry from JSON.
func NewAliasEntryFromJSON(jsonBytes []byte) (AliasEntry, error) {
	var err error

	// Validate JSON.
	schemaLoader := gojsonschema.NewStringLoader(aliasEntryJSONSchema)
	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return AliasEntry{}, errors.New("failed to validate JSON: " + err.Error())
	}
	if !result.Valid() {
		return AliasEntry{}, errors.New("invalid JSON")
	}

	// Unmarshal JSON.
	var entry AliasEntry
	err = json.Unmarshal(jsonBytes, &entry)
	if err != nil {
		return AliasEntry{}, errors.New("failed to unmarshal JSON: " + err.Error())
	}

	return entry, nil
}

func (entry AliasEntry) Map() map[string]interface{} {
	return map[string]interface{}{
		"type":   "alias",
		"target": entry.Target,
	}
}
