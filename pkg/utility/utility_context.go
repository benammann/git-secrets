package utility

import "context"

type AddChannel chan int
type DoneChannel chan bool

const (
	ContextKeyAddChannel = "AddChannel"
	ContextKeyDoneChannel = "DoneChannel"
)

func NewChannelContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyAddChannel, make(AddChannel))
	ctx = context.WithValue(ctx, ContextKeyDoneChannel, make(DoneChannel))
	return ctx
}

func GetContextChannels(ctx context.Context) (AddChannel, DoneChannel) {
	return ctx.Value(ContextKeyAddChannel).(AddChannel), ctx.Value(ContextKeyDoneChannel).(DoneChannel)
}