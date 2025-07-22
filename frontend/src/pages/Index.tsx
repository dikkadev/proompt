import { useState, useCallback, useEffect } from "react";
import { CommandPalette } from "@/components/CommandPalette";
import { PromptEditor } from "@/components/PromptEditor";
import { VariablePanel } from "@/components/VariablePanel";
import { SnippetSidebar } from "@/components/SnippetSidebar";
import { LivePreview } from "@/components/LivePreview";
import { ColorPicker } from "@/components/ColorPicker";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { 
  Search, 
  Command, 
  Settings, 
  Moon, 
  Sun,
  PanelLeftClose,
  PanelLeftOpen
} from "lucide-react";
import { getStoredAccentColor, applyThemeColor, getDefaultAccentColor } from "@/lib/colorUtils";

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
  const [isDarkMode, setIsDarkMode] = useState(true);

  // Initialize accent color on mount
  useEffect(() => {
    const storedColor = getStoredAccentColor();
    if (storedColor) {
      applyThemeColor(storedColor);
    } else {
      applyThemeColor(getDefaultAccentColor());
    }
    
    // Initialize dark mode
    document.documentElement.classList.add('dark');
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
    setIsDarkMode(!isDarkMode);
    document.documentElement.classList.toggle('dark');
  };

  const toggleSidebar = () => {
    setSidebarCollapsed(!sidebarCollapsed);
  };

  return (
    <div className={`min-h-screen w-full bg-background text-foreground ${isDarkMode ? 'dark' : ''}`}>
      {/* Global Header */}
      <header className="h-14 border-b border-border bg-card px-4 flex items-center justify-between">
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
      <div className="flex h-[calc(100vh-3.5rem)] w-full">
        {/* Left Sidebar - Snippets */}
        <div className={`border-r border-border transition-all duration-300 ${
          sidebarCollapsed ? 'w-14' : 'w-80'
        }`}>
          <SnippetSidebar
            onSnippetInsert={handleSnippetInsert}
            isCollapsed={sidebarCollapsed}
            onToggleCollapse={toggleSidebar}
          />
        </div>

        {/* Main Content Area */}
        <div className="flex-1 flex flex-col min-w-0">
          {/* Editor */}
          <div className={`${showPreview ? 'h-3/5' : 'flex-1'} border-b border-border`}>
            <PromptEditor
              onVariablesChange={handleVariablesChange}
              onPreviewToggle={setShowPreview}
              showPreview={showPreview}
            />
          </div>

          {/* Preview Panel */}
          {showPreview && (
            <div className="h-2/5">
              <LivePreview
                content="" // This would come from the editor
                variables={variables}
                variableValues={variableValues}
                isVisible={showPreview}
              />
            </div>
          )}
        </div>

        {/* Right Panel - Variables */}
        <div className="w-80 border-l border-border">
          <VariablePanel
            variables={variables}
            onVariableChange={handleVariableChange}
          />
        </div>
      </div>

      {/* Command Palette */}
      <CommandPalette
        open={commandPaletteOpen}
        onOpenChange={setCommandPaletteOpen}
      />

      {/* Color Picker */}
      <ColorPicker />

      {/* Keyboard Shortcuts Help */}
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
    </div>
  );
};

export default Index;
