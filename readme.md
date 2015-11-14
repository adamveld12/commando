# Commando

A simple bare bones CLI helper. 

Features include:
- Auto parses arguments to the right type for your handlers
- Free help text/usages output using `help | h | --help`

## How to use

```
func main(){
  app := commando.New()
  app.Add("add  sum", "adds any 2 numbers", add)
  app.Add("hello  hi", "greets you", add)
  app.Add("hot", "greets you", hotOrNot)

  app.Execute(os.Args[1:]...) 
  
  // or any of the following:
  // app.Execute("add", "2", "3")
  // app.Execute("hello", "Steve")
  // app.Execute("sum", "1", "2")
  // app.Execute("ho", "false")
}

func hotOrNot(hot bool){
  if hot {
    fmt.Println("I'm so hot")
  } else {
    fmt.Println("I'm so not :(")
  }
}

func sayHello(name string) {
  fmt.Printf("Hello, %s!\n", name)
}

func add(a, b int){
  fmt.Printf("%d\n", a + b)
}
```
