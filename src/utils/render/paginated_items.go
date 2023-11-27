package render

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
)

type Meta struct {
	Page          int    `json:"page"`
	Size          int    `json:"size"`
	TotalFiltered int    `json:"total_filtered"`
	Count         int    `json:"count"`
	LastPage      int    `json:"last_page"`
	From          int    `json:"from"`
	To            int    `json:"to"`
	Query         string `json:"query"`
	Order         string `json:"order"`
}

// Basic query params for pagination and searching
type QueryParams struct {
	Page  int    `json:"page,omitempty"`  // Offset for pagination (default: 0)
	Limit int    `json:"limit,omitempty"` // Number of items per page (default: 10)
	Query string `json:"query,omitempty"` // Query string for filtering
	Order string `json:"order,omitempty"` // Query string for filtering
}

type PaginatedResults struct {
	Meta    *Meta       `json:"meta"`
	Results interface{} `json:"results"`
}

// ParseOrderString parses a string containing comma-separated order fields and directions
// into a slice of order field structures. The input string should have a format like
// "field1:asc,field2:desc,field3:asc".
//
// It returns a slice of *models.OrderFields representing each order field and direction,
// or an error if the input string is in an invalid format.
//
// The order of the parsed fields determines their priority for sorting. Fields listed
// earlier in the input string have higher sorting priority.
//
// Example:
//
//	orderString := "name:asc,budget:desc,created_at:asc"
//	orders, err := ParseOrderString(orderString)
//
// This function ensures that the chronological order of orders is preserved, with the
// first order being the most absolute and serving as the primary sorting criterion.
func ParseOrderString(orderString string) ([]*OrderFields, error) {
	// Split the input string into individual order parts using commas
	orderParts := strings.Split(orderString, ",")
	orders := []*OrderFields{}

	for _, orderPart := range orderParts {
		// Split each order part into field and direction using a colon
		parts := strings.Split(orderPart, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid order string format: %s", orderPart)
		}

		field := parts[0]
		direction := parts[1]

		// Check if the direction is "asc" or "desc"
		if direction != "asc" && direction != "desc" {
			return nil, fmt.Errorf("invalid order direction: %s", direction)
		}

		// Create an OrderFields structure and add it to the slice of orders
		orders = append(orders, &OrderFields{
			Field:     field,
			Direction: direction,
		})
	}

	return orders, nil
}

// ParseQueryFilterParams parses a raw query string into a struct of query parameters.
// It extracts values for fields like "page", "limit", "query", and "order" from the
// query string and sets default values when not specified.
//
// The raw query string should be in a format like "?page=1&limit=10&query=search&order=name:asc".
//
// Example:
//
//	rawQuery := "?page=2&limit=20&query=keyword&order=name:desc"
//	params, err := ParseQueryFilterParams(rawQuery)
//
// This function returns a *QueryParams struct representing the parsed query parameters, or an error
// if there is an issue parsing the raw query.
//
// Example:
//
//	params := &QueryParams{
//	    Page:  2,
//	    Limit: 20,
//	    Query: "keyword",
//	    Order: "name:desc",
//	}
//	rawQuery := BuildRawQueryFromParams(params)
//	// rawQuery will be "?page=2&limit=20&query=keyword&order=name:desc"
func ParseQueryFilterParams(rawQuery string) (*QueryParams, error) {
	// Parse the raw query string into a map of values
	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return nil, err
	}

	// Create a QueryParams struct to hold the parsed values
	params := &QueryParams{}

	// Set default values for page, limit, query, and order fields when not specified
	params.Page = getIntValue(values, "page", 1)
	params.Limit = getIntValue(values, "limit", 15)
	params.Query = getStringValue(values, "query", "")
	params.Order = getStringValue(values, "order", "created_at:desc")

	return params, nil
}

// getIntValue retrieves an integer value from the provided URL values (query parameters)
// using the given key. If the value is not found, empty, or cannot be converted to an integer,
// it returns the specified defaultValue.
//
// Example:
//
//	values := url.Values{"page": {"2"}}
//	key := "page"
//	defaultValue := 1
//	result := getIntValue(values, key, defaultValue) // Result will be 2
func getIntValue(values url.Values, key string, defaultValue int) int {
	value := values.Get(key)

	if result, err := strconv.Atoi(value); err == nil && result != 0 {
		return result
	}

	return defaultValue
}

// getStringValue retrieves a string value from the provided URL values (query parameters)
// using the given key. If the value is not found or empty, it returns the specified defaultValue.
//
// Example:
//
//	values := url.Values{"query": {"search"}}
//	key := "query"
//	defaultValue := ""
//	result := getStringValue(values, key, defaultValue) // Result will be "search"
func getStringValue(values url.Values, key string, defaultValue string) string {
	value := values.Get(key)
	if value == "" {
		return defaultValue
	}

	return value
}

// GenerateMeta generates metadata for paginated results based on the total number of items,
// query parameters, and the number of items returned in the current page.
//
// It calculates and returns a *Meta structure containing information about the current page,
// page size, total filtered items, count of items on the current page, last page number, range
// of items displayed on the current page, query string, and sorting order.
//
// Example:
//
//	total_filtered := 1000
//	queryParams := &QueryParams{Page: 2, Limit: 20, Query: "search", Order: "created_at:desc"}
//	count := 20
//	meta := GenerateMeta(total_filtered, queryParams, count)
//
// The generated Meta structure provides metadata that can be used for building paginated
// responses.
func GenerateMeta(total int, queryParams *QueryParams, count int) *Meta {

	from := (queryParams.Page-1)*queryParams.Limit + 1
	to := (queryParams.Page-1)*queryParams.Limit + count

	if total == 0 {
		from = 0
	}

	return &Meta{
		Page:          queryParams.Page,
		Size:          queryParams.Limit,
		TotalFiltered: total,
		Count:         count,
		LastPage:      int(math.Ceil(float64(total) / float64(queryParams.Limit))),
		From:          from,
		To:            to,
		Query:         queryParams.Query,
		Order:         queryParams.Order,
	}
}

type OrderFields struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

type Filters struct {
	Order []OrderFields `json:"order"`
}
