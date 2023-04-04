package entries

import (
	"encoding/json"
	"errors"

	"github.com/xeipuuv/gojsonschema"
)

type IEntry interface {
	Map() map[string]interface{}
}

const jsonSchema = `
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "anyOf": [
        {
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
        },
        {
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
    ]
}
`

func NewFromJSON(jsonBytes []byte) (IEntry, error) {
	var err error

	// Validate JSON.
	schemaLoader := gojsonschema.NewStringLoader(jsonSchema)
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

	if entryType, ok := jsonMap["type"]; ok && entryType.(string) == "alias" {
		return NewAliasEntryFromJSON(jsonBytes)
	} else {
		return NewToothEntryFromJSON(jsonBytes)
	}
}
