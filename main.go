package main

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/liteldev/registrygenerator/entries"
	"github.com/liteldev/registrygenerator/logger"
)

const helpMessage = `
Usage:
  reggen [options]
  
Options:
  --h, -help                  Show help.
  --input <path>              Input directory. (default: entries)
  --output <path>             Output file. (default: index.json)`

func main() {
	// Parse flags.
	flagSet := flag.NewFlagSet("generator", flag.ExitOnError)
	flagSet.Usage = func() {
		logger.Info(helpMessage)
	}
	input := flagSet.String("input", "entries", "Input directory.")
	output := flagSet.String("output", "index.json", "Output file.")
	flagSet.Parse(os.Args[1:])

	// Check flags.
	if *input == "" {
		logger.Error("input directory is not specified")
		os.Exit(1)
	}
	if *output == "" {
		logger.Error("output file is not specified")
		os.Exit(1)
	}
	if flagSet.NArg() > 0 {
		logger.Error("unknown arguments")
		os.Exit(1)
	}

	// Load entries from input directory.
	entryMap, err := loadEntries(*input)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Generate index file.
	err = generateIndex(*output, entryMap)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("index file generated")
}

// loadEntries loads entries from input directory.
func loadEntries(input string) (map[string]entries.IEntry, error) {
	// Open input directory.
	dir, err := os.Open(input)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Read input directory.
	names, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	// Load entries.
	entryMap := make(map[string]entries.IEntry)
	for _, name := range names {
		// File name must end with ".json" extension.
		if filepath.Ext(name) != ".json" {
			continue
		}

		// Load entry.
		entry, err := loadEntry(input, name)
		if err != nil {
			return nil, err
		}

		// Add entry to map.
		entryMap[strings.TrimSuffix(name, filepath.Ext(name))] = entry
	}

	return entryMap, nil
}

// loadEntry loads entry from input directory.
func loadEntry(input, name string) (entries.IEntry, error) {
	var err error

	// Open entry file.
	file, err := os.Open(filepath.Join(input, name))
	if err != nil {
		return nil, errors.New("failed to open entry file: " + err.Error())
	}
	defer file.Close()

	// Read entry file.
	jsonBytes, err := os.ReadFile(file.Name())
	if err != nil {
		return nil, errors.New("failed to read entry file: " + err.Error())
	}

	// Create entry.
	entry, err := entries.NewFromJSON(jsonBytes)
	if err != nil {
		return nil, errors.New("failed to create entry with " + name + ": " + err.Error())
	}

	return entry, nil
}

// generateIndex generates index file.
func generateIndex(output string, entryMap map[string]entries.IEntry) error {
	var err error

	// Create output directory.
	err = os.MkdirAll(filepath.Dir(output), 0755)
	if err != nil {
		return errors.New("failed to create output directory: " + err.Error())
	}

	// Create index file.
	file, err := os.Create(output)
	if err != nil {
		return errors.New("failed to create index file: " + err.Error())
	}
	defer file.Close()

	// Make output map.
	outputObject := make(map[string]interface{})
	outputObject["format_version"] = 1
	outputObject["index"] = make(map[string]interface{})
	for identifier, entry := range entryMap {
		outputObject["index"].(map[string]interface{})[identifier] = entry.Map()
	}

	// Encode output object to JSON.
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(outputObject)
	if err != nil {
		return errors.New("failed to encode output map: " + err.Error())
	}

	return nil
}
