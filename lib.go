package spoon

// Quote quotes the string.
func Quote(unquoted string) string {
	return "`" + unquoted + "`"
}

// Semicolon adds a semicolon at the end.
func Semicolon(schema string) string {
	return schema + ";"
}
