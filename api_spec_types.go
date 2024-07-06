package allsrvc

// JSON API spec envelope types
type (
	// RespBody represents a JSON-API response body.
	// 	https://jsonapi.org/format/#document-top-level
	//
	// note: data can be either an array or a single resource object. This allows for both.
	RespBody[Attr Attrs] struct {
		Meta RespMeta    `json:"meta"`
		Errs []RespErr   `json:"errors,omitempty"`
		Data *Data[Attr] `json:"data,omitempty"`
	}
	
	// Data represents a JSON-API data response.
	//
	//	https://jsonapi.org/format/#document-top-level
	Data[Attr Attrs] struct {
		Type  string `json:"type"`
		ID    string `json:"id"`
		Attrs Attr   `json:"attributes"`
		
		// omitting the relationships here for brevity not at lvl 3 RMM
	}
	
	// Attrs can be either a document or a collection of documents.
	Attrs interface {
		any | []Attrs
	}
	
	// RespMeta represents a JSON-API meta object. The data here is
	// useful for our example service. You can add whatever non-standard
	// context that is relevant to your domain here.
	//	https://jsonapi.org/format/#document-meta
	RespMeta struct {
		TookMilli int    `json:"took_ms"`
		TraceID   string `json:"trace_id"`
	}
	
	// RespErr represents a JSON-API error object. Do note that we
	// aren't implementing the entire error type. Just the most impactful
	// bits for this workshop. Mainly, skipping Title & description separation.
	//	https://jsonapi.org/format/#error-objects
	RespErr struct {
		Status int            `json:"status,string"`
		Code   int            `json:"code"`
		Msg    string         `json:"message"`
		Source *RespErrSource `json:"source"`
	}
	
	// RespErrSource represents a JSON-API err source.
	//	https://jsonapi.org/format/#error-objects
	RespErrSource struct {
		Pointer   string `json:"pointer"`
		Parameter string `json:"parameter,omitempty"`
		Header    string `json:"header,omitempty"`
	}
	
	// ReqBody represents a JSON-API request body.
	//	https://jsonapi.org/format/#crud-creating
	ReqBody[Attr Attrs] struct {
		Data Data[Attr] `json:"data"`
	}
)
