package loki

import (
	"fmt"
	"strings"

	apiv1 "github.com/elyosemite/openwatchit/api/v1"
)

// toLogQL converts an OWL AST to a Loki LogQL stream selector string.
// Only FILTER_OP_EQ and FILTER_OP_NE are supported in v0.1.
// Multiple filters are AND-combined inside the stream selector.
//
// Examples:
//
//	level eq error              → {level="error"}
//	level eq error, service ne  → {level="error", service!="debug"}
func toLogQL(ast *apiv1.AST) (string, error) {
	filters := ast.GetFilters()
	if len(filters) == 0 {
		return "{}", nil
	}

	parts := make([]string, 0, len(filters))
	for _, f := range filters {
		expr, err := filterExpr(f)
		if err != nil {
			return "", err
		}
		parts = append(parts, expr)
	}

	return "{" + strings.Join(parts, ", ") + "}", nil
}

func filterExpr(f *apiv1.Filter) (string, error) {
	switch f.GetOp() {
	case apiv1.FilterOp_FILTER_OP_EQ:
		return fmt.Sprintf(`%s="%s"`, f.GetField(), f.GetValue()), nil
	case apiv1.FilterOp_FILTER_OP_NE:
		return fmt.Sprintf(`%s!="%s"`, f.GetField(), f.GetValue()), nil
	default:
		return "", fmt.Errorf("unsupported filter op %v for field %q", f.GetOp(), f.GetField())
	}
}
