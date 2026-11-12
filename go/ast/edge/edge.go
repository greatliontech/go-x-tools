// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package edge defines identifiers for each field of an ast.Node
// struct type that refers to another Node.
package edge

import (
	"fmt"
	"go/ast"
	"reflect"
)

// A Kind describes a field of an [ast.Node] struct.
type Kind uint8

// String returns a description of the edge kind.
func (k Kind) String() string {
	if k == Invalid {
		return "<invalid>"
	}
	info := k.info()
	return fmt.Sprintf("%v.%s", info.nodeType.Elem().Name(), info.name)
}

// NodeType returns the pointer-to-struct type of the ast.Node implementation.
func (k Kind) NodeType() reflect.Type { return k.info().nodeType }

// FieldName returns the name of the field.
func (k Kind) FieldName() string { return k.info().name }

// FieldType returns the declared type of the field.
func (k Kind) FieldType() reflect.Type { return k.info().fieldType }

// Get returns the direct child of n identified by (k, idx).
// n's type must match k.NodeType().
// idx must be a valid slice index, or -1 for a non-slice.
func (k Kind) Get(n ast.Node, idx int) ast.Node {
	info := k.info()
	if info.nodeType != reflect.TypeOf(n) {
		panic(fmt.Sprintf("%v.Get(%T): invalid node type", k, n))
	}
	v := reflect.ValueOf(n).Elem().Field(info.index)

	if v.Kind() == reflect.Slice {
		v = v.Index(idx) // asserts valid idx
	} else if idx != -1 {
		panic(fmt.Sprintf("%v, Get(%T, %d): cannot index non-slice", v, n, idx))
	}

	out, _ := v.Interface().(ast.Node) // may be nil
	return out
}

// Each [Kind] is named Type_Field, where Type is the
// [ast.Node] struct type and Field is the name of the field
const (
	Invalid Kind = iota // for nodes at the root of the traversal

	// As of Go1.26 these kinds are sorted alphabetically, but
	// numbering must be stable, so any new addition of const should
	// use a new value (be added at the end of the list).

	ArrayType_Elt
	ArrayType_Len
	AssignStmt_Lhs
	AssignStmt_Rhs
	BinaryExpr_X
	BinaryExpr_Y
	BlockStmt_List
	BranchStmt_Label
	CallExpr_Args
	CallExpr_Fun
	CaseClause_Body
	CaseClause_List
	ChanType_Value
	CommClause_Body
	CommClause_Comm
	CommentGroup_List
	CompositeLit_Elts
	CompositeLit_Type
	DeclStmt_Decl
	DeferStmt_Call
	Ellipsis_Elt
	ExprStmt_X
	FieldList_List
	Field_Comment
	Field_Doc
	Field_Names
	Field_Tag
	Field_Type
	File_Decls
	File_Doc
	File_Name
	ForStmt_Body
	ForStmt_Cond
	ForStmt_Init
	ForStmt_Post
	FuncDecl_Body
	FuncDecl_Doc
	FuncDecl_Name
	FuncDecl_Recv
	FuncDecl_Type
	FuncLit_Body
	FuncLit_Type
	FuncType_Params
	FuncType_Results
	FuncType_TypeParams
	GenDecl_Doc
	GenDecl_Specs
	GoStmt_Call
	IfStmt_Body
	IfStmt_Cond
	IfStmt_Else
	IfStmt_Init
	ImportSpec_Comment
	ImportSpec_Doc
	ImportSpec_Name
	ImportSpec_Path
	IncDecStmt_X
	IndexExpr_Index
	IndexExpr_X
	IndexListExpr_Indices
	IndexListExpr_X
	InterfaceType_Methods
	KeyValueExpr_Key
	KeyValueExpr_Value
	LabeledStmt_Label
	LabeledStmt_Stmt
	MapType_Key
	MapType_Value
	ParenExpr_X
	RangeStmt_Body
	RangeStmt_Key
	RangeStmt_Value
	RangeStmt_X
	ReturnStmt_Results
	SelectStmt_Body
	SelectorExpr_Sel
	SelectorExpr_X
	SendStmt_Chan
	SendStmt_Value
	SliceExpr_High
	SliceExpr_Low
	SliceExpr_Max
	SliceExpr_X
	StarExpr_X
	StructType_Fields
	SwitchStmt_Body
	SwitchStmt_Init
	SwitchStmt_Tag
	TypeAssertExpr_Type
	TypeAssertExpr_X
	TypeSpec_Comment
	TypeSpec_Doc
	TypeSpec_Name
	TypeSpec_Type
	TypeSpec_TypeParams
	TypeSwitchStmt_Assign
	TypeSwitchStmt_Body
	TypeSwitchStmt_Init
	UnaryExpr_X
	ValueSpec_Comment
	ValueSpec_Doc
	ValueSpec_Names
	ValueSpec_Type
	ValueSpec_Values

	maxKind
)

// Assert that the encoding fits in 7 bits,
// as the inspector relies on this.
// (We are currently at 104.)
var _ = [1 << 7]struct{}{}[maxKind]

type fieldInfo struct {
	nodeType  reflect.Type // pointer-to-struct type of ast.Node implementation
	name      string
	index     int
	fieldType reflect.Type
}

func info[N ast.Node](fieldName string) fieldInfo {
	var zero N
	nodePtrType := reflect.TypeOf(zero)
	f, ok := nodePtrType.Elem().FieldByName(fieldName)
	if !ok {
		panic(fieldName)
	}
	return fieldInfo{nodePtrType, fieldName, f.Index[0], f.Type}
}

func (k Kind) info() fieldInfo {
	switch k {
	case Invalid:
		return fieldInfo{}
	case ArrayType_Elt:
		return info[*ast.ArrayType]("Elt")
	case ArrayType_Len:
		return info[*ast.ArrayType]("Len")
	case AssignStmt_Lhs:
		return info[*ast.AssignStmt]("Lhs")
	case AssignStmt_Rhs:
		return info[*ast.AssignStmt]("Rhs")
	case BinaryExpr_X:
		return info[*ast.BinaryExpr]("X")
	case BinaryExpr_Y:
		return info[*ast.BinaryExpr]("Y")
	case BlockStmt_List:
		return info[*ast.BlockStmt]("List")
	case BranchStmt_Label:
		return info[*ast.BranchStmt]("Label")
	case CallExpr_Args:
		return info[*ast.CallExpr]("Args")
	case CallExpr_Fun:
		return info[*ast.CallExpr]("Fun")
	case CaseClause_Body:
		return info[*ast.CaseClause]("Body")
	case CaseClause_List:
		return info[*ast.CaseClause]("List")
	case ChanType_Value:
		return info[*ast.ChanType]("Value")
	case CommClause_Body:
		return info[*ast.CommClause]("Body")
	case CommClause_Comm:
		return info[*ast.CommClause]("Comm")
	case CommentGroup_List:
		return info[*ast.CommentGroup]("List")
	case CompositeLit_Elts:
		return info[*ast.CompositeLit]("Elts")
	case CompositeLit_Type:
		return info[*ast.CompositeLit]("Type")
	case DeclStmt_Decl:
		return info[*ast.DeclStmt]("Decl")
	case DeferStmt_Call:
		return info[*ast.DeferStmt]("Call")
	case Ellipsis_Elt:
		return info[*ast.Ellipsis]("Elt")
	case ExprStmt_X:
		return info[*ast.ExprStmt]("X")
	case FieldList_List:
		return info[*ast.FieldList]("List")
	case Field_Comment:
		return info[*ast.Field]("Comment")
	case Field_Doc:
		return info[*ast.Field]("Doc")
	case Field_Names:
		return info[*ast.Field]("Names")
	case Field_Tag:
		return info[*ast.Field]("Tag")
	case Field_Type:
		return info[*ast.Field]("Type")
	case File_Decls:
		return info[*ast.File]("Decls")
	case File_Doc:
		return info[*ast.File]("Doc")
	case File_Name:
		return info[*ast.File]("Name")
	case ForStmt_Body:
		return info[*ast.ForStmt]("Body")
	case ForStmt_Cond:
		return info[*ast.ForStmt]("Cond")
	case ForStmt_Init:
		return info[*ast.ForStmt]("Init")
	case ForStmt_Post:
		return info[*ast.ForStmt]("Post")
	case FuncDecl_Body:
		return info[*ast.FuncDecl]("Body")
	case FuncDecl_Doc:
		return info[*ast.FuncDecl]("Doc")
	case FuncDecl_Name:
		return info[*ast.FuncDecl]("Name")
	case FuncDecl_Recv:
		return info[*ast.FuncDecl]("Recv")
	case FuncDecl_Type:
		return info[*ast.FuncDecl]("Type")
	case FuncLit_Body:
		return info[*ast.FuncLit]("Body")
	case FuncLit_Type:
		return info[*ast.FuncLit]("Type")
	case FuncType_Params:
		return info[*ast.FuncType]("Params")
	case FuncType_Results:
		return info[*ast.FuncType]("Results")
	case FuncType_TypeParams:
		return info[*ast.FuncType]("TypeParams")
	case GenDecl_Doc:
		return info[*ast.GenDecl]("Doc")
	case GenDecl_Specs:
		return info[*ast.GenDecl]("Specs")
	case GoStmt_Call:
		return info[*ast.GoStmt]("Call")
	case IfStmt_Body:
		return info[*ast.IfStmt]("Body")
	case IfStmt_Cond:
		return info[*ast.IfStmt]("Cond")
	case IfStmt_Else:
		return info[*ast.IfStmt]("Else")
	case IfStmt_Init:
		return info[*ast.IfStmt]("Init")
	case ImportSpec_Comment:
		return info[*ast.ImportSpec]("Comment")
	case ImportSpec_Doc:
		return info[*ast.ImportSpec]("Doc")
	case ImportSpec_Name:
		return info[*ast.ImportSpec]("Name")
	case ImportSpec_Path:
		return info[*ast.ImportSpec]("Path")
	case IncDecStmt_X:
		return info[*ast.IncDecStmt]("X")
	case IndexExpr_Index:
		return info[*ast.IndexExpr]("Index")
	case IndexExpr_X:
		return info[*ast.IndexExpr]("X")
	case IndexListExpr_Indices:
		return info[*ast.IndexListExpr]("Indices")
	case IndexListExpr_X:
		return info[*ast.IndexListExpr]("X")
	case InterfaceType_Methods:
		return info[*ast.InterfaceType]("Methods")
	case KeyValueExpr_Key:
		return info[*ast.KeyValueExpr]("Key")
	case KeyValueExpr_Value:
		return info[*ast.KeyValueExpr]("Value")
	case LabeledStmt_Label:
		return info[*ast.LabeledStmt]("Label")
	case LabeledStmt_Stmt:
		return info[*ast.LabeledStmt]("Stmt")
	case MapType_Key:
		return info[*ast.MapType]("Key")
	case MapType_Value:
		return info[*ast.MapType]("Value")
	case ParenExpr_X:
		return info[*ast.ParenExpr]("X")
	case RangeStmt_Body:
		return info[*ast.RangeStmt]("Body")
	case RangeStmt_Key:
		return info[*ast.RangeStmt]("Key")
	case RangeStmt_Value:
		return info[*ast.RangeStmt]("Value")
	case RangeStmt_X:
		return info[*ast.RangeStmt]("X")
	case ReturnStmt_Results:
		return info[*ast.ReturnStmt]("Results")
	case SelectStmt_Body:
		return info[*ast.SelectStmt]("Body")
	case SelectorExpr_Sel:
		return info[*ast.SelectorExpr]("Sel")
	case SelectorExpr_X:
		return info[*ast.SelectorExpr]("X")
	case SendStmt_Chan:
		return info[*ast.SendStmt]("Chan")
	case SendStmt_Value:
		return info[*ast.SendStmt]("Value")
	case SliceExpr_High:
		return info[*ast.SliceExpr]("High")
	case SliceExpr_Low:
		return info[*ast.SliceExpr]("Low")
	case SliceExpr_Max:
		return info[*ast.SliceExpr]("Max")
	case SliceExpr_X:
		return info[*ast.SliceExpr]("X")
	case StarExpr_X:
		return info[*ast.StarExpr]("X")
	case StructType_Fields:
		return info[*ast.StructType]("Fields")
	case SwitchStmt_Body:
		return info[*ast.SwitchStmt]("Body")
	case SwitchStmt_Init:
		return info[*ast.SwitchStmt]("Init")
	case SwitchStmt_Tag:
		return info[*ast.SwitchStmt]("Tag")
	case TypeAssertExpr_Type:
		return info[*ast.TypeAssertExpr]("Type")
	case TypeAssertExpr_X:
		return info[*ast.TypeAssertExpr]("X")
	case TypeSpec_Comment:
		return info[*ast.TypeSpec]("Comment")
	case TypeSpec_Doc:
		return info[*ast.TypeSpec]("Doc")
	case TypeSpec_Name:
		return info[*ast.TypeSpec]("Name")
	case TypeSpec_Type:
		return info[*ast.TypeSpec]("Type")
	case TypeSpec_TypeParams:
		return info[*ast.TypeSpec]("TypeParams")
	case TypeSwitchStmt_Assign:
		return info[*ast.TypeSwitchStmt]("Assign")
	case TypeSwitchStmt_Body:
		return info[*ast.TypeSwitchStmt]("Body")
	case TypeSwitchStmt_Init:
		return info[*ast.TypeSwitchStmt]("Init")
	case UnaryExpr_X:
		return info[*ast.UnaryExpr]("X")
	case ValueSpec_Comment:
		return info[*ast.ValueSpec]("Comment")
	case ValueSpec_Doc:
		return info[*ast.ValueSpec]("Doc")
	case ValueSpec_Names:
		return info[*ast.ValueSpec]("Names")
	case ValueSpec_Type:
		return info[*ast.ValueSpec]("Type")
	case ValueSpec_Values:
		return info[*ast.ValueSpec]("Values")
	default:
		panic(fmt.Sprintf("invalid Kind: %d", int(k)))
	}
}
