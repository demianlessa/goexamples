package app

import (
   "errors"
   "fmt"
)

type Context interface {
}

// an object that receives a context and runs a single task
// Context should expose all necessary services so that the
// function properl yperform setup, teardown, error handling,
// synchronization, etc
type Job func(context Context) error

// an object capable of running jobs and waiting for their completion
type Runnable interface {
   Run() error
}

// an object capable of running and stopping
type Stoppable interface {
   Runnable
   Stop() error
}

// an object capable of running, stopping, and restarting
type Graceful interface {
   Stoppable
   Restart() error
}

// an application object can choose the runnable interface(s) it support
type Application interface {
   Stoppable
}

type Builder interface {
   Build() Application
   WithJobs(jobs... Job) Builder 
   WithOption1(val string) Builder 
   WithOption2(val int) Builder 
   WithOption3(val float64) Builder 
   WithOption4(val []string) Builder 
   WithOption5(val []int) Builder 
   WithOption6(val []float64) Builder 
}

type AppBuilderFunc func() Builder

var (
   app        Application
   appBuilder AppBuilderFunc
)

func GetApp() Application {
   if app == nil {
      if appBuilder == nil {
         panic(errors.New(fmt.Sprint("No application instance found and no application builder function defined.")))
      }
      app = appBuilder().Build()
   }
   return app
}

func GetAppBuilderFunc() AppBuilderFunc {
   return appBuilder
}

func SetAppBuilderFunc(builderFunc AppBuilderFunc) {
   appBuilder = builderFunc
}