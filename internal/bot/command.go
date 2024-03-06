package bot

type Command struct {
	Handle func(context *Context, args []string)
	Hint   string
}
