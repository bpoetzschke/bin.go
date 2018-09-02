package gif

func NewGiphy() Gif {
	return &giphy{}
}

type giphy struct {
}

func (g *giphy) Random(searchQuery string) (string, error) {
	return "", nil
}
