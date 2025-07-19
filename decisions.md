# Project Decisions Log

## Decision 1: No Prompt Execution
**Date**: Initial discussion
**Decision**: The app will NOT execute prompts - only manage them
**Reasoning**: User already has tools for execution (Langfuse, Anthropic workbench, etc.). Focus is purely on management/organization.

## Decision 2: Target Position on Spectrum  
**Date**: Initial discussion
**Decision**: Build something between full execution platforms (Langfuse) and generic note apps
**Reasoning**: Execution platforms have too much fluff for pure management use case, but generic note apps lack prompt-specific features.

## Decision 3: Prompt Types as Top-Level Categories
**Date**: Follow-up discussion
**Decision**: Treat "system", "user", "image", "video" etc. as the same conceptual level - the "where" to put the prompt
**Reasoning**: User clarified these are all types at the same abstraction level, not hierarchical.

## Decision 4: Templating System Architecture
**Date**: Follow-up discussion  
**Decision**: 
- Snippet directory system for reusable components
- Variables in prompts can be accessed by snippets (1 layer deep only)
- NO nested snippets (snippets cannot contain other snippets)
**Reasoning**: Prevents complexity explosion while allowing powerful composition. Clear mental model with exactly one abstraction layer.

## Decision 5: Metadata Features
**Date**: Feature refinement discussion
**Decision**: Include these metadata features:
- 5a. Model compatibility tags
- 5b. Configurable token count heuristics (user can define and choose which to use)
- 5c. Temperature/parameter suggestions
- 5d. Multiple short notes with title+body (instead of single "success rate" field)
**Reasoning**: 5b provides deep customization for power users. 5d more flexible than single rating field.

## Decision 6: Use Case as First-Class Concept
**Date**: Feature refinement discussion
**Decision**: Use case should be its own field/concept, not just a tag
**Reasoning**: Allows organizing by use case as primary dimension, feels natural for prompt organization.

## Decision 7: Git-Based Versioning (Core Feature)
**Date**: Feature refinement discussion
**Decision**: Use git in background for versioning, tightly controlled (auto-rebase, etc.)
**Reasoning**: Proper versioning without exposing git complexity to users. Merge conflicts handled automatically since we control the workflow. Must be integrated from start.

## Decision 8: Global Snippets
**Date**: Scope discussion
**Decision**: Snippets are global, not per-project
**Reasoning**: For now, keep simple. Later when user/team auth is added, this will be handled by team membership.

## Decision 9: Simple Bidirectional Prompt Links
**Date**: Multi-message discussion
**Decision**: Prompts can link to "followup" prompts, links work bidirectionally automatically
**Reasoning**: Solves navigation pain without adding complexity. No new entities, just enhanced UX.

## Decision 10: General Tagging System
**Date**: Search/organization discussion
**Decision**: 
- Prompts: arbitrary tags + model compatibility tags + use case field + title
- Snippets: arbitrary tags (no use case field)
- Notes: arbitrary tags (no use case field)
**Reasoning**: Flexible organization. Use case only makes sense for full prompts.

## Decision 11: Semantic Search (Future Feature)
**Date**: Search/organization discussion
**Decision**: Separate feature from regular search for exploratory "I want to achieve X" queries
**Reasoning**: Different use case than keyword search - for discovery when you don't know what you're looking for.

## Decision 12: Variable System Details
**Date**: Implementation planning
**Decision**: 
- Syntax: `{{variable_name}}` with optional defaults `{{var:default}}`
- Warnings for missing variables (not errors)
- No resolved storage - always compute on-demand
- String-only variables (no types)
- Snippet variables can be overridden by prompt variables (higher precedence)
**Reasoning**: Bash-style defaults are familiar. Warnings keep flexibility. String-only keeps it simple.

## Decision 13: Variable Dependency Visualization
**Date**: Implementation planning  
**Decision**: Color-coded variables in snippet browser:
- Red: Missing variables (will warn)
- Yellow: Variables with defaults
- Green: Variables provided by current prompt
**Reasoning**: Simple visual system solves dependency awareness and override detection.

## Future Features Timeline:
1. User/team authentication (medium term)
2. Semantic search (later)

## Open Questions Still to Resolve:
- Snippet organization approach