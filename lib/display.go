package lib

import "strings"

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

func SanitizeArgsForDisplay(args []string) []string {
	displayArgs := append([]string(nil), args...)

	for i := 0; i < len(displayArgs); i++ {
		switch {
		case displayArgs[i] == "--api-key" || displayArgs[i] == "-a":
			if i+1 < len(displayArgs) && !strings.HasPrefix(displayArgs[i+1], "-") {
				displayArgs[i+1] = maskAPIKeyForDisplay(displayArgs[i+1])
				i++
			}
		case strings.HasPrefix(displayArgs[i], "--api-key="):
			displayArgs[i] = "--api-key=" + maskAPIKeyForDisplay(strings.TrimPrefix(displayArgs[i], "--api-key="))
		case strings.HasPrefix(displayArgs[i], "-a="):
			displayArgs[i] = "-a=" + maskAPIKeyForDisplay(strings.TrimPrefix(displayArgs[i], "-a="))
		}
	}

	return displayArgs
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
