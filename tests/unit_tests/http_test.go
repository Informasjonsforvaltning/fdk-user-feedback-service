package unit_tests

import (
	"testing"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/util"
)

func TestSuccsessfulStatus(t *testing.T) {
	var tests = []struct {
		testName   string
		statusCode int
		expected   bool
	}{
		{"OK", 200, true},
		{"Created", 201, true},
		{"Accepted", 202, true},
		{"Non-authoritative Information", 203, true},
		{"No Content", 204, true},
		{"Reset Content", 205, true},
		{"Partial Content", 206, true},
		{"Multi-Status", 207, true},
		{"Already Reported", 208, true},
		{"IM Used", 226, true},
		{"Continue", 100, false},
		{"Switching Protocols", 101, false},
		{"Processing", 102, false},
		{"Multiple Choices", 300, false},
		{"Moved Permanently", 301, false},
		{"Found", 302, false},
		{"See Other", 303, false},
		{"Not Modified", 304, false},
		{"Use Proxy", 305, false},
		{"Temporary Redirect", 307, false},
		{"Permanent Redirect", 308, false},
		{"Bad Request", 400, false},
		{"Unauthorized", 401, false},
		{"Payment Required", 402, false},
		{"Forbidden", 403, false},
		{"Not Found", 404, false},
		{"Method Not Allowed", 405, false},
		{"Not Acceptable", 406, false},
		{"Proxy Authentication Required", 407, false},
		{"Request Timeout", 408, false},
		{"Conflict", 409, false},
		{"Gone", 410, false},
		{"Length Required", 411, false},
		{"Precondition Failed", 412, false},
		{"Payload Too Large", 413, false},
		{"Request-URI Too Long", 414, false},
		{"Unsupported Media Type", 415, false},
		{"Requested Range Not Satisfiable", 416, false},
		{"Expectation Failed", 417, false},
		{"I'm a teapot", 418, false},
		{"Misdirected Request", 421, false},
		{"Unprocessable Entity", 422, false},
		{"Locked", 423, false},
		{"Failed Dependency", 424, false},
		{"Upgrade Required", 426, false},
		{"Precondition Required", 428, false},
		{"Too Many Requests", 429, false},
		{"Request Header Fields Too Large", 431, false},
		{"Connection Closed Without Response", 444, false},
		{"Unavailable For Legal Reasons", 451, false},
		{"Client Closed Request", 499, false},
		{"Internal Server Error", 500, false},
		{"Not Implemented", 501, false},
		{"Bad Gateway", 502, false},
		{"Service Unavailable", 503, false},
		{"Gateway Timeout", 504, false},
		{"HTTP Version Not Supported", 505, false},
		{"Variant Also Negotiates", 506, false},
		{"Insufficient Storage", 507, false},
		{"Loop Detected", 508, false},
		{"Not Extended", 510, false},
		{"Network Authentication Required", 511, false},
		{"Network Connect Timeout Error", 599, false},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			actual := util.SuccsessfulStatus(test.statusCode)
			if actual != test.expected {
				t.Errorf("expected %t, got %t", test.expected, actual)
			}
		})
	}
}

func compareStringPointers(a *string, b *string) bool {
	if a == b {
		return true
	}

	if a != nil && b != nil {
		return *a == *b
	}

	return false
}

func TestParseRequestUrlPath(t *testing.T) {
	testRoute := "r"
	testEntity := "e"
	testPost := "p"

	var tests = []struct {
		testPath             string
		expectedRouteString  *string
		expectedEntityString *string
		expectedPostString   *string
	}{
		{"/r/e/p", &testRoute, &testEntity, &testPost},
		{"r/e/p", &testEntity, &testPost, nil},
		{"/r/e", &testRoute, &testEntity, nil},
		{"/r/", &testRoute, nil, nil},
		{"/r", &testRoute, nil, nil},
		{"/", nil, nil, nil},
		{"", nil, nil, nil},
	}

	for _, test := range tests {
		t.Run(test.testPath, func(t *testing.T) {
			actualRouteString, actualEntityString, actualPostString := util.ParseRequestUrlPath(test.testPath)
			if !compareStringPointers(actualRouteString, test.expectedRouteString) || !compareStringPointers(actualEntityString, test.expectedEntityString) || !compareStringPointers(actualPostString, test.expectedPostString) {
				t.Errorf("expected %v %v %v, got %v %v %v", test.expectedRouteString, test.expectedEntityString, test.expectedPostString, actualRouteString, actualPostString, actualEntityString)
			}
		})
	}
}
