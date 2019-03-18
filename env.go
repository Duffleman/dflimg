package dflimg

import (
	"encoding/json"
	"fmt"
	"os"
)

func GetEnv(key string) string {
	var v string
	switch key {
	case "salt":
		v = os.Getenv("DFL_SALT")
		if v == "" {
			return Salt
		}
	case "root_url":
		v = os.Getenv("DFL_ROOT_URL")
		if v == "" {
			return RootURL
		}
	case "pg_connection_string":
		v = os.Getenv("PG_OPTS")
		if v == "" {
			return PostgresCS
		}
	case "addr":
		v = os.Getenv("ADDR")
		if v == "" {
			return DefaultAddr
		}
	}

	return v
}

func GetUsers() map[string]string {
	v := os.Getenv("DFL_USERS")
	if v == "" {
		return Users
	}

	var users map[string]string

	err := json.Unmarshal([]byte(v), &users)
	if err != nil {
		panic(fmt.Errorf("cannot unmarshal user config: %s", err))
	}

	return users
}
