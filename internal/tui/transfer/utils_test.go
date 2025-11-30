/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package transfer

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestExpandPathAbsolute(t *testing.T) {
	absPath := filepath.Join(os.TempDir(), "testfile.txt")
	result, err := ExpandPath(absPath)
	if err != nil {
		t.Fatalf("ExpandPath() error = %v", err)
	}
	if result != absPath {
		t.Errorf("ExpandPath() = %q, want %q", result, absPath)
	}
}

func TestExpandPathRelative(t *testing.T) {
	result, err := ExpandPath("relative/path/file.txt")
	if err != nil {
		t.Fatalf("ExpandPath() error = %v", err)
	}
	if !filepath.IsAbs(result) {
		t.Errorf("ExpandPath() = %q, should be absolute", result)
	}
}

func TestExpandPathWithSpaces(t *testing.T) {
	path := "   /some/path   "
	result, err := ExpandPath(path)
	if err != nil {
		t.Fatalf("ExpandPath() error = %v", err)
	}
	if !filepath.IsAbs(result) {
		t.Errorf("ExpandPath() = %q, should be absolute", result)
	}
}

func TestExpandPathWindowsStyle(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	result, err := ExpandPath("C:\\Users\\Test")
	if err != nil {
		t.Fatalf("ExpandPath() error = %v", err)
	}
	if !filepath.IsAbs(result) {
		t.Errorf("ExpandPath() = %q, should be absolute", result)
	}
}

func TestExpandPathUnixStyle(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-specific test")
	}

	result, err := ExpandPath("/home/user/file.txt")
	if err != nil {
		t.Fatalf("ExpandPath() error = %v", err)
	}
	if result != "/home/user/file.txt" {
		t.Errorf("ExpandPath() = %q, want '/home/user/file.txt'", result)
	}
}

func TestExpandPathCurrentDir(t *testing.T) {
	result, err := ExpandPath(".")
	if err != nil {
		t.Fatalf("ExpandPath() error = %v", err)
	}

	cwd, _ := os.Getwd()
	if result != cwd {
		t.Errorf("ExpandPath('.') = %q, want %q", result, cwd)
	}
}

func TestExpandPathParentDir(t *testing.T) {
	result, err := ExpandPath("..")
	if err != nil {
		t.Fatalf("ExpandPath() error = %v", err)
	}

	if !filepath.IsAbs(result) {
		t.Errorf("ExpandPath('..') = %q, should be absolute", result)
	}
}

func TestExpandPathComplexRelative(t *testing.T) {
	result, err := ExpandPath("./subdir/../other/file.txt")
	if err != nil {
		t.Fatalf("ExpandPath() error = %v", err)
	}

	if !filepath.IsAbs(result) {
		t.Errorf("ExpandPath() = %q, should be absolute", result)
	}
}
