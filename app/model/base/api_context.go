package base

type contextKey string

const (
	// CorrelationIdContextKey holds the key used to store a correlation_id in the context.
	CorrelationIdContextKey contextKey = "CorrelationIdToken"
	// RequestHeaderContextKey holds the key used to store a request.Header in the context.
	RequestHeaderContextKey contextKey = "RequestHeaderToken"
	// SignedUserContextKey holds the key used to store a Signed User in the context.
	SignedUserContextKey contextKey = "SignedUserToken"
)
