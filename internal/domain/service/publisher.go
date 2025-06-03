package service

type PublisherKind string

const (
	PublisherKindLocal  PublisherKind = "local"
	PublisherKindSqlite PublisherKind = "sqlite"
)
