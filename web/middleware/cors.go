package middleware

// func (s *WithCORS) ServeHTTP(res http.ResponseWriter, req *http.Request) {
// 	if origin := req.Header.Get("Origin"); origin != "" {
// 	  res.Header().Set("Access-Control-Allow-Origin", origin)
// 	  res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	  res.Header().Set("Access-Control-Allow-Headers",
// 		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// 	}

// 	// Stop here for a Preflighted OPTIONS request.
// 	if req.Method == "OPTIONS" {
// 	  return
// 	}
// 	// Lets Gorilla work
// 	s.r.ServeHTTP(res, req)
//   }
