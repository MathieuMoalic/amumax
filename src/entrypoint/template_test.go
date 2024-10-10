package entrypoint

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Utility function to check if a file exists and read its content
func fileExistsAndContent(t *testing.T, path string, expectedContent string) {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Expected file %s to exist, but it does not. Error: %v", path, err)
	}

	if string(content) != expectedContent {
		t.Errorf("File %s content mismatch.\nExpected:\n%s\nGot:\n%s", path, expectedContent, string(content))
	}
}

func writeTemplateFile(t *testing.T, path string, content string) {
	parentDir := filepath.Dir(path)
	err := os.RemoveAll(parentDir)
	if err != nil {
		t.Fatalf("Failed to remove existing directory: %v", err)
	}
	// make parent directory
	err = os.Mkdir(parentDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create parent directory: %v", err)
	}
	// write templateContent to file
	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}
}

func writeParseTestClean(t *testing.T, templateContent string, expectedFiles []string, expectedContent []string, flat bool) {
	templatePath := "test_output/template"

	writeTemplateFile(t, templatePath, templateContent)

	// Parse template
	err := template(templatePath, flat)
	if err != nil {
		t.Fatalf("Error processing template: %v", err)
	}
	// Check generated files
	for i, file := range expectedFiles {
		fileExistsAndContent(t, file, expectedContent[i])
	}

	// Cleanup
	err = os.RemoveAll("test_output")
	if err != nil {
		t.Fatalf("Failed to clean up generated files. Error: %v", err)
	}
}

func TestFormat(t *testing.T) {
	templateContent := `x:="{array=[1];format=%02.0f}"`
	expectedFiles := []string{
		"test_output/01.mx3",
	}
	expectedContent := []string{
		`x:=1`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}
func TestFormat2(t *testing.T) {
	templateContent := `x:="{array=[1];format=%.1f}"`
	expectedFiles := []string{
		"test_output/1.0.mx3",
	}
	expectedContent := []string{
		`x:=1`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}
func TestFormat3(t *testing.T) {
	templateContent := `x:="{array=[1];format=%.1f}"`
	expectedFiles := []string{
		"test_output/1.0.mx3",
	}
	expectedContent := []string{
		`x:=1`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}

// Test case for simple arange expression
func TestGenerateFilesWithArange(t *testing.T) {
	templateContent := `x:="{start=0;end=1;step=1}"`
	expectedFiles := []string{
		"test_output/0.mx3",
		"test_output/1.mx3",
	}
	expectedContent := []string{
		`x:=0`,
		`x:=1`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}
func TestGenerateFilesWithArange2(t *testing.T) {
	templateContent := `x:="{prefix=qwe;start=0;end=1;step=1}"`
	expectedFiles := []string{
		"test_output/qwe0.mx3",
		"test_output/qwe1.mx3",
	}
	expectedContent := []string{
		`x:=0`,
		`x:=1`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}
func TestGenerateFilesWithArange3(t *testing.T) {
	templateContent := `x:="{prefix=qwe;start=0;end=5;step=1}"`
	expectedFiles := []string{
		"test_output/qwe0.mx3",
		"test_output/qwe1.mx3",
		"test_output/qwe2.mx3",
		"test_output/qwe3.mx3",
		"test_output/qwe4.mx3",
	}
	expectedContent := []string{
		`x:=0`,
		`x:=1`,
		`x:=2`,
		`x:=3`,
		`x:=4`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}

func TestGenerateFilesWithArangePrefixAndFormat(t *testing.T) {
	templateContent := `x:="{prefix=qwe;start=0;end=2;step=1;format=%02.0f}"`
	expectedFiles := []string{
		"test_output/qwe00.mx3",
		"test_output/qwe01.mx3",
		"test_output/qwe02.mx3",
	}
	expectedContent := []string{
		`x:=0`,
		`x:=1`,
		`x:=2`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}

// Test case for linspace (start, end, count)
func TestGenerateFilesWithLinspace(t *testing.T) {
	templateContent := `x:="{prefix=test;start=0;end=2;count=3}"`
	expectedFiles := []string{
		"test_output/test0.mx3",
		"test_output/test1.mx3",
		"test_output/test2.mx3",
	}
	expectedContent := []string{
		`x:=0`,
		`x:=1`,
		`x:=2`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}

// Test case for array
func TestGenerateFilesWithArray(t *testing.T) {
	templateContent := `x:="{prefix=array_test;array=[3.14, 2.71, 1.41]}"`
	expectedFiles := []string{
		"test_output/array_test3.14.mx3",
		"test_output/array_test2.71.mx3",
		"test_output/array_test1.41.mx3",
	}
	expectedContent := []string{
		`x:=3.14`,
		`x:=2.71`,
		`x:=1.41`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}

// Test case for array with formatting
func TestGenerateFilesWithArrayAndFormat(t *testing.T) {
	templateContent := `x:="{prefix=array_fmt;array=[10, 20, 30];format=%03.0f}"`
	expectedFiles := []string{
		"test_output/array_fmt010.mx3",
		"test_output/array_fmt020.mx3",
		"test_output/array_fmt030.mx3",
	}
	expectedContent := []string{
		`x:=10`,
		`x:=20`,
		`x:=30`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}

// Test case for step and format together
func TestGenerateFilesWithStepAndFormat(t *testing.T) {
	templateContent := `x:="{prefix=step_fmt;start=0;end=4;step=2;format=%04.0f}"`
	expectedFiles := []string{
		"test_output/step_fmt0000.mx3",
		"test_output/step_fmt0002.mx3",
		"test_output/step_fmt0004.mx3",
	}
	expectedContent := []string{
		`x:=0`,
		`x:=2`,
		`x:=4`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}

// Test case for both linspace and step together
func TestGenerateFilesWithLinspaceAndStep(t *testing.T) {
	templateContent := `x:="{prefix=lin_step;start=0;end=4;count=3}"`
	expectedFiles := []string{
		"test_output/lin_step0.mx3",
		"test_output/lin_step2.mx3",
		"test_output/lin_step4.mx3",
	}
	expectedContent := []string{
		`x:=0`,
		`x:=2`,
		`x:=4`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}

// format is %d, it should not be allowed
func TestFormatError(t *testing.T) {
	templateContent := `x:="{array=[1];format=%d}"`
	templatePath := "test_output/template"
	writeTemplateFile(t, templatePath, templateContent)

	// Parse template
	err := template(templatePath, false)
	if err.Error() != `error finding expressions: only the %f format is allowed: %d` {
		t.Fatalf("Expected error: %v", err)
	}
	// check if no file is generated
	_, err = os.ReadFile("test_output/1.mx3")
	if err == nil {
		t.Fatalf("Expected file %s to not exist, but it does", "test_output/1.mx3")
	}
	// Cleanup
	err = os.RemoveAll("test_output")
	if err != nil {
		t.Fatalf("Failed to clean up generated files. Error: %v", err)
	}
}

func TestTwoTemplateStrings(t *testing.T) {
	templateContent := `x:="{array=[1,2];format=%02.0f}"
y:="{array=[3,4];format=%.0f}"`

	expectedFiles := []string{
		"test_output/01/3.mx3",
		"test_output/01/4.mx3",
		"test_output/02/3.mx3",
		"test_output/02/4.mx3",
	}

	expectedContent := []string{
		"x:=1\ny:=3",
		"x:=1\ny:=4",
		"x:=2\ny:=3",
		"x:=2\ny:=4",
	}

	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}

func TestFlat(t *testing.T) {
	templateContent := `x:="{array=[1,2];format=%02.0f}"
y:="{array=[3,4];format=%.0f}"`

	expectedFiles := []string{
		"test_output/013.mx3",
		"test_output/014.mx3",
		"test_output/023.mx3",
		"test_output/024.mx3",
	}

	expectedContent := []string{
		"x:=1\ny:=3",
		"x:=1\ny:=4",
		"x:=2\ny:=3",
		"x:=2\ny:=4",
	}

	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, true)
}
func TestFlat2(t *testing.T) {
	templateContent := `x:="{array=[0,1]}"
y:="{array=[0,1]}"
z:="{array=[0,1]}"`

	expectedFiles := []string{
		"test_output/000.mx3",
		"test_output/001.mx3",
		"test_output/010.mx3",
		"test_output/011.mx3",
		"test_output/100.mx3",
		"test_output/101.mx3",
		"test_output/110.mx3",
		"test_output/111.mx3",
	}

	expectedContent := []string{
		"x:=0\ny:=0\nz:=0",
		"x:=0\ny:=0\nz:=1",
		"x:=0\ny:=1\nz:=0",
		"x:=0\ny:=1\nz:=1",
		"x:=1\ny:=0\nz:=0",
		"x:=1\ny:=0\nz:=1",
		"x:=1\ny:=1\nz:=0",
		"x:=1\ny:=1\nz:=1",
	}

	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, true)
}

// Test case for an empty array
func TestEmptyArray(t *testing.T) {
	templateContent := `x:="{array=[]}"`
	templatePath := "test_output/template"
	writeTemplateFile(t, templatePath, templateContent)

	err := template(templatePath, false)
	if err == nil || err.Error() != `error finding expressions: array cannot be empty` {
		t.Fatalf("Expected error for empty array, but got: %v", err)
	}
}

// Test case where start is greater than end
func TestInvalidStartEnd(t *testing.T) {
	templateContent := `x:="{start=5;end=1;step=1}"`
	templatePath := "test_output/template"
	writeTemplateFile(t, templatePath, templateContent)

	err := template(templatePath, false)
	if err == nil || err.Error() != `error finding expressions: start value should be less than end value` {
		t.Fatalf("Expected error for start > end, but got: %v", err)
	}
}

// Test case where step is zero
func TestZeroStep(t *testing.T) {
	templateContent := `x:="{start=0;end=1;step=0}"`
	templatePath := "test_output/template"
	writeTemplateFile(t, templatePath, templateContent)

	err := template(templatePath, false)
	if err == nil || err.Error() != `error finding expressions: step value should be greater than 0` {
		t.Fatalf("Expected error for zero step, but got: %v", err)
	}
}

// Test case with invalid format string
func TestInvalidFormat(t *testing.T) {
	templateContent := `x:="{array=[1];format=%q}"`
	templatePath := "test_output/template"
	writeTemplateFile(t, templatePath, templateContent)

	err := template(templatePath, false)
	if err == nil || err.Error() != `error finding expressions: only the %f format is allowed: %q` {
		t.Fatalf("Expected error for invalid format, but got: %v", err)
	}
}

// Test case with missing required fields
func TestMissingFields(t *testing.T) {
	templateContent := `x:="{end=5;step=1}"`
	templatePath := "test_output/template"
	writeTemplateFile(t, templatePath, templateContent)

	err := template(templatePath, false)
	if err == nil || err.Error() != `error finding expressions: start should be given when array is not given` {
		t.Fatalf("Expected error for missing start field, but got: %v", err)
	}
}

// Test case with conflicting fields (both array and start provided)
func TestConflictingFields(t *testing.T) {
	templateContent := `x:="{array=[1,2];start=0;end=1;step=1}"`
	templatePath := "test_output/template"
	writeTemplateFile(t, templatePath, templateContent)

	err := template(templatePath, false)
	if err == nil || err.Error() != `error finding expressions: start should not be given when array is given` {
		t.Fatalf("Expected error for conflicting fields, but got: %v", err)
	}
}

// Test case with unexpected tokens in the expression
func TestUnexpectedTokens(t *testing.T) {
	templateContent := `x:="{array=[1,2];unknown=5}"`
	templatePath := "test_output/template"
	writeTemplateFile(t, templatePath, templateContent)

	err := template(templatePath, false)
	if err == nil || err.Error() != `error finding expressions: invalid field name: unknown` {
		t.Fatalf("Expected error for unexpected tokens, but got: %v", err)
	}
}

// Test case with malformed expression
func TestMalformedExpression(t *testing.T) {
	templateContent := `x:="{array=[1,2]`
	templatePath := "test_output/template"
	writeTemplateFile(t, templatePath, templateContent)

	err := template(templatePath, false)
	if err == nil || !strings.Contains(err.Error(), `no expressions found`) {
		t.Fatalf("Expected error for malformed expression, but got: %v", err)
	}
}

// Test case with large ranges
func TestLargeRange(t *testing.T) {
	templateContent := `x:="{start=1;end=10000;step=1000}"`
	expectedFiles := []string{
		"test_output/1.mx3",
		"test_output/1001.mx3",
		"test_output/2001.mx3",
		"test_output/3001.mx3",
		"test_output/4001.mx3",
		"test_output/5001.mx3",
		"test_output/6001.mx3",
		"test_output/7001.mx3",
		"test_output/8001.mx3",
		"test_output/9001.mx3",
	}
	var expectedContent []string
	for i := 1; i <= 9001; i += 1000 {
		expectedContent = append(expectedContent, fmt.Sprintf("x:=%d", i))
	}

	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}

// Test case with special characters in prefix and suffix
func TestSpecialCharactersInPrefixSuffix(t *testing.T) {
	templateContent := `x:="{prefix=val_;suffix=_test;array=[1,2]}"`
	expectedFiles := []string{
		"test_output/val_1_test.mx3",
		"test_output/val_2_test.mx3",
	}
	expectedContent := []string{
		`x:=1`,
		`x:=2`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}

// Test case with non-numeric values in array
func TestNonNumericArrayValues(t *testing.T) {
	templateContent := `x:="{array=[a,b,c]}"`
	templatePath := "test_output/template"
	writeTemplateFile(t, templatePath, templateContent)

	err := template(templatePath, false)
	if err == nil || !strings.Contains(err.Error(), `invalid array value`) {
		t.Fatalf("Expected error for non-numeric array values, but got: %v", err)
	}
}

// Test case with only one value in linspace
func TestLinspaceSingleValue(t *testing.T) {
	templateContent := `x:="{start=5;end=5;count=1}"`
	expectedFiles := []string{
		"test_output/5.mx3",
	}
	expectedContent := []string{
		`x:=5`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}

// Test case where count is zero in linspace
func TestZeroCountLinspace(t *testing.T) {
	templateContent := `x:="{start=0;end=1;count=0}"`
	templatePath := "test_output/template"
	writeTemplateFile(t, templatePath, templateContent)

	err := template(templatePath, false)
	if err == nil || !strings.Contains(err.Error(), `error finding expressions: count value should be greater than 0`) {
		t.Fatalf("Expected error for zero count in linspace, but got: %v", err)
	}
}

// Test case with floating-point step value
func TestFloatingPointStep(t *testing.T) {
	templateContent := `x:="{start=0;end=1;step=0.5}"`
	expectedFiles := []string{
		"test_output/0.mx3",
		"test_output/0.5.mx3",
		"test_output/1.mx3",
	}
	expectedContent := []string{
		`x:=0`,
		`x:=0.5`,
		`x:=1`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}

// Test case with multiple variables and complex combinations
func TestComplexCombinations(t *testing.T) {
	templateContent := `a:="{array=[1,2]}"
b:="{start=0;end=1;step=1}"
c:="{array=[5];format=%02.0f}"`
	expectedFiles := []string{
		"test_output/1/0/05.mx3",
		"test_output/1/1/05.mx3",
		"test_output/2/0/05.mx3",
		"test_output/2/1/05.mx3",
	}
	expectedContent := []string{
		`a:=1
b:=0
c:=5`,
		`a:=1
b:=1
c:=5`,
		`a:=2
b:=0
c:=5`,
		`a:=2
b:=1
c:=5`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent, false)
}
