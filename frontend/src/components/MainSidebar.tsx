import { useState, useMemo } from "react";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import {
  ResizablePanelGroup,
  ResizablePanel,
  ResizableHandle,
} from "@/components/ui/resizable";
import {
  Search,
  Layers,
  SplitSquareHorizontal,
  FileText,
  Code2,
  Hash,
  Folder,
  ChevronDown,
  ChevronRight,
  GripVertical,
  Maximize,
  Minimize
} from "lucide-react";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
  TooltipProvider
} from "@/components/ui/tooltip";
import { usePrompts, useSnippets } from "@/lib/queries";
import { Prompt, Snippet } from "@/lib/api";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { cn } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";
import { PromptCard } from "./PromptCard";
import { SnippetCard } from "./SnippetCard";


interface MainSidebarProps {
  onSnippetInsert: (snippetName: string) => void;
  isCollapsed: boolean;
  onToggleCollapse: () => void;
}

export function MainSidebar({ 
  onSnippetInsert, 
  isCollapsed, 
  onToggleCollapse 
}: MainSidebarProps) {
  const [searchMode, setSearchMode] = useState<'combined' | 'split'>('combined');
  const [combinedQuery, setCombinedQuery] = useState('');
  const [promptQuery, setPromptQuery] = useState('');
  const [snippetQuery, setSnippetQuery] = useState('');
  const [promptViewMode, setPromptViewMode] = useState<'categories' | 'tags'>('categories');
  const [snippetViewMode, setSnippetViewMode] = useState<'categories' | 'tags'>('categories');
  const [expandedPromptCategories, setExpandedPromptCategories] = useState<Set<string>>(new Set());
  const [expandedSnippetCategories, setExpandedSnippetCategories] = useState<Set<string>>(new Set());
  const [allPromptsExpanded, setAllPromptsExpanded] = useState(false);
  const [allSnippetsExpanded, setAllSnippetsExpanded] = useState(false);

  const AUTO_EXPAND_THRESHOLD = 2; // Configurable: number of categories/tags to auto-expand

  const { data: promptsData } = usePrompts();
  const { data: snippetsData } = useSnippets();

  const prompts = useMemo(() => promptsData?.data || [], [promptsData]);
  const snippets = useMemo(() => snippetsData?.data || [], [snippetsData]);

  // Group prompts by use_case
  const groupedPromptsByUseCase = useMemo(() => {
    return prompts.reduce((acc, prompt) => {
      const category = prompt.use_case || 'Uncategorized';
      if (!acc[category]) {
        acc[category] = [];
      }
      acc[category].push(prompt);
      return acc;
    }, {} as Record<string, Prompt[]>);
  }, [prompts]);

  const sortedPromptCategories = useMemo(() => {
    return Object.keys(groupedPromptsByUseCase).sort((a, b) => 
      groupedPromptsByUseCase[b].length - groupedPromptsByUseCase[a].length
    );
  }, [groupedPromptsByUseCase]);

  // Group snippets by category (using category field for now, will map to proper snippet category later if needed)
  const groupedSnippetsByCategory = useMemo(() => {
    return snippets.reduce((acc, snippet) => {
      const category = (snippet as any).category || 'Uncategorized'; // Assuming 'category' field on snippet for now
      if (!acc[category]) {
        acc[category] = [];
      }
      acc[category].push(snippet);
      return acc;
    }, {} as Record<string, Snippet[]>);
  }, [snippets]);

  const sortedSnippetCategories = useMemo(() => {
    return Object.keys(groupedSnippetsByCategory).sort((a, b) => 
      groupedSnippetsByCategory[b].length - groupedSnippetsByCategory[a].length
    );
  }, [groupedSnippetsByCategory]);

  const togglePromptCategory = (category: string) => {
    setExpandedPromptCategories(prev => {
      const newSet = new Set(prev);
      if (newSet.has(category)) {
        newSet.delete(category);
      } else {
        newSet.add(category);
      }
      return newSet;
    });
  };

  const toggleSnippetCategory = (category: string) => {
    setExpandedSnippetCategories(prev => {
      const newSet = new Set(prev);
      if (newSet.has(category)) {
        newSet.delete(category);
      } else {
        newSet.add(category);
      }
      return newSet;
    });
  };

  const handleToggleAllPrompts = () => {
    if (expandedPromptCategories.size === sortedPromptCategories.length) {
      // If all are expanded, collapse all
      setExpandedPromptCategories(new Set());
    } else {
      // Otherwise, expand all
      setExpandedPromptCategories(new Set(sortedPromptCategories));
    }
  };

  const handleToggleAllSnippets = () => {
    if (expandedSnippetCategories.size === sortedSnippetCategories.length) {
      // If all are expanded, collapse all
      setExpandedSnippetCategories(new Set());
    } else {
      // Otherwise, expand all
      setExpandedSnippetCategories(new Set(sortedSnippetCategories));
    }
  };

  if (isCollapsed) {
    return (
      <div className="w-12 border-r border-border bg-card flex flex-col items-center py-4 gap-2">
        <Button
          variant="ghost"
          size="sm"
          onClick={onToggleCollapse}
          className="h-8 w-8 p-0"
        >
          <FileText className="h-4 w-4" />
        </Button>
      </div>
    );
  }

  const toggleSearchMode = () => {
    setSearchMode(searchMode === 'combined' ? 'split' : 'combined');
  };

  // Filtered prompts based on search query
  const filteredPrompts = useMemo(() => {
    const query = searchMode === 'combined' ? combinedQuery : promptQuery;
    if (!query) return prompts;
    return prompts.filter(prompt => 
      prompt.title.toLowerCase().includes(query.toLowerCase()) ||
      prompt.content.toLowerCase().includes(query.toLowerCase()) ||
      (prompt.use_case && prompt.use_case.toLowerCase().includes(query.toLowerCase())) ||
      (prompt.model_compatibility_tags && prompt.model_compatibility_tags.some(tag => tag.toLowerCase().includes(query.toLowerCase())))
    );
  }, [prompts, combinedQuery, promptQuery, searchMode]);

  // Filtered snippets based on search query
  const filteredSnippets = useMemo(() => {
    const query = searchMode === 'combined' ? combinedQuery : snippetQuery;
    if (!query) return snippets;
    return snippets.filter(snippet =>
      snippet.title.toLowerCase().includes(query.toLowerCase()) ||
      snippet.content.toLowerCase().includes(query.toLowerCase()) ||
      (snippet.description && snippet.description.toLowerCase().includes(query.toLowerCase())) ||
      ((snippet as any).tags && (snippet as any).tags.some((tag: string) => tag.toLowerCase().includes(query.toLowerCase()))) ||
      ((snippet as any).category && (snippet as any).category.toLowerCase().includes(query.toLowerCase()))
    );
  }, [snippets, combinedQuery, snippetQuery, searchMode]);


  return (
    <div className="h-full flex flex-col">
      <TooltipProvider delayDuration={100}>
        <Card className="flex-1 border-0 rounded-none bg-card flex flex-col">
      {/* Search Header (always visible) */}
      <div className="p-2 border-b border-border">
        <div className="flex items-center gap-1">
          <div className="flex-1">
            {searchMode === 'combined' ? (
              <div className="relative">
                <Search className="absolute left-2 top-1/2 transform -translate-y-1/2 h-3 w-3 text-muted-foreground" />
                <Input
                  placeholder="Search prompts & snippets..."
                  value={combinedQuery}
                  onChange={(e) => setCombinedQuery(e.target.value)}
                  className="h-8 pl-8 text-sm"
                />
              </div>
            ) : (
              <div className="relative">
                <Search className="absolute left-2 top-1/2 transform -translate-y-1/2 h-3 w-3 text-muted-foreground" />
                <Input
                  placeholder="Search prompts..."
                  value={promptQuery}
                  onChange={(e) => setPromptQuery(e.target.value)}
                  className="h-8 pl-8 text-sm"
                />
              </div>
            )}
          </div>
          
          <Tooltip>
            <TooltipTrigger asChild className="cursor-pointer">
              <Button
                variant="ghost"
                size="sm"
                onClick={toggleSearchMode}
                className="h-7 w-7 p-0"
              >
                {searchMode === 'combined' ? (
                  <Layers className="h-3.5 w-3.5" />
                ) : (
                  <SplitSquareHorizontal className="h-3.5 w-3.5" />
                )}
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              {searchMode === 'combined' 
                ? 'Switch to separate searches' 
                : 'Switch to combined search'}
            </TooltipContent>
          </Tooltip>
        </div>
      </div>

      {/* Split Panes */}
      <div className="flex-1 overflow-hidden">
        <ResizablePanelGroup direction="vertical">
          {/* Prompts Pane */}
          <ResizablePanel defaultSize={50} minSize={30}>
            <div className="h-full flex flex-col">
              <div className="p-2 border-b border-border">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-1">
                    {/* <FileText className="h-3.5 w-3.5 text-muted-foreground" /> */}
                    <span className="font-medium text-sm">PROMPTS</span>
                  </div>
                  <div className="flex items-center gap-0.5">
                    {/* New Toggle All Prompts Button */}
                    <Tooltip>
                      <TooltipTrigger asChild className="cursor-pointer">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={handleToggleAllPrompts}
                          className="h-6 w-6 p-0"
                        >
                          {expandedPromptCategories.size === sortedPromptCategories.length ? <Minimize className="h-3 w-3" /> : <Maximize className="h-3 w-3" />}
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>
                        {expandedPromptCategories.size === sortedPromptCategories.length ? 'Collapse all prompts' : 'Expand all prompts'}
                      </TooltipContent>
                    </Tooltip>
                    <Separator orientation="vertical" className="h-4 mx-1" />
                    <Tooltip>
                      <TooltipTrigger asChild className="cursor-pointer">
                        <Button
                          variant={promptViewMode === 'categories' ? 'secondary' : 'ghost'}
                          size="sm"
                          onClick={() => setPromptViewMode('categories')}
                          className="h-6 w-6 p-0"
                        >
                          <Folder className="h-3 w-3" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>View by categories</TooltipContent>
                    </Tooltip>
                    <Tooltip>
                      <TooltipTrigger asChild className="cursor-pointer">
                        <Button
                          variant={promptViewMode === 'tags' ? 'secondary' : 'ghost'}
                          size="sm"
                          onClick={() => setPromptViewMode('tags')}
                          className="h-6 w-6 p-0"
                        >
                          <Hash className="h-3 w-3" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>View by tags</TooltipContent>
                    </Tooltip>
                  </div>
                </div>
              </div>
              <ScrollArea className="flex-1 px-1 py-0">
                {combinedQuery || promptQuery ? (
                  // Search results for prompts
                  <div className="space-y-1 px-1">
                    {filteredPrompts.length > 0 ? (
                      filteredPrompts.map(prompt => (
                        <PromptCard key={prompt.id} prompt={prompt} className="my-0.5" />
                      ))
                    ) : (
                      <div className="text-sm text-muted-foreground py-1 px-1">No matching prompts found.</div>
                    )}
                  </div>
                ) : (
                  // Organized prompts by categories/tags
                  <div className="space-y-1 px-1">
                    {sortedPromptCategories.length > 0 ? (
                      sortedPromptCategories.map(category => (
                        <Collapsible 
                          key={category} 
                          open={expandedPromptCategories.has(category) || allPromptsExpanded || sortedPromptCategories.length <= AUTO_EXPAND_THRESHOLD} 
                          onOpenChange={() => togglePromptCategory(category)}
                          className="my-px cursor-pointer"
                        >
                          <CollapsibleTrigger asChild>
                            <Button variant="ghost" className="w-full justify-between font-normal h-7 px-1 py-1 text-sm">
                              <span className="flex items-center gap-1 min-w-0">
                                <Folder className="h-3 w-3 text-muted-foreground" />
                                <span className="flex-1 overflow-hidden whitespace-nowrap cursor-pointer">{category}</span>
                              </span>
                              <span className="flex items-center gap-1">
                                <Badge variant="outline" className="ml-1 text-xs px-0 py-0 flex items-center justify-center border-transparent bg-transparent text-muted-foreground shrink-0">
                                  {groupedPromptsByUseCase[category].length}
                                </Badge>
                                {expandedPromptCategories.has(category) ? <ChevronDown className="h-3 w-3 text-muted-foreground" /> : <ChevronRight className="h-3 w-3 text-muted-foreground" />}
                              </span>
                            </Button>
                          </CollapsibleTrigger>
                          <CollapsibleContent className="space-y-1 p-0">
                            {groupedPromptsByUseCase[category].map(prompt => (
                              <PromptCard key={prompt.id} prompt={prompt} className="ml-1 my-0.5" />
                            ))}
                          </CollapsibleContent>
                        </Collapsible>
                      ))
                    ) : (
                      <div className="text-sm text-muted-foreground py-1 px-1">No prompts available.</div>
                    )}
                  </div>
                )}
              </ScrollArea>
            </div>
          </ResizablePanel>

          <ResizableHandle withHandle />

          {/* Snippets Pane */}
          <ResizablePanel defaultSize={50} minSize={30}>
            <div className="h-full flex flex-col">
              <div className="p-2 border-b border-border">
                <div className="flex items-center justify-between mb-1">
                  <div className="flex items-center gap-1">
                    {/* <Code2 className="h-3.5 w-3.5 text-muted-foreground" /> */}
                    <span className="font-medium text-sm">SNIPPETS</span>
                  </div>
                  <div className="flex items-center gap-0.5">
                    {/* New Toggle All Snippets Button */}
                    <Tooltip>
                      <TooltipTrigger asChild className="cursor-pointer">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={handleToggleAllSnippets}
                          className="h-6 w-6 p-0"
                        >
                          {expandedSnippetCategories.size === sortedSnippetCategories.length ? <Minimize className="h-3 w-3" /> : <Maximize className="h-3 w-3" />}
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>
                        {expandedSnippetCategories.size === sortedSnippetCategories.length ? 'Collapse all snippets' : 'Expand all snippets'}
                      </TooltipContent>
                    </Tooltip>
                    <Separator orientation="vertical" className="h-4 mx-1" />
                    <Tooltip>
                      <TooltipTrigger asChild className="cursor-pointer">
                        <Button
                          variant={snippetViewMode === 'categories' ? 'secondary' : 'ghost'}
                          size="sm"
                          onClick={() => setSnippetViewMode('categories')}
                          className="h-6 w-6 p-0"
                        >
                          <Folder className="h-3 w-3" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>View by categories</TooltipContent>
                    </Tooltip>
                    <Tooltip>
                      <TooltipTrigger asChild className="cursor-pointer">
                        <Button
                          variant={snippetViewMode === 'tags' ? 'secondary' : 'ghost'}
                          size="sm"
                          onClick={() => setSnippetViewMode('tags')}
                          className="h-6 w-6 p-0"
                        >
                          <Hash className="h-3 w-3" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>View by tags</TooltipContent>
                    </Tooltip>
                  </div>
                </div>
                
                {/* Snippet search - only show in split mode */}
                {searchMode === 'split' && (
                  <div className="relative mt-1">
                    <Search className="absolute left-2 top-1/2 transform -translate-y-1/2 h-3 w-3 text-muted-foreground" />
                    <Input
                      placeholder="Search snippets..."
                      value={snippetQuery}
                      onChange={(e) => setSnippetQuery(e.target.value)}
                      className="h-8 pl-8 text-sm"
                    />
                  </div>
                )}
              </div>
              <ScrollArea className="flex-1 px-1 py-0">
                {combinedQuery || snippetQuery ? (
                  // Search results for snippets
                  <div className="space-y-1 px-1">
                    {filteredSnippets.length > 0 ? (
                      filteredSnippets.map(snippet => (
                        <SnippetCard key={snippet.id} snippet={snippet} onInsert={onSnippetInsert} className="my-px" />
                      ))
                    ) : (
                      <div className="text-sm text-muted-foreground py-1 px-1">No matching snippets found.</div>
                    )}
                  </div>
                ) : (
                  // Organized snippets by categories/tags
                  <div className="space-y-1 px-1">
                    {sortedSnippetCategories.length > 0 ? (
                      sortedSnippetCategories.map(category => (
                        <Collapsible 
                          key={category} 
                          open={expandedSnippetCategories.has(category) || allSnippetsExpanded || sortedSnippetCategories.length <= AUTO_EXPAND_THRESHOLD} 
                          onOpenChange={() => toggleSnippetCategory(category)}
                          className="my-px cursor-pointer"
                        >
                          <CollapsibleTrigger asChild>
                            <Button variant="ghost" className="w-full justify-between font-normal h-7 px-1 py-1 text-sm">
                              <span className="flex items-center gap-1 min-w-0">
                                <Folder className="h-3 w-3 text-muted-foreground" />
                                <span className="flex-1 overflow-hidden whitespace-nowrap cursor-pointer">{category}</span>
                              </span>
                              <span className="flex items-center gap-1">
                                <Badge variant="outline" className="ml-1 text-xs px-0 py-0 flex items-center justify-center border-transparent bg-transparent text-muted-foreground shrink-0">
                                  {groupedSnippetsByCategory[category].length}
                                </Badge>
                                {expandedSnippetCategories.has(category) ? <ChevronDown className="h-3 w-3 text-muted-foreground" /> : <ChevronRight className="h-3 w-3 text-muted-foreground" />}
                              </span>
                            </Button>
                          </CollapsibleTrigger>
                          <CollapsibleContent className="space-y-1 p-0">
                            {groupedSnippetsByCategory[category].map(snippet => (
                              <SnippetCard key={snippet.id} snippet={snippet} onInsert={onSnippetInsert} className="ml-1" />
                            ))}
                          </CollapsibleContent>
                        </Collapsible>
                      ))
                    ) : (
                      <div className="text-sm text-muted-foreground py-1 px-1">No snippets available.</div>
                    )}
                  </div>
                )}
              </ScrollArea>
            </div>
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>
    </Card>
    </TooltipProvider>
    </div>
  );
} 