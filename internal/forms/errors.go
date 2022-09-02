package forms

// Type errors a map with string keys and string array values
type errors map[string][]string

// Adds an error message for a given form field
func (e errors) AddErrorMessage(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns the first error message found
func (e errors) GetErrorMessage(field string) string {
	// Grab error message of associated field
	errorStringArray := e[field]

	// If the length of error array is 0
	if len(errorStringArray) == 0 {
		// return nothing
		return ""
	}

	// else return first index of error string
	return errorStringArray[0]
}
