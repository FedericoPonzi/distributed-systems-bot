package pkg

type Publisher interface {
	publishText() error
	publishImage() error
	publishLink() error
	publishVideo() error
}
