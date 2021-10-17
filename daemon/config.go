package main

type Config struct {
	Scripts map[string]map[string]interface{} `json:"scripts"`
}
