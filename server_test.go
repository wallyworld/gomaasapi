// Copyright 2013 Canonical Ltd.  This software is licensed under the
// GNU Lesser General Public License version 3 (see the file COPYING).

package gomaasapi

import (
	. "launchpad.net/gocheck"
)

func (suite *GomaasapiTestSuite) TestNewMAASParsesURL(c *C) {
	maas, err := NewMAAS("https://server.com:888/path/to/api", Client{})

	c.Check(err, IsNil)
	c.Check(maas.URL(), Equals, "https://server.com:888/path/to/api")
	jsonObj := maas.(jsonMAASObject)
	uri, err := jsonObj._URI()
	c.Check(err, IsNil)
	c.Check(uri, Equals, "/path/to/api")
	c.Check(jsonObj.baseURL, Equals, "https://server.com:888")
}