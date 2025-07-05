package patterns

var registered []Pattern

func Register(p Pattern) {
	registered = append(registered, p)
}

func All() []Pattern {
	return registered
}
