// Package discover provides YouTube ad domain discovery via browser network capture.
package discover

import (
	"context"
	"net/url"
	"os"
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

// DefaultVideoURLs are YouTube URLs to visit for ad discovery (trending + popular videos with lots of ads).
var DefaultVideoURLs = []string{
	"https://www.youtube.com/feed/trending",
	"https://www.youtube.com/watch?v=9bZkp7q19f0",  // Gangnam Style
	"https://www.youtube.com/watch?v=kJQP7kiw5Fk",  // Despacito
	"https://www.youtube.com/watch?v=RgKAFK5djSk",  // See You Again
	"https://www.youtube.com/watch?v=OPf0YbXqDm0",  // Uptown Funk
	"https://www.youtube.com/watch?v=09R8_2nJtjg",  // Sugar
	"https://www.youtube.com/watch?v=JGwWNGJdvx8",  // Shape of You
	"https://www.youtube.com/watch?v=kxopViU98Xo",  // Baby Shark
	"https://www.youtube.com/watch?v=dQw4w9WgXcQ",  // Never Gonna Give You Up
}

// Config holds configuration for the discovery client.
type Config struct {
	DurationPerVideo time.Duration // How long to capture traffic per video. 0 = 1 minute.
	VideoURLs        []string      // YouTube URLs to visit (trending, popular videos). Empty = use DefaultVideoURLs.
	Blocklist        string
	AdPatterns       []string
	ChromePath       string // Path to Chrome/Chromium binary (e.g. for CI). Empty = auto-detect.
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
	urls := cfg.VideoURLs
	if len(urls) == 0 {
		urls = DefaultVideoURLs
	}
	dur := cfg.DurationPerVideo
	if dur == 0 {
		dur = time.Minute
	}
	return &Client{
		config: Config{
			DurationPerVideo: dur,
			VideoURLs:        urls,
			Blocklist:        cfg.Blocklist,
			AdPatterns:       patterns,
			ChromePath:       cfg.ChromePath,
		},
		existing: existing,
		domains:  make(map[string]struct{}),
	}
}

// Run navigates to YouTube, visits each configured video URL, captures network traffic
// for the configured duration per video, and returns newly discovered ad-related domains.
func (c *Client) Run(ctx context.Context) ([]string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)
	if path := c.config.ChromePath; path != "" {
		opts = append(opts, chromedp.ExecPath(path))
	} else if path := os.Getenv("CHROME_PATH"); path != "" {
		opts = append(opts, chromedp.ExecPath(path))
	}

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

	// Initial load
	if err := chromedp.Run(chromeCtx,
		network.Enable(),
		chromedp.Navigate("https://www.youtube.com"),
		chromedp.Sleep(5*time.Second),
	); err != nil {
		return nil, err
	}

	// Visit each video and capture traffic
	for i, videoURL := range c.config.VideoURLs {
		if err := chromedp.Run(chromeCtx,
			chromedp.Navigate(videoURL),
			chromedp.Sleep(5*time.Second), // Let page and ads load
			chromedp.Sleep(c.config.DurationPerVideo),
		); err != nil {
			return nil, err
		}
		_ = i // avoid unused if we add logging
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
