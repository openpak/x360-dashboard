package dashboard

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHomeManifestContract(t *testing.T) {
	body, err := HomeXML()
	if err != nil {
		t.Fatal(err)
	}
	h := &Handler{Manifest: body}
	for _, locale := range []string{"en-US", "en-IE", "ja-JP", "zh-Hant-TW"} {
		r := httptest.NewRequest(http.MethodGet, "/epix/"+locale+"/dashhome.xml", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 200 {
			t.Fatalf("%s: %d", locale, w.Code)
		}
		if w.Header().Get("Content-Length") != strconv.Itoa(w.Body.Len()) {
			t.Fatal("incorrect content length")
		}
		var manifest struct {
			XMLName  xml.Name `xml:"manifest"`
			ID       string   `xml:"id"`
			Default  string   `xml:"defaultchannelid"`
			Channels []struct {
				ID         string `xml:"id"`
				Definition struct {
					Layout string `xml:"layouttemplate,attr"`
					Groups []struct {
						Slots []struct {
							Definition string `xml:"definitionpath"`
							Name       string `xml:"name"`
							Click      struct {
								Button  string `xml:"button"`
								Action  string `xml:"action"`
								Command string `xml:"cmd"`
								Param   string `xml:"param"`
							} `xml:"onclick"`
						} `xml:"slot"`
					} `xml:"slotgroup"`
				} `xml:"channeldef"`
			} `xml:"channel"`
		}
		if err := xml.Unmarshal(w.Body.Bytes(), &manifest); err != nil {
			t.Fatal(err)
		}
		if manifest.ID == "" || manifest.Default != "HOME" || len(manifest.Channels) != 3 || manifest.Channels[0].ID != manifest.Default {
			t.Fatalf("unselectable Home: %+v", manifest)
		}
		seen := map[string]bool{}
		for _, channel := range manifest.Channels {
			if seen[channel.ID] {
				t.Fatalf("duplicate channel %s", channel.ID)
			}
			seen[channel.ID] = true
			def := channel.Definition
			if def.Layout != "epix://TemplateOffline1HD.xur" || len(def.Groups) != 4 {
				t.Fatalf("invalid layout for %s: %+v", channel.ID, def)
			}
			for _, group := range def.Groups {
				if len(group.Slots) != 1 {
					t.Fatal("each layout group needs one tile")
				}
				slot := group.Slots[0]
				if slot.Definition == "slot://DiscInTray" {
					continue
				}
				if slot.Name == "" || slot.Click.Button != "A" || slot.Click.Action != "EpixCmd" || slot.Click.Command == "" {
					t.Fatalf("inert tile in %s: %+v", channel.ID, slot)
				}
				if slot.Name == "My Games" && (slot.Click.Command != "EcNavToURI" || slot.Click.Param != "gameslibrary:_RootTitle=%slot_description%;") {
					t.Fatalf("My Games must open the native installed-title library: %+v", slot)
				}
			}
		}
		for _, id := range []string{"HOME", "GAMES", "SETTINGS"} {
			if !seen[id] {
				t.Fatalf("missing channel %s", id)
			}
		}
	}
}

func TestHTTPBehaviour(t *testing.T) {
	body, err := HomeXML()
	if err != nil {
		t.Fatal(err)
	}
	h := &Handler{Manifest: body}
	for _, tc := range []struct {
		method, path string
		status       int
	}{
		{"HEAD", "/epix/en-US/dashhome.xml", 200},
		{"GET", "/healthz", 200},
		{"GET", "/epix/en-US/missing.xml", 404},
		{"GET", "/epix/not-a-locale/dashhome.xml", 404},
		{"POST", "/epix/en-US/dashhome.xml", 405},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
		if w.Code != tc.status {
			t.Fatalf("%s %s: %d", tc.method, tc.path, w.Code)
		}
		if tc.method == "HEAD" && w.Body.Len() != 0 {
			t.Fatal("HEAD must have no body")
		}
	}
}
