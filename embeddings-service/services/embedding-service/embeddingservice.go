package embeddingservice

type EmbeddingService interface {
	EmbedText(text string) float32
	SetModel() string
	GetModel() string
}
