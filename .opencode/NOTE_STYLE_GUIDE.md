# NOTE STYLE GUIDE

## Purpose of these notes
These are working notes for AI assistants working on this project. They should be:
- **Timeless**: Capture state and discoveries, not plans or todos
- **Personal**: Raw observations and genuine reactions while exploring code
- **Useful**: Help understand context when working on specific parts later
- **Honest**: Include both positive and negative findings

## Tone guidelines
- Professional but personal
- Include genuine reactions and observations
- Avoid overly dramatic language ("fucking brilliant" → "well-designed")
- Keep some emotional context ("this surprised me", "this was confusing")
- Be direct about problems without being harsh
- Show the discovery process, not just conclusions

## File naming conventions
- **ALLCAPS.md**: Persistent reference files that rarely change
- **lowercase.md**: Working notes that get updated as understanding evolves

## What to include
- Initial impressions and how they changed
- Specific code examples that illustrate points
- Technical details that matter for future work
- Confusion or surprises encountered
- Assessment of code quality and architecture
- Integration issues and their context

## What to avoid
- Formal documentation (belongs elsewhere)
- Step-by-step plans (these are notes, not task lists)
- Overly technical details that belong in code comments
- Judgmental language about previous developers
- Information that will quickly become outdated

## Example good note style
"The template system surprised me - it's more sophisticated than expected. Variable resolution with `{{var:default}}` syntax works properly, and the status tracking (provided/missing/default) is exactly what the frontend needs. But the frontend isn't using it - still doing local regex parsing instead."

This captures: surprise, technical detail, assessment, and current problem state.