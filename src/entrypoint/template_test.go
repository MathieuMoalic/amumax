package entrypoint

import (
	"os"
	"path/filepath"
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

func writeParseTestClean(t *testing.T, templateContent string, expectedFiles []string, expectedContent []string) {
	templatePath := "test_output/template"

	writeTemplateFile(t, templatePath, templateContent)

	// Parse template
	err := template(templatePath)
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
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent)
}
func TestFormat2(t *testing.T) {
	templateContent := `x:="{array=[1];format=%.1f}"`
	expectedFiles := []string{
		"test_output/1.0.mx3",
	}
	expectedContent := []string{
		`x:=1`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent)
}
func TestFormat3(t *testing.T) {
	templateContent := `x:="{array=[1];format=%.1f}"`
	expectedFiles := []string{
		"test_output/1.0.mx3",
	}
	expectedContent := []string{
		`x:=1`,
	}
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent)
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
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent)
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
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent)
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
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent)
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
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent)
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
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent)
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
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent)
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
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent)
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
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent)
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
	writeParseTestClean(t, templateContent, expectedFiles, expectedContent)
}

// format is %d, it should not be allowed
func TestFormatError(t *testing.T) {
	templateContent := `x:="{array=[1];format=%d}"`
	templatePath := "test_output/template"
	writeTemplateFile(t, templatePath, templateContent)

	// Parse template
	err := template(templatePath)
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

	writeParseTestClean(t, templateContent, expectedFiles, expectedContent)
}
