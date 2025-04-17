package dump

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func AsJSON(v any) {
	AsJSONToWriter(os.Stdout, v)
}

func AsJSONToWriter(w io.Writer, v any) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	if err := enc.Encode(v); err != nil {
		panic(fmt.Sprintf("dump.AsJSON(%T): %v", v, err))
	}
}

func AsJSONString(v any) string {
	var s strings.Builder
	AsJSONToWriter(&s, v)
	return s.String()
}
