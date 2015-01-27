package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "gopkg.in/check.v1"
)

type HTTPSuite struct {
	srv *httptest.Server
	mux *http.ServeMux
	uri *url.URL
}

var _ = Suite(&HTTPSuite{})

func (s *HTTPSuite) SetUpSuite(c *C) {
	s.mux = buildMux()
	s.srv = httptest.NewServer(s.mux)
	s.uri, _ = url.Parse(s.srv.URL)

}

func (s *HTTPSuite) TearDownSuite(c *C) {
	s.srv.Close()
}

type jsResponse struct {
	http *http.Response
	js   map[string]interface{}
}

func (s *HTTPSuite) request(c *C, path, data string) (*jsResponse, error) {
	uri := *s.uri
	uri.Path = path

	req, err := http.NewRequest("POST", uri.String(), bytes.NewBufferString(data))
	if err != nil {
		c.Fatalf("Could build request %s: %s", uri.String(), err)
	}
	req.Header.Add("Accept-Encoding", "identity")
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.Fatalf("Could not post %s: %s", uri.String(), err)
	}

	defer resp.Body.Close()

	c.Check(resp.StatusCode, Equals, 200)

	if p, err := ioutil.ReadAll(resp.Body); err != nil {
		log.Println("error from ReadAll(), read", string(p))
		return nil, err
	} else {

		r := new(jsResponse)

		r.http = resp

		err = json.Unmarshal(p, &r.js)
		return r, err
	}
}

func (s *HTTPSuite) TestNotify(c *C) {

	resp, err := s.request(c, "/api/v1/notify/example.com", "{}")
	c.Assert(err, IsNil)
	c.Check(resp.js["Error"], Equals, "No servers")

	// no DNS server here, connection refused
	servers = []string{"127.0.0.1:55"}

	resp, err = s.request(c, "/api/v1/notify/example.com", "{}")
	c.Assert(err, IsNil)
	c.Check(resp.js["Error"], Matches, "^read udp 127.0.0.1:.*: connection refused")

}
