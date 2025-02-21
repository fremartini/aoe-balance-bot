package bot

type Command struct {
	Handle func(context *Context, args []string) error
	Hint   string
	Hidden bool
}
