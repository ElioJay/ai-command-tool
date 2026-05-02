package safety

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

//go:embed defaults.toml
var defaultsToml []byte

type rulesFile struct {
	Rules []Rule `toml:"rules"`
}

type Blacklist struct {
	rules    []Rule
	compiled []*regexp.Regexp
	userPath string
}

func NewBlacklist(configDir string) (*Blacklist, error) {
	bl := &Blacklist{}
	if configDir != "" {
		bl.userPath = filepath.Join(configDir, "blacklist.toml")
	}

	var defaults rulesFile
	if err := toml.Unmarshal(defaultsToml, &defaults); err != nil {
		return nil, fmt.Errorf("加载内置黑名单失败: %w", err)
	}
	for _, r := range defaults.Rules {
		if err := bl.addRule(r); err != nil {
			return nil, err
		}
	}

	if bl.userPath != "" {
		if err := bl.loadUserRules(); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	return bl, nil
}

func (bl *Blacklist) addRule(r Rule) error {
	re, err := regexp.Compile(r.Pattern)
	if err != nil {
		return fmt.Errorf("规则 %q 正则无效: %w", r.ID, err)
	}
	bl.rules = append(bl.rules, r)
	bl.compiled = append(bl.compiled, re)
	return nil
}

func (bl *Blacklist) loadUserRules() error {
	data, err := os.ReadFile(bl.userPath)
	if err != nil {
		return err
	}
	var uf rulesFile
	if err := toml.Unmarshal(data, &uf); err != nil {
		return err
	}
	for _, r := range uf.Rules {
		r.Source = "user"
		if err := bl.addRule(r); err != nil {
			return err
		}
	}
	return nil
}

func (bl *Blacklist) Match(cmd string) (*Rule, bool) {
	cmd = strings.TrimSpace(cmd)
	for i, re := range bl.compiled {
		if re.MatchString(cmd) {
			r := bl.rules[i]
			return &r, true
		}
	}
	return nil, false
}

func (bl *Blacklist) Add(r Rule) error {
	r.Source = "user"
	if err := bl.addRule(r); err != nil {
		return err
	}
	return bl.saveUserRules()
}

func (bl *Blacklist) saveUserRules() error {
	var userRules []Rule
	for _, r := range bl.rules {
		if r.Source == "user" {
			userRules = append(userRules, r)
		}
	}
	f, err := os.OpenFile(bl.userPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(rulesFile{Rules: userRules})
}

func (bl *Blacklist) List() []Rule { return bl.rules }
