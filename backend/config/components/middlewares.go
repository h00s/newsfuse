package components

import (
	"time"

	"github.com/go-raptor/middlewares/csrf"
	"github.com/go-raptor/middlewares/limiter"
	"github.com/go-raptor/middlewares/logger"
	"github.com/go-raptor/middlewares/requestid"
	"github.com/go-raptor/middlewares/secure"
	"github.com/go-raptor/raptor/v4"
	"golang.org/x/time/rate"
)

func Middlewares() raptor.Middlewares {
	return raptor.Middlewares{
		// First, so the access log and Raptor's own error lines carry request_id.
		raptor.Use(&requestid.RequestIDMiddleware{}),
		raptor.Use(&logger.LoggerMiddleware{}),
		// Security headers on every response, the CSRF 403 and the 429s included. HSTS is on once
		// app.secure_hsts_max_age is set.
		raptor.Use(&secure.SecureMiddleware{}),
		// There are no sessions to protect, but this keeps other sites from making their visitors'
		// browsers spend summaries, spread over many addresses past the per-client limiter.
		raptor.Use(&csrf.CSRFMiddleware{}),
		// Each summary is a paid LLM call: 5 at once, then one every 12 seconds per client.
		raptor.UseOnly(limiter.NewRateLimiterMiddleware(limiter.RateLimiterConfig{
			Rate:  rate.Every(12 * time.Second),
			Burst: 5,
		}), "Stories.Summarize"),
		// A story's first read scrapes the news site.
		raptor.UseOnly(limiter.NewRateLimiterMiddleware(limiter.RateLimiterConfig{
			Rate:  1,
			Burst: 20,
		}), "Headlines.Story"),
	}
}
