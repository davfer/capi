package datum

type Parser interface {
	Parse(string) ([]DataKey, error)
	Replace(string, DataKey, string) string
}
