package arguments

import (
	"os"
	"strings"
)

type Arguments struct {
	ShortFlags []string          // can take seperately, like "-t", "-o", "-m", or all together like "-tmo", stored as "t", "m", "o"
	LongFlags  []string          // can take seperately, like "--token", "--output", "--media", stored as "token", "output", "media"
	Keys       map[string]string // --key value or --key=value
}

func (a *Arguments) GetKey(key string) string { // return the key or an empty string
	return a.Keys[key]
}

func (a *Arguments) GetKeyEither(short string, long string) string { // return the key or an empty string
	value := a.Keys[short]
	if value != "" {
		return value
	}
	return a.Keys[long]
}

func (a *Arguments) HasShortFlag(flag string) bool { // return true if the short flag is present
	for i := 0; i < len(a.ShortFlags); i++ {
		if a.ShortFlags[i] == flag {
			return true
		}
	}
	return false
}

func (a *Arguments) HasLongFlag(flag string) bool { // return true if the long flag is present
	for i := 0; i < len(a.LongFlags); i++ {
		if a.LongFlags[i] == flag {
			return true
		}
	}
	return false
}

func (a *Arguments) HasFlag(short string, long string) bool { // return true if the flag is present
	if a.HasShortFlag(short) || a.HasLongFlag(long) {
		return true
	}
	return false
}

func (a *Arguments) DebugDump() {
	for i := 0; i < len(a.ShortFlags); i++ {
		println("ShortFlags:", a.ShortFlags[i])
	}
	for i := 0; i < len(a.LongFlags); i++ {
		println("LongFlags:", a.LongFlags[i])
	}
	for k, v := range a.Keys {
		println("Keys:", k, "-", v)
	}

}

func NewArguments() *Arguments { // create an Arguments struct, process it, and return it
	a := new(Arguments)
	a.Keys = make(map[string]string) // Initialize the map here

	// process os.Args and store it in the struct
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if strings.HasPrefix(arg, "-") {
			if strings.Contains(arg, "=") {
				// Handle --key=value format
				parts := strings.SplitN(arg, "=", 2)
				key := parts[0]
				if strings.HasPrefix(key, "--") {
					key = key[2:]
				} else {
					key = key[1:]
				}
				a.Keys[key] = parts[1]
			} else if i+1 < len(os.Args) && !strings.HasPrefix(os.Args[i+1], "-") {
				// Handle --key value format
				key := arg
				if strings.HasPrefix(key, "--") {
					key = key[2:]
				} else {
					key = key[1:]
				}
				a.Keys[key] = os.Args[i+1]
				i++ // skip next arg
			} else {
				if strings.HasPrefix(arg, "--") {
					a.LongFlags = append(a.LongFlags, arg[2:])
				} else {
					a.ShortFlags = append(a.ShortFlags, arg[1:])
				}
			}
		}
	}

	return a
}
