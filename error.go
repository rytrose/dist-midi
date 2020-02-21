package util

// Must panics if there's an error.
func Must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
