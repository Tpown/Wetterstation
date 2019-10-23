package http

import "strings"

type UnparsableRequestError struct {
	s string
}

func (err *UnparsableRequestError) Error() string {
	return err.s
}

type Request struct {
	method      string            // GET oder POST
	path        string            //URL
	httpVersion string            //HTTP/1.0
	headers     map[string]string //e.g. Content-type: application/json
}

func (rq *Request) addHeader(key, val string) {
	rq.headers[key] = val
}

func (rq *Request) Path() string {
	return rq.path
}

func (rq *Request) Headers() map[string]string {
	return rq.headers
}

func (rq *Request) Method() string {
	return rq.method
}

func (rq *Request) HTTPVersion() string {
	return rq.httpVersion
}

func parseRequest(requestLines []string) (*Request, error) {
	var request Request
	request.headers = make(map[string]string)
	if len(requestLines) < 1 {
		return nil, &UnparsableRequestError{"Error Parsing the Header Line"}
	}
	headerLineFields := strings.Fields(requestLines[0])
	if len(headerLineFields) < 3 {
		return nil, &UnparsableRequestError{"Error Parsing the Header Line"}
	}
	request.method = headerLineFields[0]
	request.path = headerLineFields[1]
	request.httpVersion = headerLineFields[2]

	for i := 1; i < len(requestLines); i++ {
		tmp := strings.Split(requestLines[i], ":")
		if len(tmp) < 2 {
			return nil, &UnparsableRequestError{"Error Parsing the Headers"}
		}
		request.addHeader(strings.TrimSpace(tmp[0]), strings.TrimSpace(tmp[1]))
		//request.addHeader("Content-Type", "application/json")
	}
	return &request, nil
}
