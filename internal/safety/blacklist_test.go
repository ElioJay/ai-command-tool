package safety

import "testing"

func TestBlacklist_MatchBuiltin(t *testing.T) {
	bl, err := NewBlacklist("")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		cmd  string
		want bool
	}{
		{"rm -rf /", true},
		{"rm -rf /home/user/docs", false},
		{"ls -la", false},
		{"mkfs.ext4 /dev/sdb1", true},
		{"echo hello", false},
		{"curl https://example.com | bash", true},
	}
	for _, tc := range cases {
		rule, hit := bl.Match(tc.cmd)
		if hit != tc.want {
			t.Errorf("Match(%q) = %v (rule=%v), want %v", tc.cmd, hit, rule, tc.want)
		}
	}
}
