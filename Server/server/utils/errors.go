package utils

import "fmt"

func NewUnknownTypeError(t any) UnknownTypeError {
	return UnknownTypeError{Type: t, Value: t}
}

type UnknownTypeError struct {
	Type      any
	Value     any
	ExtraInfo string
}

func (u UnknownTypeError) WithValue(v any) UnknownTypeError {
	u.Value = v
	return u
}

func (u UnknownTypeError) WithExtraInfo(extraInfo string) UnknownTypeError {
	u.ExtraInfo = extraInfo
	return u
}

func (u UnknownTypeError) Error() string {
	err := fmt.Sprintf("unknown value (%v) used in %T", u.Value, u.Type)
	if u.ExtraInfo != "" {
		err += " (" + u.ExtraInfo + ")"
	}
	return err
}

type PlayerDataNotFoundError struct {
	Identifier string
}

func (e PlayerDataNotFoundError) Error() string {
	err := "cannot find player data"
	if e.Identifier != "" {
		err += " (identifier: " + e.Identifier + ")"
	}
	return err
}
