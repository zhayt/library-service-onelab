package main

import "strconv"

func Fooer(input int) string {
	if input%3 == 0 {
		return "Foo"
	}

	return strconv.Itoa(input)
}
