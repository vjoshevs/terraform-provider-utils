package provider

import (
	"fmt"
	"reflect"
)

func MergeMaps(dst, src reflect.Value, mergeListItems bool) error {

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
			MergeMaps(dValue, sValue, mergeListItems)
		} else if sValue.Kind() == reflect.Slice {
			dValue := dst.MapIndex(sKey)
			// if slice does not exist in dst, add empty slice
			if !dValue.IsValid() {
				dst.SetMapIndex(sKey, reflect.MakeSlice(sValue.Type(), 0, 0))
			}
			dValue = reflect.ValueOf(dst.MapIndex(sKey).Interface())
			if dValue.Kind() == reflect.Slice {
				// iterate over source slice elements and add merge with dst list
				for i := 0; i < sValue.Len(); i++ {
					MergeListItem(dst, sKey, sValue.Index(i), mergeListItems)
				}
			}
		} else {
			// else we have primitive type -> add/replace dst value
			dst.SetMapIndex(sKey, sValue)
		}
	}
	return nil
}

func MergeListItem(dst, key, src reflect.Value, mergeListItems bool) {
	dValue := reflect.ValueOf(dst.MapIndex(key).Interface())
	if src.Kind() == reflect.Interface {
		src = src.Elem()
	}
	if src.Kind() == reflect.Map && mergeListItems {
		for i := 0; i < dValue.Len(); i++ {
			match := true
			comparison := false
			uniqueSource := false
			uniqueDest := false
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
				}
				x := dValue.Index(i).Elem().MapIndex(sKey)
				if x.Kind() == reflect.Interface {
					x = x.Elem()
				}
				// check if element exists in dst map and value is the same as in src map
				if x.IsValid() && sValue.IsValid() && sValue.Interface() == x.Interface() {
					comparison = true
					continue
				}
				// if value does not exist in dst map -> continue
				if !x.IsValid() && sValue.IsValid() {
					uniqueSource = true
					continue
				}
				comparison = true
				match = false
			}
			// iterate over all dst map keys and values
			iter = dValue.Index(i).Elem().MapRange()
			for iter.Next() {
				dKey := iter.Key()
				if dKey.Kind() == reflect.Interface {
					dKey = dKey.Elem()
				}
				dValue := iter.Value()
				if dValue.Kind() == reflect.Interface {
					dValue = dValue.Elem()
				}

				if dValue.Kind() == reflect.Map || dValue.Kind() == reflect.Slice {
					// we only compare primitive types
					continue
				}
				x := src.MapIndex(dKey)
				if x.Kind() == reflect.Interface {
					x = x.Elem()
				}
				// check if element exists in src map and value is the same as in dst map
				if x.IsValid() && dValue.IsValid() && dValue.Interface() == x.Interface() {
					comparison = true
					continue
				}
				// if value does not exist in src map -> continue
				if !x.IsValid() && dValue.IsValid() {
					uniqueDest = true
					continue
				}
				comparison = true
				match = false
			}
			// Check if all primitive values have matched AND at least one comparison was done
			if match && comparison && !(uniqueSource && uniqueDest) {
				dv := reflect.ValueOf(dst.MapIndex(key).Elem().Index(i).Interface())
				MergeMaps(dv, src, mergeListItems)
				return
			}
		}

	} else {
		// check if primitive value exists in dst and if so return
		slice := dst.MapIndex(key).Elem()
		for i := 0; i < slice.Len(); i++ {
			element := slice.Index(i).Elem()
			if element.IsValid() && src.IsValid() && element.Kind() != reflect.Map && element.Kind() != reflect.Slice && element.Interface() == src.Interface() {
				return
			}
		}
	}
	dst.SetMapIndex(key, reflect.Append(dValue, src))
}
