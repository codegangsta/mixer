// Package Mixer is a small but powerful package for writing modular HTTP
// handlers in Go. Mixer attempts to solve some of the problems that
// https://github.com/go-martini/martini tried to solve, but with a smaller,
// more idiomatic scope.
//
// If you liked the idea of dependency injection in Martini, but you think it
// contained too much magic, then Mixer is a great fit.
//
// For a full guide visit http://github.com/codegangsta/mixer
//
//    package main
//
//    import (
//        "fmt"
//        "log"
//        "os"
//
//        "github.com/codegangsta/mixer"
//    )
//
//    func AppContext struct {
//        mixer.Context
//
//        Logger *log.Logger
//    }
//
//    func main() {
//        // We've changed mixer.Classic to mixer.New, and passed
//        // a context function into it, so that the compiler knows
//        // what our context is, and what data to start it with
//        m := mixer.New(func(c mixer.Context) *AppContext {
//            return &AppContext{
//                Context: c,
//                Logger: log.New(os.Stdout, "[mixer]", 0)
//            }
//        })
//
//        hello := func(c *AppContext) {
//            fmt.Fprint(c.ResponseWriter(), "Hello world")
//        }
//
//        http.ListenAndServe(":3000", m.Handler(hello))
//    }
package mixer
