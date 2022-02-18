package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

const app1v1 = `
title: App1
version: 0.0.1
maintainers:
- name: firstmaintainer app1
  email: firstmaintainer@hotmail.com
- name: secondmaintainer app1
  email: secondmaintainer@gmail.com
company: Random Inc.
website: https://website.com
source: https://github.com/random/repo
license: Apache-2.0
description: |
 ### Interesting Title
 Some application content, and description
`
const app1v2 = `
title: App1
version: 0.0.2
maintainers:
- name: firstmaintainer app1
  email: firstmaintainer@hotmail.com
- name: secondmaintainer app1
  email: secondmaintainer@gmail.com
company: Random Inc.
website: https://website.com
source: https://github.com/random/repo
license: Apache-2.0
description: |
 ### Interesting Title
 Some application content, and description
`
const app2v1 = `
title: App2
version: 0.0.1
maintainers:
- name: firstmaintainer app2
  email: firstmaintainer@hotmail.com
- name: secondmaintainer app2
  email: secondmaintainer@gmail.com
company: Random Inc.
website: https://website.com
source: https://github.com/random/repo
license: Apache-2.0
description: |
 ### Interesting Title
 Some application content, and description
`
const app3WithSpaceInTitle = `
title: App3 with space in title
version: 0.0.1
maintainers:
- name: firstmaintainer app3
  email: firstmaintainer@hotmail.com
- name: secondmaintainer app3
  email: secondmaintainer@gmail.com
company: Random Inc.
website: https://website.com
source: https://github.com/random/repo
license: Apache-2.0
description: |
 ### Interesting Title
 Some application content, and description
`

const appWithoutTitle = `
version: 0.0.1
maintainers:
- name: firstmaintainer app1
  email: firstmaintainer@hotmail.com
- name: secondmaintainer app1
  email: secondmaintainer@gmail.com
company: Random Inc.
website: https://website.com
source: https://github.com/random/repo
license: Apache-2.0
description: |
 ### Interesting Title
 Some application content, and description
`
const appWithoutVersion = `
title: App4
maintainers:
- name: firstmaintainer app1
  email: firstmaintainer@hotmail.com
- name: secondmaintainer app1
  email: secondmaintainer@gmail.com
company: Random Inc.
website: https://website.com
source: https://github.com/random/repo
license: Apache-2.0
description: |
 ### Interesting Title
 Some application content, and description
`
const appWithBadMaintainerEmail = `
title: App5
version: 0.0.1
maintainers:
- name: firstmaintainer app1
  email: firstmaintainer
company: Random Inc.
website: https://website.com
source: https://github.com/random/repo
license: Apache-2.0
description: |
 ### Interesting Title
 Some application content, and description
`
const appWithMultipleBadFields = `
title: App5
version: 0.0.1
company: Random Inc.
website: abc
source: https://github.com/random/repo
license: Apache-2.0
description: |
 ### Interesting Title
 Some application content, and description
`

type Request struct {
	method string
	url    string
	data   string
}

var scenarios = []struct {
	scenarioTitle        string
	requests             []Request
	expectedStatusCode   int
	expectedResponseBody string
}{
	// -----------------------------------------------------------
	// App creation scenarios
	// -----------------------------------------------------------
	{
		"Create a new app, response 201 and the app detail",
		[]Request{
			{"POST", "/apps", app1v1},
		},
		201,
		`{
			"Title":"App1",
			"Version":"0.0.1",
			"Maintainers":
			[
				{"Name":"firstmaintainer app1","Email":"firstmaintainer@hotmail.com"},
				{"Name":"secondmaintainer app1","Email":"secondmaintainer@gmail.com"}
			],
			"Company":"Random Inc.",
			"Website":"https://website.com",
			"Source":"https://github.com/random/repo",
			"License":"Apache-2.0",
			"Description":"### Interesting Title\nSome application content, and description\n"
		}`,
	},
	{
		"Create app with existing name but with different version, response 201 and the app detail",
		[]Request{
			{"POST", "/apps", app1v1},
			{"POST", "/apps", app1v2},
		},
		201,
		`{
			"Title":"App1",
			"Version":"0.0.2",
			"Maintainers":
			[
				{"Name":"firstmaintainer app1","Email":"firstmaintainer@hotmail.com"},
				{"Name":"secondmaintainer app1","Email":"secondmaintainer@gmail.com"}
			],
			"Company":"Random Inc.",
			"Website":"https://website.com",
			"Source":"https://github.com/random/repo",
			"License":"Apache-2.0",
			"Description":"### Interesting Title\nSome application content, and description\n"
		}`,
	},
	{
		"Create app with existing name and version, response 409 Conflict",
		[]Request{
			{"POST", "/apps", app1v1},
			{"POST", "/apps", app1v1},
		},
		409,
		`{"error":"App 'App1' with version '0.0.1' already exists."}`,
	},
	{
		"Create app without title, response 400",
		[]Request{
			{"POST", "/apps", appWithoutTitle},
		},
		400,
		`{"error":"Key: 'Meta.Title' Error:Field validation for 'Title' failed on the 'required' tag"}`,
	},
	{
		"Create app without version, response 400",
		[]Request{
			{"POST", "/apps", appWithoutVersion},
		},
		400,
		`{"error":"Key: 'Meta.Version' Error:Field validation for 'Version' failed on the 'required' tag"}`,
	},
	{
		"Create app with bad maintainer email, response 400",
		[]Request{
			{"POST", "/apps", appWithBadMaintainerEmail},
		},
		400,
		`{"error":"Key: 'Meta.Maintainers[0].Email' Error:Field validation for 'Email' failed on the 'email' tag"}`,
	},
	{
		"Create app with multiple bad fields, response 400",
		[]Request{
			{"POST", "/apps", appWithMultipleBadFields},
		},
		400,
		`{"error":"Key: 'Meta.Maintainers' Error:Field validation for 'Maintainers' failed on the 'required' tag\nKey: 'Meta.Website' Error:Field validation for 'Website' failed on the 'url' tag"}`,
	},

	// -----------------------------------------------------------
	// Retrieve one app with title and (optionally) version
	// -----------------------------------------------------------
	{
		"Show the latest version of specific app",
		[]Request{
			{"POST", "/apps", app1v1},
			{"POST", "/apps", app1v2},
			{"POST", "/apps", app2v1},
			{"GET", "/apps/App1", ""},
		},
		200,
		`{
			"Title":"App1",
			"Version":"0.0.2",
			"Maintainers":
			[
				{"Name":"firstmaintainer app1","Email":"firstmaintainer@hotmail.com"},
				{"Name":"secondmaintainer app1","Email":"secondmaintainer@gmail.com"}
			],
			"Company":"Random Inc.",
			"Website":"https://website.com",
			"Source":"https://github.com/random/repo",
			"License":"Apache-2.0",
			"Description":"### Interesting Title\nSome application content, and description\n"
		}`,
	},
	{
		"Show specific app with specific version",
		[]Request{
			{"POST", "/apps", app1v1},
			{"POST", "/apps", app1v2},
			{"POST", "/apps", app2v1},
			{"GET", "/apps/App1/versions/0.0.1", ""},
		},
		200,
		`{
			"Title":"App1",
			"Version":"0.0.1",
			"Maintainers":
			[
				{"Name":"firstmaintainer app1","Email":"firstmaintainer@hotmail.com"},
				{"Name":"secondmaintainer app1","Email":"secondmaintainer@gmail.com"}
			],
			"Company":"Random Inc.",
			"Website":"https://website.com",
			"Source":"https://github.com/random/repo",
			"License":"Apache-2.0",
			"Description":"### Interesting Title\nSome application content, and description\n"
		}`,
	},
	{
		"Show an app with spaces in title, use urlencode",
		[]Request{
			{"POST", "/apps", app3WithSpaceInTitle},
			{"GET", "/apps/App3%20with%20space%20in%20title", ""},
		},
		200,
		`{
			"Title":"App3 with space in title",
			"Version":"0.0.1",
			"Maintainers":
			[
				{"Name":"firstmaintainer app3","Email":"firstmaintainer@hotmail.com"},
				{"Name":"secondmaintainer app3","Email":"secondmaintainer@gmail.com"}
			],
			"Company":"Random Inc.",
			"Website":"https://website.com",
			"Source":"https://github.com/random/repo",
			"License":"Apache-2.0",
			"Description":"### Interesting Title\nSome application content, and description\n"
		}`,
	},
	{
		"Show a non-exist app, response 404",
		[]Request{
			{"POST", "/apps", app1v1},
			{"POST", "/apps", app1v2},
			{"POST", "/apps", app2v1},
			{"GET", "/apps/App6", ""},
		},
		404,
		`{"error":"App with title 'App6' does not exist."}`,
	},
	{
		"Show a non-exist version of an existing app, response 404",
		[]Request{
			{"POST", "/apps", app1v1},
			{"POST", "/apps", app1v2},
			{"POST", "/apps", app2v1},
			{"GET", "/apps/App1/versions/0.0.3", ""},
		},
		404,
		`{"error":"App with title 'App1' and version '0.0.3' does not exist."}`,
	},

	// -----------------------------------------------------------
	// Listing & Filtering
	// Supports both precise match and non-precise match
	// -----------------------------------------------------------
	{
		"List apps, return empty list when no apps stored",
		[]Request{
			{"GET", "/apps", ""},
		},
		200,
		`[]`,
	},
	{
		"List apps, return the details of apps when there were apps stored",
		[]Request{
			{"POST", "/apps", app1v1},
			{"POST", "/apps", app1v2},
			{"POST", "/apps", app2v1},
			{"GET", "/apps", ""},
		},
		200,
		`[
			{"Title":"App1","Version":"0.0.1","Maintainers":[{"Name":"firstmaintainer app1","Email":"firstmaintainer@hotmail.com"},{"Name":"secondmaintainer app1","Email":"secondmaintainer@gmail.com"}],"Company":"Random Inc.","Website":"https://website.com","Source":"https://github.com/random/repo","License":"Apache-2.0","Description":"### Interesting Title\nSome application content, and description\n"},
			{"Title":"App1","Version":"0.0.2","Maintainers":[{"Name":"firstmaintainer app1","Email":"firstmaintainer@hotmail.com"},{"Name":"secondmaintainer app1","Email":"secondmaintainer@gmail.com"}],"Company":"Random Inc.","Website":"https://website.com","Source":"https://github.com/random/repo","License":"Apache-2.0","Description":"### Interesting Title\nSome application content, and description\n"},
			{"Title":"App2","Version":"0.0.1","Maintainers":[{"Name":"firstmaintainer app2","Email":"firstmaintainer@hotmail.com"},{"Name":"secondmaintainer app2","Email":"secondmaintainer@gmail.com"}],"Company":"Random Inc.","Website":"https://website.com","Source":"https://github.com/random/repo","License":"Apache-2.0","Description":"### Interesting Title\nSome application content, and description\n"}
		]`,
	},
	{
		"List apps, filter with app title (multiple results)",
		[]Request{
			{"POST", "/apps", app1v1},
			{"POST", "/apps", app1v2},
			{"POST", "/apps", app2v1},
			{"GET", "/apps?title=App1", ""},
		},
		200,
		`[
			{"Title":"App1","Version":"0.0.1","Maintainers":[{"Name":"firstmaintainer app1","Email":"firstmaintainer@hotmail.com"},{"Name":"secondmaintainer app1","Email":"secondmaintainer@gmail.com"}],"Company":"Random Inc.","Website":"https://website.com","Source":"https://github.com/random/repo","License":"Apache-2.0","Description":"### Interesting Title\nSome application content, and description\n"},
			{"Title":"App1","Version":"0.0.2","Maintainers":[{"Name":"firstmaintainer app1","Email":"firstmaintainer@hotmail.com"},{"Name":"secondmaintainer app1","Email":"secondmaintainer@gmail.com"}],"Company":"Random Inc.","Website":"https://website.com","Source":"https://github.com/random/repo","License":"Apache-2.0","Description":"### Interesting Title\nSome application content, and description\n"}
		]`,
	},
	{
		"List apps, filter with app title (single result)",
		[]Request{
			{"POST", "/apps", app1v1},
			{"POST", "/apps", app1v2},
			{"POST", "/apps", app2v1},
			{"GET", "/apps?title=App2", ""},
		},
		200,
		`[
			{"Title":"App2","Version":"0.0.1","Maintainers":[{"Name":"firstmaintainer app2","Email":"firstmaintainer@hotmail.com"},{"Name":"secondmaintainer app2","Email":"secondmaintainer@gmail.com"}],"Company":"Random Inc.","Website":"https://website.com","Source":"https://github.com/random/repo","License":"Apache-2.0","Description":"### Interesting Title\nSome application content, and description\n"}
		]`,
	},
	{
		"List apps, filter with app title (no result)",
		[]Request{
			{"POST", "/apps", app1v1},
			{"POST", "/apps", app1v2},
			{"POST", "/apps", app2v1},
			{"GET", "/apps?title=App6", ""},
		},
		200,
		`[]`,
	},
	{
		"List apps, non-precise match",
		[]Request{
			{"POST", "/apps", app1v1},
			{"POST", "/apps", app1v2},
			{"POST", "/apps", app2v1},
			{"GET", "/apps?title[like]=App", ""},
		},
		200,
		`[
			{"Title":"App1","Version":"0.0.1","Maintainers":[{"Name":"firstmaintainer app1","Email":"firstmaintainer@hotmail.com"},{"Name":"secondmaintainer app1","Email":"secondmaintainer@gmail.com"}],"Company":"Random Inc.","Website":"https://website.com","Source":"https://github.com/random/repo","License":"Apache-2.0","Description":"### Interesting Title\nSome application content, and description\n"},
			{"Title":"App1","Version":"0.0.2","Maintainers":[{"Name":"firstmaintainer app1","Email":"firstmaintainer@hotmail.com"},{"Name":"secondmaintainer app1","Email":"secondmaintainer@gmail.com"}],"Company":"Random Inc.","Website":"https://website.com","Source":"https://github.com/random/repo","License":"Apache-2.0","Description":"### Interesting Title\nSome application content, and description\n"},
			{"Title":"App2","Version":"0.0.1","Maintainers":[{"Name":"firstmaintainer app2","Email":"firstmaintainer@hotmail.com"},{"Name":"secondmaintainer app2","Email":"secondmaintainer@gmail.com"}],"Company":"Random Inc.","Website":"https://website.com","Source":"https://github.com/random/repo","License":"Apache-2.0","Description":"### Interesting Title\nSome application content, and description\n"}
		]`,
	},
	{
		"List apps, filter with multiple fields",
		[]Request{
			{"POST", "/apps", app1v1},
			{"POST", "/apps", app1v2},
			{"POST", "/apps", app2v1},
			{"GET", "/apps?title=App1&version=0.0.1", ""},
		},
		200,
		`[
			{"Title":"App1","Version":"0.0.1","Maintainers":[{"Name":"firstmaintainer app1","Email":"firstmaintainer@hotmail.com"},{"Name":"secondmaintainer app1","Email":"secondmaintainer@gmail.com"}],"Company":"Random Inc.","Website":"https://website.com","Source":"https://github.com/random/repo","License":"Apache-2.0","Description":"### Interesting Title\nSome application content, and description\n"}
		]`,
	},
	{
		"List apps, filter with bad field names",
		[]Request{
			{"POST", "/apps", app1v1},
			{"POST", "/apps", app1v2},
			{"POST", "/apps", app2v1},
			{"GET", "/apps?title=App1&dummy=unknown", ""},
		},
		400,
		`{"error":"Error occurred during searching app: Unsupported rule for field 'dummy'"}[]`,
	},
	{
		"List apps, filter with bad operator",
		[]Request{
			{"POST", "/apps", app1v1},
			{"POST", "/apps", app1v2},
			{"POST", "/apps", app2v1},
			{"GET", "/apps?title=App1&version[dummy]=0.0.1", ""},
		},
		400,
		`{"error":"Unrecognized operator type 'dummy'"}`,
	},
}

func TestScenarios(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	formatUrl := func(path string) string {
		return fmt.Sprintf("%s/v1%s", ts.URL, path)
	}

	performRequest := func(req Request) (statusCode int, responseBody string) {
		var resp *http.Response
		var err error
		if req.method == "POST" {
			resp, err = http.Post(formatUrl(req.url), "application/json", strings.NewReader(req.data))
		} else {
			resp, err = http.Get(formatUrl(req.url))
		}

		if err != nil {
			t.Fatalf("Error occurred during %s %s, detail: %e", req.method, req.url, err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error while reading response body from %s %s, detail: %e", req.method, req.url, err)
		}

		return resp.StatusCode, string(body)
	}

	for _, tt := range scenarios {
		// Clean the store to make the data independentable between test cases.
		setupStore()

		t.Run(tt.scenarioTitle, func(t *testing.T) {
			var statusCode int
			var responseBody string
			// Perform all the requests in the list.
			// The last request's response is to be checked.
			// Other requests are considered as setup steps and not supposed to fail.
			for _, request := range tt.requests {
				statusCode, responseBody = performRequest(request)
			}
			if statusCode != tt.expectedStatusCode {
				t.Errorf("Expected status code to be %d but got %d", tt.expectedStatusCode, statusCode)
			}

			if responseBody != trimAndMergeToOneLine(tt.expectedResponseBody) {
				t.Errorf("Expected response body to be '%s' but got '%s'", tt.expectedResponseBody, responseBody)
			}
		})
	}
}

func trimAndMergeToOneLine(text string) string {
	lines := strings.Split(text, "\n")
	for i := 0; i < len(lines); i++ {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return strings.Join(lines, "")
}
