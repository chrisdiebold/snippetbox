# Environment Variables

If you’ve built and deployed web applications before, then you’re probably thinking what about environment variables?

If you want, you can store your configuration settings in environment variables and access them directly from your application by using the
`os.Getenv()` function like so:

```go
addr := os.Getenv("SNIPPETBOX_ADDR")
```
But this has some drawbacks compared to using command-line flags. You can’t specify a default setting (the return value from `os.Getenv()` is the empty string if the environment variable doesn’t exist), you don’t get the `-help` functionality that you do with command-line flags, and the return value from `os.Getenv()` is always a string — you don’t get automatic type conversions like you do with `flag.Int()`, `flag.Bool()` and the other command-line flag functions.

Instead, you can get the best of both worlds by passing the environment variable as a command-line flag when starting the application. For example:

```bash
export SNIPPETBOX_ADDR=":9999"
go run ./cmd/web -addr=$SNIPPETBOX_ADDR
```

## Boolean flags
For flags defined with `flag.Bool()`, omitting a value when starting the application is the same as writing `-flag=true`. The following two commands are equivalent:

```bash
go run ./example -flag=true
go run ./example -flag
```
You must explicitly use `-flag=false` if you want to set a boolean flag value to false.

## Pre-existing variables

It’s possible to parse command-line flag values into the memory addresses of pre-existing variables, using `flag.StringVar()`, `flag.IntVar()`, `flag.BoolVar()`, and similar functions for other types.

These functions are particularly useful if you want to store all your configuration settings in a single struct. As a rough example:

  ```bash
  type config struct {
    addr      string
    staticDir string
}

// ...

var cfg config

flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")

flag.Parse()
```