package base

import (
	"errors"
	"fmt"
	"github.com/bilibili/gengine/context"
	"reflect"
)

var TypeMap = map[string]string{
	"int":     "int",
	"int8":    "int8",
	"int16":   "int16",
	"int32":   "int32",
	"int64":   "int64",
	"uint":    "uint",
	"uint8":   "uint8",
	"uint16":  "uint16",
	"uint32":  "uint32",
	"uint64":  "uint64",
	"float32": "float32",
	"float64": "float64",
}

type Expression struct {
	SourceCode
	ExpressionLeft     *Expression
	ExpressionRight    *Expression
	ExpressionAtom     *ExpressionAtom
	MathExpression     *MathExpression
	LogicalOperator    string
	ComparisonOperator string
	NotOperator        string
}

func (e *Expression) AcceptExpressionAtom(atom *ExpressionAtom) error {
	if e.ExpressionAtom == nil {
		e.ExpressionAtom = atom
		return nil
	}
	return errors.New("ExpressionAtom already set twice! ")
}

func (e *Expression) AcceptMathExpression(atom *MathExpression) error {
	if e.MathExpression == nil {
		e.MathExpression = atom
		return nil
	}
	return errors.New(" Expression's MathExpression set twice! ")
}

func (e *Expression) AcceptExpression(expression *Expression) error {
	if e.ExpressionLeft == nil {
		e.ExpressionLeft = expression
		return nil
	}

	if e.ExpressionRight == nil {
		e.ExpressionRight = expression
		return nil
	}
	return errors.New("Expression already set twice! ")
}

func (e *Expression) Evaluate(dc *context.DataContext, Vars map[string]reflect.Value) (reflect.Value, error) {

	//priority to calculate single value
	var math reflect.Value
	if e.MathExpression != nil {
		evl, err := e.MathExpression.Evaluate(dc, Vars)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
		math = evl
	}

	var atom reflect.Value
	if e.ExpressionAtom != nil {
		evl, err := e.ExpressionAtom.Evaluate(dc, Vars)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
		atom = evl
	}

	var b reflect.Value
	if e.ExpressionRight == nil {
		if e.ExpressionLeft != nil {
			left, err := e.ExpressionLeft.Evaluate(dc, Vars)
			if err != nil {
				return reflect.ValueOf(nil), err
			}
			b = left
		}
	}

	// && ||  just only to be used between boolean
	if e.LogicalOperator != "" {

		lv, err := e.ExpressionLeft.Evaluate(dc, Vars)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		rv, err := e.ExpressionRight.Evaluate(dc, Vars)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		//
		flv := lv //reflect.ValueOf(lv)
		frv := rv //reflect.ValueOf(rv)

		if lv.Kind() == reflect.Bool && rv.Kind() == reflect.Bool {
			if e.LogicalOperator == "&&" {
				b = reflect.ValueOf(flv.Bool() && frv.Bool())
			}
			if e.LogicalOperator == "||" {
				b = reflect.ValueOf(flv.Bool() || frv.Bool())
			}
		} else {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, || or && can't be used between %s and %s:\n", e.LineNum, e.Column, e.Code, flv.Kind().String(), frv.Kind().String()))
		}
	}

	// == > < != >= <=  just only to be used between number and number, string and string, bool and bool
	if e.ComparisonOperator != "" {

		lv, err := e.ExpressionLeft.Evaluate(dc, Vars)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		rv, err := e.ExpressionRight.Evaluate(dc, Vars)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		//
		flv := lv //reflect.ValueOf(lv)
		frv := rv //reflect.ValueOf(rv)

		//string compare
		tlv := lv //reflect.TypeOf(lv).String()
		trv := rv //reflect.TypeOf(rv).String()
		if tlv.Kind() == reflect.String && trv.Kind() == reflect.String {
			switch e.ComparisonOperator {
			case "==":
				b = reflect.ValueOf(flv.String() == frv.String())
				break
			case "!=":
				b = reflect.ValueOf(flv.String() != frv.String())
				break
			case ">":
				b = reflect.ValueOf(flv.String() > frv.String())
				break
			case "<":
				b = reflect.ValueOf(flv.String() < frv.String())
				break
			case ">=":
				b = reflect.ValueOf(flv.String() >= frv.String())
				break
			case "<=":
				b = reflect.ValueOf(flv.String() <= frv.String())
				break
			default:
				return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, Can't be recognized ComparisonOperator: %s", e.LineNum, e.Column, e.Code, e.ComparisonOperator))
			}
			goto LAST
		}

		//data compare
		if l, ok1 := TypeMap[tlv.Kind().String()]; ok1 {
			if r, ok2 := TypeMap[trv.Kind().String()]; ok2 {
				var ll float64
				switch l {
				case "int", "int8", "int16", "int32", "int64":
					ll = float64(flv.Int())
					break
				case "uint", "uint8", "uint16", "uint32", "uint64":
					ll = float64(flv.Uint())
					break
				case "float32", "float64":
					ll = flv.Float()
					break
				}

				var rr float64
				switch r {
				case "int", "int8", "int16", "int32", "int64":
					rr = float64(frv.Int())
					break
				case "uint", "uint8", "uint16", "uint32", "uint64":
					rr = float64(frv.Uint())
					break
				case "float32", "float64":
					rr = frv.Float()
					break
				}

				switch e.ComparisonOperator {
				case "==":
					b = reflect.ValueOf(ll == rr)
					break
				case "!=":
					b = reflect.ValueOf(ll != rr)
					break
				case ">":
					b = reflect.ValueOf(ll > rr)
					break
				case "<":
					b = reflect.ValueOf(ll < rr)
					break
				case ">=":
					b = reflect.ValueOf(ll >= rr)
					break
				case "<=":
					b = reflect.ValueOf(ll <= rr)
					break
				default:
					return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, Can't be recognized ComparisonOperator: %s", e.LineNum, e.Column, e.Code, e.ComparisonOperator))
				}
			}
			goto LAST
		}

		if tlv.Kind() == reflect.Bool && trv.Kind() == reflect.Bool {
			switch e.ComparisonOperator {
			case "==":
				b = reflect.ValueOf(flv.Bool() == frv.Bool())
				break
			case "!=":
				b = reflect.ValueOf(flv.Bool() != frv.Bool())
				break
			default:
				return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, Can't be recognized ComparisonOperator: %s", e.LineNum, e.Column, e.Code, e.ComparisonOperator))
			}
			goto LAST
		}
	}

LAST:
	if e.NotOperator == "!" {

		if math != reflect.ValueOf(nil) {
			return reflect.ValueOf(!math.Bool()), nil
		}

		if atom != reflect.ValueOf(nil) {
			return reflect.ValueOf(!atom.Bool()), nil
		}

		if b != reflect.ValueOf(nil) {
			return reflect.ValueOf(!b.Bool()), nil
		}
	} else {
		if math != reflect.ValueOf(nil) {
			return math, nil
		}

		if atom != reflect.ValueOf(nil) {
			return atom, nil
		}

		if b != reflect.ValueOf(nil) {
			return b, nil
		}
	}
	return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, evaluate Expression err!", e.LineNum, e.Column, e.Code))
}
