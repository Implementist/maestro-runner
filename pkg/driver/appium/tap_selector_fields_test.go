package appium

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/devicelab-dev/maestro-runner/pkg/flow"
)

// Two "Save" buttons; only the second carries id "confirm_save". A text-only
// native query answers with the first, so id and state have to be checked in
// page source, where every selector field is matched together.
const twoSaveButtons = `<?xml version="1.0" encoding="UTF-8"?>
<hierarchy rotation="0">
  <android.widget.Button resource-id="com.app:id/draft_save" text="Save" clickable="true" enabled="true" checked="false" displayed="true" bounds="[0,100][500,200]"/>
  <android.widget.Button resource-id="com.app:id/confirm_save" text="Save" clickable="true" enabled="true" checked="true" displayed="true" bounds="[0,1000][500,1100]"/>
</hierarchy>`

func twoSaveServer(t *testing.T) (*Driver, *[]byte) {
	t.Helper()
	var body []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch p := r.URL.Path; {
		case strings.HasSuffix(p, "/source"):
			writeJSON(w, map[string]interface{}{"value": twoSaveButtons})
		case strings.HasSuffix(p, "/actions"):
			body, _ = io.ReadAll(r.Body)
			writeJSON(w, map[string]interface{}{"value": nil})
		case strings.HasSuffix(p, "/element") && r.Method == http.MethodPost:
			// A native text query finds the first "Save" (the wrong one).
			writeJSON(w, map[string]interface{}{"value": map[string]interface{}{w3cElementKey: "draft"}})
		case strings.HasSuffix(p, "/rect"):
			writeJSON(w, map[string]interface{}{"value": map[string]interface{}{"x": 0, "y": 100, "width": 500, "height": 100}})
		default:
			writeJSON(w, map[string]interface{}{"value": ""})
		}
	}))
	t.Cleanup(server.Close)
	return createTestAppiumDriver(server), &body
}

func TestTapOnTextWithIDOrStateMatchesAllFields(t *testing.T) {
	checked := true
	for name, sel := range map[string]flow.Selector{
		"text + id":      {Text: "Save", ID: "confirm_save"},
		"text + checked": {Text: "Save", Checked: &checked},
	} {
		t.Run(name, func(t *testing.T) {
			d, body := twoSaveServer(t)
			if res := d.tapOn(&flow.TapOnStep{Selector: sel}); !res.Success {
				t.Fatalf("tapOn failed: %s", res.Message)
			}
			if x, y := firstMove(t, *body); x != 250 || y != 1050 {
				t.Errorf("tapped (%v,%v), want the matching button at (250,1050)", x, y)
			}
		})
	}
}
