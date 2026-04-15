package images

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	generated "github.com/Blackhawk-Software/imgwire-go/generated"
)

type ImgwireImage struct {
	generated.ImageSchema
}

type StandardUploadResponse struct {
	Image                ImgwireImage
	UploadURL            string
	AdditionalProperties map[string]interface{}
}

type URLPreset string
type GravityType string
type ResizingType string
type OutputFormat string

const (
	PresetThumbnail URLPreset = "thumbnail"
	PresetSmall     URLPreset = "small"
	PresetMedium    URLPreset = "medium"
	PresetLarge     URLPreset = "large"
)

const (
	GravityNorth     GravityType = "no"
	GravitySouth     GravityType = "so"
	GravityEast      GravityType = "ea"
	GravityWest      GravityType = "we"
	GravityNorthEast GravityType = "noea"
	GravityNorthWest GravityType = "nowe"
	GravitySouthEast GravityType = "soea"
	GravitySouthWest GravityType = "sowe"
	GravityCenter    GravityType = "ce"
)

const (
	ResizingFit      ResizingType = "fit"
	ResizingFill     ResizingType = "fill"
	ResizingFillDown ResizingType = "fill-down"
	ResizingForce    ResizingType = "force"
	ResizingAuto     ResizingType = "auto"
)

const (
	FormatJPG  OutputFormat = "jpg"
	FormatPNG  OutputFormat = "png"
	FormatAVIF OutputFormat = "avif"
	FormatGIF  OutputFormat = "gif"
	FormatWEBP OutputFormat = "webp"
)

type URLOptions struct {
	Preset            *URLPreset
	Background        *string
	Blur              *float64
	Crop              *string
	DPR               *float64
	Enlarge           *bool
	Extend            *string
	ExtendAspectRatio *string
	Flip              *string
	Format            *OutputFormat
	Gravity           *string
	Height            *int
	KeepCopyright     *bool
	MinHeight         *int
	MinWidth          *int
	Padding           *string
	Pixelate          *float64
	Quality           *int
	ResizingType      *ResizingType
	Rotate            *int
	Sharpen           *float64
	StripColorProfile *bool
	StripMetadata     *bool
	Width             *int
	Zoom              *float64
}

type transformationEntry struct {
	canonical  string
	cacheValue string
}

func ExtendImage(image generated.ImageSchema) ImgwireImage {
	return ImgwireImage{ImageSchema: image}
}

func ExtendImagePtr(image *generated.ImageSchema) *ImgwireImage {
	if image == nil {
		return nil
	}
	extended := ExtendImage(*image)
	return &extended
}

func ExtendStandardUploadResponse(
	response *generated.StandardUploadResponseSchema,
) *StandardUploadResponse {
	if response == nil {
		return nil
	}
	return &StandardUploadResponse{
		Image:                ExtendImage(response.Image),
		UploadURL:            response.UploadUrl,
		AdditionalProperties: response.AdditionalProperties,
	}
}

func (image ImgwireImage) URL(options URLOptions) (string, error) {
	builder := urlBuilder{image: image}
	return builder.build(options)
}

type urlBuilder struct {
	image ImgwireImage
}

func (b urlBuilder) build(options URLOptions) (string, error) {
	parsed, err := url.Parse(b.image.CdnUrl)
	if err != nil {
		return "", err
	}

	path := parsed.Path
	if options.Preset != nil {
		path, err = appendPresetToPath(parsed.Path, *options.Preset)
		if err != nil {
			return "", err
		}
	}

	entries, err := parseTransformationEntries(options)
	if err != nil {
		return "", err
	}
	if len(entries) == 0 {
		parsed.Path = path
		parsed.RawQuery = ""
		return parsed.String(), nil
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].canonical < entries[j].canonical
	})
	query := url.Values{}
	for _, entry := range entries {
		query.Set(entry.canonical, entry.cacheValue)
	}

	parsed.Path = path
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func parseTransformationEntries(options URLOptions) ([]transformationEntry, error) {
	entries := make([]transformationEntry, 0, 23)

	if options.Background != nil {
		entry, err := parseBackground(*options.Background)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if options.Blur != nil {
		entry, err := parseNonNegativeNumberRule("blur", *options.Blur)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if options.Crop != nil {
		entry, err := parseCrop(*options.Crop)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if options.DPR != nil {
		entry, err := parsePositiveNumberRule("dpr", *options.DPR)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if options.Enlarge != nil && *options.Enlarge {
		entries = append(entries, transformationEntry{canonical: "enlarge", cacheValue: "true"})
	}
	if options.Extend != nil {
		entry, ok, err := parseExtendLike("extend", *options.Extend)
		if err != nil {
			return nil, err
		}
		if ok {
			entries = append(entries, entry)
		}
	}
	if options.ExtendAspectRatio != nil {
		entry, ok, err := parseExtendLike("extend_aspect_ratio", *options.ExtendAspectRatio)
		if err != nil {
			return nil, err
		}
		if ok {
			entries = append(entries, entry)
		}
	}
	if options.Flip != nil {
		entry, err := parseFlip(*options.Flip)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if options.Format != nil {
		entry, err := parseFormat(*options.Format)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if options.Gravity != nil {
		entry, err := parseGravity(*options.Gravity)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if options.Height != nil {
		entry, err := parsePositiveIntRule("height", *options.Height)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if options.KeepCopyright != nil {
		entries = append(entries, parseBooleanRule("keep_copyright", *options.KeepCopyright))
	}
	if options.MinHeight != nil {
		entry, err := parsePositiveIntRule("min-height", *options.MinHeight)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if options.MinWidth != nil {
		entry, err := parsePositiveIntRule("min-width", *options.MinWidth)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if options.Padding != nil {
		entry, err := parsePadding(*options.Padding)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if options.Pixelate != nil {
		entry, err := parsePositiveNumberRule("pixelate", *options.Pixelate)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if options.Quality != nil {
		if *options.Quality < 0 || *options.Quality > 100 {
			return nil, invalidTransformation("quality")
		}
		entries = append(entries, transformationEntry{canonical: "quality", cacheValue: strconv.Itoa(*options.Quality)})
	}
	if options.ResizingType != nil {
		entry, err := parseResizingType(*options.ResizingType)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if options.Rotate != nil {
		entry, err := parseRotate(*options.Rotate)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if options.Sharpen != nil {
		entry, err := parseNonNegativeNumberRule("sharpen", *options.Sharpen)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if options.StripColorProfile != nil {
		entries = append(entries, parseBooleanRule("strip_color_profile", *options.StripColorProfile))
	}
	if options.StripMetadata != nil {
		entries = append(entries, parseBooleanRule("strip_metadata", *options.StripMetadata))
	}
	if options.Width != nil {
		entry, err := parsePositiveIntRule("width", *options.Width)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if options.Zoom != nil {
		entry, err := parsePositiveNumberRule("zoom", *options.Zoom)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func appendPresetToPath(path string, preset URLPreset) (string, error) {
	switch preset {
	case PresetThumbnail, PresetSmall, PresetMedium, PresetLarge:
	default:
		return "", invalidTransformation("preset")
	}

	slashIndex := strings.LastIndex(path, "/")
	prefix := ""
	fileName := path
	if slashIndex >= 0 {
		prefix = path[:slashIndex+1]
		fileName = path[slashIndex+1:]
	}
	dotIndex := strings.LastIndex(fileName, ".")
	if dotIndex <= 0 || dotIndex == len(fileName)-1 {
		return "", fmt.Errorf("cannot apply an image URL preset to a CDN url without a file extension")
	}
	return fmt.Sprintf("%s%s@%s", prefix, fileName, preset), nil
}

func parsePositiveIntRule(canonical string, value int) (transformationEntry, error) {
	if value < 1 {
		return transformationEntry{}, invalidTransformation(canonical)
	}
	return transformationEntry{canonical: canonical, cacheValue: strconv.Itoa(value)}, nil
}

func parsePositiveNumberRule(canonical string, value float64) (transformationEntry, error) {
	if value <= 0 {
		return transformationEntry{}, invalidTransformation(canonical)
	}
	return transformationEntry{canonical: canonical, cacheValue: stringifyNumber(value)}, nil
}

func parseNonNegativeNumberRule(canonical string, value float64) (transformationEntry, error) {
	if value < 0 {
		return transformationEntry{}, invalidTransformation(canonical)
	}
	return transformationEntry{canonical: canonical, cacheValue: stringifyNumber(value)}, nil
}

func parseBooleanRule(canonical string, value bool) transformationEntry {
	if value {
		return transformationEntry{canonical: canonical, cacheValue: "true"}
	}
	return transformationEntry{canonical: canonical, cacheValue: "false"}
}

func parseResizingType(value ResizingType) (transformationEntry, error) {
	switch value {
	case ResizingFit, ResizingFill, ResizingFillDown, ResizingForce, ResizingAuto:
		return transformationEntry{canonical: "resizing_type", cacheValue: string(value)}, nil
	default:
		return transformationEntry{}, invalidTransformation("resizing_type")
	}
}

func parseFormat(value OutputFormat) (transformationEntry, error) {
	switch value {
	case FormatJPG, FormatPNG, FormatAVIF, FormatGIF, FormatWEBP:
		return transformationEntry{canonical: "format", cacheValue: string(value)}, nil
	default:
		return transformationEntry{}, invalidTransformation("format")
	}
}

func parseRotate(value int) (transformationEntry, error) {
	switch value {
	case 0, 90, 180, 270, 360:
		return transformationEntry{canonical: "rotate", cacheValue: strconv.Itoa(value)}, nil
	default:
		return transformationEntry{}, invalidTransformation("rotate")
	}
}

func parseBackground(value string) (transformationEntry, error) {
	parts := strings.Split(value, ":")
	if len(parts) == 3 {
		normalized := make([]string, 0, 3)
		for _, part := range parts {
			component, err := strconv.Atoi(part)
			if err != nil || component < 0 || component > 255 {
				return transformationEntry{}, invalidTransformation("background")
			}
			normalized = append(normalized, strconv.Itoa(component))
		}
		return transformationEntry{canonical: "background", cacheValue: strings.Join(normalized, ":")}, nil
	}

	hexColor := strings.TrimPrefix(value, "#")
	if len(hexColor) != 3 && len(hexColor) != 6 && len(hexColor) != 8 {
		return transformationEntry{}, invalidTransformation("background")
	}
	for _, character := range hexColor {
		if !strings.ContainsRune("0123456789abcdefABCDEF", character) {
			return transformationEntry{}, invalidTransformation("background")
		}
	}
	return transformationEntry{canonical: "background", cacheValue: strings.ToLower(hexColor)}, nil
}

func parseCrop(value string) (transformationEntry, error) {
	parts := strings.Split(value, ":")
	if len(parts) < 2 {
		return transformationEntry{}, invalidTransformation("crop")
	}
	width, err := parsePositiveStringNumber(parts[0], "crop")
	if err != nil {
		return transformationEntry{}, err
	}
	height, err := parsePositiveStringNumber(parts[1], "crop")
	if err != nil {
		return transformationEntry{}, err
	}
	gravity := "ce:0:0"
	if len(parts) > 2 {
		gravity, err = parseGravityParts(parts[2:], "crop", true)
		if err != nil {
			return transformationEntry{}, err
		}
	}
	return transformationEntry{
		canonical:  "crop",
		cacheValue: fmt.Sprintf("%s:%s:%s", width, height, gravity),
	}, nil
}

func parsePadding(value string) (transformationEntry, error) {
	parts := strings.Split(value, ":")
	if len(parts) < 1 || len(parts) > 4 {
		return transformationEntry{}, invalidTransformation("padding")
	}
	normalized := make([]string, 0, len(parts))
	for _, part := range parts {
		number, err := parseNonNegativeStringNumber(part, "padding")
		if err != nil {
			return transformationEntry{}, err
		}
		normalized = append(normalized, number)
	}
	return transformationEntry{canonical: "padding", cacheValue: strings.Join(normalized, ":")}, nil
}

func parseFlip(value string) (transformationEntry, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return transformationEntry{}, invalidTransformation("flip")
	}
	first, err := parseBoolString(parts[0], "flip")
	if err != nil {
		return transformationEntry{}, err
	}
	second, err := parseBoolString(parts[1], "flip")
	if err != nil {
		return transformationEntry{}, err
	}
	return transformationEntry{
		canonical:  "flip",
		cacheValue: fmt.Sprintf("%s:%s", serializeBool(first), serializeBool(second)),
	}, nil
}

func parseGravity(value string) (transformationEntry, error) {
	gravity, err := parseGravityParts(strings.Split(value, ":"), "gravity", true)
	if err != nil {
		return transformationEntry{}, err
	}
	return transformationEntry{canonical: "gravity", cacheValue: gravity}, nil
}

func parseExtendLike(canonical string, value string) (transformationEntry, bool, error) {
	parts := strings.Split(value, ":")
	if len(parts) < 1 {
		return transformationEntry{}, false, invalidTransformation(canonical)
	}
	extend, err := parseBoolString(parts[0], canonical)
	if err != nil {
		return transformationEntry{}, false, err
	}
	if !extend {
		return transformationEntry{}, false, nil
	}
	if len(parts) == 1 {
		return transformationEntry{canonical: canonical, cacheValue: "true"}, true, nil
	}
	gravity, err := parseGravityParts(parts[1:], canonical, false)
	if err != nil {
		return transformationEntry{}, false, err
	}
	return transformationEntry{canonical: canonical, cacheValue: "true:" + gravity}, true, nil
}

func parseGravityParts(parts []string, label string, allowSmart bool) (string, error) {
	if len(parts) != 1 && len(parts) != 2 && len(parts) != 3 {
		return "", invalidTransformation(label)
	}
	switch GravityType(parts[0]) {
	case GravityNorth, GravitySouth, GravityEast, GravityWest, GravityNorthEast,
		GravityNorthWest, GravitySouthEast, GravitySouthWest, GravityCenter:
	default:
		return "", invalidTransformation(label)
	}
	if len(parts) == 1 {
		return parts[0], nil
	}
	if len(parts) == 2 {
		if !allowSmart || parts[1] != "sm" {
			return "", invalidTransformation(label)
		}
		return parts[0] + ":sm", nil
	}
	if _, err := strconv.Atoi(parts[1]); err != nil {
		return "", invalidTransformation(label)
	}
	if _, err := strconv.Atoi(parts[2]); err != nil {
		return "", invalidTransformation(label)
	}
	return fmt.Sprintf("%s:%s:%s", parts[0], parts[1], parts[2]), nil
}

func parsePositiveStringNumber(value string, label string) (string, error) {
	parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil || parsed <= 0 {
		return "", invalidTransformation(label)
	}
	return stringifyNumber(parsed), nil
}

func parseNonNegativeStringNumber(value string, label string) (string, error) {
	parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil || parsed < 0 {
		return "", invalidTransformation(label)
	}
	return stringifyNumber(parsed), nil
}

func parseBoolString(value string, label string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "t", "1":
		return true, nil
	case "false", "f", "0":
		return false, nil
	default:
		return false, invalidTransformation(label)
	}
}

func serializeBool(value bool) string {
	if value {
		return "t"
	}
	return "f"
}

func stringifyNumber(value float64) string {
	if value == float64(int64(value)) {
		return strconv.FormatInt(int64(value), 10)
	}
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func invalidTransformation(label string) error {
	return fmt.Errorf("invalid transformation rule value for %s", label)
}
