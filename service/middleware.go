package service

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"
)

// Key to use when setting the error value.
const ErrorKey = "error"

// ErrorHandler handles writing error responses to the client.
func ErrorHandler(c *web.C, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

		err := getError(c)
		if err != nil {
			log.Error(err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.Status)
			writeJSON(map[string]string{"error": err.Error()}, w)
		}
	})
}

// getError returns an error from the given context if one is present.
func getError(c *web.C) *RequestError {
	if c.Env == nil {
		return nil
	}

	v, ok := c.Env[ErrorKey]
	if !ok {
		return nil
	}

	if err, ok := v.(*RequestError); ok {
		return err
	}

	return nil
}

// Headers injects headers defined in the configuration into every request.
func Headers(config *Config) web.MiddlewareType {
	return func(c *web.C, h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, header := range config.Headers {
				w.Header().Set(header.Name, header.Value)
			}

			h.ServeHTTP(w, r)
		})
	}
}
