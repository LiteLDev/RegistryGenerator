package entries

import (
	"encoding/json"
	"errors"

	"github.com/xeipuuv/gojsonschema"
)

type Entry struct {
	ToothPath   string
	Author      string
	Description string
	Homepage    string
	License     string
	Name        string
	Repository  string
}

const entryJSONSchema = `
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "additionalProperties": false,
    "required": [
        "format_version",
        "tooth",
        "information"
    ],
    "properties": {
        "format_version": {
            "const": 1
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
                "homepage",
                "license",
                "name",
                "repository"
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
                }
            }
        }
    }
}
`

// New creates a new entry.
func New() *Entry {
	return &Entry{}
}

// NewFromJSON creates a new entry from JSON.
func NewFromJSON(jsonBytes []byte) (*Entry, error) {
	var err error

	// Validate JSON.
	schemaLoader := gojsonschema.NewStringLoader(entryJSONSchema)
	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, errors.New("failed to validate JSON: " + err.Error())
	}
	if !result.Valid() {
		return nil, errors.New("invalid JSON")
	}

	// Unmarshal JSON.
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	if err != nil {
		return nil, errors.New("failed to unmarshal JSON: " + err.Error())
	}

	// Create entry.
	entry := New()
	entry.ToothPath = jsonMap["tooth"].(string)
	if information, ok := jsonMap["information"]; ok {
		informationMap := information.(map[string]interface{})
		if author, ok := informationMap["author"]; ok {
			entry.Author = author.(string)
		}
		if description, ok := informationMap["description"]; ok {
			entry.Description = description.(string)
		}
		if homepage, ok := informationMap["homepage"]; ok {
			entry.Homepage = homepage.(string)
		}
		if license, ok := informationMap["license"]; ok {
			entry.License = license.(string)
		}
		if name, ok := informationMap["name"]; ok {
			entry.Name = name.(string)
		}
		if repository, ok := informationMap["repository"]; ok {
			entry.Repository = repository.(string)
		}
	}

	return entry, nil
}
