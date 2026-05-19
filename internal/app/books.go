package app

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const maxBookRecords = 20

// BookRecord 表示一个最近打开的书籍工作目录。
type BookRecord struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	Author       string `json:"author"`
	LastOpenedAt string `json:"last_opened_at"`
}

type bookRegistryData struct {
	Current string       `json:"current"`
	Books   []BookRecord `json:"books"`
}

// BookRegistry 持久化最近打开的书籍工作目录。
type BookRegistry struct {
	path       string
	legacyPath string
}

// NewBookRegistry 创建书籍记录管理器。
func NewBookRegistry(novaDir string) *BookRegistry {
	return &BookRegistry{
		path:       filepath.Join(novaDir, "books.json"),
		legacyPath: legacyBookRegistryPath(),
	}
}

// Current 返回上次打开且仍存在的工作目录。
func (r *BookRegistry) Current() string {
	data := r.load()
	if data.Current == "" {
		return ""
	}
	if info, err := os.Stat(data.Current); err == nil && info.IsDir() {
		return data.Current
	}
	return ""
}

// List 返回最近打开的书籍列表。
func (r *BookRegistry) List() []BookRecord {
	data := r.load()
	books := make([]BookRecord, 0, len(data.Books))
	for _, book := range data.Books {
		if book.Path == "" {
			continue
		}
		books = append(books, book)
	}
	sort.SliceStable(books, func(i, j int) bool {
		return books[i].LastOpenedAt > books[j].LastOpenedAt
	})
	return books
}

// Touch 记录并置顶一个书籍工作目录。
func (r *BookRegistry) Touch(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("路径不是目录")
	}

	data := r.load()
	now := time.Now().Format(time.RFC3339)
	record := BookRecord{
		Name:         filepath.Base(absPath),
		Path:         absPath,
		LastOpenedAt: now,
	}
	books := []BookRecord{record}
	for _, book := range data.Books {
		if book.Path == "" || book.Path == absPath {
			continue
		}
		books = append(books, book)
		if len(books) >= maxBookRecords {
			break
		}
	}
	data.Current = absPath
	data.Books = books
	return r.save(data)
}

// Remove 移除一个书籍记录，不删除磁盘文件。
func (r *BookRegistry) Remove(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	data := r.load()
	books := make([]BookRecord, 0, len(data.Books))
	for _, book := range data.Books {
		if book.Path != absPath {
			books = append(books, book)
		}
	}
	if data.Current == absPath {
		data.Current = ""
		if len(books) > 0 {
			data.Current = books[0].Path
		}
	}
	data.Books = books
	return r.save(data)
}

func (r *BookRegistry) load() bookRegistryData {
	var data bookRegistryData
	raw, err := os.ReadFile(r.path)
	if err != nil && r.legacyPath != "" && r.legacyPath != r.path {
		raw, err = os.ReadFile(r.legacyPath)
	}
	if err != nil {
		return data
	}
	_ = json.Unmarshal(raw, &data)
	return data
}

func (r *BookRegistry) save(data bookRegistryData) error {
	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, raw, 0o644)
}

func legacyBookRegistryPath() string {
	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Join(dir, "nova", "books.json")
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".nova", "books.json")
	}
	return filepath.Join(".", ".nova-books.json")
}
