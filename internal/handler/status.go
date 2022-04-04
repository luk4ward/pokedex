package handler

import "net/http"

// HealthCheck handler for returning a 200
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
