// internal/httpResponse/error.go
package httpResponse

type Error struct {
	Code	int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`

}