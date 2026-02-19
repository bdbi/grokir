package command

import "fmt"

var commands = make(map[string]Command)

// Register adds a command to the global registry.
func Register(cmd Command) {
	commands[cmd.Name()] = cmd
}

// Get retrieves a command by name from the registry.
func Get(name string) (Command, bool) {
	cmd, ok := commands[name]
	return cmd, ok
}

// Names returns a list of all registered command names.
func Names() []string {
	var names []string
	for name := range commands {
		names = append(names, name)
	}
	return names
}

// Run executes a command by name with the given runtime and arguments.
func Run(ctx Runtime, name string, args []string) error {
	cmd, ok := Get(name)
	if !ok {
		return NewError(fmt.Sprintf("unknown command: %s", name), false)
	}
	return cmd.Run(ctx, args)
}
