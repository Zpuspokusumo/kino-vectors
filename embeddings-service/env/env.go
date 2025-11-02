package env

type ENV struct {
	QdrantAPIKEY          string
	QdrantMovieCollection string
}

func Setup() ENV {
	return ENV{
		QdrantAPIKEY:          "localhost:6334:",
		QdrantMovieCollection: "testmovie01",
	}
}
