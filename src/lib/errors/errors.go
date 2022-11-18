package errors

import "errors"

var ErrorAuth = errors.New("user is unauthenticated")
var ErrorInitDatabase = errors.New("database is not initialised")
var ErrorLoadingDatabase = errors.New("failed to load database")
var ErrorEnvironmentEnvNotFound = errors.New("could not find env variable 'environment'")
var ErrorLoadingEnv = errors.New("failed to load .env")
