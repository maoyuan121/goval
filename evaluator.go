package goval

import (
	"github.com/maja42/goval/internal"
)

// NewEvaluator creates a new evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// Evaluator is used to evaluate expression strings.
type Evaluator struct {
}

// 可以从表达式内部调用 ExpressionFunction。
//
// 返回的对象必须是下面类型中的一种： `nil`, `bool`, `int`, `float64`, `[]interface{}` or `map[string]interface{}`.
type ExpressionFunction = func(args ...interface{}) (interface{}, error)

// Evaluate the given expression string.
//
// Optionally accepts a list of variables (accessible but not modifiable from within expressions).
//
// Optionally accepts a list of expression functions (can be called from within expressions).
//
// Returns the resulting object or an error.
//
// Stateless. Can be called concurrently. If expression functions modify variables, concurrent execution requires additional synchronization.
func (e *Evaluator) Evaluate(str string, variables map[string]interface{}, functions map[string]ExpressionFunction) (result interface{}, err error) {
	return internal.Evaluate(str, variables, functions)
}
