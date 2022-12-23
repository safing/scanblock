package scanblock

import "net/http"

// ResponseWriter is used to wrap given response writers.
type ResponseWriter struct {
	http.ResponseWriter

	cacheEntry *CacheEntry
}

// WriteHeader adds custom handling to the wrapped WriterHeader method.
func (rw *ResponseWriter) WriteHeader(code int) {
	// If the request returns with a 4xx code, increase the scan requests counter.
	if code >= 400 && code < 500 {
		rw.cacheEntry.ScanRequests.Add(1)
	}

	// Continue with the original method.
	rw.ResponseWriter.WriteHeader(code)
}
