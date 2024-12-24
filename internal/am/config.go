package am

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config manages configuration settings loaded from environment variables and CLI flags.
type Config struct {
	namespace string // Namespace prefix for environment variables, e.g., MWZ, MYCVS, APP, etc.
	values    map[string]string
	flags     map[string]string
}

// LoadCfg initializes a Config instance with the specified namespace and loads the corresponding environment variables and CLI flags.
func LoadCfg(namespace string, flagDefs map[string]interface{}) *Config {
	cfg := &Config{}
	cfg.SetNamespace(namespace)
	cfg.defineFlags(flagDefs)
	flag.Parse()
	cfg.loadNamespaceEnvVars()
	cfg.loadFlags()
	return cfg
}

// SetNamespace sets the namespace for the configuration, converting it to uppercase.
func (c *Config) SetNamespace(namespace string) {
	c.namespace = strings.ToUpper(namespace)
}

// namespacePrefix returns the namespace prefix used for environment variables.
func (c *Config) namespacePrefix() string {
	return fmt.Sprintf("%s_", c.namespace)
}

// SetValues sets the configuration values directly.
func (c *Config) SetValues(values map[string]string) {
	c.values = values
}

// Get retrieves all environment variables that belong to the namespace.
// If reload is true, it re-reads the values from the environment.
func (c *Config) Get(reload ...bool) map[string]string {
	if len(reload) > 0 && reload[0] {
		return c.get(true)
	}
	return c.get(false)
}

func (c *Config) get(reload bool) map[string]string {
	if reload || len(c.values) == 0 {
		c.values = c.readNamespaceEnvVars()
	}
	merged := make(map[string]string)
	for k, v := range c.values {
		merged[k] = v
	}
	for k, v := range c.flags {
		merged[k] = v
	}
	return merged
}

// Val retrieves the value of a specific namespaced environment variable or CLI flag.
// If reload is true, it re-reads the values from the environment and CLI flags.
func (c *Config) Val(key string, reload ...bool) (value string, ok bool) {
	vals := c.get(false)
	if len(reload) > 0 && reload[0] {
		vals = c.get(true)
	}
	val, ok := vals[key]
	return val, ok
}

// ValOrDef retrieves the value of a specific namespaced environment variable or CLI flag.
// If the key is not found, it returns the provided default value.
// If reload is true, it re-reads the values from the environment and CLI flags.
func (c *Config) ValOrDef(key string, defVal string, reload ...bool) (value string) {
	vals := c.get(false)
	if len(reload) > 0 && reload[0] {
		vals = c.get(true)
	}
	val, ok := vals[key]
	if !ok {
		val = defVal
	}
	return val
}

// ValAsInt retrieves the value of a specific namespaced environment variable or CLI flag as an int.
// If the key is not found or cannot be parsed as an int, it returns the provided default value.
// If reload is true, it re-reads the values from the environment and CLI flags.
func (c *Config) ValAsInt(key string, defVal int64, reload ...bool) (value int64) {
	vals := c.get(false)
	if len(reload) > 0 && reload[0] {
		vals = c.get(true)
	}
	val, ok := vals[key]
	if !ok {
		return defVal
	}
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defVal
	}
	return i
}

// ValAsFloat retrieves the value of a specific namespaced environment variable or CLI flag as a float.
// If the key is not found or cannot be parsed as a float, it returns the provided default value.
// If reload is true, it re-reads the values from the environment and CLI flags.
func (c *Config) ValAsFloat(key string, defVal float64, reload ...bool) (value float64) {
	vals := c.get(false)
	if len(reload) > 0 && reload[0] {
		vals = c.get(true)
	}
	val, ok := vals[key]
	if !ok {
		return defVal
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defVal
	}
	return f
}

// ValAsBool retrieves the value of a specific namespaced environment variable or CLI flag as a bool.
// If the key is not found or cannot be parsed as a bool, it returns the provided default value.
// If reload is true, it re-reads the values from the environment and CLI flags.
func (c *Config) ValAsBool(key string, defVal bool, reload ...bool) (value bool) {
	vals := c.get(false)
	if len(reload) > 0 && reload[0] {
		vals = c.get(true)
	}
	val, ok := vals[key]
	if !ok {
		return defVal
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defVal
	}
	return b
}

// loadNamespaceEnvVars loads all visible environment variables that belong to the namespace.
func (c *Config) loadNamespaceEnvVars() {
	c.values = c.readNamespaceEnvVars()
}

// readNamespaceEnvVars reads all visible environment variables that belong to the namespace.
func (c *Config) readNamespaceEnvVars() map[string]string {
	nevs := make(map[string]string)
	np := c.namespacePrefix()

	for _, ev := range os.Environ() {
		if strings.HasPrefix(ev, np) {
			varval := strings.SplitN(ev, "=", 2)

			if len(varval) < 2 {
				continue
			}

			key := c.keyify(varval[0])
			nevs[key] = varval[1]
		}
	}

	return nevs
}

// keyify converts environment variable names to a dot-separated, lowercase format.
// For example, NAMESPACE_CONFIG_VALUE becomes config.value.
func (c *Config) keyify(name string) string {
	split := strings.Split(name, "_")
	if len(split) < 1 {
		return ""
	}
	// Remove namespace prefix
	wnsp := strings.Join(split[1:], ".")
	// Convert to dot-separated lowercase
	dots := strings.ToLower(strings.Replace(wnsp, "_", ".", -1))
	return dots
}

// getEnvOrDef returns the value of an environment variable or a default value if the variable is not set or is empty.
// If no default value is provided, it returns an empty string.
func getEnvOrDef(envar string, def ...string) string {
	val := os.Getenv(envar)
	if val != "" {
		return val
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// loadFlags parses CLI flags and stores them in the Config struct.
func (c *Config) loadFlags() {
	c.flags = make(map[string]string)
	flag.VisitAll(func(f *flag.Flag) {
		c.flags[f.Name] = f.Value.String()
	})
}

// defineFlags defines the CLI flags used by the application.
func (c *Config) defineFlags(flagDefs map[string]interface{}) {
	for name, defVal := range flagDefs {
		switch v := defVal.(type) {
		case string:
			flag.String(name, v, "")
		case int:
			flag.Int(name, v, "")
		case bool:
			flag.Bool(name, v, "")
		}
	}
}
