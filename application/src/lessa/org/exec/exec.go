package exec

// an object capable of running one or more tasks
type Runnable interface {
   Run() error
}

// an object capable of running and stopping one or more tasks
type Stoppable interface {
   Runnable
   Stop() error
}

// an object capable of running, stopping, and restarting one or more tasks
type Graceful interface {
   Stoppable
   Restart() error
}

// application has options
type Application struct {
   options
   done    chan bool
   sigs    chan os.Signal
}

// allow clients to install a default options builder
func InstallAppBuilderFunc() {
   model.SetAppBuilderFunc(defaultBuilder)
}

func (a Application) Run() error {

   fmt.Println()
   fmt.Println("[impl/application] Registering for specified signal types.")
   signal.Notify(a.sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

   fmt.Println("[impl/application] Setting up signal handling.")
   go a.waitForSignal()

   fmt.Println("[impl/application] Waiting for a registered signal.")
   <- a.done

   fmt.Println("[impl/application] Received and processed signal.")
   return nil
}

func (a Application) waitForSignal() {

      // waiting for a registered signal
      sig := <-a.sigs

      fmt.Println()
      fmt.Println("[impl/application] Signal received:", sig)

      // cleaning up allocated resources
      a.cleanup()

      // releasing the runnable
      a.done <- true
}

func (a Application) cleanup() error {

   fmt.Println("[impl/application] Cleaning up prior to stopping.")
   return nil
}

func (a Application) Stop() error {

   fmt.Println("[impl/application] Sending an interrupt signal.")
   a.sigs <- os.Interrupt

   fmt.Println("[impl/application] Stop completed.")
   return nil
}
