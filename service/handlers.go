package service

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/zenazn/goji/web"
)

// RequestError is an error that occurred during a request.
type RequestError struct {
	msg    string
	Status int
}

func (e *RequestError) Error() string {
	return e.msg
}

// RestdisResource represents a Restdis web service resource handler.
type RestdisResource struct {
	version string
	pool    *redis.Pool
	config  *Config
}

// Default is the default request handler.
func (r *RestdisResource) Default(c web.C, w http.ResponseWriter, req *http.Request) {
	err := writeJSON(map[string]string{"version": r.version}, w)
	if err != nil {
		c.Env["error"] = err
	}
}

// RedisCommand handles all requests for Redis commands.
func (r *RestdisResource) RedisCommand(c web.C, w http.ResponseWriter, req *http.Request) {
	command, args, err := buildRedisCommandFromURI(req.RequestURI)
	if err != nil {
		writeError(c, err, http.StatusInternalServerError)
		return
	}

	conn := r.pool.Get()
	defer conn.Close()

	res, err := conn.Do(command, args...)
	if err != nil {
		if res == nil {
			// If there's no response, it's most likely a Redis connection error.
			writeError(c, err, http.StatusInternalServerError)
			return
		}

		// Otherwise, the user entered an invalid command.
		writeError(c, err, http.StatusBadRequest)
		return
	}

	err = writeJSON(map[string]interface{}{
		"response": formatResponse(res),
	}, w)
	if err != nil {
		writeError(c, err, http.StatusInternalServerError)
	}
}

// writeError sets an error in the middleware context for the ErrorHandler middleware to handle later.
func writeError(c web.C, err error, status int) {
	c.Env[ErrorKey] = &RequestError{
		msg:    err.Error(),
		Status: status,
	}
}

func writeJSON(data interface{}, w http.ResponseWriter) error {
	j, err := json.Marshal(data)
	if err != nil {
		return &RequestError{
			msg:    err.Error(),
			Status: http.StatusInternalServerError,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)

	return nil
}

// buildRedisCommandFromURI splits a URI into a Redis command and arguments.
func buildRedisCommandFromURI(uri string) (string, redis.Args, error) {
	uri, err := url.QueryUnescape(uri)
	if err != nil {
		return "", nil, err
	}

	args := strings.Split(strings.Trim(uri, "/"), "/")
	redisCommand := args[0]
	redisArgs := redis.Args{}
	redisArgs = redisArgs.AddFlat(args[1:])

	return redisCommand, redisArgs, nil
}

// formatResponse formats the response from Redis into a usable map for converting to JSON.
func formatResponse(v interface{}) interface{} {
	switch v.(type) {
	case int64:
		return v.(int64)
	case string:
		return v.(string)
	case []byte:
		return string(v.([]byte)[:])
	case []interface{}:
		res := make([]interface{}, len(v.([]interface{})))

		for i, val := range v.([]interface{}) {
			res[i] = formatResponse(val)
		}

		return res
	default:
		return nil
	}
}
