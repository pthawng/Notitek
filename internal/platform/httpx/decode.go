package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/pthawng/notitek/internal/platform/apperr"
)

// DecodeJSON reads exactly one JSON value from the request body into dst.
//
// It requires Content-Type application/json, rejects unknown fields (the API
// schemas set additionalProperties: false) and anything after the value. Every
// failure is an *apperr.Error with field-level details and never echoes the
// input, because bodies can hold OTPs, links and contacts:
//
//   - wrong Content-Type: UNSUPPORTED_MEDIA_TYPE (415)
//   - body above the MaxBody limit: PAYLOAD_TOO_LARGE (413)
//   - empty, malformed, truncated or mistyped JSON: INVALID_REQUEST (400)
func DecodeJSON(r *http.Request, dst any) error {
	if err := requireJSON(r); err != nil {
		return err
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return decodeError(err)
	}
	// Only whitespace may follow the value.
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			return payloadTooLarge(tooLarge.Limit)
		}
		return apperr.New(apperr.InvalidRequest, "Nội dung yêu cầu phải là đúng một giá trị JSON.",
			apperr.FieldError(apperr.LocationBody, "/", apperr.ReasonMalformedJSON, "Có dữ liệu thừa sau giá trị JSON."))
	}
	return nil
}

func requireJSON(r *http.Request) error {
	ct := r.Header.Get("Content-Type")
	mt, _, err := mime.ParseMediaType(ct)
	if err == nil && mt == "application/json" {
		return nil
	}
	return apperr.New(apperr.UnsupportedMediaType, "Content-Type phải là application/json.",
		apperr.FieldError(apperr.LocationHeader, "Content-Type", apperr.ReasonInvalidValue, "Giá trị được hỗ trợ: application/json."))
}

func decodeError(err error) error {
	var (
		tooLarge *http.MaxBytesError
		syntax   *json.SyntaxError
		typed    *json.UnmarshalTypeError
	)
	switch {
	case errors.Is(err, io.EOF):
		return apperr.New(apperr.InvalidRequest, "Nội dung yêu cầu không được để trống.",
			apperr.FieldError(apperr.LocationBody, "/", apperr.ReasonRequired, "Cần có nội dung JSON."))
	case errors.As(err, &tooLarge):
		return payloadTooLarge(tooLarge.Limit)
	case errors.Is(err, io.ErrUnexpectedEOF):
		return apperr.New(apperr.InvalidRequest, "JSON không hợp lệ.",
			apperr.FieldError(apperr.LocationBody, "/", apperr.ReasonMalformedJSON, "Nội dung JSON bị cắt dở."))
	case errors.As(err, &syntax):
		return apperr.New(apperr.InvalidRequest, "JSON không hợp lệ.",
			apperr.FieldError(apperr.LocationBody, "/", apperr.ReasonMalformedJSON,
				fmt.Sprintf("Sai cú pháp JSON gần byte %d.", syntax.Offset)))
	case errors.As(err, &typed):
		return apperr.New(apperr.InvalidRequest, "Có trường sai kiểu dữ liệu.",
			apperr.FieldError(apperr.LocationBody, pointer(typed.Field), apperr.ReasonInvalidType,
				"Trường phải có kiểu "+jsonKind(typed.Type.Kind().String())+"."))
	}
	if name, ok := unknownField(err); ok {
		return apperr.New(apperr.InvalidRequest, "Có trường không được hỗ trợ.",
			apperr.FieldError(apperr.LocationBody, pointer(name), apperr.ReasonUnknownField, "Trường này không có trong hợp đồng."))
	}
	return apperr.New(apperr.InvalidRequest, "JSON không hợp lệ.",
		apperr.FieldError(apperr.LocationBody, "/", apperr.ReasonMalformedJSON, "Không đọc được nội dung JSON."))
}

func payloadTooLarge(limit int64) *apperr.Error {
	return apperr.New(apperr.PayloadTooLarge, fmt.Sprintf("Nội dung yêu cầu vượt quá %d byte.", limit),
		apperr.FieldError(apperr.LocationBody, "/", apperr.ReasonTooLong, fmt.Sprintf("Tối đa %d byte.", limit)))
}

// unknownField extracts the field name from encoding/json's unknown-field
// error, which has no exported type. Only the name is used, never a value.
func unknownField(err error) (string, bool) {
	const prefix = `json: unknown field "`
	msg := err.Error()
	if !strings.HasPrefix(msg, prefix) || !strings.HasSuffix(msg, `"`) {
		return "", false
	}
	return strings.TrimSuffix(strings.TrimPrefix(msg, prefix), `"`), true
}

// pointer converts the dotted path of encoding/json ("recipients.phone") into
// a JSON Pointer ("/recipients/phone"). encoding/json does not report array
// indexes, so a path through an array points at the field of its elements.
func pointer(dotted string) string {
	if dotted == "" {
		return "/"
	}
	parts := strings.Split(dotted, ".")
	for i, p := range parts {
		p = strings.ReplaceAll(p, "~", "~0")
		parts[i] = strings.ReplaceAll(p, "/", "~1")
	}
	return "/" + strings.Join(parts, "/")
}

// jsonKind names a Go kind the way API consumers know it.
func jsonKind(goKind string) string {
	switch {
	case strings.HasPrefix(goKind, "int"), strings.HasPrefix(goKind, "uint"), strings.HasPrefix(goKind, "float"):
		return "number"
	case goKind == "bool":
		return "boolean"
	case goKind == "slice", goKind == "array":
		return "array"
	case goKind == "map", goKind == "struct":
		return "object"
	default:
		return goKind
	}
}
