package mysql

import (
	"fmt"
	"strings"
)

// ApplyFilters dynamically applies filters to a base query.
// It returns the updated query and a slice of arguments.
func ApplyFilters(baseQuery string, filters map[string]interface{}) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	// Iterate over the filters and add them to the query
	for key, value := range filters {
		// Use a placeholder for the value
		conditions = append(conditions, fmt.Sprintf("%s = ?", key))
		args = append(args, value)
	}

	// If there are conditions, append them to the query
	if len(conditions) > 0 {
		filterQuery := strings.Join(conditions, " AND ")
		baseQuery = fmt.Sprintf("%s AND %s", baseQuery, filterQuery)
	}

	return baseQuery, args
}
