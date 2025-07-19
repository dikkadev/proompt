# Multi-Message/Chaining Brainstorm

## The Problem
How to handle scenarios where you need multiple prompts in sequence:
- Chat conversations with system + user + assistant + user...
- Prompt chains where output of one feeds into next
- Multi-step workflows
- Different roles in same conversation

## Possible Approaches

### Option A: Conversation Objects
- Create "Conversation" entity that contains multiple prompts
- Each prompt has a role (system, user, assistant)
- Order matters
- Can template across the whole conversation

### Option B: Prompt Sequences  
- Link prompts together in sequences
- Each prompt can reference previous outputs
- More flexible than conversations (not just chat)
- Could branch/fork

### Option C: Prompt Relationships
- Prompts can have "parent" and "child" relationships
- Build trees/graphs of related prompts
- Very flexible but maybe too complex

### Option D: Keep Simple
- Just manage individual prompts
- Let users copy/paste between tools for chaining
- Focus on single prompt management excellence

### Option E: Simple Bidirectional Links (USER'S IDEA)
- Prompts can link to "followup" prompts
- Links are bidirectional automatically
- Just navigation aid, no complex structure
- Minimal complexity, maximum usability gain

## User's Insight
"just have a feature where in a prompt you can 'link' (possible) followup prompts. so you can quickly go to them, and when defined in one direction they also work in the other direction"

This is brilliant because:
- No new entities or complex structures
- Just enhances navigation
- Bidirectional automatically (define once, works both ways)
- Doesn't add conceptual complexity
- Solves the real pain point (finding related prompts)