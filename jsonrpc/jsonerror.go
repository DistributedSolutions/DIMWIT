package jsonrpc

// -32700				Parse error	Invalid JSON was received by the server.
// 						An error occurred on the server while parsing the JSON text.
// -32600				Invalid Request	The JSON sent is not a valid Request object.
// -32601				Method not found	The method does not exist / is not available.
// -32602				Invalid params	Invalid method parameter(s).
// -32603				Internal error	Internal JSON-RPC error.
// -32000 to -32099		Server error	Reserved for implementation-defined server-errors.

// -32700				Parse error	Invalid JSON was received by the server.
func NewParseError(extra string, id uint32) *JSONRPCReponse {
	j := new(JSONRPCReponse)
	e := new(JSONError)
	e.Code = 32700
	e.Message = "Invalid Json"
	e.Data = extra
	j.Error = e
	j.Id = id

	return j
}

// -32600				Invalid Request	The JSON sent is not a valid Request object.
func NewInvalidRequestError(extra string, id uint32) *JSONRPCReponse {
	j := new(JSONRPCReponse)
	e := new(JSONError)
	e.Code = 32600
	e.Message = "Invalid Request"
	e.Data = extra
	j.Error = e
	j.Id = id

	return j
}

// -32601				Method not found	The method does not exist / is not available.
func NewMethodNotFoundError(extra string, id uint32) *JSONRPCReponse {
	j := new(JSONRPCReponse)
	e := new(JSONError)
	e.Code = 32601
	e.Message = "Method not found"
	e.Data = extra
	j.Error = e
	j.Id = id

	return j
}

// -32602				Invalid params	Invalid method parameter(s).
func NewInvalidParametersError(extra string, id uint32) *JSONRPCReponse {
	j := new(JSONRPCReponse)
	e := new(JSONError)
	e.Code = 32602
	e.Message = "Invalid params"
	e.Data = extra
	j.Error = e
	j.Id = id

	return j
}

// -32603				Internal error	Internal JSON-RPC error.
func NewInternalRPCSError(extra string, id uint32) *JSONRPCReponse {
	j := new(JSONRPCReponse)
	e := new(JSONError)
	e.Code = 32603
	e.Message = "Internal error"
	e.Data = extra
	j.Error = e
	j.Id = id

	return j
}

// -32000 to -32099		Server error	Reserved for implementation-defined server-errors.
func NewCustomError(extra string, id uint32, errorID uint32) *JSONRPCReponse {
	j := new(JSONRPCReponse)
	e := new(JSONError)
	e.Code = 32000 + errorID
	e.Message = "Server error"
	e.Data = extra
	j.Error = e
	j.Id = id

	return j
}
