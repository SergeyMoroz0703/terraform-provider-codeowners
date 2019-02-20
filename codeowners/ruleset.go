package codeowners

import (
	"fmt"
	"strings"
)

type File struct {
	RepositoryName  string
	RepositoryOwner string
	Branch          string
	Ruleset         Ruleset
}

type Ruleset []Rule

type Rule struct {
	Pattern   string
	Usernames []string
}

func sameStringSlice(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[string]int, len(x))
	for _, _x := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range y {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y]--
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	if len(diff) == 0 {
		return true
	}
	return false
}

func (ruleset Ruleset) Compile() []byte {
	if ruleset == nil {
		return []byte{}
	}
	output := "# automatically generated by terraform - please do not edit here\n"
	for _, rule := range ruleset {
		usernames := ""
		for _, username := range rule.Usernames {
			if !strings.Contains(username, "@") {
				usernames = fmt.Sprintf("%s@%s ", usernames, username)
			} else {
				usernames = fmt.Sprintf("%s%s ", usernames, username)
			}
		}
		output = fmt.Sprintf("%s%s %s\n", output, rule.Pattern, usernames)
	}
	return []byte(output)
}

func parseRulesFile(data string) Ruleset {

	rules := []Rule{}
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		if trimmed[0] == '#' { // ignore comments
			continue
		}
		words := strings.Split(trimmed, " ")
		if len(words) < 2 {
			continue
		}
		rule := Rule{
			Pattern: words[0],
		}
		for _, username := range words[1:] {
			if len(username) == 0 { // may be split by multiple spaces
				continue
			}
			if username[0] == '@' {
				username = username[1:]
			}
			rule.Usernames = append(rule.Usernames, username)
		}
		rules = append(rules, rule)
	}

	return rules

}
