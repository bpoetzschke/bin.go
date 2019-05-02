package word_manager

type Word struct {
	Value   string
	AddedBy string
	GifUrl  string
}

type FoundWord struct {
	Word
	FoundBy string
}
