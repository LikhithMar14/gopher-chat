package utils


import (
	"encoding/json"
	"net/http"

)



type Envelope map[string]interface{}

func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    encoder := json.NewEncoder(w)
    encoder.SetIndent("", " ") 
    if err := encoder.Encode(data); err != nil {
        return err
    }
    return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func WriteJSONError(w http.ResponseWriter, status int, message string) error {
	return WriteJSON(w, status, Envelope{"error": message})
}
