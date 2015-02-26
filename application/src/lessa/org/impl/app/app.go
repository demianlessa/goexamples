package app

import (
   "encoding/json"
   "fmt"
   model "lessa/org/app"
)

// private options, exposed indirectly via the Builder interface
// options are embedded by builder and application to promote
// implementation reuse
type options struct {
   option1 string
   option2 int
   option3 float64
   option4 []string
   option5 []int
   option6 []float64
}

// builder has options
type builder struct {
   options
}

// application has options
type application struct {
   options
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
   fmt.Print(a.String())
   return nil
}

// default options used by the default builder
func defaultoptions() options {
   return options{
      option1: "default",
      option2: 666,
      option3: 666.666,
      option4: []string{"default1","default2"},
      option5: []int{666, 999},
      option6: []float64{666.999, 999.666},
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
   }
}

// update the internal options
func (b builder) WithOption1(val string) model.Builder {
   b.option1 = val
   return b
}

// update the internal options
func (b builder) WithOption2(val int) model.Builder {
   b.option2 = val
   return b
}

// update the internal options
func (b builder) WithOption3(val float64) model.Builder {
   b.option3 = val
   return b
}

// update the internal options
func (b builder) WithOption4(val []string) model.Builder {
   b.option4 = val
   return b
}

// update the internal options
func (b builder) WithOption5(val []int) model.Builder {
   b.option5 = val
   return b
}

// update the internal options
func (b builder) WithOption6(val []float64) model.Builder {
   b.option6 = val
   return b
}
