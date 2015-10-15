package kmgSys

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"strings"
)

type IptableRule struct {
	Table string // example: "nat"
	Rule  string // example: "-A PREROUTING -s 172.20.0.0/16 -p udp -m udp --dport 53 -j REDIRECT --to-ports 53"
}

func MustSetIptableRule(rule IptableRule) {
	for _, thisRule := range MustGetIptableRuleList() {
		if thisRule.Table == rule.Table && thisRule.Rule == rule.Rule {
			return
		}
	}
	// Another app is currently holding the xtables lock. Perhaps you want to use the -w option?
	kmgCmd.MustRun("iptables -w -t " + rule.Table + " " + rule.Rule)
}

func MustGetIptableRuleList() []IptableRule {
	content := kmgCmd.MustCombinedOutput("iptables-save")
	return parseIptableSave(string(content))
}

func parseIptableSave(content string) []IptableRule {
	thisTable := ""
	output := []IptableRule{}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line[0] == '#' || line[0] == ':' {
			continue
		}
		if line[0] == '*' {
			thisTable = line[1:]
			continue
		}
		if line == "COMMIT" {
			continue
		}
		output = append(output, IptableRule{
			Table: thisTable,
			Rule:  line,
		})
	}
	return output
}
