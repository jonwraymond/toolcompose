package set

import (
	"errors"

	"github.com/jonwraymond/toolfoundation/adapter"
)

// Exposure exports a Toolset to protocol-specific formats.
type Exposure struct {
	toolset *Toolset
	adapter adapter.Adapter
}

// NewExposure creates an Exposure for the given toolset and adapter.
func NewExposure(ts *Toolset, adapter adapter.Adapter) *Exposure {
	return &Exposure{toolset: ts, adapter: adapter}
}

// Export converts all tools to the adapter's format.
func (e *Exposure) Export() ([]any, error) {
	if e.adapter == nil {
		return nil, errors.New("adapter is nil")
	}
	tools := e.toolset.Tools()
	result := make([]any, 0, len(tools))
	for _, t := range tools {
		converted, err := e.adapter.FromCanonical(t)
		if err != nil {
			return nil, err
		}
		result = append(result, converted)
	}
	return result, nil
}

// ExportWithWarnings converts tools and returns feature loss warnings and conversion errors.
// Unlike Export, this method continues on conversion errors and collects them for reporting.
// Callers should check the errors slice to detect tools that failed to convert.
func (e *Exposure) ExportWithWarnings() ([]any, []adapter.FeatureLossWarning, []error) {
	if e.adapter == nil {
		return nil, nil, []error{errors.New("adapter is nil")}
	}

	tools := e.toolset.Tools()
	result := make([]any, 0, len(tools))
	var warnings []adapter.FeatureLossWarning
	var errs []error

	for _, t := range tools {
		sourceName := t.SourceFormat
		if sourceName == "" {
			sourceName = "canonical"
		}

		warnings = append(warnings, detectSchemaFeatureLoss(t.InputSchema, sourceName, e.adapter)...)
		warnings = append(warnings, detectSchemaFeatureLoss(t.OutputSchema, sourceName, e.adapter)...)

		// Convert
		converted, err := e.adapter.FromCanonical(t)
		if err != nil {
			errs = append(errs, &ConversionError{
				ToolID: t.ID(),
				Cause:  err,
			})
			continue
		}
		result = append(result, converted)
	}
	return result, warnings, errs
}

// ConversionError represents a tool that failed to convert.
type ConversionError struct {
	ToolID string
	Cause  error
}

func (e *ConversionError) Error() string {
	return "failed to convert tool " + e.ToolID + ": " + e.Cause.Error()
}

func (e *ConversionError) Unwrap() error {
	return e.Cause
}

// detectSchemaFeatureLoss checks which features in a schema are not supported.
func detectSchemaFeatureLoss(schema *adapter.JSONSchema, sourceName string, target adapter.Adapter) []adapter.FeatureLossWarning {
	if schema == nil {
		return nil
	}

	featureUsage := map[adapter.SchemaFeature]bool{
		adapter.FeatureRef:                  schema.Ref != "",
		adapter.FeatureDefs:                 len(schema.Defs) > 0,
		adapter.FeatureAnyOf:                len(schema.AnyOf) > 0,
		adapter.FeatureOneOf:                len(schema.OneOf) > 0,
		adapter.FeatureAllOf:                len(schema.AllOf) > 0,
		adapter.FeatureNot:                  schema.Not != nil,
		adapter.FeaturePattern:              schema.Pattern != "",
		adapter.FeatureFormat:               schema.Format != "",
		adapter.FeatureAdditionalProperties: schema.AdditionalProperties != nil,
		adapter.FeatureMinimum:              schema.Minimum != nil,
		adapter.FeatureMaximum:              schema.Maximum != nil,
		adapter.FeatureMinLength:            schema.MinLength != nil,
		adapter.FeatureMaxLength:            schema.MaxLength != nil,
		adapter.FeatureEnum:                 len(schema.Enum) > 0,
		adapter.FeatureConst:                schema.Const != nil,
		adapter.FeatureDefault:              schema.Default != nil,
	}

	var warnings []adapter.FeatureLossWarning
	for feature, used := range featureUsage {
		if used && !target.SupportsFeature(feature) {
			warnings = append(warnings, adapter.FeatureLossWarning{
				Feature:     feature,
				FromAdapter: sourceName,
				ToAdapter:   target.Name(),
			})
		}
	}

	if schema.Properties != nil {
		for _, prop := range schema.Properties {
			warnings = append(warnings, detectSchemaFeatureLoss(prop, sourceName, target)...)
		}
	}
	if schema.Items != nil {
		warnings = append(warnings, detectSchemaFeatureLoss(schema.Items, sourceName, target)...)
	}
	if schema.Defs != nil {
		for _, def := range schema.Defs {
			warnings = append(warnings, detectSchemaFeatureLoss(def, sourceName, target)...)
		}
	}
	for _, s := range schema.AnyOf {
		warnings = append(warnings, detectSchemaFeatureLoss(s, sourceName, target)...)
	}
	for _, s := range schema.OneOf {
		warnings = append(warnings, detectSchemaFeatureLoss(s, sourceName, target)...)
	}
	for _, s := range schema.AllOf {
		warnings = append(warnings, detectSchemaFeatureLoss(s, sourceName, target)...)
	}
	if schema.Not != nil {
		warnings = append(warnings, detectSchemaFeatureLoss(schema.Not, sourceName, target)...)
	}

	return warnings
}
