package handler

import "strings"

func isAllowedExtension(filename string, allowed []string) bool {
    lower := strings.ToLower(filename)
    for _, ext := range allowed {
        if strings.HasSuffix(lower, ext) {
            return true
        }
    }
    return false
}
