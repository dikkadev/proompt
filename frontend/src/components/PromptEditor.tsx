import { useState, useRef, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import { Badge } from "@/components/ui/badge";
import { Play, Copy, Save, Eye, EyeOff } from "lucide-react";
import { toast } from "@/hooks/use-toast";

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
}

export function PromptEditor({ onVariablesChange, onPreviewToggle, showPreview }: PromptEditorProps) {
  const [content, setContent] = useState(
    `# Example Prompt Template

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

  // Parse variables from content
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
    setSnippets(foundSnippets);
    onVariablesChange(found);
  }, [content, onVariablesChange]);

  const handleSave = () => {
    toast({
      title: "Prompt saved",
      description: "Your prompt template has been saved successfully.",
    });
  };

  const handleCopy = () => {
    navigator.clipboard.writeText(content);
    toast({
      title: "Copied to clipboard",
      description: "Prompt content copied to clipboard.",
    });
  };

  const handleRun = () => {
    toast({
      title: "Processing prompt",
      description: "Generating preview with current variables...",
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
      <div className="flex items-center justify-between p-4 border-b border-border bg-card">
        <div className="flex items-center gap-2">
          <h2 className="text-lg font-semibold">Prompt Editor</h2>
          <div className="flex gap-1">
            {variables.map((variable) => (
              <Badge 
                key={variable.name} 
                variant="outline"
                className={`text-xs ${
                  variable.state === 'provided' ? 'border-variable-provided text-variable-provided' :
                  variable.state === 'default' ? 'border-variable-default text-variable-default' :
                  'border-variable-missing text-variable-missing'
                }`}
              >
                {variable.name}
              </Badge>
            ))}
          </div>
        </div>
        
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onPreviewToggle(!showPreview)}
            className="gap-2"
          >
            {showPreview ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
            {showPreview ? 'Hide Preview' : 'Show Preview'}
          </Button>
          <Button variant="ghost" size="sm" onClick={handleCopy} className="gap-2">
            <Copy className="h-4 w-4" />
            Copy
          </Button>
          <Button variant="ghost" size="sm" onClick={handleSave} className="gap-2">
            <Save className="h-4 w-4" />
            Save
          </Button>
          <Button onClick={handleRun} size="sm" className="gap-2 bg-primary hover:bg-primary-hover">
            <Play className="h-4 w-4" />
            Run
          </Button>
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

      {/* Quick snippet insertion */}
      {snippets.length > 0 && (
        <div className="p-4 border-t border-border bg-muted/50">
          <div className="flex items-center gap-2 text-sm text-muted-foreground mb-2">
            Referenced snippets:
          </div>
          <div className="flex flex-wrap gap-2">
            {snippets.map((snippet) => (
              <Button
                key={snippet}
                variant="secondary"
                size="sm"
                onClick={() => insertSnippet(snippet)}
                className="text-xs h-6 px-2"
              >
                @{snippet}
              </Button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}