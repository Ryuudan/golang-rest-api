package render

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"runtime"
	"time"
)

func Error(w http.ResponseWriter, r *http.Request, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, file, line, _ := runtime.Caller(1)
	errorLocation := fmt.Sprintf("%s:%d", path.Dir(file), line)
	// Create a buffer to hold the JSON encoding.

	var buf []byte
	data := map[string]interface{}{
		"error": map[string]interface{}{
			"code":             code,
			"message":          message,
			"path":             r.URL.Path,
			"method":           r.Method,
			"user_agent":       r.Header.Get("User-Agent"),
			"query_parameters": r.URL.Query(),
			"time_stamp":       time.Now(),
			"location":         errorLocation,
		},
	}

	// Encode the data into the buffer.
	buf, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the buffer's content to the response writer.
	if _, err := w.Write(buf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
