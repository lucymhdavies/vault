package api

import "fmt"

// CapabilitiesSelf returns the capabilities of the client token against a given path
// For compatibility, this is only valid when specifying a single path
func (c *Sys) CapabilitiesSelf(path string) ([]string, error) {
	return c.Capabilities(c.c.Token(), path)
}

// CapabilitiesSelfMultiple returns the capabilities of the client token against given paths
func (c *Sys) CapabilitiesSelfMultiple(path string) (map[string][]string, error) {
	return c.CapabilitiesMultiple(c.c.Token(), path)
}

// Capabilities returns the capabilities of a specified token against a given path
// For compatibility, this is only valid when specifying a single path
func (c *Sys) Capabilities(token, path string) ([]string, error) {
	body := map[string]string{
		"token": token,
		"path":  path,
	}

	reqPath := "/v1/sys/capabilities"
	if token == c.c.Token() {
		reqPath = fmt.Sprintf("%s-self", reqPath)
	}

	r := c.c.NewRequest("POST", reqPath)
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = resp.DecodeJSON(&result)
	if err != nil {
		return nil, err
	}

	if result["capabilities"] == nil {
		return nil, nil
	}
	var capabilities []string
	capabilitiesRaw, ok := result["capabilities"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("error interpreting returned capabilities")
	}
	for _, capability := range capabilitiesRaw {
		capabilities = append(capabilities, capability.(string))
	}
	return capabilities, nil
}

// CapabilitiesMultiple returns the capabilities of a specified token against given paths
func (c *Sys) CapabilitiesMultiple(token, path string) (map[string][]string, error) {
	body := map[string]string{
		"token": token,
		"path":  path,
	}

	reqPath := "/v1/sys/capabilities"
	if token == c.c.Token() {
		reqPath = fmt.Sprintf("%s-self", reqPath)
	}

	r := c.c.NewRequest("POST", reqPath)
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = resp.DecodeJSON(&result)
	if err != nil {
		return nil, err
	}

	data := result["data"].(map[string]interface{})

	capabilities := make(map[string][]string)

	for path, pathCapabilities := range data {
		capabilities[path] = []string{}

		pathCapabilitiesRaw := pathCapabilities.([]interface{})

		for _, capability := range pathCapabilitiesRaw {
			capabilities[path] = append(capabilities[path], capability.(string))
		}
	}

	return capabilities, nil
}
