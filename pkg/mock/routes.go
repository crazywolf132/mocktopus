package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
)

type RouteConfig struct {
	Routes []Route `toml:"routes"`
}

type Route struct {
	Path       string          `toml:"path"`
	Method     string          `toml:"method"`
	Conditions []Condition     `toml:"conditions"`
	Request    *RequestConfig  `toml:"request"`
	Response   *ResponseConfig `toml:"response"`
	Priority   int             `toml:"priority"`
	Callback   string          `toml:"callback"`
}

type Condition struct {
	Matcher RequestMatcher         `toml:"matcher"`
	State   map[string]interface{} `toml:"state"`
}

type RequestConfig struct {
	Headers map[string]string `toml:"headers"`
	Body    string            `toml:"body"`
}

type ResponseConfig struct {
	Status   int               `toml:"status"`
	Headers  map[string]string `toml:"headers"`
	Body     string            `toml:"body"`
	Template string            `toml:"template"`
}

// This is a TOML copy of the RequestMatcher from the server.
type RequestMatcher struct {
	Path        string            `toml:"path"`
	Method      string            `toml:"method"`
	Queries     map[string]string `toml:"queries"`
	Headers     map[string]string `toml:"headers"`
	Body        string            `toml:"body"`
	BodyPattern string            `toml:"bodyPattern"`
}

func (m *Mock) LoadConfig() error {
	if ok, _ := m.exists(); !ok {
		return fmt.Errorf("it appears the mock has been removed")
	}

	mockConfig := viper.New()
	mockConfig.SetConfigFile(filepath.Join(m.location, "config.toml"))
	mockConfig.ReadInConfig()
	m.mockConfig = mockConfig

	return nil
}

func (m *Mock) Discovery() error {
	// We are going to read the config in relation to
	// if they are using file routing... or if they are using config based routing.Discovery

	var usingJavascript bool = m.mockConfig.GetBool("useJavascript")

	startSearchLocation := m.mockConfig.GetString("routes")
	startSearchLocation = filepath.Join(m.location, startSearchLocation)

	fmt.Println("startSearchLocation: ", startSearchLocation)

	var ext string = ".toml"
	if usingJavascript {
		ext = ".js"
	}

	routes := []RouteConfig{}
	err := filepath.Walk(startSearchLocation, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ext) {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			var routesConfig RouteConfig
			if err := toml.NewDecoder(file).Decode(&routesConfig); err != nil {
				return err
			}
			routes = append(routes, routesConfig)

		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to load routes from files: %v", err)
	}

	m.Routes = routes

	return nil
}

func (m *Mock) getRoutePrefix(location, filePath string) string {
	relPath, _ := filepath.Rel(location, filePath)
	ext := filepath.Ext(relPath)
	routePrefix := "/" + strings.TrimSuffix(relPath, ext)
	routePrefix = strings.ReplaceAll(routePrefix, "\\", "/")
	return routePrefix
}

func (m *Mock) MakeRoutes() error {
	for _, routeConfig := range m.Routes {
		for _, route := range routeConfig.Routes {
			path := fmt.Sprintf("%s %s", route.Method, route.Path)
			m.Server.Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				// We will check to see if the request meets the required conditions...
				// if (m.machRequest(r, route.Conditions)) {

				// }
				if route.Response.Status != 0 {
					w.WriteHeader(route.Response.Status)
				}
				if route.Response.Headers != nil {
					for key, value := range route.Response.Headers {
						w.Header().Set(key, value)
					}
				}
				if route.Response.Body != "" {
					w.Write([]byte(route.Response.Body))
				}
			})
		}
	}
	return nil
}

func (m *Mock) matchRequest(r *http.Request, matcher RequestMatcher) bool {
	if len(matcher.Queries) > 0 {
		for k, v := range matcher.Queries {
			if r.URL.Query().Get(k) != v {
				return false
			}
		}
	}

	if len(matcher.Headers) > 0 {
		for k, v := range matcher.Headers {
			if r.Header.Get(k) != v {
				return false
			}
		}
	}

	if matcher.Body != "" {
		var body interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return false
		}
		if !reflect.DeepEqual(body, matcher.Body) {
			return false
		}
	}

	return true
}

// func (m *Mock) MakeRoutes() ([]server.Route, error) {

// 	routes := []server.Route{}

// 	for _, routeConfig := range m.Routes {
// 		for _, route := range routeConfig.Rotues {
// 			newRoute := server.Route{}
// 			for _, condition := range route.Conditions {
// 				newRoute.Matchers = append(newRoute.Matchers, server.RequestMatcher{
// 					Path:    condition.Matcher.Path,
// 					Method:  condition.Matcher.Method,
// 					Headers: condition.Matcher.Headers,
// 					Body:    condition.Matcher.Body,
// 					Queries: condition.Matcher.Queries,
// 				})
// 			}

// 			newRoute.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 				response := route.Response

// 				if response.Status != 0 {
// 					w.WriteHeader(response.Status)
// 				}

// 				if response.Headers != nil {
// 					for key, value := range response.Headers {
// 						w.Header().Set(key, value)
// 					}
// 				}

// 				if response.Body != "" {
// 					w.Write([]byte(response.Body))
// 				}
// 			})

// 			routes = append(routes, newRoute)
// 		}
// 	}

// 	return routes, nil
// }

func (m *Mock) templateResponse(route *Route, req *http.Request) (string, error) {

	pathParams := make(map[string]string)
	match := regexp.MustCompile(route.Path).FindStringSubmatch(req.URL.Path)
	if len(match) > 1 {
		subexpNames := regexp.MustCompile(route.Path).SubexpNames()
		for i, name := range subexpNames {
			if i != 0 && name != "" {
				pathParams[name] = match[i]
			}
		}
	}

	templateData := struct {
		Request     *http.Request
		PathParams  map[string]string
		QueryParams url.Values
		Headers     http.Header
		Body        string
	}{
		Request:     req,
		PathParams:  pathParams,
		QueryParams: req.URL.Query(),
		Headers:     req.Header,
	}
	fmt.Println("route.Request.Body: ", route.Request.Body)
	tmpl, err := template.New("response").Parse(route.Request.Body)
	if err != nil {
		return "", err
	}
	// Create a new writer
	var responseWriter strings.Builder
	// Execute the template and write the result to the response writer
	err = tmpl.Execute(&responseWriter, templateData)
	if err != nil {
		return "", err
	}
	// Get the final response body text from the response writer
	responseBody := responseWriter.String()
	return responseBody, nil
}
