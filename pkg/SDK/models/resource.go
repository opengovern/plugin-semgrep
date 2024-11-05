package model

type StreamSender func(Resource) error

type Resource struct {
	ID          string
	Description interface{}

	Name           string
	Type           string
	ResourceGroup  string
	Location       string
	SubscriptionID string
	TenantID       string
}

func (r Resource) UniqueID() string {
	return r.ID
}
