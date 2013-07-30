package main

import (
	"bytes"
	"encoding/json"
	"github.com/coocood/jas"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"net/url"
)

type HttpSuite struct {
	srv *httptest.Server
	mux *http.ServeMux
	uri *url.URL
}

var _ = Suite(&HttpSuite{})

func (s *HttpSuite) SetUpSuite(c *C) {
	s.mux = buildMux()
	s.srv = httptest.NewServer(s.mux)
	s.uri, _ = url.Parse(s.srv.URL)

}

func (s *HttpSuite) TearDownSuite(c *C) {
	s.srv.Close()
}

func (s *HttpSuite) request(c *C, path, data string) (*jas.Response, error) {
	uri := *s.uri
	uri.Path = path

	resp, err := http.Post(uri.String(), "application/json", bytes.NewBufferString(data))
	if err != nil {
		c.Fatalf("Could not post %s: %s", uri.String(), err)
	}

	defer resp.Body.Close()

	c.Check(resp.StatusCode, Equals, 200)

	if p, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		js := new(jas.Response)
		err = json.Unmarshal(p, &js)
		return js, err
	}
}

func (s *HttpSuite) TestNotify(c *C) {

	resp, err := s.request(c, "/api/v1/notify/example.com", "{}")
	c.Assert(err, IsNil)
	c.Check(resp.Error, Equals, "No servers")

	// no DNS server here, connection refused
	servers = []string{"127.0.0.1:55"}

	resp, err = s.request(c, "/api/v1/notify/example.com", "{}")
	c.Assert(err, IsNil)
	c.Check(resp.Error, Matches, "^read udp 127.0.0.1:.*: connection refused")

}
