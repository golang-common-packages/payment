package payment

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

// SendWithBasicAuth makes a request to the API using clientID:secret basic auth
func (c *PayPalClient) SendWithBasicAuth(req *http.Request, v interface{}) error {
	req.SetBasicAuth(c.ClientID, c.Secret)

	return c.Send(req, v)
}

// Send makes a request to the API, the response body will be
// unmarshalled into v, or if v is an io.Writer, the response will
// be written to it without decoding
func (c *PayPalClient) Send(req *http.Request, v interface{}) error {
	var (
		err  error
		resp *http.Response
		data []byte
	)

	// Set default headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en_US")

	// Default values for headers
	if req.Header.Get("Content-type") == "" {
		req.Header.Set("Content-type", "application/json")
	}
	if c.returnRepresentation {
		req.Header.Set("Prefer", "return=representation")
	}

	resp, err = c.Client.Do(req)
	c.log(req, resp)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errResp := &ErrorResponse{Response: resp}
		data, err = ioutil.ReadAll(resp.Body)

		if err == nil && len(data) > 0 {
			json.Unmarshal(data, errResp)
		}

		return errResp
	}
	if v == nil {
		return nil
	}

	if w, ok := v.(io.Writer); ok {
		io.Copy(w, resp.Body)
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// log will dump request and response to the log file
func (c *PayPalClient) log(r *http.Request, resp *http.Response) {
	if c.Log != nil {
		var (
			reqDump  string
			respDump []byte
		)

		if r != nil {
			reqDump = fmt.Sprintf("%s %s. Data: %s", r.Method, r.URL.String(), r.Form.Encode())
		}
		if resp != nil {
			respDump, _ = httputil.DumpResponse(resp, true)
		}

		c.Log.Write([]byte(fmt.Sprintf("Request: %s\nResponse: %s\n", reqDump, string(respDump))))
	}
}

// Error method implementation for ErrorResponse struct
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %s, %+v", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message, r.Details)
}