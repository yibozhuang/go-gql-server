package utils

import (
	"os"
	"strconv"

	log "github.com/yibozhuang/go-gql-server/internal/logger"
)

// MustGet will return the env or panic if it is not present
func MustGet(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicln("ENV missing, key: " + k)
	}
	return v
}

// MustGetBool will return the env as boolean or panic if it is not present
func MustGetBool(k string) bool {
	v := os.Getenv(k)
	if v == "" {
		log.Panicln("ENV missing, key: " + k)
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Panicln("ENV err: [" + k + "]\n" + err.Error())
	}
	return b
}

// MustGetInt32 will return the env as int32 or panic if it is not present
func MustGetInt32(k string) int {
	v := os.Getenv(k)
	if v == "" {
		log.MissingArg(k)
		log.Panic("ENV missing, key: " + k)
	}
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		log.MissingArg(k)
		log.Panic("ENV err: [" + k + "]" + err.Error())
	}
	return int(i)
}

// MustGetInt64 will return the env as int64 or panic if it is not present
func MustGetInt64(k string) int64 {
	v := os.Getenv(k)
	if v == "" {
		log.MissingArg(k)
		log.Panic("ENV missing, key: " + k)
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.MissingArg(k)
		log.Panic("ENV err: [" + k + "]" + err.Error())
	}
	return i
}
