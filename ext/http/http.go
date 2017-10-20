package http

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/kataras/chronos"
)

type Client struct {
	*chronos.C
	*http.Client
}

func New(max uint32, per time.Duration) *Client {
	return &Client{
		C:      chronos.New(max, per),
		Client: NewTimeoutClient(20 * time.Second),
	}
}

// Do sends an HTTP request and returns an HTTP response, following
// policy (such as redirects, cookies, auth) as configured on the
// client.
//
// An error is returned if caused by client policy (such as
// CheckRedirect), or failure to speak HTTP (such as a network
// connectivity problem). A non-2xx status code doesn't cause an
// error.
//
// If the returned error is nil, the Response will contain a non-nil
// Body which the user is expected to close. If the Body is not
// closed, the Client's underlying RoundTripper (typically Transport)
// may not be able to re-use a persistent TCP connection to the server
// for a subsequent "keep-alive" request.
//
// The request Body, if non-nil, will be closed by the underlying
// Transport, even on errors.
//
// On error, any Response can be ignored. A non-nil Response with a
// non-nil error only occurs when CheckRedirect fails, and even then
// the returned Response.Body is already closed.
//
// Generally Get, Post, or PostForm will be used instead of Do.
//
// If the server replies with a redirect, the Client first uses the
// CheckRedirect function to determine whether the redirect should be
// followed. If permitted, a 301, 302, or 303 redirect causes
// subsequent requests to use HTTP method GET
// (or HEAD if the original request was HEAD), with no body.
// A 307 or 308 redirect preserves the original HTTP method and body,
// provided that the Request.GetBody function is defined.
// The NewRequest function automatically sets GetBody for common
// standard library body types.
func (hc *Client) Do(req *http.Request) (*http.Response, error) {
	<-hc.C.Acquire() // <-- the only part that this differs, except the ReadJSON/XML.
	return hc.Client.Do(req)
}

// Head issues a HEAD to the specified URL. If the response is one of the
// following redirect codes, Head follows the redirect after calling the
// Client's CheckRedirect function:
//
//    301 (Moved Permanently)
//    302 (Found)
//    303 (See Other)
//    307 (Temporary Redirect)
//    308 (Permanent Redirect)
func (hc *Client) Head(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return hc.Do(req)
}

// Get issues a GET to the specified URL. If the response is one of the
// following redirect codes, Get follows the redirect after calling the
// Client's CheckRedirect function:
//
//    301 (Moved Permanently)
//    302 (Found)
//    303 (See Other)
//    307 (Temporary Redirect)
//    308 (Permanent Redirect)
//
// An error is returned if the Client's CheckRedirect function fails
// or if there was an HTTP protocol error. A non-2xx response doesn't
// cause an error.
//
// When err is nil, resp always contains a non-nil resp.Body.
// Caller should close resp.Body when done reading from it.
//
// To make a request with custom headers, use NewRequest and Client.Do.
func (hc *Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return hc.Do(req)
}

// Post issues a POST to the specified URL.
//
// Caller should close resp.Body when done reading from it.
//
// If the provided body is an io.Closer, it is closed after the
// request.
//
// To set custom headers, use NewRequest and Client.Do.
//
// See the Client.Do method documentation for details on how redirects
// are handled.
func (hc *Client) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return hc.Do(req)
}

// PostForm issues a POST to the specified URL,
// with data's keys and values URL-encoded as the request body.
//
// The Content-Type header is set to application/x-www-form-urlencoded.
// To set other headers, use NewRequest and DefaultClient.Do.
//
// When err is nil, resp always contains a non-nil resp.Body.
// Caller should close resp.Body when done reading from it.
//
// See the Client.Do method documentation for details on how redirects
// are handled.
func (hc *Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return hc.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

type (
	// BodyDecoder is an interface which any struct can implement in order to customize the decode action
	// from ReadJSON and ReadXML
	//
	// Trivial example of this could be:
	// type User struct { Username string }
	//
	// func (u *User) Decode(data []byte) error {
	//	  return json.Unmarshal(data, u)
	// }
	//
	// the 'ReadJSON/ReadXML(&User{})' will call the User's
	// Decode option to decode the request body
	//
	// Note: This is totally optionally, the default decoders
	// for ReadJSON is the encoding/json and for ReadXML is the encoding/xml.
	BodyDecoder interface {
		Decode(data []byte) error
	}

	// Unmarshaler is the interface implemented by types that can unmarshal any raw data
	// TIP INFO: Any v object which implements the BodyDecoder can be override the unmarshaler.
	Unmarshaler interface {
		Unmarshal(data []byte, v interface{}) error
	}

	// UnmarshalerFunc a shortcut for the Unmarshaler interface
	//
	// See 'Unmarshaler' and 'BodyDecoder' for more.
	UnmarshalerFunc func(data []byte, v interface{}) error
)

// Unmarshal parses the X-encoded data and stores the result in the value pointed to by v.
// Unmarshal uses the inverse of the encodings that Marshal uses, allocating maps,
// slices, and pointers as necessary.
func (u UnmarshalerFunc) Unmarshal(data []byte, v interface{}) error {
	return u(data, v)
}

// unmarshalBody reads the request's body and binds it to a value or pointer of any type
// Examples of usage: ReadJSON, ReadXML.
func (hc *Client) unmarshalBody(body io.ReadCloser, v interface{}, unmarshaler Unmarshaler) error {
	if body == nil {
		return errors.New("unmarshal: empty body")
	}

	rawData, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	// check if the v contains its own decode
	// in this case the v should be a pointer also,
	// but this is up to the user's custom Decode implementation*
	//
	// See 'BodyDecoder' for more
	if decoder, isDecoder := v.(BodyDecoder); isDecoder {
		return decoder.Decode(rawData)
	}

	// check if v is already a pointer, if yes then pass as it's
	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		return unmarshaler.Unmarshal(rawData, v)
	}
	// finally, if the v doesn't contains a self-body decoder and it's not a pointer
	// use the custom unmarshaler to bind the body
	return unmarshaler.Unmarshal(rawData, &v)
}

// ReadJSON will fill the "v" from the GET: "url" body's content based on the "unmarshaler".
func (hc *Client) Read(url string, v interface{}, unmarshaler Unmarshaler) error {
	resp, err := hc.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("status code %d when fetching '%s'", resp.StatusCode, url)
	}

	return hc.unmarshalBody(resp.Body, v, unmarshaler)
}

// ReadJSON will fill the "v" from the GET: "url" body's JSON content.
func (hc *Client) ReadJSON(url string, v interface{}) error {
	return hc.Read(url, v, UnmarshalerFunc(json.Unmarshal))
}

// ReadXML will fill the "v" from the GET: "url" body's XML content.
func (hc *Client) ReadXML(url string, v interface{}) error {
	return hc.Read(url, v, UnmarshalerFunc(xml.Unmarshal))
}
