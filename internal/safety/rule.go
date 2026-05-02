package safety

type Rule struct {
	ID      string `toml:"id"`
	Pattern string `toml:"pattern"`
	Reason  string `toml:"reason"`
	Source  string `toml:"source"`
}
