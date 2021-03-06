package main

import (
   "fmt"
   "lessa/org/app"
   impl "lessa/org/impl/app"
   "time"
)

func main() {
   fmt.Println("[main] Starting main.")

   fmt.Println("[main] Installing builder func...")
   impl.InstallAppBuilderFunc()

   fmt.Println("[main] Retrieving builder constructor func...")
   defBuilderFunc := app.GetAppBuilderFunc()

   fmt.Println("[main] The default builder: ", defBuilderFunc())

   fmt.Println("[main] Installing a customized builder constructor func...")
   app.SetAppBuilderFunc(func() app.Builder {
      return defBuilderFunc().
                WithJobs(func(ctx app.Context) error { 
                   fmt.Println("") 
                   fmt.Println("------------") 
                   fmt.Println("Hello World!") 
                   fmt.Println("------------") 
                   fmt.Println("") 
                   return nil
                }, func(ctx app.Context) error { 
                   fmt.Println("") 
                   fmt.Println("--------------------") 
                   fmt.Println("Hello Waiting World!") 
                   fmt.Println("--------------------") 
                   fmt.Println("") 
                   for {
                      time.Sleep(10*time.Second)
                   }
                   return nil
                }).
                WithOption1("custom").
                WithOption2(333).
                WithOption3(333.666).
                WithOption4([]string{"custom1", "custom2"}).
                WithOption5([]int{333, 666}).
                WithOption6([]float64{333.666, 666.333})
   })

   fmt.Println("[main] Retrieving builder constructor func...")
   newBuilderFunc := app.GetAppBuilderFunc()

   fmt.Println("[main] The customized builder: ", newBuilderFunc())

   fmt.Println("[main] Retrieving the application object.")
   a := app.GetApp()

   fmt.Println("[main] Running the application.")
   a.Run()

   fmt.Println("")
   fmt.Println("[main] Exiting.")
}
