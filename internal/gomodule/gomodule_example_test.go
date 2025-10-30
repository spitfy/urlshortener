package gomodule

import (
	"fmt"
	"os"
	"path/filepath"
)

func ExampleFindModuleRoot() {
	tmpDir, err := os.MkdirTemp("", "example")
	if err != nil {
		fmt.Printf("Error creating temp dir: %v\n", err)
		return
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(tmpDir)

	projectDir := filepath.Join(tmpDir, "myproject", "cmd", "app")
	_ = os.MkdirAll(projectDir, 0755)

	goModPath := filepath.Join(tmpDir, "myproject", "go.mod")
	_ = os.WriteFile(goModPath, []byte("module myproject\n"), 0644)

	result, err := FindModuleRoot(projectDir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Module root: %s\n", filepath.Base(result))
	}
	// Output: Module root: myproject
}

func ExampleFindModuleRoot_rootDirectory() {
	tmpDir, err := os.MkdirTemp("", "example")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(tmpDir)

	goModPath := filepath.Join(tmpDir, "go.mod")
	_ = os.WriteFile(goModPath, []byte("module example\n"), 0644)

	_, err = FindModuleRoot(tmpDir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Found module root successfully")
	}
	// Output:
	// Found module root successfully
}

func ExampleFindModuleRoot_errorCases() {
	// Несуществующая директория
	_, err := FindModuleRoot("/nonexistent/path")
	fmt.Printf("Error for nonexistent path: %v\n", err)

	// Пустая строка
	_, err = FindModuleRoot("")
	fmt.Printf("Error for empty path: %v\n", err)

	// Output:
	// Error for nonexistent path: go.mod not found in any parent directory
	// Error for empty path: go.mod not found in any parent directory
}
