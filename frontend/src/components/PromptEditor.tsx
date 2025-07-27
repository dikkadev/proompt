import { useState, useRef, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Copy, Save, Eye, EyeOff } from "lucide-react";
import { toast } from "@/hooks/use-toast";
import { useCreatePrompt, useUpdatePrompt } from "@/lib/queries";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

interface Variable {
  name: string;
  value?: string;
  defaultValue?: string;
  state: 'provided' | 'default' | 'missing';
}

interface PromptEditorProps {
  onVariablesChange: (variables: Variable[]) => void;
  onPreviewToggle: (show: boolean) => void;
  showPreview: boolean;
  onContentChange: (content: string) => void;
  onSnippetsChange: (snippets: string[]) => void;
  promptId?: string; // For editing existing prompts
  initialTitle?: string;
  initialContent?: string;
}

export function PromptEditor({ onVariablesChange, onPreviewToggle, showPreview, onContentChange, onSnippetsChange, promptId, initialTitle, initialContent }: PromptEditorProps) {
  const [title, setTitle] = useState(initialTitle || "Untitled Prompt");
  const [content, setContent] = useState(
    initialContent || `# Example Prompt Template

Please analyze this {{document_type:document}} and provide insights about {{topic}}.

Key areas to focus on:
- {{focus_area_1}}
- {{focus_area_2:user engagement}}
- {{focus_area_3:performance metrics}}

@analysis_guidelines

Please format your response as:
@response_template`
  );
  
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const [variables, setVariables] = useState<Variable[]>([]);
  const [snippets, setSnippets] = useState<string[]>([]);
  
  // Mutations for saving prompts
  const createPromptMutation = useCreatePrompt();
  const updatePromptMutation = useUpdatePrompt();

  // Parse variables from content and notify parent
  useEffect(() => {
    const variableRegex = /\{\{([^}]+)\}\}/g;
    const found: Variable[] = [];
    const snippetRegex = /@(\w+|\{[^}]+\})/g;
    const foundSnippets: string[] = [];
    
    let match;
    
    // Extract variables
    while ((match = variableRegex.exec(content)) !== null) {
      const [, varContent] = match;
      const [name, defaultValue] = varContent.split(':');
      
      const existing = found.find(v => v.name === name);
      if (!existing) {
        found.push({
          name: name.trim(),
          defaultValue: defaultValue?.trim(),
          state: defaultValue ? 'default' : 'missing'
        });
      }
    }
    
    // Extract snippets
    while ((match = snippetRegex.exec(content)) !== null) {
      const [, snippetName] = match;
      const cleanName = snippetName.replace(/[{}]/g, '');
      if (!foundSnippets.includes(cleanName)) {
        foundSnippets.push(cleanName);
      }
    }
    
    setVariables(found);
    
    // Use view transition for snippet changes that might affect layout
    if (foundSnippets.length !== snippets.length) {
      document.startViewTransition(() => {
        setSnippets(foundSnippets);
        onSnippetsChange(foundSnippets);
      });
    } else {
      setSnippets(foundSnippets);
      onSnippetsChange(foundSnippets);
    }
    
    onVariablesChange(found);
    onContentChange(content);
  }, [content, onVariablesChange, onContentChange, onSnippetsChange]);

  const handleSave = async () => {
    if (!title.trim()) {
      toast({
        title: "Title required",
        description: "Please enter a title for your prompt.",
        variant: "destructive",
      });
      return;
    }

    if (!content.trim()) {
      toast({
        title: "Content required", 
        description: "Please enter some content for your prompt.",
        variant: "destructive",
      });
      return;
    }

    try {
      if (promptId) {
        // Update existing prompt
        await updatePromptMutation.mutateAsync({
          id: promptId,
          title: title.trim(),
          content: content.trim(),
          type: 'user' as const, // Default to user type
        });
      } else {
        // Create new prompt
        await createPromptMutation.mutateAsync({
          title: title.trim(),
          content: content.trim(),
          type: 'user' as const, // Default to user type
          model_compatibility_tags: [],
        });
      }
    } catch (error) {
      // Error handling is done in the mutation's onError callback
      console.error('Save failed:', error);
    }
  };

  const handleCopy = () => {
    navigator.clipboard.writeText(content);
    toast({
      title: "Copied to clipboard",
      description: "Prompt content copied to clipboard.",
    });
  };



  const insertSnippet = (snippetName: string) => {
    const textarea = textareaRef.current;
    if (!textarea) return;

    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;
    const before = content.substring(0, start);
    const after = content.substring(end);
    const snippet = `@${snippetName}`;
    
    setContent(before + snippet + after);
    
    // Set cursor position after the inserted snippet
    setTimeout(() => {
      textarea.focus();
      textarea.setSelectionRange(start + snippet.length, start + snippet.length);
    }, 0);
  };

  return (
    <div className="flex flex-col h-full">
      {/* Toolbar */}
      <div className="px-4 py-2 border-b border-border bg-card">
        {/* Mobile/Narrow: Two rows */}
        <div className="lg:hidden">
          {/* Title Row */}
          <div className="flex items-center gap-2 flex-1 min-w-0 mb-2">
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Input
                    style={{viewTransitionName: 'prompt-title'}}
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    placeholder="Enter prompt title..."
                    className="font-semibold text-lg bg-transparent border-none shadow-none focus-visible:ring-0 focus-visible:ring-offset-0 px-1 h-auto py-0 flex-1 min-w-0 hover:bg-muted/20 focus:bg-muted/30 rounded transition-colors"
                  />
                </TooltipTrigger>
                <TooltipContent>
                  <p>Enter a descriptive title for your prompt</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </div>
          
          {/* Buttons Row */}
          <div className="flex items-center justify-end gap-2 flex-shrink-0">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onPreviewToggle(!showPreview)}
              className="gap-2 ghost-icon-button"
            >
              {showPreview ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
              {showPreview ? 'Hide Preview' : 'Show Preview'}
            </Button>
            <Button variant="ghost" size="sm" onClick={handleCopy} className="gap-2 ghost-icon-button">
              <Copy className="h-4 w-4" />
              Copy
            </Button>
            <Button 
              variant="ghost" 
              size="sm" 
              onClick={handleSave} 
              className="gap-2 ghost-icon-button"
              disabled={createPromptMutation.isPending || updatePromptMutation.isPending}
            >
              <Save className="h-4 w-4" />
              {promptId ? 'Update' : 'Save'}
            </Button>
          </div>
        </div>
        
        {/* Desktop/Wide: Single row */}
        <div className="hidden lg:flex lg:items-center lg:justify-between">
          <div className="flex items-center gap-2 flex-1 min-w-0">
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Input
                    style={{viewTransitionName: 'prompt-title'}}
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    placeholder="Enter prompt title..."
                    className="font-semibold text-lg bg-transparent border-none shadow-none focus-visible:ring-0 focus-visible:ring-offset-0 px-1 h-auto py-0 flex-1 min-w-0 hover:bg-muted/20 focus:bg-muted/30 rounded transition-colors"
                  />
                </TooltipTrigger>
                <TooltipContent>
                  <p>Enter a descriptive title for your prompt</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </div>
          
          <div className="flex items-center gap-2 flex-shrink-0">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onPreviewToggle(!showPreview)}
              className="gap-2 ghost-icon-button"
            >
              {showPreview ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
              {showPreview ? 'Hide Preview' : 'Show Preview'}
            </Button>
            <Button variant="ghost" size="sm" onClick={handleCopy} className="gap-2 ghost-icon-button">
              <Copy className="h-4 w-4" />
              Copy
            </Button>
            <Button 
              variant="ghost" 
              size="sm" 
              onClick={handleSave} 
              className="gap-2 ghost-icon-button"
              disabled={createPromptMutation.isPending || updatePromptMutation.isPending}
            >
              <Save className="h-4 w-4" />
              {promptId ? 'Update' : 'Save'}
            </Button>
          </div>
        </div>
      </div>

      {/* Editor */}
      <div className="flex-1 p-4">
        <Card className="h-full p-4 bg-workspace-editor">
          <Textarea
            ref={textareaRef}
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="Start typing your prompt template..."
            className="min-h-full resize-none border-none bg-transparent font-mono text-sm leading-relaxed focus-visible:ring-0 focus-visible:ring-offset-0"
            style={{ height: 'calc(100% - 2rem)' }}
          />
        </Card>
      </div>


    </div>
  );
}