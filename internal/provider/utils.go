package provider

import (
	"context"
	"fmt"
	"reflect"
)

func MergeMaps(ctx context.Context, dst, src reflect.Value) error {

	if dst.Kind() != reflect.Map || src.Kind() != reflect.Map {
		return fmt.Errorf("[ERROR] src and/or dst in MergeMaps not a Map.")
	}
	// iterate over source map keys and values
	iter := src.MapRange()
	for iter.Next() {
		sKey := iter.Key()
		if sKey.Kind() == reflect.Interface {
			sKey = sKey.Elem()
		}
		sValue := iter.Value()
		if sValue.Kind() == reflect.Interface {
			sValue = sValue.Elem()
		}

		if sValue.Kind() == reflect.Map {
			dValue := dst.MapIndex(sKey)
			// add empty map to dst if key does not exist
			if !dValue.IsValid() {
				dst.SetMapIndex(sKey, reflect.MakeMap(sValue.Type()))
			}
			// merge src map into dst map
			dValue = reflect.ValueOf(dst.MapIndex(sKey).Interface())
			MergeMaps(ctx, dValue, sValue)
		} else if sValue.Kind() == reflect.Slice {
			dValue := dst.MapIndex(sKey)
			// if slice does not exist in dst, add empty slice
			if !dValue.IsValid() {
				dst.SetMapIndex(sKey, reflect.MakeSlice(sValue.Type(), 0, 0))
			}
			dValue = reflect.ValueOf(dst.MapIndex(sKey).Interface())
			if dValue.Kind() == reflect.Slice {
				// interate over source slice elements and add merge with dst list
				for i := 0; i < sValue.Len(); i++ {
					MergeListItem(ctx, dst, sKey, sValue.Index(i))
				}
			}
		} else {
			// else we have primitive type -> add/replace dst value
			dst.SetMapIndex(sKey, sValue)
		}
	}
	return nil
}

func MergeListItem(ctx context.Context, dst, key, src reflect.Value) {
	dValue := reflect.ValueOf(dst.MapIndex(key).Interface())
	if src.Kind() == reflect.Interface {
		src = src.Elem()
	}
	if src.Kind() == reflect.Map {
		for i := 0; i < dValue.Len(); i++ {
			match := true
			comparison := false
			// iterate over all source map keys and values
			iter := src.MapRange()
			for iter.Next() {
				sKey := iter.Key()
				if sKey.Kind() == reflect.Interface {
					sKey = sKey.Elem()
				}
				sValue := iter.Value()
				if sValue.Kind() == reflect.Interface {
					sValue = sValue.Elem()
				}

				if sValue.Kind() == reflect.Map || sValue.Kind() == reflect.Slice {
					// we only compare primitive types
					continue
				} else {
					x := dValue.Index(i).Elem().MapIndex(sKey)
					if x.Kind() == reflect.Interface {
						x = x.Elem()
					}
					// check if element exists in dst map and value is the same as in src map
					if x.IsValid() && sValue.Interface() == x.Interface() {
						comparison = true
						continue
					}
					comparison = true
					match = false
				}
			}
			// Check if all primitive values have matched AND at least one comparison was done
			if match && comparison {
				dv := reflect.ValueOf(dst.MapIndex(key).Elem().Index(i).Interface())
				MergeMaps(ctx, dv, src)
				return
			}
		}

	} else {
		// check if primitive value exists in dst and if so return
		slice := dst.MapIndex(key).Elem()
		for i := 0; i < slice.Len(); i++ {
			if slice.Index(i).Elem().String() == src.String() {
				return
			}
		}
	}
	dst.SetMapIndex(key, reflect.Append(dValue, src))
}
