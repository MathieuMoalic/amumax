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
	Prefix    string
	Array     []string
	Format    string
	Suffix    string
	Original  string
	IsNumeric bool
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

func isValidFormat(format string) bool {
	// Allow %s for strings, and float formats like %f, %02.0f, %.1f, %03.0f, etc.
	re := regexp.MustCompile(`^%[-+]?[0-9]*(\.[0-9]*)?[fs]$`)
	return re.MatchString(format)
}

func validateFields(exprMap map[string]string) error {
	// Check if all the field names are valid
	validKeys := map[string]bool{
		"prefix": true, "array": true, "format": true, "suffix": true,
		"start": true, "end": true, "step": true, "count": true,
	}
	for key := range exprMap {
		if !validKeys[key] && key != "original" {
			return fmt.Errorf("invalid field name: %s", key)
		}
	}
	// If array is given, start, end, step, and count should not be given
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
		// If array is not given, start and end should be given
		if _, ok := exprMap["start"]; !ok {
			return fmt.Errorf("start should be given when array is not given")
		}
		if _, ok := exprMap["end"]; !ok {
			return fmt.Errorf("end should be given when array is not given")
		}
		// Step or count should be given
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

func parseArray(exprMap map[string]string) ([]string, error) {
	if arrStr, ok := exprMap["array"]; ok {
		arrStr = strings.Trim(arrStr, "[]")
		if arrStr == "" {
			return nil, fmt.Errorf("array cannot be empty")
		}
		arr := strings.Split(arrStr, ",")
		var array []string
		for _, valStr := range arr {
			valStr = strings.TrimSpace(valStr)
			// Remove quotes if any
			valStr = strings.Trim(valStr, "\"'")
			array = append(array, valStr)
		}
		if len(array) == 0 {
			return nil, fmt.Errorf("array cannot be empty")
		}
		return array, nil
	}
	return nil, nil
}

func parseFormat(exprMap map[string]string) (string, error) {
	if format, ok := exprMap["format"]; ok {
		if !isValidFormat(format) {
			return "", fmt.Errorf("invalid format: %s", format)
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

func parseRange(exprMap map[string]string) ([]string, error) {
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
		return []string{fmt.Sprintf("%v", start)}, nil
	}

	var numericArray []float64
	if stepStr, ok := exprMap["step"]; ok {
		step, err := strconv.ParseFloat(stepStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid step value: %v", err)
		}
		if step <= 0 {
			return nil, fmt.Errorf("step value should be greater than 0")
		}
		length := int((end-start)/step) + 1
		if length > 1000 {
			return nil, fmt.Errorf("step value is too small (more than 1000 elements)")
		}
		numericArray = arange(start, end, step)
	} else if countStr, ok := exprMap["count"]; ok {
		count, err := strconv.Atoi(countStr)
		if err != nil {
			return nil, fmt.Errorf("invalid count value: %v", err)
		}
		if count <= 0 {
			return nil, fmt.Errorf("count value should be greater than 0")
		}
		numericArray = linspace(start, end, count)
	} else {
		return nil, fmt.Errorf("invalid range specification")
	}

	array := make([]string, len(numericArray))
	for i, val := range numericArray {
		array[i] = fmt.Sprintf("%v", val)
	}
	return array, nil
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
	var array []string
	isNumeric := true
	if format == "%s" {
		isNumeric = false
		if _, ok := exprMap["array"]; ok {
			array, err = parseArray(exprMap)
			if err != nil {
				return Expression{}, err
			}
		} else {
			return Expression{}, fmt.Errorf("format %s is only valid for numeric arrays", format)
		}
	} else {
		if _, ok := exprMap["array"]; ok {
			array, err = parseArray(exprMap)
			if err != nil {
				return Expression{}, err
			}
			// Determine if the array is numeric
			for _, val := range array {
				if _, err = strconv.ParseFloat(val, 64); err != nil {
					isNumeric = false
					break
				}
			}
		} else {
			// Parse range
			array, err = parseRange(exprMap)
			if err != nil {
				return Expression{}, err
			}
		}
	}
	exp := Expression{
		Prefix:    prefix,
		Array:     array,
		Format:    format,
		Suffix:    suffix,
		IsNumeric: isNumeric,
		Original:  exprMap["original"],
	}
	return exp, nil
}

// Finds expressions in the mx3 template and parses them
func findExpressions(mx3 string) (expressions []Expression, err error) {
	regex := regexp.MustCompile(`"\{(.*?)\}"`)
	matches := regex.FindAllStringSubmatch(mx3, -1)

	for _, match := range matches {
		extracted := match[1]
		parts := strings.Split(extracted, ";")
		exprMap := make(map[string]string)
		exprMap["original"] = extracted

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
		exp, err := parseExpression(exprMap)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, exp)
	}
	if len(expressions) == 0 {
		return nil, fmt.Errorf("no expressions found")
	}
	return expressions, nil
}

// Generates the files with the processed expressions
func generateFiles(parentDir, mx3 string, expressions []Expression, flat bool) error {
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
			var formattedValue string

			// Check if format is compatible with value type
			if exp.IsNumeric {
				numValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return fmt.Errorf("error parsing numeric value '%v': %v", value, err)
				}
				if strings.HasSuffix(exp.Format, "s") {
					return fmt.Errorf("invalid format '%s' for numeric value '%v'", exp.Format, value)
				}
				formattedValue = fmt.Sprintf(exp.Format, numValue)
			} else {
				if strings.HasSuffix(exp.Format, "d") || strings.HasSuffix(exp.Format, "f") {
					return fmt.Errorf("invalid format '%s' for string value '%v'", exp.Format, value)
				}
				formattedValue = fmt.Sprintf(exp.Format, value)
			}
			// formattedValue := fmt.Sprintf(exp.Format, value)

			// Build path parts and replace placeholders
			pathParts[j] = fmt.Sprintf("%s%s%s", exp.Prefix, formattedValue, exp.Suffix)
			newMx3 = strings.Replace(newMx3, `"{`+exp.Original+`}"`, value, 1)
		}

		var joinedPath string
		if flat {
			joinedPath = strings.Join(pathParts, "")
		} else {
			joinedPath = filepath.Join(pathParts...)
		}

		fullPath := filepath.Join(parentDir, joinedPath+".mx3")
		if !flat {
			err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
			if err != nil {
				return fmt.Errorf("error creating directories: %v", err)
			}
		}

		// Write the new file
		err := os.WriteFile(fullPath, []byte(newMx3), 0644)
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
func template(path string, flat bool) (err error) {
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

	err = generateFiles(parentDir, mx3, expressions, flat)
	if err != nil {
		return fmt.Errorf("error generating files: %v", err)
	}
	return
}
