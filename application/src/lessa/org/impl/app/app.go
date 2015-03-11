package app

import (
   "encoding/json"
   "fmt"
   "os"
   "os/signal"
   //"syscall"
   "sync"
   model "lessa/org/app"
)

// private options, exposed indirectly via the Builder interface
// options are embedded by builder and application to promote
// implementation reuse
// public members on the type buys us json marshalling support
type options struct {
   Option1 string
   Option2 int
   Option3 float64
   Option4 []string
   Option5 []int
   Option6 []float64
}

// builder has options
type builder struct {
   options
   jobs    []model.Job
}

// application has options
type application struct {
   options
   jobs    []model.Job
   done    chan bool
   sigs    chan os.Signal
}

// allow clients to install a default options builder
func InstallAppBuilderFunc() {
   model.SetAppBuilderFunc(defaultBuilder)
}

// pretty print
func (o options) String() string {
   js, _ := json.MarshalIndent(o, "", "   ")
   return string(js)
}

// dummy
func (a application) Run() error {

   fmt.Println()
   fmt.Println("[impl/application] Registering for specified signal types.")
   signal.Notify(a.sigs)

   fmt.Println("[impl/application] Setting up signal handling.")
   go a.waitForSignal()

   if len(a.jobs) > 0 {
      fmt.Println("[impl/application] Spawning jobs.")
      go a.runJobs()
   }

   fmt.Println("[impl/application] Waiting for jobs or a registered signal.")
   <- a.done

   fmt.Println("[impl/application] Received and processed signal.")
   return nil
}

func (a application) runJobs() {

   wg := sync.WaitGroup{}

   for _, each := range a.jobs {
      wg.Add(1)
      go func(job model.Job) {
         // clean up correctly even if we panic
         defer func() {
            if err := recover(); err != nil {
               wg.Done()
               panic(err)
            }
         }()
         job(a)
         wg.Done()
      }(each)
   }

   wg.Wait()

   a.done <- true
}

func (a application) waitForSignal() {

   select {
      // waiting for a registered signal
      case sig := <-a.sigs:

         fmt.Println()
         fmt.Println("[impl/application] Signal received:", sig)

         // cleaning up allocated resources
         a.cleanup()

         // releasing the runnable
         a.done <- true

         // remove from select
         a.sigs = nil
   }
}

func (a application) cleanup() error {

   fmt.Println("[impl/application] Cleaning up prior to stopping.")
   return nil
}

func (a application) Stop() error {

   fmt.Println("[impl/application] Sending an interrupt signal.")
   a.sigs <- os.Interrupt

   fmt.Println("[impl/application] Stop completed.")
   return nil
}

// default options used by the default builder
func defaultoptions() options {
   return options{
      Option1: "default",
      Option2: 666,
      Option3: 666.666,
      Option4: []string{"default1","default2"},
      Option5: []int{666, 999},
      Option6: []float64{666.999, 999.666},
   }
}

// default builder
func defaultBuilder() model.Builder {
   return builder{
      options: defaultoptions(),
   }
}

// creates an application by passing a copy of the builder's options
func (b builder) Build() model.Application {
   return application {
      options: b.options,
      jobs:    b.jobs,
      done:    make(chan bool, 1),
      sigs:    make(chan os.Signal, 1),
   }
}

// update the internal jobs
func (b builder) WithJobs(jobs... model.Job) model.Builder {
   b.jobs = jobs
   return b
}

// update the internal options
func (b builder) WithOption1(val string) model.Builder {
   b.Option1 = val
   return b
}

// update the internal options
func (b builder) WithOption2(val int) model.Builder {
   b.Option2 = val
   return b
}

// update the internal options
func (b builder) WithOption3(val float64) model.Builder {
   b.Option3 = val
   return b
}

// update the internal options
func (b builder) WithOption4(val []string) model.Builder {
   b.Option4 = val
   return b
}

// update the internal options
func (b builder) WithOption5(val []int) model.Builder {
   b.Option5 = val
   return b
}

// update the internal options
func (b builder) WithOption6(val []float64) model.Builder {
   b.Option6 = val
   return b
}
