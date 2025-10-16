package toggo

// Context represents the evaluation context containing arbitrary attributes
// used for feature flag evaluation. It can hold any key-value pairs such as
// user_id, country, plan, etc.
type Context map[string]interface{}

// Get retrieves a value from the context by key.
// Returns the value and a boolean indicating whether the key exists.
func (c Context) Get(key string) (interface{}, bool) {
	val, ok := c[key]
	return val, ok
}

// GetString retrieves a string value from the context.
// Returns empty string if the key doesn't exist or value is not a string.
func (c Context) GetString(key string) string {
	val, ok := c[key]
	if !ok {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

// Set adds or updates a key-value pair in the context.
func (c Context) Set(key string, value interface{}) {
	c[key] = value
}
