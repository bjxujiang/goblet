package render

import (
	"net/http"
)

type RawRender int8

func (r *RawRender) PrepareInstance(c RenderContext) (RenderInstance, error) {
	return new(RawRenderInstance), nil
}

func (r *RawRender) Init(s RenderServer) {
}

type RawRenderInstance int8

func (r *RawRenderInstance) Render(wr http.ResponseWriter, data interface{}, status int) error {
	switch tdata := data.(type) {
	case []byte:
		wr.Write(tdata)
	}
	return nil
}