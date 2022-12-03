package middleware

import (
	"context"
	"go-klikdokter/app/model/base"
	"net"
	stdHttp "net/http"
	"strings"

	"github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
)

const (
	CorrelationIdContextKey = base.CorrelationIdContextKey
	RequestHeaderContextKey = base.RequestHeaderContextKey
	SignedUserContextKey = base.SignedUserContextKey
)

// CorrelationIdToContext create a correlation_id and add it to context.
// Needed for microservice single tracing.
func CorrelationIdToContext() http.RequestFunc {
	return func(ctx context.Context, r *stdHttp.Request) context.Context {
		correlationId := r.Header.Get("X-Correlation-ID")
		if correlationId == "" {
			correlationId = uuid.NewString()
		}

		return context.WithValue(ctx, CorrelationIdContextKey, correlationId)
	}
}

func GetIP(r *stdHttp.Request) string {
	userIP := r.Header.Get("X-FORWARDED-FOR")
	if userIP == "" {
		userIP = r.RemoteAddr
	}

	// Handle multiple IP. ex: "66.96.247.62, 10.0.191.120"
	ipMultiple := strings.Split(userIP, ",")
	if len(ipMultiple) > 0 {
		userIP = ipMultiple[0]
	}

	ip, _, err := net.SplitHostPort(userIP)
	if err != nil {
		return userIP
	}

	// Handle localhost
	if ip == "::1" {
		return "127.0.0.1"
	}

	return ip
}