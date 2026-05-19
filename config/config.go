package config

import (
	"fmt"
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

// Config 保存 Nova 的全局配置
type Config struct {
	OpenAIAPIKey        string `toml:"openai_api_key"`
	OpenAIBaseURL       string `toml:"openai_base_url"`
	OpenAIModel         string `toml:"openai_model"`
	SkillsDir           string `toml:"skills_dir"`
	NovaDir             string `toml:"nova_dir"`
	Workspace           string `toml:"workspace"`
	ResumeLastWorkspace bool   `toml:"-"`
}

// Load 加载配置，优先级：环境变量 > config.toml > 默认值
func Load() *Config {
	cfg := &Config{
		OpenAIModel:         "deepseek-v4-pro",
		NovaDir:             defaultNovaDir(),
		Workspace:           ".",
		ResumeLastWorkspace: true,
	}

	// 从 config.toml 加载（可执行文件同级目录）
	loadFromFile(cfg, "config.toml")

	// 环境变量覆盖
	overrideFromEnv(cfg)

	// 路径标准化
	if abs, err := filepath.Abs(cfg.Workspace); err == nil {
		cfg.Workspace = abs
	}
	if cfg.SkillsDir != "" {
		cfg.SkillsDir = normalizePath(cfg.SkillsDir)
	}
	if cfg.NovaDir == "" {
		cfg.NovaDir = defaultNovaDir()
	}
	cfg.NovaDir = normalizePath(cfg.NovaDir)
	return cfg
}

// loadFromFile 从 TOML 文件加载配置
func loadFromFile(cfg *Config, filename string) {
	// 依次尝试：当前目录、可执行文件同级目录
	candidates := []string{filename}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), filename))
	}

	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if err := toml.Unmarshal(data, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "警告: 解析 %s 失败: %v\n", path, err)
			continue
		}
		return
	}
}

// overrideFromEnv 用环境变量覆盖配置
func overrideFromEnv(cfg *Config) {
	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		cfg.OpenAIAPIKey = v
	}
	if v := os.Getenv("OPENAI_BASE_URL"); v != "" {
		cfg.OpenAIBaseURL = v
	}
	if v := os.Getenv("OPENAI_MODEL"); v != "" {
		cfg.OpenAIModel = v
	}
	if v := os.Getenv("NOVA_SKILLS_DIR"); v != "" {
		cfg.SkillsDir = v
	}
	if v := os.Getenv("NOVA_DIR"); v != "" {
		cfg.NovaDir = v
	}
	if v := os.Getenv("NOVA_WORKSPACE"); v != "" {
		cfg.Workspace = v
	}
}

func defaultNovaDir() string {
	return "~/.nova"
}

func normalizePath(path string) string {
	path = expandHome(path)
	if abs, err := filepath.Abs(path); err == nil {
		return abs
	}
	return path
}

func expandHome(path string) string {
	if path == "~" {
		if home, err := os.UserHomeDir(); err == nil && home != "" {
			return home
		}
		return path
	}
	if len(path) > 2 && path[:2] == "~/" {
		if home, err := os.UserHomeDir(); err == nil && home != "" {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}
