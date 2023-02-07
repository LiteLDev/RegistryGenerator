package entry

import (
	"encoding/json"
	"errors"

	"github.com/xeipuuv/gojsonschema"
)

type Entry struct {
	ToothPath   string
	Name        string
	Description string
	Author      string
	License     string
	Homepage    string
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
                "name",
                "description",
                "author",
                "license",
                "homepage"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "author": {
                    "type": "string"
                },
                "license": {
                    "type": "string"
                },
                "homepage": {
                    "type": "string"
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
		if name, ok := informationMap["name"]; ok {
			entry.Name = name.(string)
		}
		if description, ok := informationMap["description"]; ok {
			entry.Description = description.(string)
		}
		if author, ok := informationMap["author"]; ok {
			entry.Author = author.(string)
		}
		if license, ok := informationMap["license"]; ok {
			entry.License = license.(string)
		}
		if homepage, ok := informationMap["homepage"]; ok {
			entry.Homepage = homepage.(string)
		}
	}

	return entry, nil
}
