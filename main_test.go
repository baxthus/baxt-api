package main

import (
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func readBody(closer io.ReadCloser) string {
	buf := new(strings.Builder)
	_, _ = io.Copy(buf, closer)
	defer func() {
		err := closer.Close()
		if err != nil {
			panic(err)
		}
	}()
	return buf.String()
}

func TestIpRoute(t *testing.T) {
	const ExceptedStatusCode = fiber.StatusSeeOther
	const ExceptedLocation = "https://youtu.be/dQw4w9WgXcQ"

	req := httptest.NewRequest("GET", "/ip", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != ExceptedStatusCode {
		t.Fatalf("Expected status code %v, got %d. Body: %v", ExceptedStatusCode, resp.StatusCode, readBody(resp.Body))
	}

	location, err := resp.Location()
	if err != nil {
		t.Fatal(err)
	}
	if location.String() != ExceptedLocation {
		t.Fatalf("Expected redirect location %v, got %v. Body: %v", ExceptedLocation, location.String(), readBody(resp.Body))
	}

	t.Log("Success")
}
