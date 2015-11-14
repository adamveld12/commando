package commando

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Command struct {
	Names       []string
	Description string
	handler     interface{}
}

func (cc Command) Match(commandString string) bool {
	for _, str := range cc.Names {
		if commandString == str {
			return true
		}
	}
	return false
}

type CommandMux struct {
	commands []Command
}

func New() *CommandMux {
	cMux := &CommandMux{}
	cMux.Add("help h --help", "Displays usages", func() {
		fmt.Printf("%s", cMux.Usage())
	})

	return cMux
}

func (c *CommandMux) Usage() string {
	usageStr := "Usage:\n"
	for _, command := range c.commands {

		args := handlerArguments(command.handler)

		argumentList := []string{}
		for _, arg := range args {
			argumentList = append(argumentList, fmt.Sprintf("%s", arg.Name()))
		}

		argumentsStr := strings.Join(argumentList, ", ")
		commandString := strings.Join(command.Names, ", ")
		if argumentsStr != "" {
			argumentsStr = fmt.Sprintf(" [%s]", argumentsStr)
		}
		usageStr += fmt.Sprintf("%s%s\t%s\n", commandString, argumentsStr, command.Description)
	}
	return usageStr
}

func (c *CommandMux) Add(names, description string, handlerFunc interface{}) {
	funcType := reflect.TypeOf(handlerFunc)

	if funcType.Kind() != reflect.Func {
		panic("Didn't pass a func")
	}

	command := Command{strings.Split(names, " "), description, handlerFunc}
	c.commands = append(c.commands, command)
}

func (c *CommandMux) Execute(args ...string) error {
	argc := len(args)

	if argc == 0 {
		return errors.New("not enough arguments")
	}

	commandString := args[0]
	commandArgs := []string{}

	if argc > 1 {
		commandArgs = args[1:]
	}

	for _, cDef := range c.commands {
		if cDef.Match(commandString) {
			handler := cDef.handler
			return executeHandler(commandString, handler, commandArgs)
		}
	}

	return errors.New(fmt.Sprintf("\"%s\" is not a recognized command", commandString))
}

func handlerArguments(handler interface{}) []reflect.Type {
	handlerType := reflect.TypeOf(handler)

	results := []reflect.Type{}

	for i := 0; i < handlerType.NumIn(); i++ {
		paramType := handlerType.In(i)
		results = append(results, paramType)
	}

	return results
}

func executeHandler(commandString string, handler interface{}, args []string) error {
	handlerType := reflect.TypeOf(handler)

	numIn := handlerType.NumIn()

	if numIn != len(args) {
		return errors.New("Not enough arguments")
	}

	inputArgs := []reflect.Value{}

	for i, paramType := range handlerArguments(handler) {
		argStr := args[i]

		var val interface{} = argStr
		var err error

		switch paramType.Kind() {
		case reflect.Bool:
			val, err = strconv.ParseBool(argStr)
		case reflect.Float64:
			val, err = strconv.ParseFloat(argStr, 64)
		case reflect.Float32:
			val, err = strconv.ParseFloat(argStr, 64)
			val = val.(float32)
		case reflect.Int:
			val, err = strconv.Atoi(argStr)
		case reflect.Int8:
			val, err = strconv.Atoi(argStr)
			val = int8(val.(int))
		case reflect.Int16:
			val, err = strconv.Atoi(argStr)
			val = int16(val.(int))
		case reflect.Int32:
			val, err = strconv.Atoi(argStr)
			val = int32(val.(int))
		case reflect.Int64:
			val, err = strconv.Atoi(argStr)
			val = int64(val.(int))
		case reflect.Uint:
			val, err = strconv.ParseUint(argStr, 10, 64)
			val = uint(val.(uint64))
		case reflect.Uint8:
			val, err = strconv.ParseUint(argStr, 10, 8)
			val = uint8(val.(uint64))
		case reflect.Uint16:
			val, err = strconv.ParseUint(argStr, 10, 16)
			val = uint16(val.(uint64))
		case reflect.Uint32:
			val, err = strconv.ParseUint(argStr, 10, 32)
			val = uint32(val.(uint64))
		case reflect.Uint64:
			val, err = strconv.ParseUint(argStr, 10, 64)
			val = uint64(val.(uint64))
		case reflect.String:
		default:
			panic(fmt.Sprintf("%s arguments are not supported", paramType.Kind()))
		}

		if err != nil {
			return errors.New(fmt.Sprintf("\"%s\" expects %s but got %s",
				commandString, paramType.Name(), argStr))
		}

		inputArgs = append(inputArgs, reflect.ValueOf(val))
	}

	reflect.ValueOf(handler).Call(inputArgs)

	return nil
}
