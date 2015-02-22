package app

import (
   "encoding/json"
   "fmt"
   model "lessa/org/app"
)

type Options struct {
   Option1 string
   Option2 int
   Option3 float64
   Option4 []string
   Option5 []int
   Option6 []float64
}

type builder struct {
   Options
}

type application struct {
   Options
}

func InstallAppBuilderFunc() {
   model.SetAppBuilderFunc(defaultBuilder)
}

func (o Options) String() string {
   js, _ := json.MarshalIndent(o, "", "   ")
   return string(js)
}

func (a application) Run() error {
   fmt.Print(a.String())
   return nil
}

func defaultOptions() Options {
   return Options{
      Option1: "default",
      Option2: 666,
      Option3: 666.666,
      Option4: []string{"default1","default2"},
      Option5: []int{666, 999},
      Option6: []float64{666.999, 999.666},
   }
}

func defaultBuilder() model.Builder {
   return builder{
      Options: defaultOptions(),
   }
}

func (b builder) Build() model.Application {
   return application {
      Options: b.Options,
   }
}

func (b builder) WithOption1(val string) model.Builder {
   b.Option1 = val
   return b
}

func (b builder) WithOption2(val int) model.Builder {
   b.Option2 = val
   return b
}

func (b builder) WithOption3(val float64) model.Builder {
   b.Option3 = val
   return b
}

func (b builder) WithOption4(val []string) model.Builder {
   b.Option4 = val
   return b
}

func (b builder) WithOption5(val []int) model.Builder {
   b.Option5 = val
   return b
}

func (b builder) WithOption6(val []float64) model.Builder {
   b.Option6 = val
   return b
}
