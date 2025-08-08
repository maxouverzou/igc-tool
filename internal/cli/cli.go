package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"igc-tool/internal/logbook"
	"igc-tool/internal/sites"
)

// FindIGCFiles finds all IGC files from the given paths (files or directories)
// If recursive is true, it will search subdirectories as well
func FindIGCFiles(paths []string, recursive bool) ([]string, error) {
	var igcFiles []string

	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("error accessing %s: %w", path, err)
		}

		if stat.IsDir() {
			// Handle directory
			if recursive {
				err = filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if !d.IsDir() && strings.ToLower(filepath.Ext(filePath)) == ".igc" {
						igcFiles = append(igcFiles, filePath)
					}
					return nil
				})
			} else {
				entries, err := os.ReadDir(path)
				if err != nil {
					return nil, fmt.Errorf("error reading directory %s: %w", path, err)
				}
				for _, entry := range entries {
					if !entry.IsDir() && strings.ToLower(filepath.Ext(entry.Name())) == ".igc" {
						igcFiles = append(igcFiles, filepath.Join(path, entry.Name()))
					}
				}
			}
			if err != nil {
				return nil, fmt.Errorf("error walking directory %s: %w", path, err)
			}
		} else {
			// Handle regular file
			if strings.ToLower(filepath.Ext(path)) == ".igc" {
				igcFiles = append(igcFiles, path)
			} else {
				return nil, fmt.Errorf("file %s is not an IGC file", path)
			}
		}
	}

	return igcFiles, nil
}

// PrintTemplatedLogbookData prints logbook output using the provided template with TemplateData
func PrintTemplatedLogbookData(data *logbook.TemplateData, templateStr string) error {
	if data == nil {
		fmt.Println("No flight data available for logbook entry")
		return nil
	}

	tmpl, err := template.New("logbook").Parse(templateStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// LoadLandingSitesIfSpecified loads landing sites if a file is specified
func LoadLandingSitesIfSpecified(filename string) (*sites.Collection, error) {
	if filename == "" {
		return nil, nil
	}

	landingSites, err := sites.LoadLandingSites(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not load landing sites: %v\n", err)
		return nil, nil
	}

	// If no valid sites were loaded, return nil instead of empty collection
	if len(landingSites.Sites) == 0 {
		fmt.Fprintf(os.Stderr, "Warning: No valid landing sites found in %s\n", filename)
		return nil, nil
	}

	return landingSites, nil
}
