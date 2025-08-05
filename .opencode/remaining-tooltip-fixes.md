# Remaining Tooltip Fixes

## ✅ COMPLETED (Test Fix)
- **Index.tsx** (lines 200-229) - Backend Status Indicator tooltip
- **App.tsx** - Updated global TooltipProvider with proper configuration

## 🔧 REMAINING FIXES NEEDED

### HIGH PRIORITY (Multiple instances)

#### 1. MainSidebar.tsx
**File**: `frontend/src/components/MainSidebar.tsx`
**Issue**: Has its own TooltipProvider wrapper with custom delayDuration
**Line**: 217 - Remove `<TooltipProvider delayDuration={100}>` wrapper around the entire component
**Affected tooltips**: Lines 245-265, 283-297, 299-311, 312-324, 394-408, 410-422, 423-435

#### 2. LivePreview.tsx  
**File**: `frontend/src/components/LivePreview.tsx`
**Issue**: Multiple TooltipProvider wrappers around individual tooltips
**Instances to fix**:
- Lines 204-217: Warning status badge tooltip
- Lines 220-233: Missing variables status badge tooltip  
- Lines 250-262: (Need to check - there were more instances found)
- Lines 266-278: (Need to check - there were more instances found)

#### 3. VariablePanel.tsx
**File**: `frontend/src/components/VariablePanel.tsx`
**Issue**: Multiple TooltipProvider wrappers
**Instances to fix**:
- Lines 126-143: Variable-related tooltip
- Lines 157-173: Variable-related tooltip
- Lines 222-239: Variable-related tooltip

### MEDIUM PRIORITY

#### 4. PromptEditor.tsx
**File**: `frontend/src/components/PromptEditor.tsx` 
**Issue**: Two TooltipProvider wrappers
**Instances to fix**:
- Lines 190-205: Editor-related tooltip
- Lines 239-254: Editor-related tooltip

#### 5. ColorPicker.tsx
**File**: `frontend/src/components/ColorPicker.tsx`
**Issue**: Multiple TooltipProvider wrappers (exact locations need verification)
**Action needed**: Find and remove all TooltipProvider wrappers throughout the file

## PATTERN TO APPLY FOR ALL FIXES

### ❌ REMOVE THIS PATTERN:
```tsx
<TooltipProvider>
  <Tooltip>
    <TooltipTrigger asChild>
      {/* content */}
    </TooltipTrigger>
    <TooltipContent>
      {/* tooltip content */}
    </TooltipContent>
  </Tooltip>
</TooltipProvider>
```

### ✅ REPLACE WITH THIS PATTERN:
```tsx
<Tooltip>
  <TooltipTrigger asChild>
    {/* content */}
  </TooltipTrigger>
  <TooltipContent>
    {/* tooltip content */}
  </TooltipContent>
</Tooltip>
```

## IMPORT CLEANUP NEEDED
After removing TooltipProvider usage from components, clean up imports:
- Remove `TooltipProvider` from import statements if no longer used
- Keep `Tooltip`, `TooltipTrigger`, `TooltipContent` imports

## FILES ESTIMATED IMPACT
1. **MainSidebar.tsx**: ~7-8 tooltips (highest impact)
2. **LivePreview.tsx**: ~4 tooltips  
3. **VariablePanel.tsx**: ~3 tooltips
4. **PromptEditor.tsx**: ~2 tooltips
5. **ColorPicker.tsx**: Unknown count (needs investigation)

## TESTING AFTER EACH FILE
1. Hover over tooltip triggers
2. Verify tooltip appears
3. Move mouse away and verify tooltip disappears
4. Test multiple tooltips don't interfere
5. Check browser console for errors

## STATUS
- [x] Global provider configured (App.tsx)
- [x] Test fix applied (Index.tsx - Backend Status Indicator)
- [ ] MainSidebar.tsx
- [ ] LivePreview.tsx
- [ ] VariablePanel.tsx
- [ ] PromptEditor.tsx
- [ ] ColorPicker.tsx 