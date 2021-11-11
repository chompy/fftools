package main

import (
	"errors"
)

var (
	ErrRequestTooLarge  = errors.New("request was too large")
	ErrRequestNotFound  = errors.New("request not found")
	ErrRequestTimeout   = errors.New("request timed out")
	ErrResponseTooLarge = errors.New("response was too large")
)
