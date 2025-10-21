package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// adding middleware handling security header according to OWASP guidance
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-Xss-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

// middleware loggging request
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

// middleware panic recovery
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close") // trigger: http server auto close current connection after send it
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// middleware to require authentication for certain routes
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//	check if user is authenticated
		if !app.isAuthenticated(r) {
			//add the path that user trying to access to their session data to remember
			app.sessionManager.Put(r.Context(), "redirectPathAfterLogin", r.URL.Path)
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		//	add "no-store" cache-control header to prevent caching of sensitive info
		w.Header().Set("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// Create a NoSurf middleware function which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly attributes set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

// middleware to authenticate user based on session value
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if no authenticatedUserID in session, call next handler
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}
		// check if user exists in db
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		// if user exists, set context value to true
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
