package gif

type Gif interface {
	Random(searchQuery string) (string, bool, error)
}
