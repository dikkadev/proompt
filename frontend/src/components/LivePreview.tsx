import { useState, useEffect } from "react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { StatusBadge } from "@/components/ui/status-badge";
import { LoadingState } from "@/components/ui/loading-state";
import { 
  Eye, 
  Copy, 
  Download, 
  AlertTriangle,
  Code,
  FileText,
  Hash,
  WifiOff
} from "lucide-react";
import { toast } from "@/hooks/use-toast";
import { useTemplatePreviewMutation } from "@/lib/queries";
import { debounce } from "@/lib/utils";

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
  const [previewVariables, setPreviewVariables] = useState<Array<{
    name: string;
    defaultValue?: string;
    hasDefault: boolean;
    status: 'provided' | 'default' | 'missing';
  }>>([]);
  const [warnings, setWarnings] = useState<string[]>([]);
  
  const templatePreviewMutation = useTemplatePreviewMutation();

  // Debounced preview function to avoid too many API calls
  const debouncedPreview = debounce(async (contentToPreview: string, vars: Record<string, string>) => {
    if (!contentToPreview.trim() || !isVisible) {
      setResolvedContent('');
      setPreviewVariables([]);
      setWarnings([]);
      return;
    }

    try {
      const result = await templatePreviewMutation.mutateAsync({
        content: contentToPreview,
        variables: vars
      });

      setResolvedContent(result.resolved_content);
      setPreviewVariables(result.variables.map(v => ({
        name: v.name,
        defaultValue: v.default_value,
        hasDefault: v.has_default,
        status: v.status
      })));
      setWarnings(result.warnings || []);
    } catch (error) {
      // Error is already handled by the mutation's onError
      setResolvedContent('');
      setPreviewVariables([]);
      setWarnings([]);
    }
  }, 300);

  useEffect(() => {
    if (isVisible && content) {
      debouncedPreview(content, variableValues);
    }
    
    return () => {
      debouncedPreview.cancel();
    };
  }, [content, variableValues, isVisible]);

  const handleCopy = () => {
    navigator.clipboard.writeText(resolvedContent);
    toast({
      title: "Copied to clipboard",
      description: "Preview content copied to clipboard.",
    });
  };

  const handleDownload = () => {
    const timestamp = new Date().toISOString().slice(0, 19).replace(/:/g, '-');
    const filename = `prompt-export-${timestamp}.md`;
    
    const blob = new Blob([resolvedContent], { type: 'text/markdown' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    
    toast({
      title: "Downloaded",
      description: `Preview exported as ${filename}`,
    });
  };

  // Function to highlight unresolved variables in the content
  const renderContentWithHighlights = (content: string) => {
    const variableRegex = /\{\{([^}]+)\}\}/g;
    const parts = [];
    let lastIndex = 0;
    let match;

    while ((match = variableRegex.exec(content)) !== null) {
      // Add text before the variable
      if (match.index > lastIndex) {
        parts.push(content.slice(lastIndex, match.index));
      }
      
      // Add the highlighted variable
      parts.push(
        <span 
          key={`var-${match.index}`}
          className="bg-variable-missing/10 text-variable-missing px-1 py-0.5 rounded text-xs font-medium border border-variable-missing/30"
          title="Missing variable value"
        >
          {match[0]}
        </span>
      );
      
      lastIndex = match.index + match[0].length;
    }
    
    // Add remaining text
    if (lastIndex < content.length) {
      parts.push(content.slice(lastIndex));
    }
    
    return parts.length > 1 ? parts : content;
  };

  const isLoading = templatePreviewMutation.isPending;
  const hasWarnings = warnings.length > 0;
  const hasUnresolvedVariables = previewVariables.some(v => v.status === 'missing');

  if (!isVisible) {
    return null;
  }

  return (
    <div className="flex flex-col h-full bg-workspace-preview overflow-hidden">
      {/* Header */}
      <div className="px-4 py-2 border-b border-border bg-card flex-shrink-0">
        {/* Mobile/Narrow: Two rows */}
        <div className="lg:hidden">
          {/* Title Row */}
          <div className="flex items-center justify-between mb-2">
            <div className="flex items-center gap-2">
              <Eye className="h-4 w-4 text-muted-foreground" />
              <h3 className="font-medium">Live Preview</h3>
              {isLoading && (
                <LoadingState size="sm" inline message="Generating..." />
              )}
            </div>
            
            <div className="flex items-center gap-2">
              <Button
                variant="ghost"
                size="sm"
                onClick={handleCopy}
                disabled={!resolvedContent || isLoading}
                className="gap-2 ghost-icon-button"
              >
                <Copy className="h-4 w-4" />
                Copy
              </Button>
              <Button
                variant="ghost"
                size="sm"
                onClick={handleDownload}
                disabled={!resolvedContent || isLoading}
                className="gap-2 ghost-icon-button"
              >
                <Download className="h-4 w-4" />
                Download
              </Button>
            </div>
          </div>
          
          {/* Pills Row */}
          <div className="flex items-center gap-2">
            {hasWarnings && (
              <StatusBadge status="warning">
                <AlertTriangle className="h-3 w-3" />
                {warnings.length} warning{warnings.length !== 1 ? 's' : ''}
              </StatusBadge>
            )}
            
            {hasUnresolvedVariables && (
              <StatusBadge status="missing">
                <AlertTriangle className="h-3 w-3" />
                Missing variables
              </StatusBadge>
            )}
          </div>
        </div>
        
        {/* Desktop/Wide: Single row */}
        <div className="hidden lg:flex lg:items-center lg:justify-between">
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2">
              <Eye className="h-4 w-4 text-muted-foreground" />
              <h3 className="font-medium">Live Preview</h3>
            </div>
            
            {isLoading && (
              <LoadingState size="sm" inline message="Generating..." />
            )}
            
            {hasWarnings && (
              <StatusBadge status="warning">
                <AlertTriangle className="h-3 w-3" />
                {warnings.length} warning{warnings.length !== 1 ? 's' : ''}
              </StatusBadge>
            )}
            
            {hasUnresolvedVariables && (
              <StatusBadge status="missing">
                <AlertTriangle className="h-3 w-3" />
                Missing variables
              </StatusBadge>
            )}
          </div>
          
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={handleCopy}
              disabled={!resolvedContent || isLoading}
              className="gap-2 ghost-icon-button"
            >
              <Copy className="h-4 w-4" />
              Copy
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={handleDownload}
              disabled={!resolvedContent || isLoading}
              className="gap-2 ghost-icon-button"
            >
              <Download className="h-4 w-4" />
              Download
            </Button>
          </div>
        </div>
      </div>

      {/* Warnings Section */}
      {hasWarnings && (
        <div className="p-3 bg-warning/10 border-b border-warning/20 flex-shrink-0">
          <div className="flex items-start gap-2">
            <AlertTriangle className="h-4 w-4 text-warning mt-0.5 flex-shrink-0" />
            <div className="space-y-1">
              <p className="text-sm font-medium text-warning">Warnings</p>
              <ul className="text-xs text-warning/80 space-y-1">
                {warnings.map((warning, index) => (
                  <li key={index}>• {warning}</li>
                ))}
              </ul>
            </div>
          </div>
        </div>
      )}

             {/* Variable Status Bar */}
       {previewVariables.length > 0 && (
         <div className="p-3 bg-muted/30 border-b border-border flex-shrink-0">
           <div className="flex items-center gap-2 text-xs text-muted-foreground mb-2">
             <Hash className="h-3 w-3" />
             Variables ({previewVariables.length})
           </div>
           <div className="flex flex-wrap gap-1">
             {previewVariables.map((variable) => (
               <Badge
                 key={variable.name}
                 variant="outline"
                 className={`text-xs px-2 py-0.5 ${
                   variable.status === 'provided' 
                     ? 'border-variable-provided/30 text-variable-provided bg-variable-provided/10'
                     : variable.status === 'default'
                     ? 'border-variable-default/30 text-variable-default bg-variable-default/10'
                     : 'border-variable-missing/30 text-variable-missing bg-variable-missing/10'
                 }`}
                 title={
                   variable.status === 'provided' 
                     ? `Provided: ${variableValues[variable.name]}`
                     : variable.status === 'default'
                     ? `Using default: ${variable.defaultValue}`
                     : 'Missing value'
                 }
               >
                 {variable.name}
               </Badge>
             ))}
           </div>
         </div>
       )}

      {/* Content */}
      <div className="flex-1 relative overflow-hidden">
        <ScrollArea className="h-full">
          <div className="p-4">
            {isLoading ? (
              <div className="space-y-2">
                {[...Array(6)].map((_, i) => (
                  <div 
                    key={i}
                    className="h-4 bg-muted rounded animate-pulse"
                    style={{ width: `${60 + Math.random() * 40}%` }}
                  />
                ))}
              </div>
            ) : resolvedContent ? (
              <div className="prose prose-sm max-w-none dark:prose-invert">
                <pre className="whitespace-pre-wrap font-mono text-sm leading-relaxed text-foreground bg-transparent border-none p-0">
                  {renderContentWithHighlights(resolvedContent)}
                </pre>
              </div>
            ) : content ? (
              <div className="text-center text-muted-foreground py-8">
                <FileText className="h-8 w-8 mx-auto mb-2 opacity-50" />
                <p className="text-sm">Failed to generate preview</p>
                <p className="text-xs">Check your template syntax and backend connection</p>
              </div>
            ) : (
              <div className="text-center text-muted-foreground py-8">
                <FileText className="h-8 w-8 mx-auto mb-2 opacity-50" />
                <p className="text-sm">Start typing to see a live preview...</p>
              </div>
            )}
          </div>
        </ScrollArea>
      </div>

      {/* Footer Stats */}
      {resolvedContent && (
        <div className="flex items-center justify-between p-3 border-t border-border bg-card text-xs text-muted-foreground flex-shrink-0">
          <div className="flex items-center gap-4">
            <span>{resolvedContent.length} characters</span>
            <span>{resolvedContent.split('\n').length} lines</span>
            <span>~{Math.ceil(resolvedContent.length / 4)} tokens</span>
          </div>
          
                                <div className="flex items-center gap-2">
             {previewVariables.filter(v => v.status === 'provided').length > 0 && (
               <div className="flex items-center gap-1">
                 <div className="w-2 h-2 rounded-full bg-variable-provided"></div>
                 <span>{previewVariables.filter(v => v.status === 'provided').length} provided</span>
               </div>
             )}
             {previewVariables.filter(v => v.status === 'default').length > 0 && (
               <div className="flex items-center gap-1">
                 <div className="w-2 h-2 rounded-full bg-variable-default"></div>
                 <span>{previewVariables.filter(v => v.status === 'default').length} default</span>
               </div>
             )}
             {previewVariables.filter(v => v.status === 'missing').length > 0 && (
               <div className="flex items-center gap-1">
                 <div className="w-2 h-2 rounded-full bg-variable-missing"></div>
                 <span>{previewVariables.filter(v => v.status === 'missing').length} missing</span>
               </div>
             )}
           </div>
        </div>
      )}
    </div>
  );
}