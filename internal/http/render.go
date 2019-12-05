package http

import (
	"io"
)

// RenderHTML renders the HTTP template specified by its file name. It also updates the HTTP code and sets the
// Content-Type as "text/html".
func RenderHTML(ctx *Context, code int, name string, obj H) {
	newObj := H{
		"t": func(key string, args ...interface{}) string {
			return T(ctx, key, args...)
		},
	}

	for k, v := range obj {
		if _, ok := newObj[k]; !ok {
			newObj[k] = v
		}
	}

	ctx.HTML(code, name, newObj)
}

// RenderASCIIJSON serializes the given struct as JSON into the response body with unicode to ASCII string. It also
// sets the Content-Type as "application/json".
func RenderASCIIJSON(ctx *Context, code int, obj H) {
	ctx.AsciiJSON(code, obj)
}

// RenderIndentedJSON serializes the given struct as pretty JSON (indented + endlines) into the response body. It also
// sets the Content-Type as "application/json".
//
// WARNING: We recommend to use this only for debug build since printing pretty JSON is more CPU and bandwidth
// consuming. For release build, use RenderJSON() instead.
func RenderIndentedJSON(ctx *Context, code int, obj H) {
	ctx.IndentedJSON(code, obj)
}

// RenderJSON serializes the given struct as JSON into the response body. It also sets the Content-Type as
// "application/json".
func RenderJSON(ctx *Context, code int, obj H) {
	ctx.JSON(code, obj)
}

// RenderJSONP serializes the given struct as JSON into the response body. It add padding to response body to request
// data from a server residing in a different domain than the client. It also sets the Content-Type as
// "application/javascript".
func RenderJSONP(ctx *Context, code int, obj H) {
	ctx.JSONP(code, obj)
}

// RenderPureJSON serializes the given struct as JSON into the response body unlike JSON which does not replace special
// html characters with their unicode entities.
func RenderPureJSON(ctx *Context, code int, obj H) {
	ctx.PureJSON(code, obj)
}

// RenderSecureJSON serializes the given struct as Secure JSON into the response body. By default, it prepends
// "while(1)," to response body if the given struct is array values. It also sets the Content-Type as "application/json".
func RenderSecureJSON(ctx *Context, code int, obj H) {
	ctx.SecureJSON(code, obj)
}

// RenderXML serializes the given struct as XML into the response body. It also sets the Content-Type as
// "application/xml".
func RenderXML(ctx *Context, code int, obj H) {
	ctx.XML(code, obj)
}

// RenderYAML serializes the given struct as YAML into the response body. It also sets the Content-Type as
// "application/yaml".
func RenderYAML(ctx *Context, code int, obj H) {
	ctx.YAML(code, obj)
}

// RenderProtoBuf serializes the given struct as ProtoBuf into the response body. It also sets the Content-Type as
// "application/protobuf".
func RenderProtoBuf(ctx *Context, code int, obj H) {
	ctx.ProtoBuf(code, obj)
}

// RenderString writes the given string into the response body. It also sets the Content-Type as "text/plain".
func RenderString(ctx *Context, code int, format string, values ...interface{}) {
	ctx.String(code, format, values...)
}

// RenderData writes some data into the body stream and updates the HTTP code.
func RenderData(ctx *Context, code int, contentType string, data []byte) {
	ctx.Data(code, contentType, data)
}

// RenderDataFromReader writes the specified reader into the body stream and updates the HTTP code.
func RenderDataFromReader(ctx *Context, code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string) {
	ctx.DataFromReader(code, contentLength, contentType, reader, extraHeaders)
}

// RenderFile writes the specified file into the body stream in a efficient way.
func RenderFile(ctx *Context, filepath string) {
	ctx.File(filepath)
}

// RenderFileAttachment writes the specified file into the body stream in an efficient way. On the client side, the
// file will typically be downloaded with the given filename.
func RenderFileAttachment(ctx *Context, filepath, filename string) {
	ctx.FileAttachment(filepath, filename)
}

// RenderSSEvent writes a Server-Sent Event into the body stream.
func RenderSSEvent(ctx *Context, name string, message interface{}) {
	ctx.SSEvent(name, message)
}

// RenderStream sends a streaming response and returns a boolean to indicate "Is client disconnected in middle of stream".
func RenderStream(ctx *Context, step func(w io.Writer) bool) bool {
	return ctx.Stream(step)
}

// Redirect returns a HTTP redirect to the specific location.
func Redirect(ctx *Context, code int, location string) {
	ctx.Redirect(code, location)
}
