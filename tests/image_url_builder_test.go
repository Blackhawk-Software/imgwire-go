package tests

import (
	"testing"

	imgwire "github.com/Blackhawk-Software/imgwire-go"
	generated "github.com/Blackhawk-Software/imgwire-go/generated"
)

func makeImage() imgwire.Image {
	return imgwire.Image{
		ImageSchema: generated.ImageSchema{
			CdnUrl:           "https://cdn.imgwire.dev/example.jpg",
			CustomMetadata:   map[string]generated.CustomMetadataValue{},
			ExifData:         map[string]interface{}{},
			Extension:        "jpg",
			Height:           100,
			Id:               "img_1",
			MimeType:         "image/jpeg",
			OriginalFilename: "example.jpg",
			SizeBytes:        100,
			Status:           "READY",
			Width:            100,
		},
	}
}

func TestImageURLBuildsTransformedURL(t *testing.T) {
	image := makeImage()
	preset := imgwire.PresetThumbnail
	background := "#ffffff"
	height := 150
	width := 150
	rotate := 90

	url, err := image.URL(imgwire.ImageURLOptions{
		Preset:     &preset,
		Background: &background,
		Height:     &height,
		Rotate:     &rotate,
		Width:      &width,
	})
	if err != nil {
		t.Fatalf("build image url: %v", err)
	}

	expected := "https://cdn.imgwire.dev/example.jpg@thumbnail?background=ffffff&height=150&rotate=90&width=150"
	if url != expected {
		t.Fatalf("unexpected url %q", url)
	}
}

func TestImageURLOmitsFalseEnlargeAndNormalizesBooleans(t *testing.T) {
	image := makeImage()
	enlarge := false
	stripMetadata := true

	url, err := image.URL(imgwire.ImageURLOptions{
		Enlarge:       &enlarge,
		StripMetadata: &stripMetadata,
	})
	if err != nil {
		t.Fatalf("build image url: %v", err)
	}

	expected := "https://cdn.imgwire.dev/example.jpg?strip_metadata=true"
	if url != expected {
		t.Fatalf("unexpected url %q", url)
	}
}

func TestImageURLRejectsInvalidTransformationValue(t *testing.T) {
	image := makeImage()
	rotate := 45

	_, err := image.URL(imgwire.ImageURLOptions{Rotate: &rotate})
	if err == nil {
		t.Fatalf("expected invalid rotation error")
	}
}
