package serializer

import (
	"github.com/bytedance/sonic"
	"github.com/labstack/echo/v4"
)

type SerializerType int

const (
	SerializerTypeSonic SerializerType = iota
)

func New() echo.JSONSerializer {
	defaultConfig := DefaultConfig()
	return &sonicJSONSerializer{defaultConfig.SonicConfig.Froze()}
}

func NewWithConfig(config Config) echo.JSONSerializer {
	return &sonicJSONSerializer{api: config.SonicConfig.Froze()}
}

func Type() SerializerType {
	return SerializerTypeSonic
}

// sonicJSONSerializer implements JSON encoding using github.com/bytedance/sonic
type sonicJSONSerializer struct {
	api sonic.API
}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func (s sonicJSONSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	enc := s.api.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func (s sonicJSONSerializer) Deserialize(c echo.Context, i interface{}) error {
	err := s.api.NewDecoder(c.Request().Body).Decode(i)
	return err
}
