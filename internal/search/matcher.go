package search

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

//information to match a filename sob yat store thakibo
type fileMatcher struct {
	options Options
	regularExpression *regexp.Regexp 		//only used when regex flag is enabled
}

//creates a matcher
func newMatcher(options Options) (*fileMatcher, error) {
	matcher := &fileMatcher{
		options: options,
	}

	if options.UseRegex {
		pattern := options.Pattern

		// (?i) makes a regular expression case-insensitive
		if options.IgnoreCase {
			pattern = "(?i)" + pattern
		}

		compiledRegex, err := regexp.Compile(pattern)

		if err != nil {
			return nil, fmt.Errorf(
				"invalid regular expression: %w",
				err,
			)
		}

		matcher.regularExpression = compiledRegex
	}

	return matcher, nil
}

// matches checks whether a filename matches the search pattern
func (matcher *fileMatcher) matches(name string) bool {
	options := matcher.options

	if options.UseRegex {
		return matcher.regularExpression.MatchString(name)
	}

	if options.Extension {
		return matcher.matchExtension(name)
	}

	// Support patterns
	if strings.Contains(options.Pattern, "*") || strings.Contains(options.Pattern, "?") {

		return matcher.matchWildcard(name)
	}

	if options.Exact {
		return matcher.matchExact(name)
	}

	return matcher.matchPartial(name)
}

// matchExtension compares file extensions.
func (matcher *fileMatcher) matchExtension(name string) bool {
	extension := filepath.Ext(name)
	
	// Remove the dot so that the user can search "pdf".
	extension = strings.TrimPrefix(extension, ".")

	pattern := strings.TrimPrefix(
		matcher.options.Pattern,
		".",
	)

	if matcher.options.IgnoreCase {
		return strings.EqualFold(extension, pattern)
	}

	return extension == pattern
}

// matchExact checks the complete filename
func (matcher *fileMatcher) matchExact(name string) bool {
	if matcher.options.IgnoreCase {
		return strings.EqualFold(name, matcher.options.Pattern)
	}

	return name == matcher.options.Pattern
}

// matchPartial checks whether the pattern exists inside the name.
func (matcher *fileMatcher) matchPartial(name string) bool {
	pattern := matcher.options.Pattern

	if matcher.options.IgnoreCase {
		name = strings.ToLower(name)
		pattern = strings.ToLower(pattern)
	}

	return strings.Contains(name, pattern)
}

// matchWildcard handles patterns such as "*.pdf".
func (matcher *fileMatcher) matchWildcard(name string) bool {
	pattern := matcher.options.Pattern

	if matcher.options.IgnoreCase {
		name = strings.ToLower(name)
		pattern = strings.ToLower(pattern)
	}

	matched, err := filepath.Match(pattern, name)

	if err != nil {
		return false
	}

	return matched
}
