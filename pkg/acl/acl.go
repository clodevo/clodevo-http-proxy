package acl

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/clodevo/raven-proxy/pkg/utils"
	"github.com/valyala/fasthttp"
)

// Assuming LogLevelNone, LogLevelDebug, LogLevelTrace are defined in utils package
// If not, you'll need to adjust the log level handling accordingly.

type List struct {
	Whitelist []string `json:"Whitelist"`
	Blacklist []string `json:"Blacklist"`
}

type ACLManager struct {
	TenantLists      map[string]*List
	compiledPatterns map[string]*regexp.Regexp
	aclDataPath      string
	logger           *utils.Logger
}

func NewACLManager(aclDataPath string, logger *utils.Logger) *ACLManager {
	return &ACLManager{
		TenantLists:      make(map[string]*List),
		compiledPatterns: make(map[string]*regexp.Regexp),
		aclDataPath:      aclDataPath,
		logger:           logger,
	}
}

func (a *ACLManager) LoadTenantLists(tenantName string) *List {
	list := &List{}
	filePath := filepath.Join(a.aclDataPath, tenantName+".json")
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		a.logger.Debug("Error reading list file for tenant %s: %v", tenantName, err)
		return list
	}

	err = json.Unmarshal(fileContent, list)
	if err != nil {
		a.logger.Debug("Error parsing list file for tenant %s: %v", tenantName, err)
		return list
	}

	a.TenantLists[tenantName] = list
	a.logger.Trace("Loaded ACL list for tenant: %s", tenantName)
	return list
}

func (a *ACLManager) IsRequestAllowed(ctx *fasthttp.RequestCtx, tenantName string) bool {
	hostWithPort := string(ctx.Host())
	host, port, _ := net.SplitHostPort(hostWithPort)
	if host == "" {
		host = hostWithPort
	}

	if list, exists := a.TenantLists[tenantName]; exists {
		for _, b := range list.Blacklist {
			if a.matchesPattern(host, port, b) {
				a.logger.Debug("Request to %s blocked by blacklist rule: %s", hostWithPort, b)
				return false
			}
		}
		for _, w := range list.Whitelist {
			if a.matchesPattern(host, port, w) {
				a.logger.Debug("Request to %s allowed by whitelist rule: %s", hostWithPort, w)
				return true
			}
		}
		a.logger.Trace("Evaluating request to %s against ACL rules", hostWithPort)
		return false
	}
	a.logger.Trace("No ACL rules defined for tenant %s, defaulting to block", tenantName)
	return false
}

func (a *ACLManager) matchesPattern(host, port, pattern string) bool {
	patternHost, patternPort, _ := net.SplitHostPort(pattern)
	if patternHost == "" {
		patternHost = pattern
	}

	if patternPort == "" || patternPort == port {
		regex := a.compilePattern(patternHost)
		match := regex.MatchString(host)
		a.logger.Trace("Matching host %s against pattern %s: %t", host, regex.String(), match)
		return match
	}
	return false
}

func (a *ACLManager) wildcardToRegex(pattern string) string {
	pattern = strings.Replace(pattern, "*", ".*", -1)
	pattern = strings.Replace(pattern, ".", "\\.", -1)  // Escape actual dots for regex
	pattern = strings.Replace(pattern, "\\.*", ".*", 1) // Replace the first occurrence of '\.*' with '.*'

	// Make the regex case-insensitive and ensure it matches the entire host
	return "(?i)^" + pattern + "$"
}

func (a *ACLManager) compilePattern(pattern string) *regexp.Regexp {
	if compiled, exists := a.compiledPatterns[pattern]; exists {
		return compiled
	}
	regexPattern := a.wildcardToRegex(pattern)
	regex := regexp.MustCompile(regexPattern)
	a.compiledPatterns[pattern] = regex
	return regex
}
