package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/cmiski/sawako/shared/contextx"
)

func NewReverseProxy(
	serviceName string,
	targetURL string,
) (http.Handler, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf(
			"proxy: parse target url: %w",
			err,
		)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := reverseProxy.Director

	reverseProxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host

		if requestID := contextx.GetRequestID(
			req.Context(),
		); requestID != "" {
			req.Header.Set(
				"X-Request-ID",
				requestID,
			)
		}
	}

	reverseProxy.ErrorHandler = func(
		w http.ResponseWriter,
		r *http.Request,
		err error,
	) {
		log.Printf(
			"proxy %s error: %v",
			serviceName,
			err,
		)

		w.Header().Set(
			"Content-Type",
			"application/json",
		)
		w.WriteHeader(http.StatusBadGateway)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "service unavailable",
		})
	}

	return reverseProxy, nil
}
