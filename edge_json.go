package stigg

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/Yamashou/gqlgenc/graphqljson"
)

// unmarshalEdgeData decodes an Edge API payload into a generated result type.
//
// The Edge API returns its own canonical shape, which carries keys the gqlgenc selection sets never
// requested — notably __typename on objects whose generated structs have no field for it. Stock
// graphqljson rejects any key it cannot place, so unplaceable keys are pruned first. Decoding itself
// stays with graphqljson because it is what resolves `graphql:"... on X"` union members; plain
// encoding/json would leave every union member zero-valued.
func unmarshalEdgeData(data json.RawMessage, v interface{}) error {
	var tree interface{}
	if err := json.Unmarshal(data, &tree); err != nil {
		return fmt.Errorf("failed to decode data: %w", err)
	}

	pruned, err := json.Marshal(prune(tree, reflect.TypeOf(v)))
	if err != nil {
		return fmt.Errorf("failed to re-encode pruned data: %w", err)
	}

	return graphqljson.UnmarshalData(pruned, v)
}

// prune drops the keys graphqljson would refuse to place, walking the decoded JSON alongside the
// target type. A key is kept whenever any candidate struct at this position could accept it.
func prune(node interface{}, target reflect.Type) interface{} {
	switch value := node.(type) {
	case map[string]interface{}:
		fields := fieldsFor(target)
		if fields == nil {
			// An opaque target (map, interface{}, json.RawMessage) takes the subtree verbatim, so
			// nothing here is prunable and everything must be preserved.
			return value
		}

		for key, child := range value {
			fieldType, ok := fields.lookup(key)
			if !ok {
				delete(value, key)
				continue
			}
			value[key] = prune(child, fieldType)
		}

		return value
	case []interface{}:
		element := elementType(target)
		for i, child := range value {
			value[i] = prune(child, element)
		}

		return value
	default:
		return node
	}
}

// edgeFields is the set of keys graphqljson would accept for one struct position, flattened across
// the struct's inline-fragment members. Matching mirrors graphqljson's hasGraphQLName: a `graphql`
// tag matches case-sensitively, an untagged field matches its own name case-insensitively.
type edgeFields struct {
	tagged   map[string]reflect.Type
	untagged map[string]reflect.Type
}

func (f *edgeFields) lookup(key string) (reflect.Type, bool) {
	if fieldType, ok := f.tagged[key]; ok {
		return fieldType, true
	}

	fieldType, ok := f.untagged[strings.ToLower(key)]

	return fieldType, ok
}

// Generated result types are fixed at build time, so the accepted-key set per type never changes.
// Without this cache the set is rebuilt for every object in every response.
var edgeFieldCache sync.Map // reflect.Type -> *edgeFields

func fieldsFor(target reflect.Type) *edgeFields {
	target = structType(target)
	if target == nil || target.Kind() != reflect.Struct {
		return nil
	}

	if cached, ok := edgeFieldCache.Load(target); ok {
		return cached.(*edgeFields)
	}

	fields := &edgeFields{tagged: map[string]reflect.Type{}, untagged: map[string]reflect.Type{}}
	collectFields(target, fields, map[reflect.Type]bool{})
	edgeFieldCache.Store(target, fields)

	return fields
}

func collectFields(target reflect.Type, into *edgeFields, seen map[reflect.Type]bool) {
	target = structType(target)
	if target == nil || target.Kind() != reflect.Struct || seen[target] {
		return
	}
	seen[target] = true

	for i := 0; i < target.NumField(); i++ {
		field := target.Field(i)
		if field.PkgPath != "" {
			continue
		}

		tag, tagged := field.Tag.Lookup("graphql")
		if !tagged {
			key := strings.ToLower(field.Name)
			if _, exists := into.untagged[key]; !exists {
				into.untagged[key] = field.Type
			}
			continue
		}

		// An inline-fragment member has no name of its own; its fields belong to this position.
		if strings.HasPrefix(strings.TrimSpace(tag), "...") {
			collectFields(field.Type, into, seen)
			continue
		}

		name := graphQLName(tag)
		if _, exists := into.tagged[name]; !exists {
			into.tagged[name] = field.Type
		}
	}
}

// graphQLName extracts the field name from a `graphql` tag the way graphqljson does, dropping any
// argument list or alias suffix.
func graphQLName(tag string) string {
	tag = strings.TrimSpace(tag)
	if i := strings.Index(tag, "("); i != -1 {
		tag = tag[:i]
	}
	if i := strings.Index(tag, ":"); i != -1 {
		tag = tag[:i]
	}

	return strings.TrimSpace(tag)
}

// structType strips pointers to reach the type an object decodes into.
func structType(target reflect.Type) reflect.Type {
	for target != nil && target.Kind() == reflect.Ptr {
		target = target.Elem()
	}

	return target
}

// elementType returns the type one JSON array level down.
func elementType(target reflect.Type) reflect.Type {
	target = structType(target)
	if target != nil && (target.Kind() == reflect.Slice || target.Kind() == reflect.Array) {
		return target.Elem()
	}

	return target
}
