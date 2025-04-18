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

func NewConfig() *Config {
	return &Config{
		namespace: "AQM",
		values:    make(map[string]string),
		flags:     make(map[string]string),
	}
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
func (cfg *Config) SetNamespace(namespace string) {
	cfg.namespace = strings.ToUpper(namespace)
}

// namespacePrefix returns the namespace prefix used for environment variables.
func (cfg *Config) namespacePrefix() string {
	return fmt.Sprintf("%s_", cfg.namespace)
}

// SetValues sets the configuration values directly.
func (cfg *Config) SetValues(values map[string]string) {
	cfg.values = values
}

// Get retrieves all environment variables that belong to the namespace.
// If reload is true, it re-reads the values from the environment.
func (cfg *Config) Get(reload ...bool) map[string]string {
	if len(reload) > 0 && reload[0] {
		return cfg.get(true)
	}
	return cfg.get(false)
}

func (cfg *Config) get(reload bool) map[string]string {
	if reload || len(cfg.values) == 0 {
		cfg.values = cfg.readNamespaceEnvVars()
	}
	merged := make(map[string]string)
	for k, v := range cfg.values {
		merged[k] = v
	}
	for k, v := range cfg.flags {
		merged[k] = v
	}
	return merged
}

// ByteSliceVal retrieves the value of a specific namespaced environment variable or CLI flag as a byte slice.
// If the key is not found, it returns an empty byte slice.
// If reload is true, it re-reads the values from the environment and CLI flags.
func (cfg *Config) ByteSliceVal(key string, reload ...bool) []byte {
	val, ok := cfg.StrVal(key, reload...)
	if !ok {
		return []byte{}
	}
	return []byte(val)
}

// StrVal retrieves the value of a specific namespaced environment variable or CLI flag.
// If reload is true, it re-reads the values from the environment and CLI flags.
func (cfg *Config) StrVal(key string, reload ...bool) (value string, ok bool) {
	vals := cfg.get(false)
	if len(reload) > 0 && reload[0] {
		vals = cfg.get(true)
	}
	val, ok := vals[key]
	return val, ok
}

// StrValOrDef retrieves the value of a specific namespaced environment variable or CLI flag.
// If the key is not found, it returns the provided default value.
// If reload is true, it re-reads the values from the environment and CLI flags.
func (cfg *Config) StrValOrDef(key string, defVal string, reload ...bool) (value string) {
	vals := cfg.get(false)
	if len(reload) > 0 && reload[0] {
		vals = cfg.get(true)
	}
	val, ok := vals[key]
	if !ok {
		val = defVal
	}
	return val
}

// IntVal retrieves the value of a specific namespaced environment variable or CLI flag as an int.
// If the key is not found or cannot be parsed as an int, it returns the provided default value.
// If reload is true, it re-reads the values from the environment and CLI flags.
func (cfg *Config) IntVal(key string, defVal int64, reload ...bool) (value int64) {
	vals := cfg.get(false)
	if len(reload) > 0 && reload[0] {
		vals = cfg.get(true)
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

// FloatVal retrieves the value of a specific namespaced environment variable or CLI flag as a float.
// If the key is not found or cannot be parsed as a float, it returns the provided default value.
// If reload is true, it re-reads the values from the environment and CLI flags.
func (cfg *Config) FloatVal(key string, defVal float64, reload ...bool) (value float64) {
	vals := cfg.get(false)
	if len(reload) > 0 && reload[0] {
		vals = cfg.get(true)
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

// BoolVal retrieves the value of a specific namespaced environment variable or CLI flag as a bool.
// If the key is not found or cannot be parsed as a bool, it returns the provided default value.
// If reload is true, it re-reads the values from the environment and CLI flags.
func (cfg *Config) BoolVal(key string, defVal bool, reload ...bool) (value bool) {
	vals := cfg.get(false)
	if len(reload) > 0 && reload[0] {
		vals = cfg.get(true)
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
func (cfg *Config) loadNamespaceEnvVars() {
	cfg.values = cfg.readNamespaceEnvVars()
}

// readNamespaceEnvVars reads all visible environment variables that belong to the namespace.
func (cfg *Config) readNamespaceEnvVars() map[string]string {
	nevs := make(map[string]string)
	np := cfg.namespacePrefix()

	for _, ev := range os.Environ() {
		if strings.HasPrefix(ev, np) {
			varval := strings.SplitN(ev, "=", 2)

			if len(varval) < 2 {
				continue
			}

			key := cfg.keyify(varval[0])
			nevs[key] = varval[1]
		}
	}

	return nevs
}

// keyify converts environment variable names to a dot-separated, lowercase format.
// For example, NAMESPACE_CONFIG_VALUE becomes config.value.
func (cfg *Config) keyify(name string) string {
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
func (cfg *Config) loadFlags() {
	cfg.flags = make(map[string]string)
	flag.VisitAll(func(f *flag.Flag) {
		cfg.flags[f.Name] = f.Value.String()
	})
}

// defineFlags defines the CLI flags used by the application.
func (cfg *Config) defineFlags(flagDefs map[string]interface{}) {
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

func (cfg *Config) WebAddr() string {
	host := cfg.StrValOrDef(Key.ServerWebHost, "localhost")
	port := cfg.StrValOrDef(Key.ServerWebPort, "8080")
	return host + ":" + port
}

func (cfg *Config) APIAddr() string {
	host := cfg.StrValOrDef(Key.ServerAPIHost, "localhost")
	port := cfg.StrValOrDef(Key.ServerAPIPort, "8081")
	return host + ":" + port
}

// Debug prints the configuration values in a readable format.
func (cfg *Config) Debug() {
	fmt.Println("Configuration values:")
	for k, v := range cfg.Get() {
		fmt.Printf("%s: %s\n", k, v)
	}
}
