package response

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Renderer handles a http request returning an error
type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request) error
}

// RendererFunc allows regular methods to act as Renderers
type RendererFunc func(w http.ResponseWriter, r *http.Request) error

// Render calls the wrapped method.
func (fn RendererFunc) Render(w http.ResponseWriter, r *http.Request) error {
	return fn(w, r)
}

// RenderFunc returns a Renderer for a HTTP request
type RenderFunc func(w http.ResponseWriter, r *http.Request) Renderer

// Render returns a http.HandlerFunc that executes the given render func
func Render(fn RenderFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r).Render(w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// JSON provides a renderer that write the given json body
// if a nil body is given no data will be written on the response body
func JSON(status int, body interface{}) Renderer {
	return RendererFunc(func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)

		if body != nil {
			var buff bytes.Buffer

			enc := json.NewEncoder(&buff)
			enc.SetEscapeHTML(true)

			if err := enc.Encode(body); err != nil {
				return errors.Wrap(err, "failed to encode json body")
			}

			if _, err := w.Write(buff.Bytes()); err != nil {
				return errors.Wrap(err, "failed to write json body")
			}
		}

		return nil
	})
}

// Error returns a renderer for given error
func Error(status int, err error) Renderer {
	return JSON(status, map[string]string{
		"Error": err.Error(),
	})
}

// InternalServerError returns a 500 renderer, the given error is logged
func InternalServerError(err error) Renderer {
	return RendererFunc(func(w http.ResponseWriter, r *http.Request) error {
		log.Info().Err(err).Send()
		return Error(http.StatusInternalServerError, err).Render(w, r)
	})
}

// NotFound returns a 404 renderer
func NotFound(err error) Renderer {
	return RendererFunc(func(w http.ResponseWriter, r *http.Request) error {
		return Error(http.StatusNotFound, err).Render(w, r)
	})
}
