# Implementation Steps to Fix Tooltip Issue

## Overview
The tooltip issue is caused by multiple `TooltipProvider` instances creating conflicts. The solution is to use only the global provider and remove all individual providers.

## Step-by-Step Implementation

### 1. Update App.tsx (Global Provider Configuration)
**File**: `frontend/src/App.tsx`
**Change**: Add proper configuration to the global TooltipProvider

```tsx
// BEFORE:
<TooltipProvider>

// AFTER:
<TooltipProvider delayDuration={100} skipDelayDuration={300}>
```

### 2. Fix MainSidebar.tsx (Priority 1)
**File**: `frontend/src/components/MainSidebar.tsx`
**Change**: Remove the wrapping TooltipProvider around line 217

**Before:**
```tsx
<TooltipProvider delayDuration={100}>
  <Card className="flex-1 border-0 rounded-none bg-card flex flex-col">
    {/* Component content with tooltips */}
  </Card>
</TooltipProvider>
```

**After:**
```tsx
<Card className="flex-1 border-0 rounded-none bg-card flex flex-col">
  {/* Component content with tooltips */}
</Card>
```

**Also remove import**: Remove `TooltipProvider` from the import statement if it's only used for the wrapper.

### 3. Fix LivePreview.tsx (Priority 1)
**File**: `frontend/src/components/LivePreview.tsx`
**Lines to modify**: 204-217, 220-233, 250-262, 266-278

**Pattern to replace:**
```tsx
// REMOVE THIS PATTERN:
<TooltipProvider>
  <Tooltip>
    <TooltipTrigger asChild>
      {/* trigger content */}
    </TooltipTrigger>
    <TooltipContent>
      {/* tooltip content */}
    </TooltipContent>
  </Tooltip>
</TooltipProvider>

// REPLACE WITH:
<Tooltip>
  <TooltipTrigger asChild>
    {/* trigger content */}
  </TooltipTrigger>
  <TooltipContent>
    {/* tooltip content */}
  </TooltipContent>
</Tooltip>
```

### 4. Fix VariablePanel.tsx (Priority 1)
**File**: `frontend/src/components/VariablePanel.tsx`
**Lines to modify**: 126-143, 157-173, 222-239

Apply the same pattern as Step 3.

### 5. Fix PromptEditor.tsx (Priority 2)
**File**: `frontend/src/components/PromptEditor.tsx`
**Lines to modify**: 190-205, 239-254

Apply the same pattern as Step 3.

### 6. Fix ColorPicker.tsx (Priority 2)
**File**: `frontend/src/components/ColorPicker.tsx`
**Action**: Find and remove all TooltipProvider wrappers throughout the file.

## Validation Steps

After making changes to each file:

1. **Test tooltip appearance**: Hover over elements with tooltips
2. **Test tooltip disappearance**: Move mouse away from trigger elements
3. **Test multiple tooltips**: Ensure they don't interfere with each other
4. **Check console**: Verify no React/JavaScript errors
5. **Test keyboard navigation**: Tab through tooltip triggers

## Quick Test Script
```bash
# Run the frontend to test changes
cd frontend
bun run dev
```

## Common Issues to Watch For

1. **Import cleanup**: Remove unused `TooltipProvider` imports
2. **Indentation**: Maintain proper code formatting when removing providers
3. **Component boundaries**: Ensure no orphaned closing tags
4. **TypeScript errors**: Check for any type issues after changes

## Files Priority Order
1. `App.tsx` (global config)
2. `MainSidebar.tsx` (has custom provider config)
3. `LivePreview.tsx` (multiple instances)
4. `VariablePanel.tsx` (multiple instances)
5. `PromptEditor.tsx` (two instances)
6. `ColorPicker.tsx` (multiple instances)

## Expected Result
- Tooltips appear on hover
- Tooltips disappear when mouse leaves trigger
- No persistent tooltips stuck on screen
- Consistent timing across all tooltips
- Better performance and cleaner code 