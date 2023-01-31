package gominitrader

import "net/http"

type AuthenticationTransport struct {
	http.RoundTripper
	X_SECURITY_TOKEN string
	CST              string
}

func (t *AuthenticationTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-SECURITY-TOKEN", t.X_SECURITY_TOKEN)
	req.Header.Set("CST", t.CST)
	return t.RoundTripper.RoundTrip(req)
}
