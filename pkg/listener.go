package pkg

type LineOutputListener interface {
	// receives lines of the stdout, encouraged to process asyncronously
	Out(line string)
	// receives lines of the stdout, encouraged to process asyncronously
	Err(line string)
	// a chance to report the end result of the execution + wait for pending async tasks
	Done(err error)
}

func Compose(listeners ...LineOutputListener) LineOutputListener {
	return &composedListener{listeners}
}

type composedListener struct {
	listeners []LineOutputListener
}

func (c *composedListener) Out(line string) {
	for _, l := range c.listeners {
		l.Out(line)
	}
}

func (c *composedListener) Err(line string) {
	for _, l := range c.listeners {
		l.Err(line)
	}
}

func (c *composedListener) Done(err error) {
	for _, l := range c.listeners {
		l.Done(err)
	}
}
