package fileops

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// RenameOptions stores the batch rename settings.
type RenameOptions struct {
	Prefix       string
	Suffix       string
	ReplaceText  string
	Replacement  string
	RegexPattern string
}

// RenameResult stores one completed rename operation.
type RenameResult struct {
	OldPath string
	NewPath string
}

// Rename renames one file or directory.
func Rename(
	sourcePath string,
	newName string,
	overwrite bool,
) error {
	if newName == "" {
		return fmt.Errorf("new name cannot be empty")
	}

	sourceInfo, err := os.Stat(sourcePath)

	if err != nil {
		return fmt.Errorf(
			"could not access source %s: %w",
			sourcePath,
			err,
		)
	}

	// Prevent paths such as another-folder/new.txt
	if filepath.Base(newName) != newName {
		return fmt.Errorf(
			"new name should not contain a directory path",
		)
	}

	parentDirectory := filepath.Dir(sourcePath)

	destinationPath := filepath.Join(
		parentDirectory,
		newName,
	)

	err = validateDifferentPaths(
		sourcePath,
		destinationPath,
	)

	if err != nil {
		return err
	}

	if PathExists(destinationPath) {
		if !overwrite {
			return fmt.Errorf(
				"an item named %s already exists",
				newName,
			)
		}

		err = RemoveExistingDestination(destinationPath)

		if err != nil {
			return err
		}
	}

	err = os.Rename(
		sourcePath,
		destinationPath,
	)

	if err != nil {
		return fmt.Errorf(
			"could not rename %s: %w",
			sourceInfo.Name(),
			err,
		)
	}

	return nil
}

// BatchRename renames all items inside one directory.
func BatchRename(
	directoryPath string,
	options RenameOptions,
) ([]RenameResult, error) {
	var results []RenameResult

	directoryInfo, err := os.Stat(directoryPath)

	if err != nil {
		return results, fmt.Errorf(
			"could not access directory %s: %w",
			directoryPath,
			err,
		)
	}

	if !directoryInfo.IsDir() {
		return results, fmt.Errorf(
			"%s is not a directory",
			directoryPath,
		)
	}

	if !hasRenameOperation(options) {
		return results, fmt.Errorf("provide a prefix, suffix, replacement, or regex pattern")
	}

	var compiledRegex *regexp.Regexp

	if options.RegexPattern != "" {
		compiledRegex, err = regexp.Compile(options.RegexPattern)

		if err != nil {
			return results, fmt.Errorf("invalid rename regular expression: %w", err)
		}
	}

	entries, err := os.ReadDir(directoryPath)

	if err != nil {
		return results, fmt.Errorf(
			"could not read directory %s: %w",
			directoryPath,
			err,
		)
	}

	// First calculate every new name.
	// This helps us check for errors before renaming.
	type renameOperation struct {
		oldPath string
		newPath string
	}

	var operations []renameOperation
	targetNames := make(map[string]bool)

	for _, entry := range entries {
		oldName := entry.Name()

		newName := createNewName(
			oldName,
			entry.IsDir(),
			options,
			compiledRegex,
		)

		// Do not add unchanged names.
		if oldName == newName {
			continue
		}

		if newName == "" {
			return results, fmt.Errorf(
				"rename operation created an empty name for %s",
				oldName,
			)
		}

		oldPath := filepath.Join(
			directoryPath,
			oldName,
		)

		newPath := filepath.Join(
			directoryPath,
			newName,
		)

		if targetNames[newName] {
			return results, fmt.Errorf(
				"multiple items would receive the name %s",
				newName,
			)
		}

		targetNames[newName] = true

		if PathExists(newPath) {
			return results, fmt.Errorf(
				"cannot rename %s because %s already exists",
				oldName,
				newName,
			)
		}

		operations = append(
			operations,
			renameOperation{
				oldPath: oldPath,
				newPath: newPath,
			},
		)
	}

	// Rename each item after checking the operations.
	for _, operation := range operations {
		err = os.Rename(
			operation.oldPath,
			operation.newPath,
		)

		if err != nil {
			return results, fmt.Errorf(
				"could not rename %s: %w",
				operation.oldPath,
				err,
			)
		}

		results = append(
			results,
			RenameResult{
				OldPath: operation.oldPath,
				NewPath: operation.newPath,
			},
		)
	}

	return results, nil
}

// hasRenameOperation checks whether the user provided
// at least one batch rename option.
func hasRenameOperation(options RenameOptions) bool {
	return options.Prefix != "" ||
		options.Suffix != "" ||
		options.ReplaceText != "" ||
		options.RegexPattern != ""
}

// createNewName creates a new name for one item.
func createNewName(
	oldName string,
	isDirectory bool,
	options RenameOptions,
	compiledRegex *regexp.Regexp,
) string {
	newName := oldName

	// Replace normal text.
	if options.ReplaceText != "" {
		newName = strings.ReplaceAll(
			newName,
			options.ReplaceText,
			options.Replacement,
		)
	}

	// Replace text using regular expressions.
	if compiledRegex != nil {
		newName = compiledRegex.ReplaceAllString(
			newName,
			options.Replacement,
		)
	}

	// Add the prefix before the complete name.
	if options.Prefix != "" {
		newName = options.Prefix + newName
	}

	// Add the suffix before the extension.
	if options.Suffix != "" {
		newName = addSuffix(
			newName,
			options.Suffix,
			isDirectory,
		)
	}

	return newName
}

// addSuffix adds a suffix before the file extension
func addSuffix(
	name string,
	suffix string,
	isDirectory bool,
) string {
	if isDirectory {
		return name + suffix
	}

	extension := filepath.Ext(name)

	// A file may not have an extension.
	if extension == "" {
		return name + suffix
	}

	nameWithoutExtension := strings.TrimSuffix(
		name,
		extension,
	)

	return nameWithoutExtension + suffix + extension
}
