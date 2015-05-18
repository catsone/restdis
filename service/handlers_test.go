package service

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

type commandResult struct {
	command string
	args    redis.Args
}

var commandTests = []struct {
	uri      string
	expected commandResult
}{
	{"/GET/foo", commandResult{"GET", redis.Args{"foo"}}},
	{"/SET/foo/bar", commandResult{"SET", redis.Args{"foo", "bar"}}},
	{"/GETJOB/FROM/foo/1000/", commandResult{"GETJOB", redis.Args{"FROM", "foo", "1000"}}},
	{"/INFO", commandResult{"INFO", redis.Args{}}},
}

func TestBuildRedisCommandFromURI(t *testing.T) {
	for _, tt := range commandTests {
		command, args := buildRedisCommandFromURI(tt.uri)
		if command != tt.expected.command {
			t.Errorf("Expected %v, got %v", tt.expected.command, command)
		}

		for i, arg := range args {
			if arg != tt.expected.args[i] {
				t.Errorf("Expected %v, got %v", tt.expected.args[i], arg)
			}
		}
	}
}
