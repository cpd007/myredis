package helper

import (
	"errors"
	"fmt"
)

func ToStringArray(d any) ([]string, error) {

	ds, ok := d.([]any)
	if !ok {
		return nil, errors.New("Assertion failed for any")
	}
	sa := make([]string, 0, len(ds))

	for i := range ds {
		s, ok := ds[i].(string)
		if !ok {
			return nil, fmt.Errorf("Type assertion failed for: %v", ds[i])
		}

		sa = append(sa, s)
	}

	return sa, nil
}
