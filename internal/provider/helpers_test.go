// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/prempador/go-defectdojo"
	"github.com/stretchr/testify/require"
)

func TestUnitBasetypesInt64ValueToDefectdojoNullableInt32ValueNull(t *testing.T) {
	value := basetypes.NewInt64Null()
	dNullInt := defectdojo.NewNullableInt32(nil)

	result := basetypesInt64ValueToDefectdojoNullableInt32(value)

	require.Equal(t, *dNullInt, result)
}

func TestUnitBasetypesInt64ValueToDefectdojoNullableInt32ValueUnknown(t *testing.T) {
	value := basetypes.NewInt64Unknown()
	dNullInt := defectdojo.NewNullableInt32(nil)

	result := basetypesInt64ValueToDefectdojoNullableInt32(value)

	require.Equal(t, *dNullInt, result)
}

func TestUnitBasetypesInt64ValueToDefectdojoNullableInt32Value(t *testing.T) {
	i64 := int64(1)
	i32 := int32(i64)

	value := basetypes.NewInt64Value(i64)
	dNullInt := defectdojo.NewNullableInt32(&i32)

	result := basetypesInt64ValueToDefectdojoNullableInt32(value)

	require.Equal(t, *dNullInt, result)
}

func TestUnitBasetypesInt64ValueToInt32PointerNull(t *testing.T) {
	value := basetypes.NewInt64Null()
	var i32 *int32

	result := basetypesInt64ValueToInt32Pointer(value)

	require.Equal(t, i32, result)
}

func TestUnitBasetypesInt64ValueToInt32PointerUnknown(t *testing.T) {
	value := basetypes.NewInt64Unknown()
	var i32 *int32

	result := basetypesInt64ValueToInt32Pointer(value)

	require.Equal(t, i32, result)
}

func TestUnitBasetypesInt64ValueToInt32PointerValue(t *testing.T) {
	var i32 *int32
	i := int32(1)
	i32 = &i
	value := basetypes.NewInt64Value(int64(i))

	result := basetypesInt64ValueToInt32Pointer(value)

	require.Equal(t, i32, result)
}

func TestUnitBasetypesStringValueToDefectdojoNullableStringNull(t *testing.T) {
	value := basetypes.NewStringNull()
	dNullString := defectdojo.NewNullableString(nil)

	result := basetypesStringValueToDefectdojoNullableString(value)

	require.Equal(t, *dNullString, result)
}

func TestUnitBasetypesStringValueToDefectdojoNullableStringUnknown(t *testing.T) {
	value := basetypes.NewStringUnknown()
	dNullString := defectdojo.NewNullableString(nil)

	result := basetypesStringValueToDefectdojoNullableString(value)

	require.Equal(t, *dNullString, result)
}

func TestUnitBasetypesStringValueToDefectdojoNullableStringValue(t *testing.T) {
	s := "asdf"
	value := basetypes.NewStringValue(s)
	dNullString := defectdojo.NewNullableString(&s)

	result := basetypesStringValueToDefectdojoNullableString(value)

	require.Equal(t, *dNullString, result)
}
