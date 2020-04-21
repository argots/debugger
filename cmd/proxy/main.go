//nolint: gochecknoglobals
// chromedp-proxy provides a cli utility that will proxy requests from a Chrome
// DevTools Protocol client to a browser instance.
//
// chromedp-proxy is particularly useful for recording events/data from
// Selenium (ChromeDriver), Chrome DevTools in the browser, or for debugging
// remote application instances compatible with the devtools protocol.
//
// Please see README.md for more information on using chromedp-proxy.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	flagListen  = flag.String("l", "localhost:9223", "listen address")
	flagRemote  = flag.String("r", "localhost:9222", "remote address")
	flagNoLog   = flag.Bool("n", false, "disable logging to file")
	flagLogMask = flag.String("log", "logs/cdp-%s.log", "log file mask")
)

const (
	incomingBufferSize = 10 * 1024 * 1024
	outgoingBufferSize = 25 * 1024 * 1024
)

var wsUpgrader = &websocket.Upgrader{
	ReadBufferSize:  incomingBufferSize,
	WriteBufferSize: outgoingBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var wsDialer = &websocket.Dialer{
	ReadBufferSize:  outgoingBufferSize,
	WriteBufferSize: incomingBufferSize,
}

func main() { //nolint:funlen
	flag.Parse()

	mux := http.NewServeMux()
	simplep := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: *flagRemote})
	mux.Handle("/json", simplep)
	mux.Handle("/json/", simplep)
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		id := path.Base(req.URL.Path)
		if id == "" {
			simplep.ServeHTTP(res, req)
			return
		}
		f, logger := createLog(id)
		if f != nil {
			defer f.Close()
		}
		logger.Printf("---------- connection from %s ----------", req.RemoteAddr)

		ver, err := checkVersion()
		if err != nil {
			msg := fmt.Sprintf("version error, got: %v", err)
			logger.Println(msg)
			http.Error(res, msg, http.StatusInternalServerError)
			return
		}
		logger.Printf("endpoint %s reported: %s", *flagRemote, string(ver))

		endpoint := "ws://" + *flagRemote + path.Join(path.Dir(req.URL.Path), id)

		// connect outgoing websocket
		logger.Printf("connecting to %s", endpoint)
		out, pres, err := wsDialer.Dial(endpoint, nil)
		if err != nil {
			msg := fmt.Sprintf("could not connect to %s, got: %v", endpoint, err)
			logger.Println(msg)
			http.Error(res, msg, http.StatusInternalServerError)
			return
		}
		defer pres.Body.Close()
		defer out.Close()

		logger.Printf("connected to %s", endpoint)

		// connect incoming websocket
		logger.Printf("upgrading connection on %s", req.RemoteAddr)
		in, err := wsUpgrader.Upgrade(res, req, nil)
		if err != nil {
			msg := fmt.Sprintf("could not upgrade websocket from %s, got: %v", req.RemoteAddr, err)
			logger.Println(msg)
			http.Error(res, msg, http.StatusInternalServerError)
			return
		}
		defer in.Close()
		logger.Printf("upgraded connection on %s", req.RemoteAddr)

		ctxt, cancel := context.WithCancel(context.Background())
		defer cancel()

		errc := make(chan error, 1)
		go proxyWS(ctxt, logger, "<-", in, out, errc)
		go proxyWS(ctxt, logger, "->", out, in, errc)
		<-errc
		logger.Printf("---------- closing %s ----------", req.RemoteAddr)
	})

	log.Fatal(http.ListenAndServe(*flagListen, mux))
}

// proxyWS proxies in and out messages for a websocket connection, logging the
// message to the logger with the passed prefix. Any error encountered will be
// sent to errc.
func proxyWS(ctxt context.Context, logger *log.Logger, prefix string, in, out *websocket.Conn, errc chan error) {
	var mt int
	var buf []byte
	var err error

	for {
		select {
		default:
			mt, buf, err = in.ReadMessage()
			if err != nil {
				errc <- err
				return
			}

			logger.Println(prefix, string(buf))

			err = out.WriteMessage(mt, buf)
			if err != nil {
				errc <- err
				return
			}

		case <-ctxt.Done():
			return
		}
	}
}

// checkVersion retrieves the version information for the remote endpoint, and
// formats it appropriately.
func checkVersion() ([]byte, error) {
	cl := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+*flagRemote+"/json/version", nil)
	if err != nil {
		return nil, err
	}

	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var v map[string]string
	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, fmt.Errorf("expected json result")
	}

	return body, nil
}

var (
	idCleanRE = regexp.MustCompile(`[^a-zA-Z0-9_\-\.]`)
)

// createLog creates a log for the specified id based on flags.
func createLog(id string) (io.Closer, *log.Logger) {
	var f io.Closer
	var w io.Writer = os.Stdout
	if !*flagNoLog && *flagLogMask != "" {
		filename := *flagLogMask
		if strings.Contains(*flagLogMask, "%s") {
			filename = fmt.Sprintf(*flagLogMask, idCleanRE.ReplaceAllString(id, ""))
		}
		l, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		f, w = l, io.MultiWriter(os.Stdout, l)
	}
	return f, log.New(w, "", log.LstdFlags)
}
