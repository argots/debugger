// Package debugger implements a general purpose debugger server.
//
// The current implementation supports the ability to build a Chrome
// DevTools compatible server.
//
// There are plans to support VSCode at some point.
//
// Usage:
//
// See cmd/fake/fake.go for a fake debugger server.
package debugger

import cdp "github.com/chromedp/cdproto"

// Driver is the actual implementation that drivers need to implement.
type Driver interface {
	Name() string
	List() ([]*Target, error)
	Network() Network
	Dispatcher
}

// Dispatcher dispatches a message to the driver. The response is also
// just filled in directly.
type Dispatcher interface {
	Dispatch(d Driver, msg *cdp.Message)
}

// Target is a debugger target.  Use Init() to fill in the required fields.
type Target struct {
	// these fields are filled in by calling Init(prefix, ID)
	ID                string `json:"id"`
	FrontEndURL       string `json:"devtoolsFrontendUrl"`
	FrontEndURLCompat string `json:"devtoolsFrontendUrlCompat"`
	DebuggerURL       string `json:"webSocketDebuggerUrl"`

	// Optional fields
	Description string `json:"description"`
	FaviconURL  string `json:"faviconUrl"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

// Init must be called with a prefix like "http://localhost:9222/"
func (t *Target) Init(httpPrefix, id string) {
	t.ID = id
	wsPrefix := "ws" + httpPrefix[4:] + "/" + id
	wsEqualsPrefix := "ws=" + httpPrefix[5:] + "/" + id

	//nolint:lll
	t.FrontEndURL = `chrome-devtools://devtools/bundled/js_app.html?experiments=true&v8only=true&` + wsEqualsPrefix

	//nolint:lll
	t.FrontEndURLCompat = `chrome-devtools://devtools/bundled/inspector.html?experiments=true&v8only=true&` + wsEqualsPrefix
	t.DebuggerURL = wsPrefix
}
