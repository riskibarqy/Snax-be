package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type JSONOpts struct {
	W           http.ResponseWriter
	Status      int
	Data        any    // Optional, for success
	Err         error  // Optional, for error
	UserMessage string // Optional, for user-safe error
	SkipLogging bool   // Optional, for silencing logs (e.g., 404)
}

func JSONResponse(opts JSONOpts) {
	opts.W.Header().Set("Content-Type", "application/json")
	opts.W.WriteHeader(opts.Status)

	// Error response
	if opts.Err != nil {
		if !opts.SkipLogging {
			log.Printf("[ERROR] %v\n", opts.Err) // Replace with Uptrace later
		}

		_ = json.NewEncoder(opts.W).Encode(domain.APIResponse{
			Error: &domain.APIError{
				Message: opts.UserMessage,
				Detail:  opts.Err.Error(),
				Code:    opts.Status,
			},
		})
		return
	}

	// Success response
	_ = json.NewEncoder(opts.W).Encode(domain.APIResponse{
		Data: opts.Data,
	})
}

func RespondError(w http.ResponseWriter, status int, msg string, err error) {
	if err == nil {
		err = errors.New(msg)
	}
	JSONResponse(JSONOpts{
		W:           w,
		Status:      status,
		Err:         err,
		UserMessage: msg,
	})
}
