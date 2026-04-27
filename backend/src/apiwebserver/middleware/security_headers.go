package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders sets a conservative set of HTTP security headers on every
// response. It's intentionally minimal — these are the headers that bring real
// risk reduction with no behavioral side-effects:
//
//   - X-Content-Type-Options: blocks MIME-sniff drift attacks.
//   - X-Frame-Options:        prevents clickjacking via <iframe>.
//   - Referrer-Policy:        avoids leaking the full URL on outbound clicks.
//   - X-XSS-Protection: 0     explicitly opt out of the legacy IE filter
//     (modern browsers ignore it; the official guidance is to disable it
//     because the filter itself has been a source of vulns).
//
// HSTS is intentionally NOT set here — it should be set by the TLS terminator
// (nginx / load balancer) so it's only active over HTTPS. CSP is also not set
// because a tight CSP needs to be tuned per-frontend and is easy to break.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("X-XSS-Protection", "0")
		c.Next()
	}
}
