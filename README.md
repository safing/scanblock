# Scan Block Traefik Plugin

Traefik plugin that blocks scanning IPs by counting 4xx status codes until a limit is hit.

Can also play games with scanners.

### Config

```
// MinScanRequests defines the minimum 4xx responses to observe before
// blocking an IP.
MinScanRequests uint64

// MinTotalRequests defines the minimum requests to observe before blocking
// an IP.
MinTotalRequests uint64

// MinScanPercent defines the minimum percent of 4xx responses of total
// requests before blocking an IP.
MinScanPercent float64

// BlockPrivate defines if private IP ranges (RFC1918, RFC4193) should be
// blocked too.
BlockPrivate bool

// PlayGames defines if the the plugin should respond with random 4xx status
// codes or even kill the connection sometimes.
PlayGames bool

// BlockSeconds defines for how many seconds an IP should be blocked.
BlockSeconds int

// RememberSeconds defines for how many seconds information about an IP
// should be cached after it was last seen.
RememberSeconds int
```
