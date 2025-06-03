package service

type Theme string
type MediaMakerKind string

const (
	ThemeDark          Theme          = "dark"
	ThemeLight         Theme          = "light"
	MediaMakerKindHTML MediaMakerKind = "html"
)

type MediaParams struct {
	Theme Theme
}

type Listener interface {
	Start() error
	AddMaker(maker MediaMaker) error
}

type MediaMaker interface {
	AddPublisher(publisher Publisher) error
	Make(data any) (any, error)
}

type Publisher interface {
	Publish(data any) (any, error)
}

type MakeFlux struct {
	Maker     []MediaMaker
	Publisher []Publisher
}

type MediaCore struct {
	CoreConfig CoreConfig
	Listener   Listener
	Flux       []MakeFlux
}

type CoreConfig struct {
	Listeners []struct {
		Name   string       `yaml:"name"`
		Kind   ListenerKind `yaml:"kind"`
		Makers []string     `yaml:"makers"`
	} `yaml:"listeners"`
	Makers []struct {
		Name       string         `yaml:"name"`
		Kind       MediaMakerKind `yaml:"kind"`
		Publishers []string       `yaml:"publishers"`
	} `yaml:"makers"`
	Publishers []struct {
		Name string        `yaml:"name"`
		Kind PublisherKind `yaml:"kind"`
	} `yaml:"publishers"`
}

func NewMediaCore() {

}
