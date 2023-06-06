package menu

import (
	"github.com/c-bata/go-prompt"
)

func createSuggestions(options []string) []prompt.Suggest {
	var suggestions []prompt.Suggest
	for _, option := range options {
		suggestions = append(suggestions, prompt.Suggest{Text: option})
	}
	return suggestions
}

func completer(options []string) prompt.Completer {
	return func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(createSuggestions(options), d.GetWordBeforeCursor(), true)
	}
}

func Select(values []string, text string) string {
	options := []prompt.Option{
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSuggestionTextColor(prompt.White),
	}

	return prompt.Input(text, completer(values), options...)
}
