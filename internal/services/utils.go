package services

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func extractParams(expr string) []string {
	if strings.HasPrefix(expr, "=-") {
		expr = "=" + expr[2:]
	}

	// Match all sequences of word characters
	re := regexp.MustCompile(`\b\w+\b`)
	matches := re.FindAllString(expr, -1)

	uniqueMatches := make(map[string]bool)
	for _, match := range matches {
		index := strings.Index(expr, match)
		if index == 0 || (!unicode.IsLetter(rune(expr[index-1])) && !unicode.IsDigit(rune(expr[index-1]))) {
			if _, err := strconv.Atoi(match); err != nil { // If not purely numeric
				uniqueMatches[match] = true
			}
		}
	}

	params := make([]string, 0, len(uniqueMatches))
	for param := range uniqueMatches {
		params = append(params, param)
	}

	return params
}

func replaceParamsWithValues(expr string, params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}

	// Sort the keys based on length, longest first
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})

	// Replace based on sorted keys
	for _, k := range keys {
		expr = strings.ReplaceAll(expr, k, params[k])
	}

	return expr
}

func calculate(expr string) (float64, error) {
	trimmedExpr := strings.TrimPrefix(cleanExpression(expr), "=")
	expression, err := parser.ParseExpr(trimmedExpr)
	if err != nil {
		return 0, err
	}
	return eval(expression)
}

func eval(expr ast.Expr) (float64, error) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		x, err := eval(e.X)
		if err != nil {
			return 0, err
		}
		y, err := eval(e.Y)
		if err != nil {
			return 0, err
		}
		switch e.Op {
		case token.ADD:
			return x + y, nil
		case token.SUB:
			return x - y, nil
		case token.MUL:
			return x * y, nil
		case token.QUO:
			if y == 0 {
				return 0, errors.New("division by zero")
			}
			return x / y, nil
		default:
			return 0, fmt.Errorf("unsupported binary operator: %v", e.Op)
		}
	case *ast.ParenExpr:
		return eval(e.X)
	case *ast.BasicLit:
		if e.Kind == token.INT || e.Kind == token.FLOAT {
			return strconv.ParseFloat(e.Value, 64)
		}
	case *ast.UnaryExpr:
		// For unary minus
		if e.Op == token.SUB {
			value, err := eval(e.X)
			if err != nil {
				return 0, err
			}
			return -value, nil
		}
		return 0, fmt.Errorf("unary operation %v not supported", e.Op)
	default:
		return 0, fmt.Errorf("expression type %T not supported", e)
	}
	return 0, nil
}

func cleanExpression(expr string) string {
	re := regexp.MustCompile(`(\+|\-)+`)
	expr = re.ReplaceAllStringFunc(expr, func(s string) string {
		if strings.Count(s, "-")%2 == 0 {
			return "+"
		}
		return "-"
	})
	return strings.Replace(expr, "=+", "=", 1)
}

func isValid(str string) bool {
	_, intErr := strconv.Atoi(str)
	_, floatErr := strconv.ParseFloat(str, 64)

	if intErr != nil && floatErr != nil {
		if !strings.HasPrefix(str, "=") {
			return false
		}
	}
	forbidden := []string{"++", "--", "//", "**"}

	for _, combo := range forbidden {
		if strings.Contains(str, combo) {
			return false
		}
	}
	return true
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
