# Servemux Features

Request URL paths are automatically sanitized. If the request path contains any `.` or `..` elements or repeated slashes, the user will automatically be redirected to an equivalent clean URL. For example, if a user makes a request to `/foo//bar/../baz/./` they will automatically be sent a `301` Permanent Redirect to `/foo/baz/` instead.

If a subtree path has been registered and a request is received for that subtree path without a trailing slash, then the user will automatically be sent a `301` Permanent Redirect to the subtree path with the slash added. For example, if you
have registered the subtree path `/foo/`, then any request to `/foo` will be redirected to `/foo/`.

## Hostname matching
It’s possible to include host names in your route patterns. This can be useful when you want to redirect all HTTP requests to a canonical URL, or if your application is acting as the back end for multiple sites or services. For example:

```go
mux := http.NewServeMux()
mux.HandleFunc("foo.example.org/", fooHandler)
mux.HandleFunc("bar.example.org/", barHandler)
mux.HandleFunc("/baz", bazHandler)
```

When it comes to pattern matching, any host-specific patterns will be checked first and if there is a match the request will be dispatched to the corresponding handler. Only when there isn’t a host-specific match found will the non-host-specific patterns also be checked.