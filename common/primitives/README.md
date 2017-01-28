# Primitive types
Wrapper of common types or custom types that allows structures using these types to be easily marshaled into and unmarshaled out of a binary block. This is used for storing in a database.

## Must Implment

```Go
MarshalBinary() ([]byte, error) {
	// Return the object as a byte array that is able
	// to be unmarshaled back into the type.
}

UnmarshalBinary(data []byte) error {
	_, err := UnmarshalBinaryData(data)
	return err
}

UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	// newData returned is the given data minus the bytes used
	// E.g: If a type used 20 bytes, the data given could be 
	//      100 bytes. This function will return the remaining
	//      80 bytes.
}
```

## Can Implement

Any function that could prove useful to interacting with the object.

```Go
String() string
NewOBJECTNAME() *OBJECT
```

## Write unit tests

Don't forget to write these ;)