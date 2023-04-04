package entries

import (
	"encoding/json"
	"errors"

	"github.com/xeipuuv/gojsonschema"
)

type ToothEntry struct {
	Type        string
	Tooth       string
	Author      string
	Description string
	Homepage    string
	License     string
	Name        string
	Repository  string
	Tags        []string
}

const toothEntryJSONSchema = `
{
    "$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"properties": {
		"format_version": {
			"const": 1
		},
		"type": {
			"type": "string",
			"const": "tooth"
		},
		"tooth": {
			"type": "string"
		},
		"information": {
			"type": "object",
			"additionalProperties": false,
			"required": [
				"author",
				"description",
				"name"
			],
			"properties": {
				"author": {
					"type": "string",
					"pattern": "^[a-zA-Z0-9-]+$"
				},
				"description": {
					"type": "string",
					"minLength": 1
				},
				"homepage": {
					"type": "string",
					"pattern": "^https?:\/\/.+$"
				},
				"license": {
					"type": "string",
					"pattern": "^[a-zA-Z0-9-+.]*$"
				},
				"name": {
					"type": "string",
					"minLength": 1
				},
				"repository": {
					"type": "string",
					"pattern": "^github\\.com\/[a-zA-Z0-9-]+\/[a-zA-Z0-9-_.]+$"
				},
				"tags": {
					"type": "array",
					"items": {
						"type": "string",
						"pattern": "^[a-z0-9-]+$"
					}
				}
			}
		}
	},
	"required": [
		"format_version",
		"tooth",
		"information"
	],
	"additionalProperties": false
}
`

// NewToothEntryFromJSON creates a new entry from JSON.
func NewToothEntryFromJSON(jsonBytes []byte) (ToothEntry, error) {
	var err error

	// Validate JSON.
	schemaLoader := gojsonschema.NewStringLoader(toothEntryJSONSchema)
	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return ToothEntry{}, errors.New("failed to validate JSON: " + err.Error())
	}
	if !result.Valid() {
		return ToothEntry{}, errors.New("invalid JSON")
	}

	// Unmarshal JSON.
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	if err != nil {
		return ToothEntry{}, errors.New("failed to unmarshal JSON: " + err.Error())
	}

	// Create entry.
	entry := ToothEntry{}
	entry.Tooth = jsonMap["tooth"].(string)

	information := jsonMap["information"].(map[string]interface{})
	entry.Author = information["author"].(string)
	entry.Description = information["description"].(string)
	entry.Name = information["name"].(string)

	if homepage, ok := information["homepage"]; ok {
		entry.Homepage = homepage.(string)
	}

	if license, ok := information["license"]; ok {
		entry.License = license.(string)
	}

	if repository, ok := information["repository"]; ok {
		entry.Repository = repository.(string)
	}

	entry.Tags = []string{}
	if tags, ok := information["tags"]; ok {
		for _, tag := range tags.([]interface{}) {
			entry.Tags = append(entry.Tags, tag.(string))
		}
	}

	return entry, nil
}

func (entry ToothEntry) Map() map[string]interface{} {
	return map[string]interface{}{
		"type":        "tooth",
		"tooth":       entry.Tooth,
		"author":      entry.Author,
		"description": entry.Description,
		"homepage":    entry.Homepage,
		"license":     entry.License,
		"name":        entry.Name,
		"repository":  entry.Repository,
		"tags":        entry.Tags,
	}
}
