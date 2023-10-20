package services

import (
	"math"
	"reflect"
	"testing"
)

func Test_extractParams(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Handle '=-' prefix",
			args: args{expr: "=-x + 2"},
			want: []string{"x"},
		},
		{
			name: "Mixed words and numbers",
			args: args{expr: "=2-x1"},
			want: []string{"x1"},
		},
		{
			name: "No valid params",
			args: args{expr: "=123 + 456"},
			want: []string{},
		},
		{
			name: "Only numbers",
			args: args{expr: "123456"},
			want: []string{},
		},
		{
			name: "Empty expression",
			args: args{expr: ""},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractParams(tt.args.expr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_replaceParamsWithValues(t *testing.T) {
	type args struct {
		expr   string
		params map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Different lengths, no prefix",
			args: args{
				expr:   "x+y",
				params: map[string]string{"x": "10", "y": "20"},
			},
			want: "10+20",
		},
		{
			name: "Numbers with params",
			args: args{
				expr:   "2*x+y-3",
				params: map[string]string{"x": "10", "y": "20"},
			},
			want: "2*10+20-3",
		},
		{
			name: "One key is a prefix of another",
			args: args{
				expr:   "xy+x",
				params: map[string]string{"xy": "15", "x": "10"},
			},
			want: "15+10",
		},
		{
			name: "Key containing number",
			args: args{
				expr:   "1x+x1+1",
				params: map[string]string{"1x": "15", "x1": "10"},
			},
			want: "15+10+1",
		},
		{
			name: "No replacements",
			args: args{
				expr:   "1+2",
				params: map[string]string{},
			},
			want: "1+2",
		},
		{
			name: "Empty expression",
			args: args{
				expr:   "",
				params: map[string]string{},
			},
			want: "",
		},
		{
			name: "Empty params map",
			args: args{
				expr:   "x+y",
				params: map[string]string{},
			},
			want: "x+y",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceParamsWithValues(tt.args.expr, tt.args.params); got != tt.want {
				t.Errorf("replaceParamsWithValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculate(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Addition",
			args: args{expr: "=1+2"},
			want: 3,
		},
		{
			name: "Subtraction",
			args: args{expr: "=7-4"},
			want: 3,
		},
		{
			name: "Multiplication",
			args: args{expr: "=3*3"},
			want: 9,
		},
		{
			name: "Division",
			args: args{expr: "=8/2"},
			want: 4,
		},
		{
			name:    "Division by zero",
			args:    args{expr: "=8/0"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Unsupported binary operator",
			args:    args{expr: "=8^0"},
			want:    0,
			wantErr: true,
		},
		{
			name: "Nested Expressions",
			args: args{expr: "=(3+2)*(2+3)"},
			want: 25,
		},
		{
			name: "Order of Operations",
			args: args{expr: "=2+3*4"},
			want: 14,
		},
		{
			name: "Unary Operation",
			args: args{expr: "=-3+4"},
			want: 1,
		},
		{
			name: "Complex Nested Parentheses",
			args: args{expr: "=(((3+2)*2)+4)/(2+1)"},
			want: 4.666666666666667,
		},
		{
			name: "Unary Operation on Nested Expression",
			args: args{expr: "=-((3+2)*2)"},
			want: -10,
		},
		{
			name: "Large Arithmetic Operation",
			args: args{expr: "=1000000*1000000+1000000"},
			want: 1000001000000,
		},
		{
			name: "Float Calculation",
			args: args{expr: "=3.5*2"},
			want: 7,
		},
		{
			name: "Another Float Calculation",
			args: args{expr: "=10.5/2"},
			want: 5.25,
		},
		{
			name: "Multiplication with Zero",
			args: args{expr: "=0*4"},
			want: 0,
		},
		{
			name:    "Invalid Sequence with Operator at the End",
			args:    args{expr: "=5*"},
			want:    0,
			wantErr: true,
		},
		{
			name: "Expression with Spaces",
			args: args{expr: "= 5 + 2"},
			want: 7,
		},
		{
			name: "Unary Minus with Zero",
			args: args{expr: "=-0"},
			want: 0,
		},
		{
			name: "Addition with Zero",
			args: args{expr: "=5+0"},
			want: 5,
		},
		{
			name: "Smallest Nonzero Float64",
			args: args{expr: "=4.940656458412465441765687928682213723651e-324"},
			want: 4.940656458412465441765687928682213723651e-324,
		},
		{
			name: "Maximum Float64 Addition",
			args: args{expr: "=1.797693134862315708145274237317043567981e+308 + -1.797693134862315708145274237317043567981e+308"},
			want: 0,
		},
		{
			name:    "Maximum Float64 Multiplication",
			args:    args{expr: "=1.797693134862315708145274237317043567981e+308 * 2"},
			want:    math.Inf(1),
			wantErr: false, // Because the result would be Infinity
		},
		{
			name:    "Minimum Float64 Multiplication",
			args:    args{expr: "=-1.797693134862315708145274237317043567981e+308 * 2"},
			want:    math.Inf(-1),
			wantErr: false, // Because the result would be Infinity
		},
		{
			name:    "Division by Near Zero",
			args:    args{expr: "=5/0.00000000000000000001"},
			want:    5e+20,
			wantErr: false, // Depending on how your function handles it.
		},
		{
			name: "Subtraction of Very Close Values",
			args: args{expr: "=1.00000000000001 - 1"},
			want: 9.992007221626409e-15,
		},
		{
			name: "Expression with Trailing Whitespaces",
			args: args{expr: "=5+2      "},
			want: 7,
		},
		{
			name: "Expression with Multiple Whitespaces in Between",
			args: args{expr: "=5   +   2"},
			want: 7,
		},
		{
			name:    "Just Operators Without Numbers",
			args:    args{expr: "=-"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculate(tt.args.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanExpression(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "No replacements no signs",
			args: args{expr: "=x+y"},
			want: "=x+y",
		},
		{
			name: "No replacements - present",
			args: args{expr: "=-x+y"},
			want: "=-x+y",
		},
		{
			name: "Empty expression",
			args: args{expr: ""},
			want: "",
		},
		{
			name: "-- replacement",
			args: args{expr: "=--x+y"},
			want: "=x+y",
		},
		{
			name: "=+ replacement",
			args: args{expr: "=+x+y"},
			want: "=x+y",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanExpression(tt.args.expr); got != tt.want {
				t.Errorf("cleanExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValid(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid integer",
			args: args{str: "1234"},
			want: true,
		},
		{
			name: "Valid float",
			args: args{str: "123.45"},
			want: true,
		},
		{
			name: "String starts with =",
			args: args{str: "=x+123"},
			want: true,
		},
		{
			name: "String neither starts with = nor a valid number",
			args: args{str: "x+123"},
			want: false,
		},
		{
			name: "String with forbidden pattern ++",
			args: args{str: "=x++123"},
			want: false,
		},
		{
			name: "String with forbidden pattern --",
			args: args{str: "=x--123"},
			want: false,
		},
		{
			name: "String with forbidden pattern //",
			args: args{str: "=x//123"},
			want: false,
		},
		{
			name: "String with forbidden pattern **",
			args: args{str: "=x**123"},
			want: false,
		},
		{
			name: "Empty string",
			args: args{str: ""},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValid(tt.args.str); got != tt.want {
				t.Errorf("isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	type args struct {
		slice []string
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Slice contains value",
			args: args{
				slice: []string{"a", "bb", "ccc"},
				value: "bb",
			},
			want: true,
		},
		{
			name: "Slice does not contain value",
			args: args{
				slice: []string{"a", "bb", "ccc"},
				value: "m",
			},
			want: false,
		},
		{
			name: "Empty slice",
			args: args{
				slice: []string{},
				value: "m",
			},
			want: false,
		},
		{
			name: "Slice contains similar but not exact value",
			args: args{
				slice: []string{"a", "bb", "ccc"},
				value: "bbs",
			},
			want: false,
		},
		{
			name: "Slice contains multiple occurrences of the value",
			args: args{
				slice: []string{"a", "bb", "bb", "ccc"},
				value: "bb",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.slice, tt.args.value); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
