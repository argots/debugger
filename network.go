package debugger

import nw "github.com/chromedp/cdproto/network"

// See https://github.com/chromedp/cdproto/blob/master/network/network.go
type Network interface {
	ClearBrowserCache() error
	ClearBrowserCookies() error
	DeleteCookies(req *nw.DeleteCookiesParams) error
	Disable() error
	EmulateNetworkConditions(req *nw.EmulateNetworkConditionsParams) error
	Enable() (*nw.EnableParams, error)
	GetAllCookies() (*nw.GetAllCookiesReturns, error)
	GetCertificate(req *nw.GetCertificateParams) (*nw.GetCertificateReturns, error)
	GetCookies(req *nw.GetCookiesParams) (*nw.GetCookiesReturns, error)
	GetResponseBody(req *nw.GetResponseBodyParams) (*nw.GetResponseBodyReturns, error)
	GetRequestPostData(req *nw.GetRequestPostDataParams) (*nw.GetRequestPostDataReturns, error)

	/*** GetResponseBodyForInterception has some streaming
		/*** business, take it out for now.
		GetResponseBodyForInterception(req
		/*** *nw.GetResponseBodyForInterceptionParams)
		/*** *nw.GetResponseBodyForInterceptionReturns
	        ***/

	ReplayXHR(req *nw.ReplayXHRParams) error
	SearchInResponseBody(req *nw.SearchInResponseBodyParams) (*nw.SearchInResponseBodyReturns, error)
	SetBlockedURLs(req *nw.SetBlockedURLSParams) error
	SetBypassServiceWorker(req *nw.SetBypassServiceWorkerParams) error
	SetCacheDisabled(req *nw.SetCacheDisabledParams) error
	SetCookie(req *nw.SetCookieParams) (*nw.SetCookieReturns, error)
	SetCookies(req *nw.SetCookiesParams) error
	SetDataSizeLimitsForTest(req *nw.SetDataSizeLimitsForTestParams) error
	SetExtraHTTPHeaders(req *nw.SetExtraHTTPHeadersParams) error
}

var _ Network = FakeNetwork{}

type FakeNetwork struct{}

func (FakeNetwork) ClearBrowserCache() error {
	return nil
}

func (FakeNetwork) ClearBrowserCookies() error {
	return nil
}

func (FakeNetwork) DeleteCookies(req *nw.DeleteCookiesParams) error {
	return nil
}

func (FakeNetwork) Disable() error {
	return nil
}

func (FakeNetwork) EmulateNetworkConditions(req *nw.EmulateNetworkConditionsParams) error {
	return nil
}

func (FakeNetwork) Enable() (*nw.EnableParams, error) {
	return &nw.EnableParams{}, nil
}

func (FakeNetwork) GetAllCookies() (*nw.GetAllCookiesReturns, error) {
	return &nw.GetAllCookiesReturns{}, nil
}

func (FakeNetwork) GetCertificate(req *nw.GetCertificateParams) (*nw.GetCertificateReturns, error) {
	return &nw.GetCertificateReturns{}, nil
}

func (FakeNetwork) GetCookies(req *nw.GetCookiesParams) (*nw.GetCookiesReturns, error) {
	return &nw.GetCookiesReturns{}, nil
}

func (FakeNetwork) GetResponseBody(req *nw.GetResponseBodyParams) (*nw.GetResponseBodyReturns, error) {
	return &nw.GetResponseBodyReturns{}, nil
}

func (FakeNetwork) GetRequestPostData(req *nw.GetRequestPostDataParams) (*nw.GetRequestPostDataReturns, error) {
	return &nw.GetRequestPostDataReturns{}, nil
}

func (FakeNetwork) ReplayXHR(req *nw.ReplayXHRParams) error {
	return nil
}

func (FakeNetwork) SearchInResponseBody(req *nw.SearchInResponseBodyParams) (*nw.SearchInResponseBodyReturns, error) {
	return &nw.SearchInResponseBodyReturns{}, nil
}

func (FakeNetwork) SetBlockedURLs(req *nw.SetBlockedURLSParams) error {
	return nil
}

func (FakeNetwork) SetBypassServiceWorker(req *nw.SetBypassServiceWorkerParams) error {
	return nil
}

func (FakeNetwork) SetCacheDisabled(req *nw.SetCacheDisabledParams) error {
	return nil
}

func (FakeNetwork) SetCookie(req *nw.SetCookieParams) (*nw.SetCookieReturns, error) {
	return &nw.SetCookieReturns{}, nil
}

func (FakeNetwork) SetCookies(req *nw.SetCookiesParams) error {
	return nil
}

func (FakeNetwork) SetDataSizeLimitsForTest(req *nw.SetDataSizeLimitsForTestParams) error {
	return nil
}

func (FakeNetwork) SetExtraHTTPHeaders(req *nw.SetExtraHTTPHeadersParams) error {
	return nil
}
