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
	Tags        []string
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

		entry.Author = informationMap["author"].(string)
		entry.Description = informationMap["description"].(string)
		entry.Name = informationMap["name"].(string)

		if homepage, ok := informationMap["homepage"]; ok {
			entry.Homepage = homepage.(string)
		}

		if license, ok := informationMap["license"]; ok {
			entry.License = license.(string)
		}

		if repository, ok := informationMap["repository"]; ok {
			entry.Repository = repository.(string)
		}

		if tags, ok := informationMap["tags"]; ok {
			for _, tag := range tags.([]interface{}) {
				entry.Tags = append(entry.Tags, tag.(string))
			}
		} else {
			entry.Tags = []string{}
		}
	}

	return entry, nil
}
