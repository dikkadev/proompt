# Complete Remaining Tooltip Fixes List

## ✅ COMPLETED (Test Fix)
- **App.tsx** - Global TooltipProvider configured with `delayDuration={100} skipDelayDuration={300}`
- **Index.tsx** (lines ~200-229) - Backend Status Indicator tooltip FIXED ✅

## 🔧 REMAINING FIXES (Confirmed by grep search)

### CRITICAL: Index.tsx has ANOTHER TooltipProvider!
**File**: `frontend/src/pages/Index.tsx`
**Line**: 167 - `<TooltipProvider>` (MISSED THIS ONE!)
**Status**: ❌ NEEDS IMMEDIATE ATTENTION

### HIGH PRIORITY

#### 1. MainSidebar.tsx
**File**: `frontend/src/components/MainSidebar.tsx`
**Line**: 216 - `<TooltipProvider delayDuration={100}>`
**Status**: ❌ NOT FIXED YET

#### 2. LivePreview.tsx
**File**: `frontend/src/components/LivePreview.tsx`
**Lines with TooltipProvider**:
- Line 204: `<TooltipProvider>`
- Line 220: `<TooltipProvider>`
- Line 250: `<TooltipProvider>`
- Line 266: `<TooltipProvider>`
**Status**: ❌ NOT FIXED YET (4 instances)

#### 3. VariablePanel.tsx
**File**: `frontend/src/components/VariablePanel.tsx`
**Lines with TooltipProvider**:
- Line 126: `<TooltipProvider>`
- Line 157: `<TooltipProvider>`
- Line 222: `<TooltipProvider>`
**Status**: ❌ NOT FIXED YET (3 instances)

#### 4. PromptEditor.tsx
**File**: `frontend/src/components/PromptEditor.tsx`
**Lines with TooltipProvider**:
- Line 190: `<TooltipProvider>`
- Line 239: `<TooltipProvider>`
**Status**: ❌ NOT FIXED YET (2 instances)

### LOWER PRIORITY

#### 5. UI Component Library
**File**: `frontend/src/components/ui/sidebar.tsx`
**Line**: 131 - `<TooltipProvider delayDuration={0}>`
**Note**: This is in the UI library, may be intentional for sidebar-specific behavior
**Status**: ⚠️ INVESTIGATE - May need special handling

## SUMMARY OF REMAINING WORK
- **Index.tsx**: 1 more TooltipProvider to fix (URGENT - missed in first pass)
- **MainSidebar.tsx**: 1 TooltipProvider wrapper to remove
- **LivePreview.tsx**: 4 TooltipProvider instances to fix
- **VariablePanel.tsx**: 3 TooltipProvider instances to fix  
- **PromptEditor.tsx**: 2 TooltipProvider instances to fix
- **UI Sidebar**: 1 instance to investigate

**Total remaining**: ~11-12 TooltipProvider instances to fix

## NEXT STEPS FOR USER
1. **Test the current fix first**: Check if the backend status indicator tooltip now disappears properly
2. **If confirmed working**: Apply fixes to the remaining files using the same pattern
3. **Priority order**: Fix Index.tsx line 167 first (missed in initial fix), then MainSidebar.tsx, then others

## TEST THE CURRENT FIX
The backend status indicator in the top bar (connection status with green/red/yellow dot) should now:
- ✅ Show tooltip on hover
- ✅ Hide tooltip when mouse moves away
- ✅ Not get stuck on screen 