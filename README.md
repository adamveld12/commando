# Commando

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/adamveld12/commando)
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/adamveld12/commando)
[![Gocover](http://gocover.io/_badge/github.com/adamveld12/commando)](http://gocover.io/github.com/adamveld12/commando)
[![Go Report Card](https://goreportcard.com/badge/github.com/adamveld12/commando)](https://goreportcard.com/report/github.com/adamveld12/commando)


[![wercker status](https://app.wercker.com/status/8fb474f2adee226f89217e8d669464b9/m "wercker status")](https://app.wercker.com/project/bykey/8fb474f2adee226f89217e8d669464b9)

A simple bare bones CLI helper.

Features include:

- Auto parses arguments to the right type for your handlers
- Free help text/usages output using `help | h | --help`

## How to use

```
func main(){
  app := commando.New()
  app.Add(add, "adds any 2 numbers", "add", "sum")
  app.Add(sayhello, "greets you", "hello", "hi")
  app.Add(hotOrNot,  "do you think this is hot?", "hot")

  app.Execute(os.Args[1:]...)

  // or any of the following:
  // app.Execute("add", "2", "3")
  // app.Execute("hello", "Steve")
  // app.Execute("sum", "1", "2")
  // app.Execute("ho", "false")
}

func hotOrNot(hot bool){
  if hot {
    fmt.Println("this cli is so hot")
  } else {
    fmt.Println("this cli is so not :(")
  }
}

func sayHello(name string) {
  fmt.Printf("Hello, %s!\n", name)
}

func add(a, b int){
  fmt.Printf("%d\n", a + b)
}
```

## License

MIT
