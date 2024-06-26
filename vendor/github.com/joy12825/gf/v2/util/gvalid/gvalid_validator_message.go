// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/joy12825/gf.

package gvalid

import (
	"context"

	"github.com/joy12825/gf/v2/util/gvalid/internal/builtin"
)

// getErrorMessageByRule retrieves and returns the error message for specified rule.
// It firstly retrieves the message from custom message map, and then checks i18n manager,
// it returns the default error message if it's not found in neither custom message map nor i18n manager.
func (v *Validator) getErrorMessageByRule(ctx context.Context, ruleKey string, customMsgMap map[string]string) string {
	content := customMsgMap[ruleKey]
	if content != "" {
		// I18n translation.
		i18nContent := v.i18nManager.GetContent(ctx, content)
		if i18nContent != "" {
			return i18nContent
		}
		return content
	}

	// Retrieve default message according to certain rule.
	content = v.i18nManager.GetContent(ctx, ruleMessagePrefixForI18n+ruleKey)
	if content == "" {
		content = defaultErrorMessages[ruleKey]
	}
	// Builtin rule message.
	if content == "" {
		if builtinRule := builtin.GetRule(ruleKey); builtinRule != nil {
			content = builtinRule.Message()
		}
	}
	// If there's no configured rule message, it uses default one.
	if content == "" {
		content = v.i18nManager.GetContent(ctx, ruleMessagePrefixForI18n+internalDefaultRuleName)
	}
	// If there's no configured rule message, it uses default one.
	if content == "" {
		content = defaultErrorMessages[internalDefaultRuleName]
	}
	return content
}
