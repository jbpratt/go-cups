package cups

import (
	"reflect"
	"testing"
)

func TestGetMimeTypes(t *testing.T) {
	want := []string{"application/pdf", "application/postscript"}

	got, err := GetMimeTypes("mime.types")
	if err != nil {
		t.Fatalf("error: failed to read file: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetMimeTypes() = %v; want: %v", got, want)
	}

}
