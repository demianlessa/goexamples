# Application Proof-of-Concept

This simple project illustrates the abstraction of applications into an
Application interface with flexible startup options, following a strong
encapsulation mechanism where neither implementations nor global variables
are exposed to the client code.

   1. Definition of an Application interface type embedding a Runnable interface.

```Go
           type Runnable interface {
              Run() error
           }

           type Application interface {
              Runnable
           }
```

   2. Definition of a Builder interface type to allow specification of application
      options, without exposing implementation details. For this example, we use a
      builder with six options.

```Go
           type Builder interface {
              Build() Application
              WithOption1(val string) Builder 
              WithOption2(val int) Builder 
              WithOption3(val float64) Builder 
              WithOption4(val []string) Builder 
              WithOption5(val []int) Builder 
              WithOption6(val []float64) Builder 
           }
```

      This provides a simple and fluent mechanism for setting application options.

   3. It is obvious that the Builder interface enables encapsulation of Application 
      object implementations. In order to encapsualte Builder object implementations 
      themselves, a functional approach is used. The general idea is simple-- define
      a function type that returns a Builder instance (i.e., a constructor function
      for Builder objects) and let implementors of Builder and Application objects
      "install" their Builder constructor functions.

```Go
           type AppBuilderFunc func() Builder

           var (
              app        Application
              appBuilder AppBuilderFunc
           )

           func GetAppBuilderFunc() AppBuilderFunc {
              return appBuilder
           }

           func SetAppBuilderFunc(builderFunc AppBuilderFunc) {
              appBuilder = builderFunc
           }
```

      Again, it is obvious that by calling the SetAppBuilderFunc implementors are
      able to specify how an application is built without disclosing implementation
      details.

   4. Global application instances are obtained from a central point, using 
      the GetApp function.

```Go
           func GetApp() Application {
              if app == nil {
                 if appBuilder == nil {
                    panic(errors.New(fmt.Sprint("No application instance found and no application builder function defined.")))
                 }
                 app = appBuilder().Build()
              }
              return app
           }
```

   5. We assume that the implementor of Builder and Application objects provide 
      sensible default building options. However, there are many occasions in which
      we need to change those, as during testing, when overriding startup options
      using runtime environments or command line arguments, etc. The patterns used
      above support this in a relatively simple manner, as illustrated below.

```Go
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
```

      The customization pattern is straightforward: get the installed Builder constructor
      function, use it to create a customized builder constructor function and install it
      globally. 

      Note that we install the implementation Builder explicitly. This would be typically
      done in the init() function of the Builder implementation package, allowing us to 
      remove all references to the implementation package in main.go. Unfortunately, this 
      would cause the implementation package not to be compiled/linked into our app! Thus,
      we must either install it explicitly (like in the above example) or use the Go idiom
      "include-package-for-side-effects". This idiom is typically used when a registration 
      mechanism is used, such as our case of registering a Builder constructor function.

Hope you enjoy the code. Any comments/criticisms are welcome. :)