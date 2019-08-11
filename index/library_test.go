package index

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var tr = track{
	ID:          "ID",
	Name:        "Name",
	Album:       "Album",
	AlbumArtist: "AlbumArtist",
	Artist:      "Artist",
	Composer:    "Composer",
	Genre:       "Genre",
	Location:    "Location",
	Kind:        "Kind",

	TotalTime:   1,
	Year:        2,
	DiscNumber:  3,
	TrackNumber: 4,
	TrackCount:  5,
	DiscCount:   6,
	BitRate:     7,

	DateAdded:    time.Now(),
	DateModified: time.Now(),
}

func TestTrack(t *testing.T) {
	stringFields := []string{"ID", "Name", "Album", "AlbumArtist", "Artist", "Composer", "Genre", "Location", "Kind"}
	for _, f := range stringFields {
		got := tr.GetString(f)
		if got != f {
			t.Errorf("tr.GetString(%#v) = %#v, expected %#v", f, got, f)
		}
	}

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic from invalid field")
			}
		}()

		y := tr.GetString("Year")
		t.Errorf("expected panic from GetString, got: %v", y)
	}()

	stringsFields := []string{"AlbumArtist", "Artist", "Composer"}
	for _, f := range stringsFields {
		got := tr.GetStrings(f)
		expected := []string{f}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("tr.GetStrings(%#v) = %#v, expected %#v", f, got, expected)
		}
	}

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic from invalid field")
			}
		}()

		y := tr.GetStrings("Name")
		t.Errorf("expected panic from GetStrings, got: %v", y)
	}()

	intFields := []string{"TotalTime", "Year", "DiscNumber", "TrackNumber", "TrackCount", "DiscCount", "BitRate"}
	for i, f := range intFields {
		got := tr.GetInt(f)
		expected := i + 1
		if got != expected {
			t.Errorf("tr.GetInt(%#v) = %d, expected %d", f, got, expected)
		}
	}

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic from invalid field")
			}
		}()

		y := tr.GetInt("Name")
		t.Errorf("expected panic from GetInt, got: %v", y)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic from invalid field")
			}
		}()

		y := tr.GetTime("Name")
		t.Errorf("expected panic from GetTime, got: %v", y)
	}()
}

type testLibrary struct {
	tr *track
}

func (t testLibrary) Tracks() []Track {
	return []Track{t.tr}
}

func (t testLibrary) Track(identifier string) (Track, bool) {
	return t.tr, true
}

func TestConvert(t *testing.T) {
	tl := testLibrary{
		tr: &tr,
	}

	l := Convert(tl, "ID")

	got := l.Tracks()
	expected := tl.Tracks()

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("l.Tracks() = %v, expected: %v", got, expected)
	}

	id := "ID"
	gotTrack, _ := l.Track(id)
	expectedTrack, _ := tl.Track(id)
	if !reflect.DeepEqual(gotTrack, expectedTrack) {
		t.Errorf("l.Track(%#v) = %#v, expected: %#v", id, gotTrack, expectedTrack)
	}
}

func TestLibraryEncodeDecode(t *testing.T) {
	tl := testLibrary{
		tr: &tr,
	}

	l := Convert(tl, "ID")
	buf := &bytes.Buffer{}
	err := WriteTo(l, buf)
	if err != nil {
		t.Errorf("unexpected error in WriteTo: %v", err)
	}

	got, err := ReadFrom(buf)
	if err != nil {
		t.Errorf("unexpected error in ReadFrom: %v", err)
	}

	gotTracks := got.Tracks()
	expectedTracks := l.Tracks()

	if len(gotTracks) != len(expectedTracks) {
		t.Errorf("expected %d tracks, got: %d", len(expectedTracks), len(gotTracks))
	}

	// TODO(dhowden): Remove this mess!
	gotTrack := gotTracks[0].(*track)
	expectedTrack := expectedTracks[0].(*track)

	gotTrack.DateAdded = gotTrack.DateAdded.Local()
	gotTrack.DateModified = gotTrack.DateModified.Local()

	if diff := cmp.Diff(expectedTrack, gotTrack); diff != "" {
		t.Errorf("Encode -> Decode inconsistent; diff\n%s", diff)
	}
}
