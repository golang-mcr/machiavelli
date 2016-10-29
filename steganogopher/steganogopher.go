package steganogopher

// Encoder is the signature of the Encode function. Specs TBD
type Encoder func(interface{}, interface{}) interface{}

// Decoder is the signature of the Decode function. Specs TBD.
type Decoder func(interface{}, interface{}) interface{}
