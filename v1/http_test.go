package core

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDurationUnmarshalSupportsString(t *testing.T) {
	var d Duration
	if err := json.Unmarshal([]byte(`"5s"`), &d); err != nil {
		t.Fatal(err)
	}

	if got := d.Or(0); got != 5*time.Second {
		t.Fatalf("expected 5s, got %s", got)
	}
}

func TestDurationUnmarshalSupportsNumericSeconds(t *testing.T) {
	var d Duration
	if err := json.Unmarshal([]byte(`10`), &d); err != nil {
		t.Fatal(err)
	}

	if got := d.Or(0); got != 10*time.Second {
		t.Fatalf("expected 10s, got %s", got)
	}
}

func TestResponseBytesUsesStructuredError(t *testing.T) {
	payload := responseBytes(
		"error",
		nil,
		&APIError{
			Code:    "bad_request",
			Message: "bad request",
		},
		Metadata{RequestID: "req-1"},
	)

	var out Response
	if err := json.Unmarshal(payload, &out); err != nil {
		t.Fatal(err)
	}

	if out.Error == nil {
		t.Fatal("expected structured error")
	}
	if out.Error.Code != "bad_request" {
		t.Fatalf("expected error code bad_request, got %q", out.Error.Code)
	}
	if out.Metadata.RequestID != "req-1" {
		t.Fatalf("expected request id req-1, got %q", out.Metadata.RequestID)
	}
}
