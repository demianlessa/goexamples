package app

import (
   "errors"
   "fmt"
)

// an object capable of running a task (single responsibility)
type Runnable interface {
   Run() error
}

// an object capable of running and stopping a task (single responsibility)
type Stoppable interface {
   Runnable
   Stop() error
}

// an object capable of running, stopping, and restarting a task (single responsibility)
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
