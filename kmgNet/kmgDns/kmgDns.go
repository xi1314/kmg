package kmgDns

import "strings"

// 例如 www.google.com -> []string{www.google.com,google.com,com}
func GetSupperDomainList(domain string) []string {
	output := []string{}
	parts := strings.Split(domain, ".")
	for i := 0; i < len(parts); i++ {
		output = append(output, strings.Join(parts[i:], "."))
	}
	return output
}
