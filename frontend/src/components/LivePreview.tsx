import { useState, useEffect } from "react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { 
  Eye, 
  Copy, 
  Download, 
  RefreshCw,
  AlertTriangle,
  Code,
  FileText
} from "lucide-react";
import { toast } from "@/hooks/use-toast";

interface Variable {
  name: string;
  value?: string;
  defaultValue?: string;
  state: 'provided' | 'default' | 'missing';
}

interface LivePreviewProps {
  content: string;
  variables: Variable[];
  variableValues: Record<string, string>;
  isVisible: boolean;
}

export function LivePreview({ content, variables, variableValues, isVisible }: LivePreviewProps) {
  const [resolvedContent, setResolvedContent] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [unresolvedVariables, setUnresolvedVariables] = useState<string[]>([]);

  useEffect(() => {
    if (!isVisible) return;
    
    setIsLoading(true);
    
    // Simulate API call delay
    const timer = setTimeout(() => {
      let resolved = content;
      const unresolved: string[] = [];
      
      // Replace variables with values
      const variableRegex = /\{\{([^}]+)\}\}/g;
      resolved = resolved.replace(variableRegex, (match, varContent) => {
        const [name, defaultValue] = varContent.split(':');
        const trimmedName = name.trim();
        
        if (variableValues[trimmedName]) {
          return variableValues[trimmedName];
        } else if (defaultValue) {
          return defaultValue.trim();
        } else {
          unresolved.push(trimmedName);
          return `{{${trimmedName}}}`;
        }
      });
      
      // Replace snippets with mock content
      const snippetRegex = /@(\w+|\{[^}]+\})/g;
      resolved = resolved.replace(snippetRegex, (match, snippetName) => {
        const cleanName = snippetName.replace(/[{}]/g, '');
        return `[Snippet: ${cleanName}]`;
      });
      
      setResolvedContent(resolved);
      setUnresolvedVariables(unresolved);
      setIsLoading(false);
    }, 300);
    
    return () => clearTimeout(timer);
  }, [content, variableValues, isVisible]);

  const handleCopy = () => {
    navigator.clipboard.writeText(resolvedContent);
    toast({
      title: "Copied to clipboard",
      description: "Preview content copied to clipboard.",
    });
  };

  const handleDownload = () => {
    const blob = new Blob([resolvedContent], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'prompt-preview.txt';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    
    toast({
      title: "Downloaded",
      description: "Preview saved as prompt-preview.txt",
    });
  };

  const handleRefresh = () => {
    setIsLoading(true);
    setTimeout(() => setIsLoading(false), 500);
  };

  if (!isVisible) {
    return null;
  }

  const completionPercentage = variables.length > 0 
    ? Math.round(((variables.length - unresolvedVariables.length) / variables.length) * 100)
    : 100;

  return (
    <Card className="h-full flex flex-col bg-workspace-preview">
      {/* Header */}
      <div className="p-4 border-b border-border">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2">
            <Eye className="h-5 w-5 text-primary" />
            <h3 className="font-semibold">Live Preview</h3>
            <Badge variant="secondary" className="ml-2">
              {completionPercentage}% Complete
            </Badge>
          </div>
          
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={handleRefresh}
              disabled={isLoading}
              className="gap-2"
            >
              <RefreshCw className={`h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={handleCopy}
              disabled={isLoading}
              className="gap-2"
            >
              <Copy className="h-4 w-4" />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={handleDownload}
              disabled={isLoading}
              className="gap-2"
            >
              <Download className="h-4 w-4" />
            </Button>
          </div>
        </div>

        {/* Status Indicators */}
        <div className="flex items-center justify-between text-sm">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-1">
              <div className="w-2 h-2 rounded-full bg-variable-provided"></div>
              <span className="text-muted-foreground">
                {variables.length - unresolvedVariables.length} resolved
              </span>
            </div>
            {unresolvedVariables.length > 0 && (
              <div className="flex items-center gap-1">
                <div className="w-2 h-2 rounded-full bg-variable-missing"></div>
                <span className="text-muted-foreground">
                  {unresolvedVariables.length} missing
                </span>
              </div>
            )}
          </div>
          
          {isLoading && (
            <div className="flex items-center gap-2 text-muted-foreground">
              <RefreshCw className="h-3 w-3 animate-spin" />
              <span className="text-xs">Updating...</span>
            </div>
          )}
        </div>
      </div>

      {/* Warnings */}
      {unresolvedVariables.length > 0 && (
        <div className="p-3 bg-warning/10 border-b border-warning/20">
          <div className="flex items-start gap-2">
            <AlertTriangle className="h-4 w-4 text-warning flex-shrink-0 mt-0.5" />
            <div className="text-sm">
              <p className="text-warning font-medium mb-1">
                Unresolved Variables
              </p>
              <p className="text-muted-foreground text-xs">
                The following variables need values: {unresolvedVariables.join(', ')}
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Preview Content */}
      <div className="flex-1 relative">
        <ScrollArea className="h-full">
          <div className="p-4">
            {isLoading ? (
              <div className="space-y-2">
                {[...Array(6)].map((_, i) => (
                  <div 
                    key={i}
                    className="h-4 bg-muted rounded animate-pulse-soft"
                    style={{ width: `${60 + Math.random() * 40}%` }}
                  />
                ))}
              </div>
            ) : (
              <div className="prose prose-sm max-w-none">
                <pre className="whitespace-pre-wrap font-mono text-sm leading-relaxed text-foreground bg-transparent border-none p-0">
                  {resolvedContent || 'Preview will appear here as you type...'}
                </pre>
              </div>
            )}
          </div>
        </ScrollArea>

        {/* Loading Overlay */}
        {isLoading && (
          <div className="absolute inset-0 bg-background/50 flex items-center justify-center">
            <div className="flex items-center gap-2 text-muted-foreground">
              <RefreshCw className="h-4 w-4 animate-spin" />
              <span className="text-sm">Generating preview...</span>
            </div>
          </div>
        )}
      </div>

      {/* Footer Stats */}
      <div className="p-3 border-t border-border bg-muted/50">
        <div className="flex justify-between items-center text-xs text-muted-foreground">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-1">
              <FileText className="h-3 w-3" />
              <span>{resolvedContent.split(' ').length} words</span>
            </div>
            <div className="flex items-center gap-1">
              <Code className="h-3 w-3" />
              <span>{resolvedContent.length} characters</span>
            </div>
          </div>
          
          <div className="flex items-center gap-1">
            <div className="w-1 h-1 rounded-full bg-primary animate-pulse"></div>
            <span>Live</span>
          </div>
        </div>
      </div>
    </Card>
  );
}