package engine

import (
	"fmt"
	"strconv"
	"time"

	"github.com/nao1215/aiondb/engine/parser/core"
)

// Operator compares 2 values and return a boolean
type Operator func(leftValue Value, rightValue Value) bool

// NewOperator initializes the operator matching the Token number
func NewOperator(token core.TokenID, lexeme string) (Operator, error) {
	switch token {
	case core.TokenIDEquality:
		return equalityOperator, nil
	case core.TokenIDDistinctness:
		return distinctnessOperator, nil
	case core.TokenIDLeftDiple:
		return lessThanOperator, nil
	case core.TokenIDRightDiple:
		return greaterThanOperator, nil
	case core.TokenIDLessOrEqual:
		return lessOrEqualOperator, nil
	case core.TokenIDGreaterOrEqual:
		return greaterOrEqualOperator, nil
	}
	return nil, fmt.Errorf("operator '%s' does not exist", lexeme)
}

// toDate converts a value to a date
func toDate(t interface{}) (time.Time, error) {
	switch t := t.(type) {
	case string:
		d, err := core.ParseDate(t)
		if err != nil {
			return time.Time{}, fmt.Errorf("cannot parse date %v", t)
		}
		return *d, nil
	default:
		return time.Time{}, fmt.Errorf("unexpected internal type %T", t)
	}
}

// toFloat converts a value to a float
func toFloat(t interface{}) (float64, error) {
	switch t := t.(type) {
	case float64:
		return t, nil
	case int64:
		return float64(t), nil
	case int:
		return float64(t), nil
	case string:
		return strconv.ParseFloat(t, 64)
	default:
		return 0, fmt.Errorf("unexpected internal type %T", t)
	}
}

// greaterThanOperator compares 2 values and return true if the left value is greater than the right value
func greaterThanOperator(leftValue Value, rightValue Value) bool {
	var left, right float64
	var rValue interface{}

	if rightValue.v != nil {
		rValue = rightValue.v
	} else {
		rValue = rightValue.lexeme
	}

	var leftDate time.Time
	var isDate bool
	left, err := toFloat(leftValue.v)
	if err != nil {
		leftDate, err = toDate(leftValue.v)
		if err != nil {
			return false
		}
		isDate = true
	}

	if !isDate {
		right, err = toFloat(rValue)
		if err != nil {
			return false
		}
		return left > right
	}

	rightDate, err := toDate(rValue)
	if err != nil {
		return false
	}
	return leftDate.After(rightDate)
}

// lessOrEqualOperator compares 2 values and return true if the left value is less or equal than the right value
func lessOrEqualOperator(leftValue Value, rightValue Value) bool {
	return lessThanOperator(leftValue, rightValue) || equalityOperator(leftValue, rightValue)
}

// greaterOrEqualOperator compares 2 values and return true if the left value is greater or equal than the right value
func greaterOrEqualOperator(leftValue Value, rightValue Value) bool {
	return greaterThanOperator(leftValue, rightValue) || equalityOperator(leftValue, rightValue)
}

// lessThanOperator compares 2 values and return true if the left value is less than the right value
func lessThanOperator(leftValue Value, rightValue Value) bool {
	var left, right float64

	var rValue interface{}
	if rightValue.v != nil {
		rValue = rightValue.v
	} else {
		rValue = rightValue.lexeme
	}

	var leftDate time.Time
	var isDate bool

	left, err := toFloat(leftValue.v)
	if err != nil {
		leftDate, err = toDate(leftValue.v)
		if err != nil {
			return false
		}
		isDate = true
	}

	if !isDate {
		right, err = toFloat(rValue)
		if err != nil {
			return false
		}
		return left < right
	}

	rightDate, err := toDate(rValue)
	if err != nil {
		return false
	}
	return leftDate.Before(rightDate)
}

// equalityOperator checks if given value are equal
func equalityOperator(leftValue Value, rightValue Value) bool {
	return fmt.Sprintf("%v", leftValue.v) == rightValue.lexeme
}

// distinctnessOperator checks if given value are distinct
func distinctnessOperator(leftValue Value, rightValue Value) bool {
	return fmt.Sprintf("%v", leftValue.v) != rightValue.lexeme
}

// TrueOperator always returns true
func trueOperator(_ Value, _ Value) bool {
	return true
}

// inOperator checks if the left value is in the right value
// Right value should be a slice of string
func inOperator(leftValue Value, rightValue Value) bool {
	values, ok := rightValue.v.([]string)
	if !ok {
		return false
	}
	for i := range values {
		if fmt.Sprintf("%v", leftValue.v) == values[i] {
			return true
		}
	}
	return false
}

// notInOperator checks if the left value is not in the right value
func notInOperator(leftValue Value, rightValue Value) bool {
	return !inOperator(leftValue, rightValue)
}

// isNullOperator checks if the left value is null
func isNullOperator(leftValue Value, _ Value) bool {
	return leftValue.v == nil
}

// isNotNullOperator checks if the left value is not null
func isNotNullOperator(leftValue Value, _ Value) bool {
	return leftValue.v != nil
}
