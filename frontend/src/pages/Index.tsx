import { useState, useCallback, useEffect } from "react";
import { CommandPalette } from "@/components/CommandPalette";
import { PromptEditor } from "@/components/PromptEditor";
import { VariablePanel } from "@/components/VariablePanel";
import { SnippetSidebar } from "@/components/SnippetSidebar";
import { LivePreview } from "@/components/LivePreview";
import { ColorPicker } from "@/components/ColorPicker";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { ResizablePanelGroup, ResizablePanel, ResizableHandle } from "@/components/ui/resizable";
import { 
  Search, 
  Command, 
  Settings, 
  Moon, 
  Sun,
  PanelLeftClose,
  PanelLeftOpen,
  Wifi,
  WifiOff,
  Loader2
} from "lucide-react";
import { getStoredAccentColor, applyThemeColor, getDefaultAccentColor, initializeTheme, setStoredTheme } from "@/lib/colorUtils";
import { useHealthCheck } from "@/lib/queries";

interface Variable {
  name: string;
  value?: string;
  defaultValue?: string;
  state: 'provided' | 'default' | 'missing';
}

const Index = () => {
  const [commandPaletteOpen, setCommandPaletteOpen] = useState(false);
  const [variables, setVariables] = useState<Variable[]>([]);
  const [variableValues, setVariableValues] = useState<Record<string, string>>({});
  const [showPreview, setShowPreview] = useState(true);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [isDarkMode, setIsDarkMode] = useState(true); // Will be updated in useEffect
  const [promptContent, setPromptContent] = useState("");
  const [snippets, setSnippets] = useState<string[]>([]);
  
  // Health check
  const { 
    data: healthData, 
    isError: healthError, 
    isLoading: healthLoading,
    isFetching: healthFetching,
    error: healthErrorData
  } = useHealthCheck();

  // Initialize accent color and theme on mount
  useEffect(() => {
    const storedColor = getStoredAccentColor();
    if (storedColor) {
      applyThemeColor(storedColor);
    } else {
      applyThemeColor(getDefaultAccentColor());
    }
    
    // Initialize theme from cookie (defaults to dark mode)
    const currentTheme = initializeTheme();
    setIsDarkMode(currentTheme === 'dark');
  }, []);

  const handleVariablesChange = useCallback((newVariables: Variable[]) => {
    setVariables(newVariables);
  }, []);

  const handleVariableChange = useCallback((name: string, value: string) => {
    setVariableValues(prev => ({ ...prev, [name]: value }));
  }, []);

  const handleSnippetInsert = useCallback((snippetName: string) => {
    console.log(`Inserting snippet: ${snippetName}`);
    // This would be handled by the PromptEditor
  }, []);

  const toggleTheme = () => {
    const newTheme = isDarkMode ? 'light' : 'dark';
    setIsDarkMode(!isDarkMode);
    document.documentElement.classList.toggle('dark');
    setStoredTheme(newTheme);
  };

  const toggleSidebar = () => {
    setSidebarCollapsed(!sidebarCollapsed);
  };

  return (
    <div className={`h-screen w-full bg-background text-foreground overflow-hidden ${isDarkMode ? 'dark' : ''}`}>
      {/* Global Header */}
      <header className="h-14 border-b border-border bg-card px-4 flex items-center justify-between flex-shrink-0">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-lg bg-gradient-primary flex items-center justify-center">
              <span className="text-white font-bold text-sm">P</span>
            </div>
            <h1 className="text-lg font-semibold">Proompt</h1>
          </div>
          
          <Separator orientation="vertical" className="h-6" />
          
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={toggleSidebar}
              className="gap-2"
            >
              {sidebarCollapsed ? <PanelLeftOpen className="h-4 w-4" /> : <PanelLeftClose className="h-4 w-4" />}
              {sidebarCollapsed ? 'Show' : 'Hide'} Snippets
            </Button>
            
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setCommandPaletteOpen(true)}
              className="gap-2"
            >
              <Search className="h-4 w-4" />
              Search
              <kbd className="pointer-events-none inline-flex h-4 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-xs font-medium text-muted-foreground opacity-100">
                <Command className="h-2 w-2" />K
              </kbd>
            </Button>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {/* Backend Status Indicator */}
          <div 
            className="flex items-center gap-2 text-xs"
            title={
              healthError 
                ? `Backend offline: ${healthErrorData?.message || 'Connection failed'}` 
                : healthData 
                ? `Backend connected - Status: ${healthData.status}, Version: ${healthData.version}`
                : healthLoading 
                ? 'Checking backend connection...'
                : 'Connecting to backend...'
            }
          >
            <div className={`w-2 h-2 rounded-full ${
              healthError ? 'bg-red-500' : healthData ? 'bg-green-500' : 'bg-yellow-500'
            }`} />
            <span className="text-muted-foreground">
              {healthError ? 'Offline' : healthData ? 'Connected' : 'Connecting...'}
            </span>
          </div>
          
          <Button
            variant="ghost"
            size="sm"
            onClick={toggleTheme}
            className="gap-2"
          >
            {isDarkMode ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
          </Button>
          
          <Button variant="ghost" size="sm">
            <Settings className="h-4 w-4" />
          </Button>
        </div>
      </header>

      {/* Main Layout */}
      <div className="h-[calc(100vh-3.5rem)] w-full overflow-hidden">
        <ResizablePanelGroup direction="horizontal" className="h-full">
          {/* Left Sidebar - Snippets */}
          <ResizablePanel 
            defaultSize={20} 
            minSize={5} 
            maxSize={30}
            className="border-r border-border"
          >
            <SnippetSidebar
              onSnippetInsert={handleSnippetInsert}
              isCollapsed={sidebarCollapsed}
              onToggleCollapse={toggleSidebar}
            />
          </ResizablePanel>

          <ResizableHandle withHandle />

          {/* Main Content Area */}
          <ResizablePanel defaultSize={60} minSize={40} className="flex flex-col min-w-0 overflow-hidden">
          {showPreview ? (
            <ResizablePanelGroup direction="vertical" className="flex-1">
              {/* Editor */}
              <ResizablePanel defaultSize={60} minSize={30}>
                <PromptEditor
                  onVariablesChange={handleVariablesChange}
                  onPreviewToggle={setShowPreview}
                  showPreview={showPreview}
                  onContentChange={setPromptContent}
                  onSnippetsChange={setSnippets}
                />
              </ResizablePanel>
              
              <ResizableHandle withHandle />
              
              {/* Preview Panel */}
              <ResizablePanel defaultSize={40} minSize={20}>
                <LivePreview
                  content={promptContent}
                  variables={variables}
                  variableValues={variableValues}
                  isVisible={showPreview}
                />
              </ResizablePanel>
            </ResizablePanelGroup>
          ) : (
            <div className="flex-1 overflow-hidden">
              <PromptEditor
                onVariablesChange={handleVariablesChange}
                onPreviewToggle={setShowPreview}
                showPreview={showPreview}
                onContentChange={setPromptContent}
                onSnippetsChange={setSnippets}
              />
            </div>
          )}
          </ResizablePanel>

          <ResizableHandle withHandle />

          {/* Right Panel - Variables */}
          <ResizablePanel 
            defaultSize={20} 
            minSize={15} 
            maxSize={30}
            className="border-l border-border overflow-hidden"
          >
            <VariablePanel
              variables={variables}
              onVariableChange={handleVariableChange}
              snippets={snippets}
            />
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>

      {/* Command Palette */}
      <CommandPalette
        open={commandPaletteOpen}
        onOpenChange={setCommandPaletteOpen}
      />

      {/* Color Picker */}
      <ColorPicker />

      {/* TODO: Redesign keyboard shortcuts helper - currently too generic, needs better UX integration */}
      {/* 
      <div className="fixed bottom-4 right-4 text-xs text-muted-foreground bg-card border border-border rounded-lg px-3 py-2 shadow-lg">
        <div className="flex items-center gap-2">
          <kbd className="pointer-events-none inline-flex h-4 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-xs font-medium text-muted-foreground">
            <Command className="h-2 w-2" />K
          </kbd>
          <span>Command Palette</span>
          <span className="text-muted-foreground/50">•</span>
          <kbd className="pointer-events-none inline-flex h-4 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-xs font-medium text-muted-foreground">
            ?
          </kbd>
          <span>Help</span>
        </div>
      </div>
      */}
    </div>
  );
};

export default Index;
