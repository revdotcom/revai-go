package revai

import "fmt"

// ErrBadStatusCode is returned when the API returns a non 2XX error code
type ErrBadStatusCode struct {
	OriginalBody string
	Code         int
}

func (e *ErrBadStatusCode) Error() string {
	return fmt.Sprintf("Invalid status code: %d. Response: %s", e.Code, e.OriginalBody)
}
