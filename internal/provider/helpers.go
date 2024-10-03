// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/prempador/go-defectdojo"
)

// basetypesInt64ValueToDefectdojoNullableInt32 converts a basetypes.Int64Value to a defectdojo.NullableInt32.
func basetypesInt64ValueToDefectdojoNullableInt32(value basetypes.Int64Value) defectdojo.NullableInt32 {
	if value.IsNull() || value.IsUnknown() {
		return *defectdojo.NewNullableInt32(nil)
	}

	v := int32(value.ValueInt64())
	return *defectdojo.NewNullableInt32(&v)
}

// basetypesInt64ValueToInt32Pointer converts a basetypes.Int64Value to a *int32.
func basetypesInt64ValueToInt32Pointer(value basetypes.Int64Value) *int32 {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	v := int32(value.ValueInt64())
	return &v
}

// int32PointerToBasetypesInt64Value converts a *int32 to a basetypes.Int64Value.
func int32PointerToBasetypesInt64Value(value *int32) basetypes.Int64Value {
	if value == nil {
		return basetypes.NewInt64Null()
	}

	return basetypes.NewInt64Value(int64(*value))
}

// basetypesStringValueToDefectdojoNullableString converts a basetypes.StringValue to a defectdojo.NullableString.
// we need to convert some fields like this because if they are unknown,
// defectdojo is treating an empty string as a string that needs to be validated.
func basetypesStringValueToDefectdojoNullableString(value basetypes.StringValue) defectdojo.NullableString {
	if value.IsNull() || value.IsUnknown() {
		return *defectdojo.NewNullableString(nil)
	}

	v := value.ValueString()
	return *defectdojo.NewNullableString(&v)
}
