package loki

import (
	"testing"

	apiv1 "github.com/elyosemite/openwatchit/api/v1"
)

func TestToLogQL_NoFilters(t *testing.T) {
	ast := &apiv1.AST{}
	got, err := toLogQL(ast)
	if err != nil {
		t.Fatal(err)
	}
	if got != "{}" {
		t.Errorf("got %q, want {}", got)
	}
}

func TestToLogQL_EqFilter(t *testing.T) {
	ast := &apiv1.AST{
		Filters: []*apiv1.Filter{
			{Field: "level", Op: apiv1.FilterOp_FILTER_OP_EQ, Value: "error"},
		},
	}
	got, err := toLogQL(ast)
	if err != nil {
		t.Fatal(err)
	}
	want := `{level="error"}`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestToLogQL_NeFilter(t *testing.T) {
	ast := &apiv1.AST{
		Filters: []*apiv1.Filter{
			{Field: "level", Op: apiv1.FilterOp_FILTER_OP_NE, Value: "debug"},
		},
	}
	got, err := toLogQL(ast)
	if err != nil {
		t.Fatal(err)
	}
	want := `{level!="debug"}`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestToLogQL_MultipleFilters(t *testing.T) {
	ast := &apiv1.AST{
		Filters: []*apiv1.Filter{
			{Field: "level", Op: apiv1.FilterOp_FILTER_OP_EQ, Value: "error"},
			{Field: "service", Op: apiv1.FilterOp_FILTER_OP_NE, Value: "debug-svc"},
		},
	}
	got, err := toLogQL(ast)
	if err != nil {
		t.Fatal(err)
	}
	want := `{level="error", service!="debug-svc"}`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestToLogQL_UnsupportedOp(t *testing.T) {
	ast := &apiv1.AST{
		Filters: []*apiv1.Filter{
			{Field: "duration", Op: apiv1.FilterOp_FILTER_OP_UNSPECIFIED, Value: "500"},
		},
	}
	_, err := toLogQL(ast)
	if err == nil {
		t.Fatal("expected error for unsupported op, got nil")
	}
}
