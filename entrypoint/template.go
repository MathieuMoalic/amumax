package entrypoint

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Expression struct {
	Prefix   string
	Array    []float64
	Format   string
	Suffix   string
	Original string
}

func (e Expression) PrettyPrint() {
	fmt.Printf("Prefix: %s\nArray: %v\nFormat: %s\nSuffix: %s\nOriginal: %s\n", e.Prefix, e.Array, e.Format, e.Suffix, e.Original)
}

// Parses template file and returns its content as a string
func parseTemplate(templatePath string) (string, error) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("error reading template file: %v", err)
	}
	return string(content), nil
}

// Generates a range similar to numpy.arange
func arange(start, end, step float64) []float64 {
	var result []float64
	for i := start; i <= end; i += step {
		result = append(result, i)
	}
	return result
}

// Generates a range similar to numpy.linspace
func linspace(start, end float64, count int) []float64 {
	var result []float64
	if count == 1 {
		return []float64{start}
	}
	step := (end - start) / float64(count-1)
	for i := 0; i < count; i++ {
		result = append(result, start+float64(i)*step)
	}
	return result
}

func isFloatFormat(format string) bool {
	// Regular expression to match formats like %f, %02.0f, %.1f, %03.0f, etc.
	re := regexp.MustCompile(`^%[-+]?[0-9]*(\.[0-9]*)?f$`)
	return re.MatchString(format)
}

func validateFields(exprMap map[string]string) error {
	// check if all the field names are valid
	for key := range exprMap {
		if key != "prefix" && key != "array" && key != "format" && key != "suffix" && key != "start" && key != "end" && key != "step" && key != "count" {
			return fmt.Errorf("invalid field name: %s", key)
		}
	}
	// if array is given, start, end, step and count should not be given
	if _, ok := exprMap["array"]; ok {
		if _, ok := exprMap["start"]; ok {
			return fmt.Errorf("start should not be given when array is given")
		}
		if _, ok := exprMap["end"]; ok {
			return fmt.Errorf("end should not be given when array is given")
		}
		if _, ok := exprMap["step"]; ok {
			return fmt.Errorf("step should not be given when array is given")
		}
		if _, ok := exprMap["count"]; ok {
			return fmt.Errorf("count should not be given when array is given")
		}
	} else {
		// if array is not given, start, end, step or count should be given
		if _, ok := exprMap["start"]; !ok {
			return fmt.Errorf("start should be given when array is not given")
		}
		if _, ok := exprMap["end"]; !ok {
			return fmt.Errorf("end should be given when array is not given")
		}
		// step or count should be given
		if _, ok := exprMap["step"]; !ok {
			if _, ok := exprMap["count"]; !ok {
				return fmt.Errorf("step or count should be given when array is not given")
			}
		}
	}
	return nil
}

func parsePrefix(exprMap map[string]string) (string, error) {
	if prefix, ok := exprMap["prefix"]; ok {
		return prefix, nil
	}
	return "", nil
}

func parseArray(exprMap map[string]string) ([]float64, error) {
	if arrStr, ok := exprMap["array"]; ok {
		arrStr = strings.Trim(arrStr, "[]")
		// splits the string into an array of strings
		arr := strings.Split(arrStr, ",")
		var array []float64
		for _, valStr := range arr {
			// remove whitespace
			valStr = strings.TrimSpace(valStr)
			val, err := strconv.ParseFloat(string(valStr), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid array value: %v", err)
			}
			array = append(array, val)
		}
		return array, nil
	}
	return nil, nil
}

func parseFormat(exprMap map[string]string) (string, error) {
	if format, ok := exprMap["format"]; ok {
		if !isFloatFormat(format) {
			return "", fmt.Errorf("only the %%f format is allowed: %s", format)
		}
		return format, nil
	}
	return "%v", nil // Default format
}

func parseSuffix(exprMap map[string]string) (string, error) {
	if suffix, ok := exprMap["suffix"]; ok {
		return suffix, nil
	}
	return "", nil
}

func parseRange(exprMap map[string]string) ([]float64, error) {
	start, err := strconv.ParseFloat(exprMap["start"], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid start value: %v", err)
	}
	end, err := strconv.ParseFloat(exprMap["end"], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid end value: %v", err)
	}
	if start > end {
		return nil, fmt.Errorf("start value should be less than end value")
	}
	if start == end {
		return []float64{start}, nil
	}

	if stepStr, ok := exprMap["step"]; ok {
		step, err := strconv.ParseFloat(stepStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid step value: %v", err)
		}
		return arange(start, end, step), nil
	} else if countStr, ok := exprMap["count"]; ok {
		count, err := strconv.Atoi(countStr)
		if err != nil {
			return nil, fmt.Errorf("invalid count value: %v", err)
		}
		return linspace(start, end, count), nil
	}
	return nil, fmt.Errorf("invalid range specification")
}

func parseExpression(exprMap map[string]string) (Expression, error) {
	prefix, err := parsePrefix(exprMap)
	if err != nil {
		return Expression{}, err
	}
	format, err := parseFormat(exprMap)
	if err != nil {
		return Expression{}, err
	}
	suffix, err := parseSuffix(exprMap)
	if err != nil {
		return Expression{}, err
	}
	var array []float64
	if exprMap["array"] == "" {
		array, err = parseRange(exprMap)
		if err != nil {
			return Expression{}, err
		}
	} else {
		array, err = parseArray(exprMap)
		if err != nil {
			return Expression{}, err
		}
	}
	return Expression{
		Prefix: prefix,
		Array:  array,
		Format: format,
		Suffix: suffix,
	}, nil
}

// Finds expressions in the mx3 template and parses them
func findExpressions(mx3 string) (expressions []Expression, err error) {
	regex := regexp.MustCompile(`"\{(.*?)\}"`)
	matches := regex.FindAllStringSubmatch(mx3, -1)

	for _, match := range matches {
		extracted := match[1]
		parts := strings.Split(extracted, ";")
		exprMap := make(map[string]string)

		// Parse the expression parts
		for _, part := range parts {
			if strings.Contains(part, "=") {
				keyVal := strings.SplitN(part, "=", 2)
				exprMap[keyVal[0]] = keyVal[1]
			} else {
				return nil, fmt.Errorf("invalid expression part: %s", part)
			}
		}

		if err = validateFields(exprMap); err != nil {
			return nil, err
		}
		exp := Expression{}
		exp.Prefix, err = parsePrefix(exprMap)
		if err != nil {
			return nil, err
		}
		exp, err := parseExpression(exprMap)
		if err != nil {
			return nil, err
		}
		exp.Original = extracted
		expressions = append(expressions, exp)
	}
	return expressions, nil
}

// Generates the files with the processed expressions
func generateFiles(parentDir, mx3 string, expressions []Expression) error {
	// Generate combinations of all arrays using cartesian product
	combinationCount := 1
	for _, exp := range expressions {
		combinationCount *= len(exp.Array)
	}

	indices := make([]int, len(expressions))
	for i := 0; i < combinationCount; i++ {
		pathParts := make([]string, len(expressions))
		newMx3 := mx3

		for j, exp := range expressions {
			value := exp.Array[indices[j]]
			formattedValue := fmt.Sprintf(exp.Format, value)
			// formattedValue := fmt.Sprintf(exp.Format, int(value)) // Convert float64 to int
			pathParts[j] = fmt.Sprintf("%s%s%s", exp.Prefix, formattedValue, exp.Suffix)
			// Replace the placeholder in the mx3 template
			newMx3 = strings.ReplaceAll(newMx3, `"{`+exp.Original+`}"`, fmt.Sprintf("%v", value))
		}

		// Construct file path
		fullPath := filepath.Join(parentDir, strings.Join(pathParts, "/")+".mx3")
		err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directories: %v", err)
		}

		// Write the new file
		err = os.WriteFile(fullPath, []byte(newMx3), 0644)
		if err != nil {
			return fmt.Errorf("error writing file %s: %v", fullPath, err)
		}

		// Update the indices for next combination
		for j := len(indices) - 1; j >= 0; j-- {
			indices[j]++
			if indices[j] < len(expressions[j].Array) {
				break
			}
			indices[j] = 0
		}
	}

	return nil
}

// Main function for handling the template logic
func template(path string) (err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}

	parentDir := filepath.Dir(path)

	mx3, err := parseTemplate(path)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	expressions, err := findExpressions(mx3)
	if err != nil {
		return fmt.Errorf("error finding expressions: %v", err)
	}

	err = generateFiles(parentDir, mx3, expressions)
	if err != nil {
		return fmt.Errorf("error generating files: %v", err)
	}
	return
}
