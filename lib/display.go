package lib

const fullyRedactedValue = "<redacted>"
const apiKeyVisiblePrefixLength = 16
const apiKeyMaskSuffix = "****************"

func maskAPIKeyForDisplay(value string) string {
	if len(value) <= apiKeyVisiblePrefixLength {
		if value == "" {
			return ""
		}
		return fullyRedactedValue
	}

	return value[:apiKeyVisiblePrefixLength] + apiKeyMaskSuffix
}

func redactSessionIDForDisplay(value string) string {
	if value == "" {
		return ""
	}

	return fullyRedactedValue
}

func (p *Profile) Display() *Profile {
	if p == nil {
		return nil
	}

	display := *p
	display.APIKey = maskAPIKeyForDisplay(display.APIKey)
	display.SessionId = redactSessionIDForDisplay(display.SessionId)
	return &display
}

func (p *Profiles) Display() *Profiles {
	if p == nil {
		return nil
	}

	display := &Profiles{
		Profiles: make(map[string]*Profile, len(p.Profiles)),
	}

	for name, profile := range p.Profiles {
		display.Profiles[name] = profile.Display()
	}

	return display
}
