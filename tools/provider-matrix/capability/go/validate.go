package capability

import (
	"fmt"
	"reflect"
)

// Validate checks the fixed result vocabulary and source membership, not the
// truth of a format owner's interpretation. Association stays with the host.
func Validate(request DefinitionRequest, inventory Inventory) error {
	limits, err := EffectiveLimits(request.Limits)
	if err != nil {
		return err
	}
	if request.Source == nil || request.Root.Package == "" || request.Root.Name == "" {
		return fmt.Errorf("missing extraction source or root")
	}
	root, err := request.Source.Declaration(request.Root)
	if err != nil || root.Symbol != request.Root {
		return fmt.Errorf("extraction root is not in source capture: %v", err)
	}
	v := validator{source: request.Source, limits: limits}
	if err := v.span(root.Span); err != nil {
		return err
	}
	if inventory.Extent != Decoded && inventory.Extent != PartlyDecoded {
		return fmt.Errorf("missing or invalid extraction extent")
	}
	keys := map[string]bool{}
	for _, f := range inventory.Features {
		if f.Key == "" || keys[f.Key] {
			return fmt.Errorf("empty or duplicate feature key %q", f.Key)
		}
		keys[f.Key] = true
		if err := v.record(f.At); err != nil {
			return err
		}
		if err := field(&v, f.Roles, true, false); err != nil {
			return err
		}
		roles := map[string]bool{}
		for _, r := range f.Roles.Value {
			if r.Name == "" || roles[r.Name] {
				return fmt.Errorf("empty or duplicate role %q in %q", r.Name, f.Key)
			}
			roles[r.Name] = true
			if err := v.record(r.At); err != nil {
				return err
			}
			if err := field(&v, r.Min, false, false); err != nil {
				return err
			}
			if err := field(&v, r.Max, false, false); err != nil {
				return err
			}
			if r.Min.State == Known && r.Min.Value < 0 {
				return fmt.Errorf("negative minimum for %q", r.Name)
			}
			if r.Max.State == Known {
				if r.Max.Value.Unbounded && r.Max.Value.Value != 0 || !r.Max.Value.Unbounded && r.Max.Value.Value < 0 {
					return fmt.Errorf("invalid maximum for %q", r.Name)
				}
				if r.Min.State == Known && !r.Max.Value.Unbounded && r.Max.Value.Value < r.Min.Value {
					return fmt.Errorf("maximum below minimum for %q", r.Name)
				}
			}
			if err := field(&v, r.Kinds, true, false); err != nil {
				return err
			}
			if err := v.strings(r.Kinds.Value); err != nil {
				return err
			}
		}
		if err := field(&v, f.Attributes, true, false); err != nil {
			return err
		}
		attributes := map[string]bool{}
		for _, a := range f.Attributes.Value {
			if a.Name == "" || attributes[a.Name] {
				return fmt.Errorf("empty or duplicate attribute %q in %q", a.Name, f.Key)
			}
			attributes[a.Name] = true
			if err := v.record(a.At); err != nil {
				return err
			}
			if err := field(&v, a.Required, false, false); err != nil {
				return err
			}
			if err := field(&v, a.Values, true, true); err != nil {
				return err
			}
			if a.Values.Value.Unrestricted && (a.Values.State != Known || len(a.Values.Value.Values) > 0) {
				return fmt.Errorf("invalid unrestricted attribute %q", a.Name)
			}
			if err := v.strings(a.Values.Value.Values); err != nil {
				return err
			}
		}
		if err := field(&v, f.References, true, false); err != nil {
			return err
		}
		if err := v.strings(f.References.Value); err != nil {
			return err
		}
	}
	for _, limitation := range inventory.Limitations {
		v.uncertain = true
		if err := v.record(limitation.At); err != nil {
			return err
		}
		if limitation.Message == "" || len(limitation.Message) > limits.MaxDiagnosticBytes {
			return fmt.Errorf("missing or oversized limitation")
		}
	}
	if inventory.Extent == Decoded && v.uncertain {
		return fmt.Errorf("decoded inventory contains uncertain fields")
	}
	if inventory.Extent == PartlyDecoded && !v.uncertain {
		return fmt.Errorf("partial inventory has no localized limitation")
	}
	return nil
}

type validator struct {
	source    SourceView
	limits    Limits
	records   int
	uncertain bool
}

func (v *validator) span(span Span) error { return ValidateSpan(v.source, span) }

func (v *validator) record(span Span) error {
	v.records++
	if v.records > v.limits.MaxRecords {
		return fmt.Errorf("extraction record limit exceeded")
	}
	return v.span(span)
}

func (v *validator) strings(values []string) error {
	seen := map[string]bool{}
	for _, value := range values {
		v.records++
		if v.records > v.limits.MaxRecords {
			return fmt.Errorf("extraction record limit exceeded")
		}
		if seen[value] {
			return fmt.Errorf("duplicate structural value %q", value)
		}
		seen[value] = true
	}
	return nil
}

func field[T any](v *validator, f Field[T], partial, opaque bool) error {
	if err := v.record(f.At); err != nil {
		return err
	}
	switch f.State {
	case Known:
		if f.Reason != "" {
			return fmt.Errorf("known field carries an uncertainty reason")
		}
	case Partial:
		if !partial {
			return fmt.Errorf("partial scalar field")
		}
	case Opaque:
		if !opaque {
			return fmt.Errorf("opaque non-predicate field")
		}
		fallthrough
	case Unknown:
		if !reflect.ValueOf(f.Value).IsZero() {
			return fmt.Errorf("unknown or opaque field carries a value")
		}
	default:
		return fmt.Errorf("missing or invalid field knowledge")
	}
	if f.State != Known {
		v.uncertain = true
		if f.Reason == "" || len(f.Reason) > v.limits.MaxDiagnosticBytes {
			return fmt.Errorf("missing or oversized field reason")
		}
	}
	return nil
}
