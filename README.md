# Mixer: Classy HTTP Handlers in Go

Mixer is a small but powerful package for writing modular HTTP handlers in Go.
Mixer attempts to solve some of the problems that
[Martini](https://github.com/go-martini/martini) tried to solve, but with a
smaller, more idiomatic scope.


## Features

* Type safe
* Brutally simple to use
* Non intrusive design
* Zero external dependencies
* Compatible with most HTTP routers out there, choose your favorite!
* Dependency injection without reflection, type assertions, or other magic


## Getting Started

Pull down this package with go get (go 1.18 or greater is required):

```sh
go get github.com/codegangsta/mixer
```


### Hello world!

After installing your package you can create a simple http handler in Mixer:

```go
package main

import (
    "fmt"

    "github.com/codegangsta/mixer"
)

func main() {
    m := mixer.Classic()
    
    hello := func(c mixer.Context) {
        fmt.Fprint(c.ResponseWriter(), "Hello world")
    }
    
    http.ListenAndServe(":3000", m.Handler(hello))
}
```

This sure is simple, but it's not really useful. Once your http handlers start
growing, they are going to need to depend on things like databases, loggers,
sessions, users and all kinds of other functionality that `mixer.Context`
doesn't have. This is where the **Custom Context** comes in.


### Adding your own context

Let's create an `AppContext` struct to hold the context for our app. Beyond the basic stuff in `mixer.Context`, it's also got a `log.Logger` instance.

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/codegangsta/mixer"
)

func AppContext struct {
    mixer.Context
    
    Logger *log.Logger
}

func main() {
    // We've changed mixer.Classic to mixer.New, and passed
    // a context function into it, so that the compiler knows
    // what our context is, and what data to start it with
    m := mixer.New(func(c mixer.Context) *AppContext {
        return &AppContext{
            Context: c,
            Logger: log.New(os.Stdout, "[mixer]", 0)
        }
    })
    
    hello := func(c *AppContext) {
        fmt.Fprint(c.ResponseWriter(), "Hello world")
    }
    
    http.ListenAndServe(":3000", m.Handler(hello))
}
```

Pretty cool right? This means we can easily add global and request scoped data
to our context struct, and the handlers will get them in a fully type-safe way.
No `interface{}`, no reflection and no type assertions.


### Adding Before and After Hooks

Mixer also supports `Before` and `After` hooks, in case you'd like to run some
logic or map some dependencies outside of the context func.

If we add a `User` type to our `AppContext`:

```go
func AppContext struct {
    mixer.Context
    
    Logger *log.Logger
    User *User
}
```

We can add that user in a before func, maybe it comes from a cookie:

```go
m.Before(func(c *AppContext) {
    cookie, err := c.Request().Cookie("user-id")
    if err != nil {
        // Return early by writing out a status
        c.ResponseWriter().WriteHeader(500)
        return
    }
    
    c.User = lookupUser(cookie.Value)
})
```

## FAQ

### Is this a framework

No. This is not a web framework. I'd hardly even call it a package. Mixer is
more of a pattern for getting global and request scoped values to your handlers without much fuss.

Mixer is designed to plug into your existing web stack, as long as it supports
the `http.Handler` interface for its handlers. (Which is should)


### Where's my router?

Like [Negroni](https://github.com/codegangsta/negroni), Mixer is a BYOR (Bring
your own Router) library. Pick one that is your favorite, and append Mixer
handlers to it using the `mixer.Handler` utility function.


### What router do you recommend?

I'm really enjoying [https://github.com/go-chi/chi](https://github.com/go-chi/chi) these days. Here's a simple example of how Mixer works with Chi:

```go
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
    "github.com/codegangsta/mixer"
)

type Context {
    mixer.Context
    
    // Add your dependencies here
}

func main() {
	r := chi.NewRouter()
    m := mixer.New(func(c mixer.Context) *Context {
        return &Context{c}
    })
    
	r.Get("/", mixer.Handler(func(c *Context) {
        fmt.Fprint(c.ResponseWriter(), "Hello world")
	}))
	http.ListenAndServe(":3000", r)
}

```


### How does this work without reflection?

Mixer uses a very straightforward generics implementation to accomplish it's
custom context feature. (Seriously, read the code, it's so freaking simple)


### What's the difference between this and `context.Context`

`context.Context` has its uses, but since it was introduced I've seen it be
abused for many things. I believe that type assertions should be used
sparingly, and that `context.Context` isn't designed to hold both global and
request scoped values in a very efficient way.


### How fast is this?

I don't know. Someone putting together a benchmark would be cool. But I imagine
it's not very slow. Speed will depend on the number of allocations you create
when setting up your context.


### I don't want to embed `mixer.Context` into my own context?

You don't have to. If you want to keep your context pure of any third party
nonsense, you can do that too. Heck, you can even make your contact an
Interface if that is your sort of thing.


## Contributing

I don't know how much more we could add to this that would be valuable, but if you have ideas that's great. Submit a Pull Request with a new feature, bug fix, or link to your project that plays nicely with Mixer.

## About

Mixer is distributed under the MIT license, see LICENSE for more details.

Mixer is obsessively designed by none other than the [Code Gangsta](https://github.com/codegangsta)
