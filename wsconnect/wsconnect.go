package wsconnect

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/alexbrainman/sspi/negotiate"
	"github.com/gorilla/websocket"
)

// Connect establishes a websocket connection through a Windows proxy using automatic SSPI authentication.
// It uses the current Windows user's context for authentication, similar to curl's "-U :" functionality.
func Connect(targetURL string, proxyURL string) (*websocket.Conn, error) {
	if targetURL == "" || proxyURL == "" {
		return nil, fmt.Errorf("targetURL and proxyURL cannot be empty")
	}

	wsURLParsed, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse WebSocket URL '%s': %w", targetURL, err)
	}

	proxyURLParsed, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse proxy URL '%s': %w", proxyURL, err)
	}

	creds, err := negotiate.AcquireCurrentUserCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to acquire Windows credentials: %w", err)
	}
	defer creds.Release()

	secContext, token, err := negotiate.NewClientContext(creds, "HTTP/"+proxyURLParsed.Host)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSPI security context for '%s': %w", proxyURLParsed.Host, err)
	}
	defer secContext.Release()

	// Configure dialer with timeout
	dialer := websocket.Dialer{
		Proxy:            http.ProxyURL(proxyURLParsed),
		HandshakeTimeout: 10 * time.Second,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}

	headers := http.Header{}
	headers.Add("Proxy-Authorization", "Negotiate "+base64.StdEncoding.EncodeToString(token))
	headers.Add("Origin", wsURLParsed.Scheme+"://"+wsURLParsed.Host)

	conn, resp, err := dialer.Dial(targetURL, headers)
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("WebSocket connection failed with status %d: %w\nResponse Headers: %v",
				resp.StatusCode, err, resp.Header)
		}
		return nil, fmt.Errorf("failed to establish WebSocket connection: %w", err)
	}

	return conn, nil
}
