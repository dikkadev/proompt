import { useState } from "react";
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
} from "lucide-react";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

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

  return (
    <div className="h-full flex flex-col">
      <Card className="flex-1 border-0 rounded-none bg-card flex flex-col">
      {/* Search Header */}
      <div className="p-3 border-b border-border">
        <div className="flex items-center gap-2 mb-2">
          <div className="flex-1">
            {searchMode === 'combined' ? (
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search prompts & snippets..."
                  value={combinedQuery}
                  onChange={(e) => setCombinedQuery(e.target.value)}
                  className="pl-10"
                />
              </div>
            ) : (
              <div className="space-y-2">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <Input
                    placeholder="Search prompts..."
                    value={promptQuery}
                    onChange={(e) => setPromptQuery(e.target.value)}
                    className="pl-10"
                  />
                </div>
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <Input
                    placeholder="Search snippets..."
                    value={snippetQuery}
                    onChange={(e) => setSnippetQuery(e.target.value)}
                    className="pl-10"
                  />
                </div>
              </div>
            )}
          </div>
          
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="ghost"
                size="sm"
                onClick={toggleSearchMode}
                className="h-8 w-8 p-0"
              >
                {searchMode === 'combined' ? (
                  <Layers className="h-4 w-4" />
                ) : (
                  <SplitSquareHorizontal className="h-4 w-4" />
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
          <ResizablePanel defaultSize={50} minSize={20}>
            <div className="h-full flex flex-col">
              <div className="p-3 border-b border-border">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <FileText className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium text-sm">PROMPTS</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <Tooltip>
                      <TooltipTrigger asChild>
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
                      <TooltipTrigger asChild>
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
              <div className="flex-1 p-3 overflow-auto">
                {/* TODO: Implement prompt lists/search results */}
                <div className="text-sm text-muted-foreground">
                  Prompt organization coming soon...
                </div>
              </div>
            </div>
          </ResizablePanel>

          <ResizableHandle withHandle />

          {/* Snippets Pane */}
          <ResizablePanel defaultSize={50} minSize={20}>
            <div className="h-full flex flex-col">
              <div className="p-3 border-b border-border">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Code2 className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium text-sm">SNIPPETS</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <Tooltip>
                      <TooltipTrigger asChild>
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
                      <TooltipTrigger asChild>
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
              </div>
              <div className="flex-1 p-3 overflow-auto">
                {/* Mock snippet cards with drag affordance */}
                <div className="space-y-2">
                  <SnippetCard
                    snippet={{
                      id: '1',
                      title: 'analysis_guidelines',
                      description: 'Standard guidelines for content analysis',
                      tags: ['analysis', 'guidelines', 'methodology'],
                      category: 'Analysis'
                    }}
                    onInsert={onSnippetInsert}
                  />
                  <SnippetCard
                    snippet={{
                      id: '2',
                      title: 'response_template',
                      description: 'Standard response format template',
                      tags: ['template', 'format', 'structure'],
                      category: 'Templates'
                    }}
                    onInsert={onSnippetInsert}
                  />
                  <SnippetCard
                    snippet={{
                      id: '3',
                      title: 'code_review_checklist',
                      description: 'Comprehensive code review checklist',
                      tags: ['code', 'review', 'checklist'],
                      category: 'Development'
                    }}
                    onInsert={onSnippetInsert}
                  />
                </div>
              </div>
            </div>
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>
    </Card>
    </div>
  );
}

// Draggable snippet card component with visual grab affordance
function SnippetCard({ 
  snippet, 
  onInsert 
}: { 
  snippet: {
    id: string;
    title: string;
    description: string;
    tags: string[];
    category: string;
  };
  onInsert: (name: string) => void;
}) {
  return (
    <Card className="p-3 cursor-grab active:cursor-grabbing hover:shadow-md transition-all duration-200 border-2 border-transparent hover:border-primary/20 bg-gradient-to-r from-background to-background hover:from-primary/5 hover:to-primary/10">
      <div className="flex items-start gap-2">
        {/* Drag handle icon */}
        <div className="flex-shrink-0 mt-0.5">
          <GripVertical className="h-4 w-4 text-muted-foreground/60" />
        </div>
        
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <Code2 className="h-3 w-3 text-muted-foreground flex-shrink-0" />
            <span className="font-mono text-sm font-medium truncate">
              @{snippet.title}
            </span>
          </div>
          
          {snippet.description && (
            <p className="text-xs text-muted-foreground mb-2 line-clamp-2">
              {snippet.description}
            </p>
          )}
          
          <div className="flex items-center justify-between text-xs">
            <div className="flex flex-wrap gap-1 min-w-0">
              {snippet.tags.slice(0, 2).map((tag) => (
                <span 
                  key={tag}
                  className="px-1 py-0.5 bg-secondary/60 text-secondary-foreground rounded text-xs"
                >
                  #{tag}
                </span>
              ))}
              {snippet.tags.length > 2 && (
                <span className="text-muted-foreground">
                  +{snippet.tags.length - 2}
                </span>
              )}
            </div>
            
            <span className="text-muted-foreground ml-2 flex-shrink-0">
              {snippet.category}
            </span>
          </div>
        </div>
      </div>
    </Card>
  );
} 