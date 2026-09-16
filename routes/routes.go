package routes

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json") // Anything that is written after the WriteHeader call will be sent as the response body
		// so any headers modifications should be done before the WriteHeader call
		w.WriteHeader(http.StatusOK) // the order of the WriteHeader and Write methods is important; WriteHeader should be called before Write

		jsonResponse := `{"status": 200, "message": "Server Running Healthy..."}`

		w.Write([]byte(jsonResponse))
	})
}
