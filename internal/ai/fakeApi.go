package ai

import (
	"math/rand"
	"time"
)

type FakeApi struct {
	responses []string
	rng       *rand.Rand
	stream    bool
	delay     time.Duration
}

func InitializeFakeApi(stream bool, delay time.Duration) *FakeApi {
	return &FakeApi{
		responses: []string{
			"TO-USER: Hello",
			"TO-USER: How can I help you",
			"TO-USER: Goodbye",
			"TO-SERVICE: getTasks",
		},
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
		stream: stream,
		delay:  delay,
	}
}

func (fakeApi *FakeApi) ChatComplete(conv *Conversation) <-chan ChatPart {
	ch := make(chan ChatPart)
	go func() {
		defer close(ch)
		value := fakeApi.responses[fakeApi.rng.Intn(len(fakeApi.responses))]
		if fakeApi.stream {
			for _, char := range value {
				time.Sleep(fakeApi.delay)
				ch <- ChatPart{Value: string(char), Err: nil}
			}
		} else {
			time.Sleep(fakeApi.delay)
			ch <- ChatPart{Value: value, Err: nil}
		}

	}()
	return ch
}

func (fakeapi *FakeApi) Metadata() map[string]any {
	return map[string]any{
		"apiName": "fake",
		"delay":   fakeapi.delay.Milliseconds(),
		"stream":  fakeapi.stream,
	}
}
