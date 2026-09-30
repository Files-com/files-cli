package lib

import (
	"strings"

	"github.com/spf13/cobra"
)

// Command annotations read by command discovery. Cobra ignores annotation keys
// it does not define, so these never change how a command runs.
const (
	commandAnnotationEffect           = "files-cli-effect"
	commandAnnotationResponse         = "files-cli-response"
	commandAnnotationResponseWithFlag = "files-cli-response-with-flag"
)

// CommandEffect is a hint about what the Files.com API request behind a
// command does. It is guidance for choosing a command, not a guarantee, and
// nothing enforces it.
type CommandEffect string

const (
	// EffectReadOnly marks a command that only reads data.
	EffectReadOnly CommandEffect = "read_only"
	// EffectMutating marks a command that changes data or starts an action and
	// is not known to be destructive. That does not make it additive.
	EffectMutating CommandEffect = "mutating"
	// EffectDestructive marks a command that can delete, overwrite, or replace
	// existing data.
	EffectDestructive CommandEffect = "destructive"
)

// The setters are not inlined: generated code calls them from hundreds of
// command constructors, and inlining the map setup into each call site added
// about 160 KB to the binary.

// SetCommandEffect records the effect hint of a command.
//
//go:noinline
func SetCommandEffect(cmd *cobra.Command, effect CommandEffect) {
	setCommandAnnotation(cmd, commandAnnotationEffect, string(effect))
}

// GetCommandEffect returns the effect recorded by SetCommandEffect.
func GetCommandEffect(cmd *cobra.Command) (CommandEffect, bool) {
	effect, ok := cmd.Annotations[commandAnnotationEffect]
	return CommandEffect(effect), ok
}

// SetCommandResponse records the API entity a command prints on success: one
// record, a list of records when list is true, or nothing when entity is "".
//
//go:noinline
func SetCommandResponse(cmd *cobra.Command, entity string, list bool) {
	setCommandAnnotation(cmd, commandAnnotationResponse, encodeResponse(entity, list))
}

// GetCommandResponse returns the response recorded by SetCommandResponse.
func GetCommandResponse(cmd *cobra.Command) (entity string, list bool, ok bool) {
	value, ok := cmd.Annotations[commandAnnotationResponse]
	entity, list = decodeResponse(value)
	return entity, list, ok
}

// SetCommandResponseWithFlag records the entity a command can print instead
// when the named flag is set.
//
//go:noinline
func SetCommandResponseWithFlag(cmd *cobra.Command, flag string, entity string, list bool) {
	setCommandAnnotation(cmd, commandAnnotationResponseWithFlag, flag+"="+encodeResponse(entity, list))
}

// GetCommandResponseWithFlag returns the response recorded by
// SetCommandResponseWithFlag.
func GetCommandResponseWithFlag(cmd *cobra.Command) (flag string, entity string, list bool, ok bool) {
	value, ok := cmd.Annotations[commandAnnotationResponseWithFlag]
	flag, response, _ := strings.Cut(value, "=")
	entity, list = decodeResponse(response)
	return flag, entity, list, ok
}

func setCommandAnnotation(cmd *cobra.Command, key string, value string) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations[key] = value
}

func encodeResponse(entity string, list bool) string {
	if list {
		return "[]" + entity
	}
	return entity
}

func decodeResponse(value string) (entity string, list bool) {
	entity, list = strings.CutPrefix(value, "[]")
	return entity, list
}
