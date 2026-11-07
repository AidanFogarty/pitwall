package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/cdproto/storage"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

const (
	loginURL          = "https://account.formula1.com/#/en/login"
	sessionCookieName = "login-session"
	pollInterval      = 500 * time.Millisecond
)

var ErrLoginIncomplete = errors.New("the browser was closed before sign in completed")

func Login(ctx context.Context, browserPath string) (*Session, error) {
	profileDir, err := os.MkdirTemp("", "pitwall-login-*")
	if err != nil {
		return nil, fmt.Errorf("creating browser profile: %w", err)
	}
	defer os.RemoveAll(profileDir)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx,
		chromedp.ExecPath(browserPath),
		chromedp.UserDataDir(profileDir),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false),
	)
	defer cancelAlloc()

	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()

	if _, err := chromedp.Run(browserCtx, chromedp.Navigate(loginURL)); err != nil {
		return nil, fmt.Errorf("opening login page: %w", err)
	}

	session, err := waitForSession(browserCtx)
	if errors.Is(err, context.Canceled) && ctx.Err() == nil {
		return nil, ErrLoginIncomplete
	}
	return session, err
}

func waitForSession(ctx context.Context) (*Session, error) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}

		open, err := hasOpenWindow(ctx)
		if err != nil {
			return nil, fmt.Errorf("lost connection to the browser: %w", err)
		}
		if !open {
			return nil, ErrLoginIncomplete
		}

		session, found, err := findSession(ctx)
		if err != nil {
			return nil, fmt.Errorf("lost connection to the browser: %w", err)
		}
		if found {
			return session, nil
		}
	}
}

func hasOpenWindow(ctx context.Context) (bool, error) {
	result, err := chromedp.CallBrowser(ctx, target.GetTargets, target.GetTargetsParams{})
	if err != nil {
		return false, err
	}

	for _, info := range result.TargetInfos {
		if info.Type == "page" {
			return true, nil
		}
	}
	return false, nil
}

func findSession(ctx context.Context) (*Session, bool, error) {
	result, err := chromedp.CallBrowser(ctx, storage.GetCookies, storage.GetCookiesParams{})
	if err != nil {
		return nil, false, err
	}

	for _, cookie := range result.Cookies {
		if cookie.Name != sessionCookieName {
			continue
		}

		session, err := sessionFromCookie(cookie.Value)
		if err != nil || session.IsExpired(time.Now()) {
			continue
		}
		return session, true, nil
	}

	return nil, false, nil
}
