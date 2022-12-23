package scanblock

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func (sb *ScanBlock) block(w http.ResponseWriter, r *http.Request) {
	// Attempt to clean the cache.
	removedEntries := sb.cache.CleanEntries(
		time.Duration(sb.config.RememberSeconds) * time.Second,
	)
	if removedEntries > 0 {
		fmt.Fprintf(os.Stdout, "scanblock plugin %q purged %d cache entries\n", sb.name, removedEntries)
	}

	// If we are not playing any games, reply with a plain error message.
	if !sb.config.PlayGames {
		http.Error(w, "blocked by scanblock", http.StatusTooManyRequests)
		return
	}

	// Always wait a little.
	select {
	case <-r.Context().Done():
	case <-time.After(getRandomWaitDuration()):
	}

	// Select random game until one works.
	start := rand.Intn(len(games)) //nolint:gosec // Not for security.
	for {
		// Try game and return if successful.
		if games[start](w, r) {
			return
		}

		// Select next game and wrap to start.
		start = (start + 1) % len(games)
	}
}

var games = []func(w http.ResponseWriter, r *http.Request) (ok bool){
	playClose,
	play4xx,
	play4xxEmoji,
	play4xxMessage,
}

func playClose(w http.ResponseWriter, r *http.Request) (ok bool) {
	// Check if the response writer supports hijacking.
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return false
	}

	// Hijack connection.
	conn, _, err := hijacker.Hijack()
	if err != nil {
		return false
	}
	// Close connection.
	_ = conn.Close()

	return true
}

func play4xx(w http.ResponseWriter, r *http.Request) (ok bool) {
	w.WriteHeader(getRandom4xxStatusCode())
	return true
}

func play4xxEmoji(w http.ResponseWriter, r *http.Request) (ok bool) {
	http.Error(w, "😛", getRandom4xxStatusCode())
	return true
}

func play4xxMessage(w http.ResponseWriter, r *http.Request) (ok bool) {
	http.Error(w, "Let's play a game.", getRandom4xxStatusCode())
	return true
}

const (
	waitDurationMin = 10 * time.Second
	waitDurationMax = 25 * time.Second
)

func getRandomWaitDuration() time.Duration {
	return time.Duration(
		rand.Intn(int(waitDurationMax-waitDurationMin)), //nolint:gosec // Not for security.
	) + waitDurationMin
}

func getRandom4xxStatusCode() int {
	return statusCodes[rand.Intn(len(statusCodes))] //nolint:gosec // Not for security.
}

var statusCodes = []int{
	http.StatusBadRequest,
	http.StatusUnauthorized,
	http.StatusPaymentRequired,
	http.StatusForbidden,
	http.StatusNotFound,
	http.StatusMethodNotAllowed,
	http.StatusNotAcceptable,
	http.StatusProxyAuthRequired,
	http.StatusRequestTimeout,
	http.StatusConflict,
	http.StatusGone,
	http.StatusLengthRequired,
	http.StatusPreconditionFailed,
	http.StatusRequestEntityTooLarge,
	http.StatusRequestURITooLong,
	http.StatusUnsupportedMediaType,
	http.StatusRequestedRangeNotSatisfiable,
	http.StatusExpectationFailed,
	http.StatusTeapot,
	http.StatusMisdirectedRequest,
	http.StatusUnprocessableEntity,
	http.StatusLocked,
	http.StatusFailedDependency,
	http.StatusTooEarly,
	http.StatusUpgradeRequired,
	http.StatusPreconditionRequired,
	http.StatusTooManyRequests,
	http.StatusRequestHeaderFieldsTooLarge,
	http.StatusUnavailableForLegalReasons,
}
