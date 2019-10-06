package format

// TODO:...
func BuildFormatsRepo() []Formatter {
	fRepo := make([]Formatter, 0, 10)

	fRepo = append(fRepo, Keywords{true})                     // to be configurable
	fRepo = append(fRepo, WsAroundKeywords{true, true, true}) // to be configurable

	return fRepo
}
