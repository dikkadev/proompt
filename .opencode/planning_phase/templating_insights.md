# Templating System Insights

## User's Templating Vision
- Strong templating system with snippet directory
- Snippets can be inserted into prompts
- Variables in prompts can be accessed by snippets (1 layer deep)
- NO nested snippets (snippets can't contain other snippets)
- This creates exactly one abstraction layer

## My Analysis
This is actually brilliant constraint design:
- Prevents complexity explosion (no infinite nesting)
- Still allows powerful composition
- Variables flow down but not sideways
- Clear mental model: Prompt -> Variables -> Snippets use those variables

## Technical Implications
- Need snippet management system
- Variable resolution system
- Dependency tracking (which snippets use which variables)
- Preview system to show resolved output

## Questions
- How to handle variable conflicts?
- Should snippets be global or per-project?
- How to organize snippet directory?