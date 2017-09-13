package flag

import (
	"strings"

	"code.cloudfoundry.org/cli/types"

	"github.com/cloudfoundry/bytefmt"
	flags "github.com/jessevdk/go-flags"
)

const (
	ALLOWED_UNITS = "mg"
)

type Megabytes struct {
	types.NullByteSize
}

func (m *Megabytes) UnmarshalFlag(val string) error {
	if val == "" {
		return nil
	}

	size, err := bytefmt.ToMegabytes(val)

	if err != nil ||
		!strings.ContainsAny(strings.ToLower(val), ALLOWED_UNITS) ||
		strings.Contains(strings.ToLower(val), ".") {
		return &flags.Error{
			Type:    flags.ErrRequired,
			Message: `Byte quantity must be an integer with a unit of measurement like M, MB, G, or GB`,
		}
	}

	m.Value = size
	m.IsSet = true
	m.IsBytes = false
	return nil
}
