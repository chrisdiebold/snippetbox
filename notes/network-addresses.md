# Network Addresses

The TCP network address that you pass to `http.ListenAndServe()` should be in the format `"host:port"`. If you omit the host (like we did with `":4000"`) then the server will listen on all available network interfaces on your computer. Generally, you only need to specify a host in the address if your computer has multiple network interfaces and you want to listen on
just one of them.

In other Go projects or documentation you might sometimes see network addresses written using named ports like `":http"` or `":http-alt"` instead of a number. If you use a named port then the `http.ListenAndServe()` function will attempt to look up the relevant port number from your `/etc/services` file when starting the server, returning an error if a match can’t be found.