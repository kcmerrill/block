package block

import "strings"

func isDirIgnored(name, osPathname string, ignored []string) bool {
	for _, ignoreDir := range ignored {
		if string(ignoreDir[0]) == "/" && strings.HasPrefix(osPathname, ignoreDir) {
			return true
		}
		if strings.Contains(osPathname, ignoreDir) {
			return true
		}
	}

	// hidden folders
	if len(name) >= 3 && string(name[0]) == "." {
		return true
	}

	return false
}
