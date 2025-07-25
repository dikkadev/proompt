import { useState } from "react";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { useSnippets } from "@/lib/queries";
import { 
  FileText, 
  Search, 
  Plus, 
  Tag,
  Code,
  MessageSquare,
  Settings,
  Star,
  ChevronDown,
  ChevronRight,
  Hash,
  Filter
} from "lucide-react";

interface Snippet {
  id: string;
  title: string;
  content: string;
  description?: string;
  created_at: string;
  updated_at: string;
  git_ref?: string;
}

interface SnippetSidebarProps {
  onSnippetInsert: (snippetName: string) => void;
  isCollapsed: boolean;
  onToggleCollapse: () => void;
}

type ViewMode = 'categories' | 'tags';

// Mock snippets removed - now using real API data
const mockSnippets: Snippet[] = [
  {
    id: '1',
    name: 'analysis_guidelines',
    content: `When analyzing the content, please follow these guidelines:

1. **Clarity**: Ensure your analysis is clear and easy to understand
2. **Evidence**: Support your points with specific examples from the content
3. **Structure**: Organize your response in logical sections
4. **Objectivity**: Maintain a neutral, analytical tone`,
    description: 'Standard guidelines for content analysis',
    tags: ['analysis', 'guidelines', 'methodology'],
    category: 'Analysis',
    isFavorite: true
  },
  {
    id: '2',
    name: 'response_template',
    content: `## Summary
[Brief overview of the main points]

## Key Findings
- Finding 1: [Description]
- Finding 2: [Description]
- Finding 3: [Description]

## Recommendations
1. [Recommendation with rationale]
2. [Recommendation with rationale]

## Conclusion
[Final thoughts and next steps]`,
    description: 'Standard response format template',
    tags: ['template', 'format', 'structure'],
    category: 'Templates'
  },
  {
    id: '3',
    name: 'code_review_checklist',
    content: `Review the following aspects:

🔍 **Code Quality**
- Is the code readable and well-documented?
- Are naming conventions consistent?
- Is the code properly structured?

🛡️ **Security**
- Are there any security vulnerabilities?
- Is input validation implemented?

⚡ **Performance**
- Are there any performance bottlenecks?
- Is the code efficient?`,
    description: 'Comprehensive code review checklist',
    tags: ['code', 'review', 'checklist', 'quality'],
    category: 'Development'
  },
  {
    id: '4',
    name: 'meeting_summary',
    content: `## Meeting Summary: {{meeting_title}}
**Date**: {{date}}
**Attendees**: {{attendees}}

### Key Decisions
- {{decision_1}}
- {{decision_2}}

### Action Items
- [ ] {{action_1}} (Owner: {{owner_1}}, Due: {{due_1}})
- [ ] {{action_2}} (Owner: {{owner_2}}, Due: {{due_2}})

### Next Steps
{{next_steps}}`,
    description: 'Template for meeting summaries',
    tags: ['meeting', 'summary', 'template', 'collaboration'],
    category: 'Communication'
  },
  {
    id: '5',
    name: 'debugging_approach',
    content: `## Debugging Approach

1. **Reproduce the Issue**
   - Can you consistently reproduce the problem?
   - What are the exact steps to trigger it?

2. **Gather Information**
   - Check error logs and stack traces
   - Review recent changes
   - Note environment details

3. **Isolate the Problem**
   - Use binary search approach
   - Test in isolation
   - Remove variables one by one

4. **Apply Fix**
   - Make minimal changes
   - Test thoroughly
   - Document the solution`,
    description: 'Systematic approach to debugging issues',
    tags: ['debugging', 'troubleshooting', 'methodology', 'process'],
    category: 'Development'
  },
  {
    id: '6',
    name: 'user_feedback_response',
    content: `Thank you for your feedback! I appreciate you taking the time to share your thoughts.

## What I understand:
{{user_concern_summary}}

## Next steps:
- {{action_item_1}}
- {{action_item_2}}

I'll keep you updated on our progress. Please don't hesitate to reach out if you have any other concerns.

Best regards,
{{your_name}}`,
    description: 'Professional response template for user feedback',
    tags: ['feedback', 'response', 'customer-service', 'communication'],
    category: 'Communication'
  },
  {
    id: '7',
    name: 'api_documentation',
    content: `## {{endpoint_name}} API

### Endpoint
\`{{method}} {{endpoint_url}}\`

### Description
{{endpoint_description}}

### Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| {{param_1}} | {{type_1}} | {{required_1}} | {{description_1}} |

### Response
\`\`\`json
{{response_example}}
\`\`\`

### Error Codes
- \`400\`: {{error_400_description}}
- \`401\`: {{error_401_description}}
- \`404\`: {{error_404_description}}`,
    description: 'Template for API endpoint documentation',
    tags: ['api', 'documentation', 'template', 'reference'],
    category: 'Development'
  },
  {
    id: '8',
    name: 'project_proposal',
    content: `# Project Proposal: {{project_name}}

## Executive Summary
{{brief_overview}}

## Problem Statement
{{problem_description}}

## Proposed Solution
{{solution_overview}}

## Timeline
- **Phase 1**: {{phase_1}} ({{timeline_1}})
- **Phase 2**: {{phase_2}} ({{timeline_2}})
- **Phase 3**: {{phase_3}} ({{timeline_3}})

## Resources Required
- {{resource_1}}
- {{resource_2}}
- {{resource_3}}

## Success Metrics
{{success_criteria}}`,
    description: 'Template for project proposals',
    tags: ['proposal', 'project', 'planning', 'business'],
    category: 'Planning'
  }
];

const categories = ['All', 'Analysis', 'Templates', 'Development', 'Communication', 'Planning'];

export function SnippetSidebar({ onSnippetInsert, isCollapsed, onToggleCollapse }: SnippetSidebarProps) {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('All');
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [viewMode, setViewMode] = useState<ViewMode>('categories');
  const [expandedCategories, setExpandedCategories] = useState<Set<string>>(new Set(['Favorites']));
  // Fetch snippets from API
  const { data: snippetsData, isLoading, isError } = useSnippets();
  const snippets = snippetsData?.items || [];

  // For now, disable tag functionality since API snippets don't have tags
  const allTags: string[] = [];
  const tagCounts: Record<string, number> = {};

  const filteredSnippets = snippets.filter(snippet => {
    const matchesSearch = searchQuery === '' ||
                         snippet.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         snippet.description?.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         snippet.content.toLowerCase().includes(searchQuery.toLowerCase());
    
    // For now, ignore category and tag filtering since API doesn't provide these
    return matchesSearch;
  });

  const favoriteSnippets = filteredSnippets.filter(s => s.isFavorite);
  const regularSnippets = filteredSnippets.filter(s => !s.isFavorite);

  // Group snippets by category
  const snippetsByCategory = regularSnippets.reduce((acc, snippet) => {
    if (!acc[snippet.category]) {
      acc[snippet.category] = [];
    }
    acc[snippet.category].push(snippet);
    return acc;
  }, {} as Record<string, Snippet[]>);

  // Group snippets by tags (snippets can appear in multiple groups)
  const snippetsByTag = allTags.reduce((acc, tag) => {
    const snippetsWithTag = regularSnippets.filter(s => s.tags.includes(tag));
    if (snippetsWithTag.length > 0) {
      acc[tag] = snippetsWithTag;
    }
    return acc;
  }, {} as Record<string, Snippet[]>);

  const toggleCategory = (category: string) => {
    const newExpanded = new Set(expandedCategories);
    if (newExpanded.has(category)) {
      newExpanded.delete(category);
    } else {
      newExpanded.add(category);
    }
    setExpandedCategories(newExpanded);
  };

  const toggleTag = (tag: string) => {
    if (selectedTags.includes(tag)) {
      setSelectedTags(selectedTags.filter(t => t !== tag));
    } else {
      setSelectedTags([...selectedTags, tag]);
    }
  };

  const clearFilters = () => {
    setSelectedTags([]);
    setSelectedCategory('All');
    setSearchQuery('');
  };

  const getCategoryIcon = (category: string) => {
    switch (category) {
      case 'Analysis': return Code;
      case 'Templates': return FileText;
      case 'Development': return Settings;
      case 'Communication': return MessageSquare;
      case 'Planning': return Hash;
      default: return FileText;
    }
  };

  if (isCollapsed) {
    return (
      <Card className="w-full h-full flex flex-col items-center py-4 bg-workspace-sidebar">
        <Button
          variant="ghost"
          size="sm"
          onClick={onToggleCollapse}
          className="mb-4 p-2 h-8 w-8"
          aria-label="Expand snippets sidebar"
        >
          <FileText className="h-4 w-4" />
        </Button>
        
        <div className="flex flex-col gap-2">
          {favoriteSnippets.slice(0, 3).map((snippet) => (
            <Button
              key={snippet.id}
              variant="ghost"
              size="sm"
              onClick={() => onSnippetInsert(snippet.title)}
              className="p-2 h-8 w-8"
              title={snippet.title}
            >
              <Star className="h-3 w-3 text-yellow-500" />
            </Button>
          ))}
        </div>
      </Card>
    );
  }

  return (
    <Card className="w-full h-full flex flex-col bg-workspace-sidebar">
      {/* Header */}
      <div className="p-4 border-b border-border">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2">
            <FileText className="h-5 w-5 text-primary" />
            <h3 className="font-semibold">Snippets</h3>
          </div>
          <div className="flex items-center gap-1">
            <Button variant="ghost" size="sm" className="h-6 w-6 p-0">
              <Plus className="h-3 w-3" />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={onToggleCollapse}
              className="h-6 w-6 p-0"
              aria-label="Collapse snippets sidebar"
            >
              <FileText className="h-3 w-3" />
            </Button>
          </div>
        </div>

        {/* Search */}
        <div className="relative mb-3">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-3 w-3 text-muted-foreground" />
          <Input
            placeholder="Search snippets..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-9 h-8 text-sm"
          />
        </div>

        {/* View Mode Toggle */}
        <div className="flex items-center gap-1 mb-3">
          <Button
            variant={viewMode === 'categories' ? "secondary" : "ghost"}
            size="sm"
            onClick={() => setViewMode('categories')}
            className="h-6 px-2 text-xs"
          >
            <Settings className="h-3 w-3 mr-1" />
            Categories
          </Button>
          <Button
            variant={viewMode === 'tags' ? "secondary" : "ghost"}
            size="sm"
            onClick={() => setViewMode('tags')}
            className="h-6 px-2 text-xs"
          >
            <Tag className="h-3 w-3 mr-1" />
            Tags
          </Button>
        </div>

        {/* Filters */}
        {(selectedTags.length > 0 || selectedCategory !== 'All') && (
          <div className="flex items-center gap-2 mb-3 p-2 bg-muted/50 rounded-md">
            <Filter className="h-3 w-3 text-muted-foreground" />
            <div className="flex flex-wrap gap-1 flex-1">
              {selectedCategory !== 'All' && (
                <Badge variant="secondary" className="text-xs">
                  {selectedCategory}
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setSelectedCategory('All')}
                    className="h-3 w-3 p-0 ml-1 hover:bg-destructive/10"
                  >
                    ×
                  </Button>
                </Badge>
              )}
              {selectedTags.map(tag => (
                <Badge key={tag} variant="secondary" className="text-xs">
                  #{tag}
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => toggleTag(tag)}
                    className="h-3 w-3 p-0 ml-1 hover:bg-destructive/10"
                  >
                    ×
                  </Button>
                </Badge>
              ))}
            </div>
            <Button
              variant="ghost"
              size="sm"
              onClick={clearFilters}
              className="h-6 px-2 text-xs"
            >
              Clear
            </Button>
          </div>
        )}

        {/* Quick Category Filter (when in categories mode) */}
        {viewMode === 'categories' && (
          <div className="flex flex-wrap gap-1">
            {categories.map((category) => (
              <Button
                key={category}
                variant={selectedCategory === category ? "secondary" : "ghost"}
                size="sm"
                onClick={() => setSelectedCategory(category)}
                className="h-6 px-2 text-xs"
              >
                {category}
              </Button>
            ))}
          </div>
        )}

        {/* Tag Cloud (when in tags mode) */}
        {viewMode === 'tags' && (
          <div className="space-y-2">
            <div className="flex flex-wrap gap-1">
              {allTags.map((tag) => (
                <Button
                  key={tag}
                  variant={selectedTags.includes(tag) ? "secondary" : "ghost"}
                  size="sm"
                  onClick={() => toggleTag(tag)}
                  className="h-6 px-2 text-xs hover:bg-primary/10"
                >
                  #{tag}
                  <span className="ml-1 text-muted-foreground">
                    {tagCounts[tag]}
                  </span>
                </Button>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Snippet List */}
      <ScrollArea className="flex-1">
        <div className="p-4 space-y-4">
          {/* Loading State */}
          {isLoading && (
            <div className="text-center py-8">
              <div className="text-sm text-muted-foreground">Loading snippets...</div>
            </div>
          )}
          
          {/* Error State */}
          {isError && (
            <div className="text-center py-8">
              <div className="text-sm text-red-600 mb-2">Failed to load snippets</div>
              <div className="text-xs text-muted-foreground">Check your connection to the backend server</div>
            </div>
          )}
          
          {/* Empty State */}
          {!isLoading && !isError && snippets.length === 0 && (
            <div className="text-center py-8">
              <FileText className="h-8 w-8 mx-auto mb-2 text-muted-foreground opacity-50" />
              <div className="text-sm text-muted-foreground">No snippets found</div>
            </div>
          )}
          
          {/* Content */}
          {!isLoading && !isError && snippets.length > 0 && (
            <>
          {/* Favorites */}
          {favoriteSnippets.length > 0 && (
            <Collapsible
              open={expandedCategories.has('Favorites')}
              onOpenChange={() => toggleCategory('Favorites')}
            >
              <CollapsibleTrigger asChild>
                <Button variant="ghost" className="w-full justify-start p-0 h-auto">
                  <div className="flex items-center gap-2 py-2">
                    {expandedCategories.has('Favorites') ? 
                      <ChevronDown className="h-3 w-3" /> : 
                      <ChevronRight className="h-3 w-3" />
                    }
                    <Star className="h-3 w-3 text-yellow-500" />
                    <h4 className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                      Favorites ({favoriteSnippets.length})
                    </h4>
                  </div>
                </Button>
              </CollapsibleTrigger>
              <CollapsibleContent className="space-y-2 mt-2">
                {favoriteSnippets.map((snippet) => (
                  <SnippetCard 
                    key={snippet.id} 
                    snippet={snippet} 
                    onInsert={onSnippetInsert}
                    getCategoryIcon={getCategoryIcon}
                    selectedTags={selectedTags}
                    onTagClick={toggleTag}
                  />
                ))}
              </CollapsibleContent>
            </Collapsible>
          )}

          {favoriteSnippets.length > 0 && (
            (viewMode === 'categories' && Object.keys(snippetsByCategory).length > 0) ||
            (viewMode === 'tags' && Object.keys(snippetsByTag).length > 0)
          ) && (
            <Separator />
          )}

          {/* Snippets grouped by Categories or Tags */}
          {viewMode === 'categories' ? (
            // Organized by categories
            Object.entries(snippetsByCategory).map(([category, categorySnippets]) => (
              <Collapsible
                key={category}
                open={expandedCategories.has(category)}
                onOpenChange={() => toggleCategory(category)}
              >
                <CollapsibleTrigger asChild>
                  <Button variant="ghost" className="w-full justify-start p-0 h-auto">
                    <div className="flex items-center gap-2 py-2">
                      {expandedCategories.has(category) ? 
                        <ChevronDown className="h-3 w-3" /> : 
                        <ChevronRight className="h-3 w-3" />
                      }
                      {(() => {
                        const Icon = getCategoryIcon(category);
                        return <Icon className="h-3 w-3 text-muted-foreground" />;
                      })()}
                      <h4 className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                        {category} ({categorySnippets.length})
                      </h4>
                    </div>
                  </Button>
                </CollapsibleTrigger>
                <CollapsibleContent className="space-y-2 mt-2">
                  {categorySnippets.map((snippet) => (
                    <SnippetCard 
                      key={snippet.id} 
                      snippet={snippet} 
                      onInsert={onSnippetInsert}
                      getCategoryIcon={getCategoryIcon}
                      selectedTags={selectedTags}
                      onTagClick={toggleTag}
                    />
                  ))}
                </CollapsibleContent>
              </Collapsible>
            ))
          ) : (
            // Organized by tags (snippets can appear multiple times)
            Object.entries(snippetsByTag).map(([tag, tagSnippets]) => (
              <Collapsible
                key={tag}
                open={expandedCategories.has(tag)}
                onOpenChange={() => toggleCategory(tag)}
              >
                <CollapsibleTrigger asChild>
                  <Button variant="ghost" className="w-full justify-start p-0 h-auto">
                    <div className="flex items-center gap-2 py-2">
                      {expandedCategories.has(tag) ? 
                        <ChevronDown className="h-3 w-3" /> : 
                        <ChevronRight className="h-3 w-3" />
                      }
                      <Tag className="h-3 w-3 text-muted-foreground" />
                      <h4 className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                        #{tag} ({tagSnippets.length})
                      </h4>
                    </div>
                  </Button>
                </CollapsibleTrigger>
                <CollapsibleContent className="space-y-2 mt-2">
                  {tagSnippets.map((snippet) => (
                    <SnippetCard 
                      key={`${tag}-${snippet.id}`} 
                      snippet={snippet} 
                      onInsert={onSnippetInsert}
                      getCategoryIcon={getCategoryIcon}
                      selectedTags={selectedTags}
                      onTagClick={toggleTag}
                    />
                  ))}
                </CollapsibleContent>
              </Collapsible>
            ))
          )}

          {filteredSnippets.length === 0 && (
            <div className="text-center text-muted-foreground py-8">
              <FileText className="h-8 w-8 mx-auto mb-2 opacity-50" />
              <p className="text-sm">No snippets found</p>
              {(searchQuery || selectedTags.length > 0 || selectedCategory !== 'All') && (
                <div className="space-y-1">
                  <p className="text-xs">Try adjusting your filters:</p>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={clearFilters}
                    className="h-6 text-xs"
                  >
                    Clear all filters
                  </Button>
                </div>
              )}
            </div>
          )}
          </>
          )}
        </div>
      </ScrollArea>
    </Card>
  );
}

function SnippetCard({ 
  snippet, 
  onInsert, 
  getCategoryIcon,
  selectedTags,
  onTagClick
}: { 
  snippet: Snippet; 
  onInsert: (name: string) => void;
  getCategoryIcon: (category: string) => any;
  selectedTags: string[];
  onTagClick: (tag: string) => void;
}) {
  const Icon = getCategoryIcon(snippet.category);
  
  return (
    <div
      className="group p-3 rounded-lg border border-border hover:border-primary/50 cursor-pointer transition-all hover:shadow-xs bg-card hover:bg-primary/5"
      onClick={() => onInsert(snippet.title)}
    >
      <div className="flex items-start justify-between mb-2">
        <div className="flex items-center gap-2 min-w-0 flex-1">
          <Icon className="h-3 w-3 text-muted-foreground shrink-0" />
          <span className="font-mono text-sm font-medium truncate">
            @{snippet.title}
          </span>
          {snippet.isFavorite && (
            <Star className="h-3 w-3 text-yellow-500 shrink-0" />
          )}
        </div>
      </div>
      
      {snippet.description && (
        <p className="text-xs text-muted-foreground mb-2 line-clamp-2">
          {snippet.description}
        </p>
      )}
      
      <div className="flex items-center justify-between">
        <div className="flex flex-wrap gap-1">
          {snippet.tags.slice(0, 3).map((tag) => (
            <Badge 
              key={tag} 
              variant={selectedTags.includes(tag) ? "default" : "secondary"} 
              className="text-xs px-1 py-0 h-4 cursor-pointer hover:bg-primary/20 transition-colors"
              onClick={(e) => {
                e.stopPropagation();
                onTagClick(tag);
              }}
            >
              #{tag}
            </Badge>
          ))}
          {snippet.tags.length > 3 && (
            <Badge variant="outline" className="text-xs px-1 py-0 h-4">
              +{snippet.tags.length - 3}
            </Badge>
          )}
        </div>
        
        <Badge variant="outline" className="text-xs px-1 py-0 h-4">
          {snippet.category}
        </Badge>
      </div>
    </div>
  );
}