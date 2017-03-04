package commando

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var errNotSupported = errors.New("not supported")

var parseFuncs = []argParser{
	parseBaseTypes,
	parseFloatTypes,
	parseIntTypes,
	parseUintTypes,
}

type argParser func(string, reflect.Kind) (interface{}, error)

type command struct {
	Names       []string
	Description string
	handler     interface{}
}

func (cc command) match(commandString string) bool {
	for _, str := range cc.Names {
		if commandString == str {
			return true
		}
	}
	return false
}

// CommandMux manages command registration
type CommandMux struct {
	commands []command
}

// New creates a new CommandMux instance
func New() *CommandMux {
	cMux := &CommandMux{}
	cMux.Add("help h --help", "Displays usages", func() {
		fmt.Printf("%s", cMux.Usage())
	})

	return cMux
}

// Usage returns a string with the command, argument types and description for each handler
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

// Add registers a new handler with the specified names (separated by spaces), description and handler func
func (c *CommandMux) Add(names, description string, handlerFunc interface{}) {
	funcType := reflect.TypeOf(handlerFunc)

	if funcType.Kind() != reflect.Func {
		panic("Didn't pass a func")
	}

	cmd := command{strings.Split(names, " "), description, handlerFunc}
	c.commands = append(c.commands, cmd)
}

// Execute executes a command (specified in arg[0]) passing in the remaining arguments as parameters (arg[1:])
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
		if cDef.match(commandString) {
			handler := cDef.handler
			return executeHandler(commandString, handler, commandArgs)
		}
	}

	return fmt.Errorf("\"%s\" is not a recognized command", commandString)
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
		val, err := parseArgument(args[i], paramType)
		if err != nil {
			return err
		}

		inputArgs = append(inputArgs, reflect.ValueOf(val))
	}

	reflect.ValueOf(handler).Call(inputArgs)
	return nil
}

func parseArgument(arg string, paramType reflect.Type) (interface{}, error) {
	kind := paramType.Kind()
	var val interface{}
	var err error
	for _, parser := range parseFuncs {
		val, err = parser(arg, kind)
		if err == errNotSupported {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("expected %s but got %s", paramType.Name(), arg)
		} else {
			return val, nil
		}
	}

	return nil, fmt.Errorf("no idea how to parse %+v", paramType.Name())
}

func parseBaseTypes(arg string, k reflect.Kind) (val interface{}, err error) {
	switch k {
	case reflect.String:
		val = arg
	case reflect.Bool:
		val, err = strconv.ParseBool(arg)
	default:
		err = errNotSupported
	}

	return

}

func parseFloatTypes(arg string, k reflect.Kind) (val interface{}, err error) {
	val, err = strconv.ParseFloat(arg, 64)
	switch k {
	case reflect.Float64:
	case reflect.Float32:
		val = val.(float32)
	default:
		err = errNotSupported
	}

	return
}

func parseUintTypes(arg string, k reflect.Kind) (val interface{}, err error) {
	val, err = strconv.ParseUint(arg, 10, 64)
	valuint := val.(uint64)
	switch k {
	case reflect.Uint:
		val = uint(valuint)
	case reflect.Uint8:
		val = uint8(valuint)
	case reflect.Uint16:
		val = uint16(valuint)
	case reflect.Uint32:
		val = uint32(valuint)
	case reflect.Uint64:
	default:
		err = errNotSupported
	}

	return
}

func parseIntTypes(arg string, k reflect.Kind) (val interface{}, err error) {
	val, err = strconv.ParseInt(arg, 10, 64)
	valint := val.(int64)
	switch k {
	case reflect.Int:
		val = int(valint)
	case reflect.Int8:
		val = int8(valint)
	case reflect.Int16:
		val = int16(valint)
	case reflect.Int32:
		val = int32(valint)
	case reflect.Int64:
	default:
		err = errNotSupported
	}

	return val, err
}
