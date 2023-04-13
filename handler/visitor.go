package handler

import (
	"fmt"
	"time"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/pingcap/tidb/parser/ast"
	td "github.com/pingcap/tidb/parser/test_driver"
)

type Visitor struct {
	Result  *mysql.Result
	Error   error
	handler *Handler
}

func NewVisitor(h *Handler) *Visitor {
	return &Visitor{
		Result:  &mysql.Result{},
		Error:   nil,
		handler: h,
	}
}

func putOneColumn(names *[]string, values *[][]interface{}, field string, value []interface{}) {
	*names = append(*names, field)
	vc := make([]interface{}, 0)
	vc = append(vc, value...)
	*values = append(*values, vc)
}

func (v *Visitor) handleSelectStmt(n *ast.SelectStmt) error {
	if n.Fields == nil || n.Fields.Fields == nil || n.Fields.Fields[0].Expr == nil {
		return fmt.Errorf("no fields")
	}
	field := n.Fields.Fields[0]
	switch expr := field.Expr.(type) {
	case *ast.VariableExpr:
		var value interface{}
		if expr.IsGlobal {
			value = GlobalVariables.GetVariable(expr.Name)
		} else {
			value = v.handler.GetVariable(expr.Name)
		}
		names := make([]string, 0)
		values := make([][]interface{}, 0)
		putOneColumn(&names, &values, field.Text(), []interface{}{value})
		v.Result.Resultset, v.Error = mysql.BuildSimpleResultset(names, values, false)
		return nil
	case *ast.FuncCallExpr:
		switch expr.FnName.L {
		case "unix_timestamp":
			names := make([]string, 0)
			values := make([][]interface{}, 0)
			now := time.Now().Unix()
			putOneColumn(&names, &values, field.Text(), []interface{}{int32(now)})
			v.Result.Resultset, v.Error = mysql.BuildSimpleResultset(names, values, false)
			return nil
		default:
			return fmt.Errorf("unsupported function: %s", expr.FnName.O)
		}
	default:
		return fmt.Errorf("unsupported expression: %T", expr)
	}
}

func (v *Visitor) handleSetStmt(n *ast.SetStmt) error {
	if len(n.Variables) == 0 {
		return fmt.Errorf("no variables")
	}
	for _, variable := range n.Variables {
		var value interface{}
		switch expr := variable.Value.(type) {
		case *ast.VariableExpr:
			if expr.IsGlobal {
				value = GlobalVariables.GetVariable(expr.Name)
			} else {
				value = v.handler.GetVariable(expr.Name)
			}
		case *td.ValueExpr:
			value = expr.Datum.GetValue()
		default:
		}
		if variable.IsGlobal {
			GlobalVariables.SetVariable(variable.Name, value)
		} else {
			v.handler.SetVariable(variable.Name, value)
		}
	}
	return nil
}

func (v *Visitor) Enter(n ast.Node) (ast.Node, bool) {
	switch x := n.(type) {
	case *ast.SelectStmt:
		v.Error = v.handleSelectStmt(x)
		return x, true
	case *ast.SetStmt:
		v.Error = v.handleSetStmt(x)
	default:
		v.Error = fmt.Errorf("unsupported statement: %T", x)
	}
	return n, true
}

func (v *Visitor) Leave(n ast.Node) (ast.Node, bool) {
	return n, true
}

func (v *Visitor) Clean() {
	v.Result = &mysql.Result{}
	v.Error = nil
}
