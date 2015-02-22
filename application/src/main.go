package main

import (
   "fmt"
   "lessa/org/app"
   impl "lessa/org/impl/app"
)

func main() {
   fmt.Println("Starting application.")

   fmt.Println("Installing builder func...")
   impl.InstallAppBuilderFunc()

   fmt.Println("Retrieving builder func...")
   defBuilderFunc := app.GetAppBuilderFunc()

   fmt.Println("Here's what the default builder looks like", defBuilderFunc())

   fmt.Println("Installing a customized builder func...")
   app.SetAppBuilderFunc(func() app.Builder {
      return defBuilderFunc().
                WithOption1("custom").
                WithOption2(333).
                WithOption3(333.666).
                WithOption4([]string{"custom1", "custom2"}).
                WithOption5([]int{333, 666}).
                WithOption6([]float64{333.666, 666.333})
   })

   fmt.Println("Retrieving builder func...")
   newBuilderFunc := app.GetAppBuilderFunc()

   fmt.Println("Here's what the customized builder looks like", newBuilderFunc())

   fmt.Println("Creating application...")
   a := app.GetApp()

   fmt.Println("Running the application...")
   a.Run()
}
