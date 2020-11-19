package adapter

type ShippingCreateInput struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type ShippingUpdateInput struct {
	Name string `json:"name"`
}

type ShippingProviderSearchOptions struct {
	Name string `json:"name" bson:"name"`
}

type ShippingAdapter interface {
	DecodeCreateInput([]byte) (ShippingCreateInput, error)
	DecodeUpdateInput([]byte) (ShippingUpdateInput, error)
}
