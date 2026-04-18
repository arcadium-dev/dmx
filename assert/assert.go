//  Copyright 2026 arcadium.dev <info@arcadium.dev>
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package assert // import "arcadium.dev/dmx/assert"

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Contains(t *testing.T, actual, expected string) {
	t.Helper()
	if !strings.Contains(actual, expected) {
		t.Errorf("\nExpected: %s\nActual:   %s", expected, actual)
	}
}

func NotContains(t *testing.T, actual, expected string) {
	t.Helper()
	if strings.Contains(actual, expected) {
		t.Errorf("Unexpected value: %+v", actual)
	}
}

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("\nExpected: %+v\nActual:   %+v", expected, actual)
	}
}

func NotEqual[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual == expected {
		t.Errorf("Expected different values: %+v", actual)
	}
}

func Compare[T any](t *testing.T, actual, expected T, opts ...cmp.Option) {
	t.Helper()
	if !cmp.Equal(actual, expected, opts...) {
		t.Errorf("\nExpected: %+v\nActual:   %+v", expected, actual)
	}
}

func Error(t *testing.T, err error, expected string) {
	t.Helper()
	if err == nil {
		t.Fatal("Expected an error")
	}
	if expected != err.Error() {
		t.Errorf("\nExpected error: %s\nActual error:   %s", expected, err)
	}
}

func ErrorIs(t *testing.T, actual, expected error) {
	t.Helper()
	if actual == nil {
		t.Fatal("Expected an error")
	}
	if !errors.Is(actual, expected) {
		t.Errorf("\nExpected error: %s\nActual error:   %s", expected, actual)
	}
}

func Nil(t *testing.T, object any) {
	t.Helper()
	if object == nil {
		return
	}

	value := reflect.ValueOf(object)
	switch value.Kind() {
	case
		reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Map,
		reflect.Pointer, reflect.Slice, reflect.UnsafePointer:

		if value.IsNil() {
			return
		}
	}

	t.Errorf("Unexpected non-nil value: %+v", object)
}

func NotNil(t *testing.T, object any) {
	t.Helper()
	if object != nil {
		return
	}

	value := reflect.ValueOf(object)
	switch value.Kind() {
	case
		reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Map,
		reflect.Pointer, reflect.Slice, reflect.UnsafePointer:

		if !value.IsNil() {
			return
		}
	}

	t.Errorf("Unexpected nil value: %+v", object)
}

func True(t *testing.T, v bool) {
	t.Helper()
	if !v {
		t.Error("Expected value to be true")
	}
}

func False(t *testing.T, v bool) {
	t.Helper()
	if v {
		t.Error("Expected value to be false")
	}
}

func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error("Expected value to be nil")
	}
}

func ErrorContains(t *testing.T, err error, value string) {
	t.Helper()
	if err != nil && !strings.Contains(err.Error(), value) {
		t.Errorf("Error message does not contain expected value")
	}
}
