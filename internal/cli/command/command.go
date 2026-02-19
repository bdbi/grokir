package command

// Command defines the interface that all CLI commands must implement.
type Command interface {
	Name() string
	Usage() string
	Run(ctx Runtime, args []string) error
}
