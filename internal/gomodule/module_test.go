package gomodule

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestFindModuleRoot(t *testing.T) {
	tests := []struct {
		name        string
		fs          fstest.MapFS
		startDir    string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name: "go.mod in current directory",
			fs: fstest.MapFS{
				"project/go.mod": &fstest.MapFile{Mode: 0644},
			},
			startDir: "project",
			want:     "project",
			wantErr:  false,
		},
		{
			name: "go.mod in parent directory",
			fs: fstest.MapFS{
				"project/go.mod":          &fstest.MapFile{Mode: 0644},
				"project/cmd/app/main.go": &fstest.MapFile{Mode: 0644},
			},
			startDir: "project/cmd/app",
			want:     "project",
			wantErr:  false,
		},
		{
			name: "go.mod in root directory",
			fs: fstest.MapFS{
				"go.mod":                  &fstest.MapFile{Mode: 0644},
				"project/cmd/app/main.go": &fstest.MapFile{Mode: 0644},
			},
			startDir: "project/cmd/app",
			want:     ".",
			wantErr:  false,
		},
		{
			name: "go.mod not found",
			fs: fstest.MapFS{
				"project/main.go": &fstest.MapFile{Mode: 0644},
			},
			startDir:    "project",
			wantErr:     true,
			errContains: "go.mod not found",
		},
		{
			name:        "empty directory",
			fs:          fstest.MapFS{},
			startDir:    "empty",
			wantErr:     true,
			errContains: "go.mod not found",
		},
		{
			name:        "non-existent directory",
			fs:          fstest.MapFS{},
			startDir:    "/nonexistent",
			wantErr:     true,
			errContains: "go.mod not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем временную файловую систему
			tmpDir := t.TempDir()

			// Создаем структуру файлов согласно тесту
			for path, file := range tt.fs {
				fullPath := filepath.Join(tmpDir, path)
				_ = os.MkdirAll(filepath.Dir(fullPath), 0755)
				if file.Mode.IsRegular() {
					err := os.WriteFile(fullPath, []byte("module test\n"), 0644)
					if err != nil {
						return
					}
				}
			}

			startPath := filepath.Join(tmpDir, tt.startDir)
			_ = os.MkdirAll(startPath, 0755)

			got, err := FindModuleRoot(startPath)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("Expected error, got nil")
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Fatalf("Error message %q should contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			expectedPath := filepath.Join(tmpDir, tt.want)
			if got != expectedPath {
				t.Fatalf("Expected %q, got %q", expectedPath, got)
			}
		})
	}
}

func TestFindModuleRoot_EdgeCases(t *testing.T) {
	t.Run("root directory without go.mod", func(t *testing.T) {
		tmpDir := t.TempDir()

		_, err := FindModuleRoot(tmpDir)
		if err == nil {
			t.Fatal("Expected error for root directory without go.mod")
		}
	})

	t.Run("empty string path", func(t *testing.T) {
		_, err := FindModuleRoot("")
		if err == nil {
			t.Fatal("Expected error for empty path")
		}
	})

	t.Run("file instead of directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		_ = os.WriteFile(filePath, []byte("test"), 0644)

		_, err := FindModuleRoot(filePath)
		if err == nil {
			t.Fatal("Expected error for file path")
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

// Benchmark тест
func BenchmarkFindModuleRoot(b *testing.B) {
	tmpDir := b.TempDir()

	// Создаем глубокую структуру с go.mod в корне
	deepDir := filepath.Join(tmpDir, "level1", "level2", "level3", "level4")
	_ = os.MkdirAll(deepDir, 0755)
	_ = os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test\n"), 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindModuleRoot(deepDir)
	}
}
