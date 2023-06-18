// Copyright 2022 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

// Code generated by eval_gen.go. DO NOT EDIT.
// Regenerate this file with either of the following commands:
//
//	./dev generate go
//	go generate ./pkg/sql/sem/tree
//
// If you use the dev command and you have added a new tree expression, like
// tree.XYZ in a new file, you may get the confusing error: undefined: XYZ.
// Run './dev generate bazel' to fix this.
package tree

import "context"

// ExprEvaluator is used to evaluate TypedExpr expressions.
type ExprEvaluator interface {
	EvalAllColumnsSelector(context.Context, *AllColumnsSelector) (Datum, error)
	EvalAndExpr(context.Context, *AndExpr) (Datum, error)
	EvalArray(context.Context, *Array) (Datum, error)
	EvalArrayFlatten(context.Context, *ArrayFlatten) (Datum, error)
	EvalBinaryExpr(context.Context, *BinaryExpr) (Datum, error)
	EvalCaseExpr(context.Context, *CaseExpr) (Datum, error)
	EvalCastExpr(context.Context, *CastExpr) (Datum, error)
	EvalCoalesceExpr(context.Context, *CoalesceExpr) (Datum, error)
	EvalCollateExpr(context.Context, *CollateExpr) (Datum, error)
	EvalColumnAccessExpr(context.Context, *ColumnAccessExpr) (Datum, error)
	EvalColumnItem(context.Context, *ColumnItem) (Datum, error)
	EvalComparisonExpr(context.Context, *ComparisonExpr) (Datum, error)
	EvalDefaultVal(context.Context, *DefaultVal) (Datum, error)
	EvalFuncExpr(context.Context, *FuncExpr) (Datum, error)
	EvalIfErrExpr(context.Context, *IfErrExpr) (Datum, error)
	EvalIfExpr(context.Context, *IfExpr) (Datum, error)
	EvalIndexedVar(context.Context, *IndexedVar) (Datum, error)
	EvalIndirectionExpr(context.Context, *IndirectionExpr) (Datum, error)
	EvalIsNotNullExpr(context.Context, *IsNotNullExpr) (Datum, error)
	EvalIsNullExpr(context.Context, *IsNullExpr) (Datum, error)
	EvalIsOfTypeExpr(context.Context, *IsOfTypeExpr) (Datum, error)
	EvalNotExpr(context.Context, *NotExpr) (Datum, error)
	EvalNullIfExpr(context.Context, *NullIfExpr) (Datum, error)
	EvalOrExpr(context.Context, *OrExpr) (Datum, error)
	EvalParenExpr(context.Context, *ParenExpr) (Datum, error)
	EvalPlaceholder(context.Context, *Placeholder) (Datum, error)
	EvalRangeCond(context.Context, *RangeCond) (Datum, error)
	EvalRoutineExpr(context.Context, *RoutineExpr) (Datum, error)
	EvalSubquery(context.Context, *Subquery) (Datum, error)
	EvalTuple(context.Context, *Tuple) (Datum, error)
	EvalTupleStar(context.Context, *TupleStar) (Datum, error)
	EvalTypedDummy(context.Context, *TypedDummy) (Datum, error)
	EvalUnaryExpr(context.Context, *UnaryExpr) (Datum, error)
	EvalUnqualifiedStar(context.Context, UnqualifiedStar) (Datum, error)
	EvalUnresolvedName(context.Context, *UnresolvedName) (Datum, error)
}

// Eval is part of the TypedExpr interface.
func (node *AllColumnsSelector) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalAllColumnsSelector(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *AndExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalAndExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *Array) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalArray(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *ArrayFlatten) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalArrayFlatten(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *BinaryExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalBinaryExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *CaseExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalCaseExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *CastExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalCastExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *CoalesceExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalCoalesceExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *CollateExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalCollateExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *ColumnAccessExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalColumnAccessExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *ColumnItem) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalColumnItem(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *ComparisonExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalComparisonExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *DArray) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DBitArray) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DBool) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DBox2D) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DBytes) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DCollatedString) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DDate) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DDecimal) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DEncodedKey) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DEnum) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DFloat) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DGeography) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DGeometry) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DIPAddr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DInt) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DInterval) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DJSON) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DOid) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DOidWrapper) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DString) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DTSQuery) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DTSVector) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DTime) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DTimeTZ) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DTimestamp) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DTimestampTZ) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DTuple) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DUuid) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DVoid) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}

// Eval is part of the TypedExpr interface.
func (node *DefaultVal) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalDefaultVal(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *FuncExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalFuncExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *IfErrExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalIfErrExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *IfExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalIfExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *IndexedVar) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalIndexedVar(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *IndirectionExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalIndirectionExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *IsNotNullExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalIsNotNullExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *IsNullExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalIsNullExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *IsOfTypeExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalIsOfTypeExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *NotExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalNotExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *NullIfExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalNullIfExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *OrExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalOrExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *ParenExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalParenExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *Placeholder) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalPlaceholder(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *RangeCond) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalRangeCond(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *RoutineExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalRoutineExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *Subquery) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalSubquery(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *Tuple) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalTuple(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *TupleStar) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalTupleStar(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *TypedDummy) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalTypedDummy(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *UnaryExpr) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalUnaryExpr(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node UnqualifiedStar) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalUnqualifiedStar(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node *UnresolvedName) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return v.EvalUnresolvedName(ctx, node)
}

// Eval is part of the TypedExpr interface.
func (node dNull) Eval(ctx context.Context, v ExprEvaluator) (Datum, error) {
	return node, nil
}
