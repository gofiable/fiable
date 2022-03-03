package minima

/**
* Minima is a free and open source software under Mit license

Copyright (c) 2021 gominima

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

* Authors @apoorvcodes @megatank58
* Maintainers @Panquesito7 @savioxavier @Shubhaankar-Sharma @apoorvcodes @megatank58
* Thank you for showing interest in minima and for this beautiful community
*/

import (
	"net/http"
)

/**
 * @info The Outgoing header structure
 * @property {http.Request} [req] The net/http request instance
 * @property {http.ResponseWriter} [res] The net/http response instance
 * @property {bool} [body] Whether body has been sent or not
 * @property {int} [status] response status code
*/
type OutgoingHeader struct {
	req *http.Request
	res http.ResponseWriter
}

var statusCodes = map[string]int{
	"OK":                         200,
	"Created":                    201,
	"Accepted":                   202,
	"No Content":                 204,
	"Reset Content":              205,
	"Partial Content":            206,
	"Moved Permanently":          301,
	"Found":                      302,
	"Not Modified":               304,
	"Use Proxy":                  305,
	"Switch Proxy":               306,
	"Temporary Redirect":         307,
	"Permanent Redirect":         308,
	"Bad Request":                400,
	"Unauthorized":               401,
	"Forbidden":                  403,
	"NOT FOUND":                  404,
	"Method Not Allowed":         405,
	"Payload Too Large":          413,
	"URI Too Long":               414,
	"Internal Server Error":      500,
	"Not Implemented":            501,
	"Bad Gateway":                502,
	"Service Unavailable":       503,
	"Gateway Timeout":            504,
	"HTTP Version Not Supported": 505,
}

/**
 * @info Make a new default request header instance
 * @param {http.Request} [req] The net/http request instance
 * @param {http.ResponseWriter} [res] The net/http response instance
 * @returns {OutgoingHeader}
*/
func NewResHeader(res http.ResponseWriter, req *http.Request) *OutgoingHeader {
	return &OutgoingHeader{req, res}
}

/**
 * @info Sets and new header to response
 * @param {string} [key] Key of the new header
 * @param {string} [value] Value of the new header
 * @returns {OutgoingHeader}
*/
func (h *OutgoingHeader) Set(key string, value string) *OutgoingHeader {
	h.res.Header().Set(key, value)
	return h
}

/**
 * @info Gets the header from response headers
 * @param {string} [key] Key of the header
 * @returns {string}
*/
func (h *OutgoingHeader) Get(key string) string {
	return h.res.Header().Get(key)
}

/**
 * @info Deletes header from respose
 * @param {string} [key] Key of the header
 * @returns {OutgoingHeader}
*/
func (h *OutgoingHeader) Del(key string) *OutgoingHeader {
	h.res.Header().Del(key)
	return h
}

/**
 * @info Clones all headers from response
 * @returns {OutgoingHeader}
*/
func (h *OutgoingHeader) Clone() http.Header {
	return h.res.Header().Clone()
}

/**
 * @info Sets content lenght
 * @param {string} [len] The lenght of the content
 * @returns {OutgoingHeader}
*/
func (h *OutgoingHeader) Setlength(len string) *OutgoingHeader {
	h.Set("Content-length", len)
	return h
}

/**
 * @info Sets response status
 * @param {int} [code] The status code for the response
 * @returns {OutgoingHeader}
*/
func (h *OutgoingHeader) Status(code int) *OutgoingHeader {
	h.res.WriteHeader(code)
	return h
}

/**
 * @info Sends good stack of base headers
 * @returns {}
*/
func (h *OutgoingHeader) BaseHeaders() {
	h.Set("transfer-encoding", "chunked")
	h.Set("connection", "keep-alive")
}

/**
 * @info Flushes and writes header to route
 * @returns {bool}
*/
func (h *OutgoingHeader) Flush() bool {
	if h.Get("Content-Type") == "" {
		h.Set("Content-Type", "text/html;charset=utf-8")
	}

	if f, ok := h.res.(http.Flusher); ok {
		f.Flush()
	}

	return true
}
