// Copyright 2013 Canonical Ltd.  This software is licensed under the
// GNU Lesser General Public License version 3 (see the file COPYING).

package gomaasapi

import (
	. "launchpad.net/gocheck"
)

func (suite *GomaasapiTestSuite) TestJoinURLsAppendsPathToHost(c *C) {
	c.Check(JoinURLs("http://example.com/", "foo"), Equals, "http://example.com/foo")
}

func (suite *GomaasapiTestSuite) TestJoinURLsAddsSlashIfNeeded(c *C) {
	c.Check(JoinURLs("http://example.com", "bar"), Equals, "http://example.com/bar")
}

func (suite *GomaasapiTestSuite) TestJoinURLsNormalizesDoubleSlash(c *C) {
	c.Check(JoinURLs("http://example.com/", "/szot"), Equals, "http://example.com/szot")
}
