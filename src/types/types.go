package types

type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	}
	Subscription string `json:"subscription"`
}
