package resources

import (
	"fmt"

	generated "github.com/imgwire/imgwire-go/generated"
)

func toCustomMetadata(input map[string]any) (map[string]generated.CustomMetadataValue, error) {
	output := make(map[string]generated.CustomMetadataValue, len(input))
	for key, value := range input {
		metadataValue, err := toCustomMetadataValue(value)
		if err != nil {
			return nil, fmt.Errorf("custom_metadata[%s]: %w", key, err)
		}
		output[key] = metadataValue
	}
	return output, nil
}

func toCustomMetadataValue(input any) (generated.CustomMetadataValue, error) {
	switch value := input.(type) {
	case string:
		return generated.CustomMetadataValue{String: &value}, nil
	case bool:
		return generated.CustomMetadataValue{Bool: &value}, nil
	case int:
		converted := int32(value)
		return generated.CustomMetadataValue{Int32: &converted}, nil
	case int32:
		return generated.CustomMetadataValue{Int32: &value}, nil
	case int64:
		converted := int32(value)
		return generated.CustomMetadataValue{Int32: &converted}, nil
	case float32:
		return generated.CustomMetadataValue{Float32: &value}, nil
	case float64:
		converted := float32(value)
		return generated.CustomMetadataValue{Float32: &converted}, nil
	default:
		return generated.CustomMetadataValue{}, fmt.Errorf("unsupported type %T", input)
	}
}
