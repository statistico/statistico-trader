package errors

import "fmt"

type DuplicationError struct {
	Message string
}

func (d *DuplicationError) Error() string {
	return fmt.Sprintf("Duplication error: %s", d.Message)
}
