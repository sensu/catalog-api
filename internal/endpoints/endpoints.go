package endpoints

type APIEndpoint interface {
	GetOutputPath() string
	GetData() interface{}
}
