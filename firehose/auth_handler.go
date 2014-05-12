package firehose

import (
  "log"
  "encoding/base64"
  "net/http"
  "strings"
  "errors"
)

type AuthHandler struct {
  app http.Handler
  token string
}

func (handler *AuthHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

  unauthorized, reason := handler.authorize(req)

  if unauthorized == true {
    log.Println("Unauthorised:", reason)
    http.Error(rw, reason, http.StatusUnauthorized)
    return
  }

  handler.app.ServeHTTP(rw, req)
}

func (handler *AuthHandler) authorize(req *http.Request) (unauthorized bool, reason string) {
  scheme, credentials, _ := ParseRequest(req)

  switch scheme {
  case "Bearer":
    if credentials != handler.token {
      unauthorized = true
      reason = "Invalid Access Token"
    }
  case "Basic":
    basic, err := NewBasic(credentials)
    if err != nil {
      unauthorized = true
      reason = "Malformed Basic Authorization crdentials"
    }
    if basic.Password != handler.token {
      unauthorized = true
      reason = "Invalid credentials"
    }
  default:
    // try the 'access_token' query param
    req.ParseForm()
    if len(req.Form["access_token"]) == 0 || req.Form["access_token"][0] != handler.token {
      unauthorized = true
      reason = "Mising or invalid access_token"
    }
  }

  return
}

func NewAuthHandler(app http.Handler, token string) (handler *AuthHandler) {

  handler = &AuthHandler{
    app: app,
    token: token,
  }

  return
}

// ParseRequest extracts an "Authorization" header from a request and returns
// its scheme and credentials.
func ParseRequest(r *http.Request) (scheme, credentials string, err error) {
  h, ok := r.Header["Authorization"]
  if !ok || len(h) == 0 {
          return "", "", errors.New("The authorization header is not set.")
  }
  return Parse(h[0])
}

// Parse parses an "Authorization" header and returns its scheme and
// credentials.
func Parse(value string) (scheme, credentials string, err error) {
  parts := strings.SplitN(value, " ", 2)
  if len(parts) == 2 {
          return parts[0], parts[1], nil
  }
  return "", "", errors.New("The authorization header is malformed.")
}

// NewBasic parses credentials from a "basic" http authentication scheme.
func NewBasic(credentials string) (*Basic, error) {
  if b, err := base64.StdEncoding.DecodeString(credentials); err == nil {
    parts := strings.Split(string(b), ":")
    if len(parts) == 2 {
      return &Basic{
        Username: parts[0],
        Password: parts[1],
      }, nil
    }
  }
  return nil, errors.New("The basic authentication header is malformed.")
}

// Basic stores username and password for the "basic" http authentication
// scheme. Reference:
//
//    http://tools.ietf.org/html/rfc2617#section-2
type Basic struct {
  Username string
  Password string
}