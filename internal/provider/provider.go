package provider

import "context"

type Message struct {
	Role    string
	Content string
}

type Chunk struct {
	Delta string
	Done  bool
	Err   error
}

type Provider interface {
	Name() string
	Stream(ctx context.Context, msgs []Message) (<-chan Chunk, error)
}
