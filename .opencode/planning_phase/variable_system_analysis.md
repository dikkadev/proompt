# Variable System Analysis

## The Precedence Issue
User's solution is smart: snippet variables can be overridden by prompt variables.

Example:
- Snippet defines: `{{tone:professional}}`
- Prompt defines: `{{tone:casual}}`
- Result: "casual" wins

This keeps snippets self-contained (they work standalone) but allows customization.

## Visualization Challenge
Need to show:
- Which variables a snippet expects/uses
- Which are provided by current prompt context
- Which will use defaults
- Which are missing (warnings)

Could be color-coded or with icons in snippet browser.

## Potential Issues
1. Variable name conflicts between snippets
2. Users might not realize they're overriding snippet defaults
3. Debugging why a snippet behaves differently in different contexts

The precedence rule helps with #1 and #3, but #2 might need UI design to address.