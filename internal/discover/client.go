// Package discover provides YouTube ad domain discovery via browser network capture.
package discover

import (
	"context"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// DefaultAdPatterns are known ad-related domain substrings used for filtering.
var DefaultAdPatterns = []string{
	"googlevideo", "doubleclick", "googlesyndication", "googleadservices",
	"innovid", "moatads", "fwmrm", "adform", "serving-sys", "tubemogul",
	"2mdn", "imasdk", "googleadapis", "adservice", "ads.youtube", "ad.youtube",
}

// Config holds configuration for the discovery client.
type Config struct {
	Duration   time.Duration
	Blocklist  string
	AdPatterns []string
}

// Client captures network traffic from YouTube and extracts ad-related domains.
type Client struct {
	config   Config
	existing map[string]bool
	domains  map[string]struct{}
	mu       sync.Mutex
}

// NewClient creates a new discovery client.
func NewClient(cfg Config, existing map[string]bool) *Client {
	patterns := cfg.AdPatterns
	if len(patterns) == 0 {
		patterns = DefaultAdPatterns
	}
	return &Client{
		config: Config{
			Duration:   cfg.Duration,
			Blocklist:  cfg.Blocklist,
			AdPatterns: patterns,
		},
		existing: existing,
		domains:  make(map[string]struct{}),
	}
}

// Run navigates to YouTube, captures network traffic for the configured duration,
// and returns newly discovered ad-related domains not in the existing set.
func (c *Client) Run(ctx context.Context) ([]string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	chromeCtx, chromeCancel := chromedp.NewContext(allocCtx,
		chromedp.WithErrorf(func(string, ...interface{}) {}), // Suppress CDP unmarshal errors
	)
	defer chromeCancel()

	chromedp.ListenTarget(chromeCtx, func(ev interface{}) {
		if ev, ok := ev.(*network.EventRequestWillBeSent); ok {
			host := extractHost(ev.Request.URL)
			if host == "" {
				return
			}
			if isAdRelated(host, c.config.AdPatterns) && !c.existing[host] {
				c.mu.Lock()
				c.domains[host] = struct{}{}
				c.mu.Unlock()
			}
		}
	})

	tasks := chromedp.Tasks{
		network.Enable(),
		chromedp.Navigate("https://www.youtube.com"),
		chromedp.Sleep(5 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return chromedp.Navigate("https://www.youtube.com/watch?v=dQw4w9WgXcQ").Do(ctx)
		}),
		chromedp.Sleep(c.config.Duration),
	}

	if err := chromedp.Run(chromeCtx, tasks); err != nil {
		return nil, err
	}

	var result []string
	c.mu.Lock()
	for d := range c.domains {
		result = append(result, d)
	}
	c.mu.Unlock()

	sort.Strings(result)
	return result, nil
}

func extractHost(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	host := u.Hostname()
	if idx := strings.Index(host, ":"); idx >= 0 {
		host = host[:idx]
	}
	return strings.ToLower(host)
}

func isAdRelated(host string, patterns []string) bool {
	hostLower := strings.ToLower(host)
	for _, p := range patterns {
		if strings.Contains(hostLower, strings.TrimSpace(p)) {
			return true
		}
	}
	return false
}
