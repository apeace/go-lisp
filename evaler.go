package lisp

import "fmt"

var Env map[string]interface{}

func init() {
	Env = make(map[string]interface{})
}

func Eval(expr interface{}) (interface{}, error) {
	fmt.Printf("EVAL: %v\n", expr)
	switch expr.(type) {
	case int: // Int
		return expr, nil
	case string:
		if val, ok := Env[expr.(string)]; ok {
			return val, nil
		} else {
			return "", fmt.Errorf("Unknown symbol: %v", expr)
		}
	case []interface{}:
		tokens := expr.([]interface{})
		t := tokens[0]
		if _, ok := t.([]interface{}); ok {
			return Eval(t)
		} else if t == "quote" { // Quote
			return tokens[1:], nil
		} else if t == "define" { // Define
			Env[tokens[1].(string)] = tokens[2]
			return tokens[2], nil
		} else if t == "set!" { // Set!
			key := tokens[1].(string)
			if _, ok := Env[key]; ok {
				Env[key] = tokens[2]
				return tokens[2], nil
			} else {
				return nil, fmt.Errorf("Can only set! variable that is previously defined")
			}
		} else if t == "if" { // If
			if tokens[1] == "true" && len(tokens) > 2 {
				return Eval(tokens[2])
			} else if len(tokens) > 3 {
				return Eval(tokens[3])
			}
			return "nil", nil
		} else if t == "begin" { // Begin
			var r interface{}
			var err error
			for _, val := range tokens[1:] {
				r, err = Eval(val)
				if err != nil {
					return nil, err
				}
			}
			return r, nil
		} else if t == "+" { // Addition
			var sum int
			for _, i := range tokens[1:] {
				j, err := Eval(i)
				if err != nil {
					return nil, err
				}
				v, ok := j.(int)
				if ok {
					sum += int(v)
				} else {
					return nil, fmt.Errorf("Cannot only add numbers: %v", i)
				}
			}
			return sum, nil
		}
	}
	return "", fmt.Errorf("Unknown data type: %v", expr)
}
