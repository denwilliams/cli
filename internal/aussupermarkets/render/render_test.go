package render_test

import (
	"bytes"
	"encoding/json"
	"math"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/denwilliams/cli/internal/aussupermarkets/render"
	"github.com/denwilliams/cli/internal/aussupermarkets/search"
)

func fixture() search.Report {
	observed := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	location := search.Location{Channel: search.InStore, Postcode: "0800", StoreID: "test-store"}
	return search.Report{Query: "milk", Sort: search.SortPrice, MaxResults: 2, Groups: []search.Group{
		{Retailer: search.Coles, Status: "ok", RequestedLocation: location, Location: location, Source: search.Source{Name: "synthetic fixture", URL: "https://example.invalid/feed", ObservedAt: &observed}, Coverage: search.Coverage{Scope: "sample", Truncated: true}, CandidateCount: 3, Limited: true, Products: []search.Product{
			{ID: "a", Name: "Fixture Milk", Price: &search.Money{MinorUnits: 1002, Currency: "AUD"}, PriceBasis: "pack", PackSize: "2L", UnitPrice: &search.UnitPrice{Amount: search.Money{MinorUnits: 501, Currency: "AUD"}, Unit: "1L"}, Availability: "unknown", URL: "https://example.invalid/a", ObservedAt: &observed},
			{ID: "b", Name: "Unknown Milk", Price: nil},
		}},
		{Retailer: search.ALDI, Status: "unconfigured", Error: "no permitted feed", RequestedLocation: search.Location{Channel: search.InStore}, Location: search.Location{Channel: search.InStore}, Coverage: search.Coverage{Scope: "unknown"}, Products: []search.Product{}},
	}}
}

func TestGoldenOutputs(t *testing.T) {
	for _, tc := range []struct {
		name   string
		render func(*bytes.Buffer, search.Report) error
	}{
		{"human", func(b *bytes.Buffer, r search.Report) error { return render.Human(b, r) }},
		{"json", func(b *bytes.Buffer, r search.Report) error { return render.JSON(b, r) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer
			if err := tc.render(&b, fixture()); err != nil {
				t.Fatal(err)
			}
			want, err := os.ReadFile("testdata/" + tc.name + ".golden")
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(b.Bytes(), want) {
				t.Fatalf("output changed:\n%s", b.String())
			}
		})
	}
}

func TestJSONExactMoneyNullAndMetadata(t *testing.T) {
	r := fixture()
	r.Groups[0].Products[0].Price.MinorUnits = math.MaxInt64
	var b bytes.Buffer
	if err := render.JSON(&b, r); err != nil {
		t.Fatal(err)
	}
	var decoded search.Report
	if err := json.Unmarshal(b.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r, decoded) || !strings.Contains(b.String(), `"minor_units": 9223372036854775807`) || !strings.Contains(b.String(), `"price": null`) {
		t.Fatal(b.String())
	}
}

func TestMoneyFormatting(t *testing.T) {
	for _, tc := range []struct {
		cents int64
		want  string
	}{{0, "AUD 0.00"}, {1, "AUD 0.01"}, {1002, "AUD 10.02"}, {math.MaxInt64, "AUD 92233720368547758.07"}} {
		if got := render.Money(search.Money{MinorUnits: tc.cents, Currency: "AUD"}); got != tc.want {
			t.Fatalf("%s != %s", got, tc.want)
		}
	}
}

func TestHumanEscapesUntrustedTextAndLabelsConditionalPrices(t *testing.T) {
	r := fixture()
	p := &r.Groups[0].Products[0]
	p.Name = "bad\n\x1b[31m"
	p.PriceBasis = "member"
	p.Conditions = "member only\n"
	r.Groups[1].Error = "failure\x1b[0m"
	var b bytes.Buffer
	if err := render.Human(&b, r); err != nil {
		t.Fatal(err)
	}
	if strings.ContainsRune(b.String(), '\x1b') || !strings.Contains(b.String(), `bad\n\x1b[31m`) || !strings.Contains(b.String(), "(member)") || !strings.Contains(b.String(), "conditions=") {
		t.Fatal(b.String())
	}
}

func TestConditionalPackRetainsPriceAndConditionsInBothFormats(t *testing.T) {
	r := fixture()
	p := &r.Groups[0].Products[0]
	p.Conditions = "members only"
	var human, machine bytes.Buffer
	if err := render.Human(&human, r); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(human.String(), "AUD 10.02 (pack)") || !strings.Contains(human.String(), `conditions="members only"`) {
		t.Fatal(human.String())
	}
	if err := render.JSON(&machine, r); err != nil {
		t.Fatal(err)
	}
	var decoded search.Report
	if err := json.Unmarshal(machine.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r, decoded) {
		t.Fatal("conditional pack metadata changed")
	}
}
