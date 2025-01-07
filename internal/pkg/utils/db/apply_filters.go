package mysql

import (
	"fmt"
	"strings"
)

type Filter struct {
	Operator string
	Value    interface{}
}

// ApplyFilters dynamically applies filters with various operators to a base query.
// It returns the updated query and a slice of arguments.
func ApplyFilters(baseQuery string, filters map[string]Filter) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	// Iterate over the filters and add them to the query
	for key, filter := range filters {
		// Append the condition with the specified operator
		conditions = append(conditions, fmt.Sprintf("%s %s ?", key, filter.Operator))
		args = append(args, filter.Value)
	}

	// If there are conditions, append them to the query
	if len(conditions) > 0 {
		filterQuery := strings.Join(conditions, " AND ")
		baseQuery = fmt.Sprintf("%s AND %s", baseQuery, filterQuery)
	}

	return baseQuery, args
}
