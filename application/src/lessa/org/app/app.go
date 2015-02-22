package app

import (
   "errors"
   "fmt"
)

type Runnable interface {
   Run() error
}

type Application interface {
   Runnable
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
