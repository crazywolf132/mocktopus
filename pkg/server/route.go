package server

// type Route struct {
// 	Matchers   []RequestMatcher
// 	Handler    http.Handler
// 	Middleware []Middleware
// }

// type RequestMatcher struct {
// 	Path      string
// 	Method    string
// 	Queries   map[string]string
// 	Headers   map[string]string
// 	Body      interface{}
// 	Predicate func(r *http.Request) bool
// }

// func (s *Server) AddRoute(route Route) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	handler := route.Handler
// 	for _, mw := range route.Middleware {
// 		handler = mw(handler)
// 	}

// 	s.router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
// 		for _, matcher := range route.Matchers {
// 			if s.matchRequest(r, matcher) {
// 				handler.ServeHTTP(w, r)
// 				return
// 			}
// 		}
// 		http.NotFound(w, r)
// 	})
// }

// func (s *Server) matchRequest(r *http.Request, matcher RequestMatcher) bool {
// 	if matcher.Path != "" && !regexp.MustCompile(matcher.Path).MatchString(r.URL.Path) {
// 		return false
// 	}
// 	if matcher.Method != "" && matcher.Method != r.Method {
// 		return false
// 	}

// 	if len(matcher.Queries) > 0 {
// 		for k, v := range matcher.Queries {
// 			if r.URL.Query().Get(k) != v {
// 				return false
// 			}
// 		}
// 	}

// 	if len(matcher.Headers) > 0 {
// 		for k, v := range matcher.Headers {
// 			if r.Header.Get(k) != v {
// 				return false
// 			}
// 		}
// 	}

// 	if matcher.Body != nil {
// 		var body interface{}
// 		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
// 			return false
// 		}
// 		if !reflect.DeepEqual(body, matcher.Body) {
// 			return false
// 		}
// 	}
// 	if matcher.Predicate != nil && !matcher.Predicate(r) {
// 		return false
// 	}
// 	return true
// }
