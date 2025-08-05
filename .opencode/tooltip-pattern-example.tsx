// WRONG PATTERN (causes tooltip issues):
// Each tooltip creates its own provider
<TooltipProvider>
  <Tooltip>
    <TooltipTrigger asChild>
      <Button variant="ghost" size="sm">
        <Settings className="h-4 w-4" />
      </Button>
    </TooltipTrigger>
    <TooltipContent>
      <p>Settings</p>
    </TooltipContent>
  </Tooltip>
</TooltipProvider>

// CORRECT PATTERN (use global provider):
// Just use Tooltip components directly
<Tooltip>
  <TooltipTrigger asChild>
    <Button variant="ghost" size="sm">
      <Settings className="h-4 w-4" />
    </Button>
  </TooltipTrigger>
  <TooltipContent>
    <p>Settings</p>
  </TooltipContent>
</Tooltip>

// EXAMPLE: Multiple tooltips in the same component
<div className="flex items-center gap-2">
  <Tooltip>
    <TooltipTrigger asChild>
      <Button variant="ghost" size="sm">
        <Settings className="h-4 w-4" />
      </Button>
    </TooltipTrigger>
    <TooltipContent>
      <p>Settings</p>
    </TooltipContent>
  </Tooltip>

  <Tooltip>
    <TooltipTrigger asChild>
      <Button variant="ghost" size="sm">
        <Search className="h-4 w-4" />
      </Button>
    </TooltipTrigger>
    <TooltipContent>
      <p>Search</p>
    </TooltipContent>
  </Tooltip>

  <Tooltip>
    <TooltipTrigger asChild>
      <Button variant="ghost" size="sm">
        <Filter className="h-4 w-4" />
      </Button>
    </TooltipTrigger>
    <TooltipContent>
      <p>Filter</p>
    </TooltipContent>
  </Tooltip>
</div>

// NOTE: All these tooltips will work correctly with the global TooltipProvider
// No individual TooltipProvider needed! 