package handler

import (
	"errors"
	"pokedex/pkg/http/response"
)

// Sentinel errors
const (
	ErrPokemonNotFound Error = iota + 1
)

// Error is a sentinel error
type Error uint

func (e Error) Error() string {
	switch e {
	case ErrPokemonNotFound:
		return "pokemon not found"
	}

	return "unknown error"
}

// ErrorRenderer returns a response.Renderer for handling errors
func ErrorRenderer(err error) response.Renderer {
	switch {
	case
		errors.Is(err, ErrPokemonNotFound):
		return response.NotFound(err)
	}

	return response.InternalServerError(err)
}
