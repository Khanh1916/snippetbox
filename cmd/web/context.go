package main

type contextKey string //custom key type

// custom variable - unique key to store and retrieve authen status without naming collisions
const isAuthenticatedContextKey = contextKey("isAuthenticated")
