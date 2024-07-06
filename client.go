package allsrvc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	
	"github.com/jsteenb2/errors"
)

const (
	resourceTypeFoo = "foo"
)

var (
	ErrIDRequired = errors.Kind("id is required")
)

type ClientHTTP struct {
	addr string
	c    *http.Client
}

func NewClientHTTP(addr string, c *http.Client) *ClientHTTP {
	return &ClientHTTP{
		addr: addr,
		c:    c,
	}
}

// Foo API types
type (
	FooCreateAttrs struct {
		Name string `json:"name"`
		Note string `json:"note"`
	}
	
	// ResourceFooAttrs are the attributes of a foo resource.
	ResourceFooAttrs struct {
		Name      string `json:"name"`
		Note      string `json:"note"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
)

func (c *ClientHTTP) CreateFoo(ctx context.Context, attrs FooCreateAttrs) (RespBody[ResourceFooAttrs], error) {
	req, err := jsonReq(ctx, "POST", c.fooPath(""), newFooData("", attrs))
	if err != nil {
		return RespBody[ResourceFooAttrs]{}, errors.Wrap(err, "create foo")
	}
	
	resp, err := c.doFooReq(req)
	return resp, errors.Wrap(err)
}

func (c *ClientHTTP) ReadFoo(ctx context.Context, id string) (RespBody[ResourceFooAttrs], error) {
	if id == "" {
		return RespBody[ResourceFooAttrs]{}, errors.Wrap(ErrIDRequired)
	}
	
	req, err := http.NewRequestWithContext(ctx, "GET", c.fooPath(id), nil)
	if err != nil {
		return RespBody[ResourceFooAttrs]{}, errors.Wrap(err)
	}
	
	resp, err := c.doFooReq(req)
	return resp, errors.Wrap(err)
}

type FooUpdAttrs struct {
	Name *string `json:"name"`
	Note *string `json:"note"`
}

func (c *ClientHTTP) UpdateFoo(ctx context.Context, id string, attrs FooUpdAttrs) (RespBody[ResourceFooAttrs], error) {
	req, err := jsonReq(ctx, "PATCH", c.fooPath(id), newFooData(id, attrs))
	if err != nil {
		return RespBody[ResourceFooAttrs]{}, errors.Wrap(err)
	}
	
	resp, err := c.doFooReq(req)
	return resp, errors.Wrap(err)
}

func (c *ClientHTTP) DelFoo(ctx context.Context, id string) (RespBody[any], error) {
	if id == "" {
		return RespBody[any]{}, errors.Wrap(ErrIDRequired)
	}
	
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.fooPath(id), nil)
	if err != nil {
		return RespBody[any]{}, errors.Wrap(err)
	}
	
	resp, err := doReq[any](c.c, req)
	return resp, errors.Wrap(err)
}

func (c *ClientHTTP) doFooReq(req *http.Request) (RespBody[ResourceFooAttrs], error) {
	resp, err := doReq[ResourceFooAttrs](c.c, req)
	return resp, errors.Wrap(err)
}

func (c *ClientHTTP) fooPath(id string) string {
	u := c.addr + "/v1/foos"
	if id == "" {
		return u
	}
	return u + "/" + id
}

func newFooData[Attr Attrs](id string, attrs Attr) Data[Attr] {
	return Data[Attr]{
		Type:  resourceTypeFoo,
		ID:    id,
		Attrs: attrs,
	}
}

// jsonReq here uses generics to provide feedback to developers when they provide some other field.
// This improves the feedback loop working with these methods. If they copy pasta wrong and provide
// ReqBody[Attr] instead, this will reject that.
func jsonReq[Attr Attrs](ctx context.Context, method, path string, v Data[Attr]) (*http.Request, error) {
	reqBody := ReqBody[Attr]{Data: v}
	
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(reqBody); err != nil {
		return nil, errors.Wrap(err, "failed to json encode request body")
	}
	
	req, err := http.NewRequestWithContext(ctx, method, path, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	
	return req, nil
}

func doReq[Attr Attrs](c *http.Client, req *http.Request) (RespBody[Attr], error) {
	resp, err := c.Do(req)
	if err != nil {
		return *new(RespBody[Attr]), errors.Wrap(err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	
	// put an upper limit on how much data we'll read off the wire. A little bit of protection
	// from a resource exhaustion attack.
	b, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1MB limit
	if err != nil {
		return *new(RespBody[Attr]), errors.Wrap(err, "failed to read response body")
	}
	
	if resp.Header.Get("Content-Type") != "application/json" {
		return *new(RespBody[Attr]), errors.Wrap(err, "invalid content type received", errors.KVs("content", string(b)))
	}
	
	var respBody RespBody[Attr]
	if err := json.Unmarshal(b, &respBody); err != nil {
		return *new(RespBody[Attr]), errors.Wrap(err)
	}
	
	return respBody, nil
}
