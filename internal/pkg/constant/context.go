package constant

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

const UserContextKey ContextKey = "user"
