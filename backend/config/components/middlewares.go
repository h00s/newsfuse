package components

import (
	"time"

	"github.com/go-raptor/middlewares/csrf"
	"github.com/go-raptor/middlewares/limiter"
	"github.com/go-raptor/middlewares/logger"
	"github.com/go-raptor/raptor/v4"
	"golang.org/x/time/rate"
)

func Middlewares() raptor.Middlewares {
	return raptor.Middlewares{
		raptor.Use(&logger.LoggerMiddleware{}),
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
