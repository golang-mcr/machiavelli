package steganogopher

type Crypter interface {
	Encode(interface{}, interface{}) interface{}
	Decode(interface{}) interface{}
}
