package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/denwilliams/cli/internal/aussupermarkets/cli"
	"github.com/denwilliams/cli/internal/aussupermarkets/providers"
	"github.com/denwilliams/cli/internal/aussupermarkets/search"
)

type fake struct {
	id search.Retailer
	fn func(search.Request) (search.Result, error)
}

func (f fake) Retailer() search.Retailer                                         { return f.id }
func (f fake) Search(_ context.Context, r search.Request) (search.Result, error) { return f.fn(r) }

func TestInvalidFlags(t *testing.T) {
	for _, args := range [][]string{nil, {" "}, {"milk", "--sort", "name"}, {"milk", "-sort", "name"}, {"milk", "bread"}, {"--unknown", "milk"}, {"--max-results", "0", "milk"}, {"--max-results", "-1", "milk"}, {"--max-results", "abc", "milk"}, {"--sort", "cost", "milk"}, {"--format", "xml", "milk"}, {"--timeout", "0s", "milk"}, {"--timeout", "-1s", "milk"}, {"--timeout", "what", "milk"}, {"--postcode", "800", "milk"}, {"--postcode", "abcd", "milk"}, {"--providers", "", "milk"}, {"--providers", "coles,", "milk"}, {"--providers", "coles,coles", "milk"}, {"--providers", "other", "milk"}, {"--store", "coles=", "milk"}, {"--store", "other=1", "milk"}, {"--store", "coles=1", "--store", "coles=2", "milk"}, {"--providers", "aldi", "--store", "coles=1", "milk"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			var out, errOut bytes.Buffer
			code := cli.Run(context.Background(), args, &out, &errOut, providers.All())
			if code != 2 || out.Len() != 0 || errOut.Len() == 0 {
				t.Fatalf("code=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
			}
		})
	}
}

func TestHelpDoesNotSearch(t *testing.T) {
	for _, arg := range []string{"--help", "-h"} {
		var out, errOut bytes.Buffer
		code := cli.Run(context.Background(), []string{arg}, &out, &errOut, []search.Provider{fake{search.Coles, func(search.Request) (search.Result, error) { t.Error("help searched"); return search.Result{}, nil }}})
		if code != 0 || !strings.Contains(errOut.String(), "no live in-store feeds") || !strings.Contains(errOut.String(), "max-results") {
			t.Fatalf("%d %s", code, errOut.String())
		}
	}
}

func TestProductionJSONAndHumanState(t *testing.T) {
	for _, format := range []string{"json", "human"} {
		var out, errOut bytes.Buffer
		code := cli.Run(context.Background(), []string{"--format", format, "--postcode", "0800", "--store", "coles=provided", "milk"}, &out, &errOut, providers.All())
		if code != 1 || errOut.Len() != 0 {
			t.Fatalf("%d %s", code, errOut.String())
		}
		if format == "json" {
			var report search.Report
			if err := json.Unmarshal(out.Bytes(), &report); err != nil {
				t.Fatal(err)
			}
			if len(report.Groups) != 3 || report.Groups[0].RequestedLocation.StoreID != "provided" {
				t.Fatal(report)
			}
			for _, g := range report.Groups {
				if g.Status != "unconfigured" || g.Products == nil || len(g.Products) != 0 {
					t.Fatal(g)
				}
			}
		} else {
			if !strings.Contains(out.String(), "unconfigured") || strings.Contains(out.String(), "AUD ") {
				t.Fatal(out.String())
			}
		}
	}
}

func TestEnabledSelectionContextAndEmptySuccess(t *testing.T) {
	selected := func(id search.Retailer) search.Provider {
		return fake{id, func(req search.Request) (search.Result, error) {
			if req.Query != "milk 2L" || req.Location.Postcode != "0800" || req.Location.Channel != search.InStore {
				t.Errorf("bad request: %+v", req)
			}
			if id == search.Coles && req.Location.StoreID != "chosen" {
				t.Error("store lost")
			}
			if id == search.ALDI && req.Location.StoreID != "" {
				t.Error("invented store")
			}
			loc := req.Location
			if loc.StoreID == "" {
				loc.StoreID = "resolved-test-store"
			}
			return search.Result{Location: loc, Source: search.Source{Name: "synthetic fixture", URL: "https://example.invalid/feed"}, Coverage: search.Coverage{Scope: "complete_catalogue"}}, nil
		}}
	}
	ps := []search.Provider{selected(search.Coles), fake{search.Woolworths, func(search.Request) (search.Result, error) {
		t.Error("disabled provider called")
		return search.Result{}, nil
	}}, selected(search.ALDI)}
	var out, errOut bytes.Buffer
	code := cli.Run(context.Background(), []string{"--providers", "aldi,coles", "--format", "json", "--postcode", "0800", "--store", "coles=chosen", "milk 2L"}, &out, &errOut, ps)
	var r search.Report
	if err := json.Unmarshal(out.Bytes(), &r); err != nil {
		t.Fatal(err)
	}
	if code != 0 || len(r.Groups) != 2 || r.Groups[0].Retailer != search.Coles || r.Groups[1].Retailer != search.ALDI || r.Failed() {
		t.Fatalf("%d %s %+v", code, errOut.String(), r)
	}
}

func TestPartialFailureExitAndLiteralDashQuery(t *testing.T) {
	ps := []search.Provider{fake{search.Coles, func(req search.Request) (search.Result, error) {
		if req.Query != "-milk" {
			t.Error(req.Query)
		}
		loc := req.Location
		loc.StoreID = "resolved-test-store"
		return search.Result{Location: loc, Source: search.Source{Name: "synthetic fixture", URL: "https://example.invalid/feed"}, Coverage: search.Coverage{Scope: "sample"}}, nil
	}}, fake{search.ALDI, func(search.Request) (search.Result, error) { return search.Result{}, search.ErrUnavailable }}}
	var out, errOut bytes.Buffer
	code := cli.Run(context.Background(), []string{"--providers", "coles,aldi", "--format", "json", "--", "-milk"}, &out, &errOut, ps)
	var r search.Report
	if err := json.Unmarshal(out.Bytes(), &r); err != nil {
		t.Fatal(err)
	}
	if code != 1 || r.Groups[0].Status != "empty" || r.Groups[1].Status != "unavailable" {
		t.Fatal(code, r)
	}
}

type brokenWriter struct{}

func (brokenWriter) Write([]byte) (int, error) { return 0, errors.New("broken output") }
func TestOutputFailure(t *testing.T) {
	for _, format := range []string{"human", "json"} {
		var errOut bytes.Buffer
		code := cli.Run(context.Background(), []string{"--format", format, "milk"}, brokenWriter{}, &errOut, providers.All())
		if code != 1 || !strings.Contains(errOut.String(), "broken output") {
			t.Fatal(code, errOut.String())
		}
	}
}
func TestInvalidRegistry(t *testing.T) {
	if code := cli.Run(context.Background(), []string{"milk"}, io.Discard, io.Discard, []search.Provider{nil}); code != 2 {
		t.Fatal(code)
	}
}
