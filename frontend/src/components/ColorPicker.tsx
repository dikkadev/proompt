import { useState, useEffect, useRef } from "react";
import { HexColorPicker } from "react-colorful";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Palette, Check } from "lucide-react";
import { applyThemeColor, getStoredAccentColor, getDefaultAccentColor } from "@/lib/colorUtils";

const presetColors = [
  '#1e9fa2', // Default teal
  '#3b82f6', // Blue
  '#8b5cf6', // Purple
  '#ef4444', // Red
  '#f59e0b', // Amber
  '#10b981', // Emerald
  '#f97316', // Orange
  '#6366f1', // Indigo
  '#ec4899', // Pink
  '#84cc16', // Lime
];

export function ColorPicker() {
  const [currentColor, setCurrentColor] = useState(getDefaultAccentColor());
  const [isOpen, setIsOpen] = useState(false);
  const [tempColor, setTempColor] = useState(currentColor);
  const debounceRef = useRef<NodeJS.Timeout>();

  useEffect(() => {
    // Load stored color on mount
    const stored = getStoredAccentColor();
    if (stored) {
      setCurrentColor(stored);
      setTempColor(stored);
      applyThemeColor(stored);
    }
  }, []);

  const handleColorChange = (newColor: string) => {
    setTempColor(newColor);
    
    // Debounce the theme application for performance
    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }
    
    debounceRef.current = setTimeout(() => {
      applyThemeColor(newColor);
    }, 100);
  };

  const handleApply = () => {
    setCurrentColor(tempColor);
    applyThemeColor(tempColor);
    setIsOpen(false);
  };

  const handlePresetClick = (color: string) => {
    setTempColor(color);
    handleColorChange(color);
  };

  const handleCancel = () => {
    setTempColor(currentColor);
    applyThemeColor(currentColor);
    setIsOpen(false);
  };

  return (
    <div className="fixed bottom-4 left-4 z-50">
      <Popover open={isOpen} onOpenChange={setIsOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="ghost"
            size="sm"
            className="h-4 w-4 p-0 rounded-full border border-border/50 hover:border-primary/50 transition-all hover:scale-110"
            style={{ backgroundColor: currentColor }}
            aria-label="Choose accent color"
          >
          </Button>
        </PopoverTrigger>
        
        <PopoverContent 
          className="w-64 p-4" 
          side="top" 
          align="start"
          sideOffset={8}
        >
          <Card className="p-4 space-y-4">
            <div className="space-y-2">
              <h4 className="font-medium text-sm">Accent Color</h4>
              <p className="text-xs text-muted-foreground">
                Choose a color to personalize your interface
              </p>
            </div>
            
            {/* Color Picker */}
            <div className="space-y-3">
              <HexColorPicker 
                color={tempColor} 
                onChange={handleColorChange}
                style={{ width: '100%', height: '120px' }}
              />
              
              {/* Current Color Display */}
              <div className="flex items-center gap-2 text-sm">
                <div 
                  className="w-6 h-6 rounded border border-border"
                  style={{ backgroundColor: tempColor }}
                />
                <code className="font-mono text-xs bg-muted px-2 py-1 rounded">
                  {tempColor.toUpperCase()}
                </code>
              </div>
            </div>
            
            {/* Preset Colors */}
            <div className="space-y-2">
              <h5 className="text-xs font-medium text-muted-foreground">Presets</h5>
              <div className="grid grid-cols-5 gap-2">
                {presetColors.map((color) => (
                  <button
                    key={color}
                    className="w-8 h-8 rounded border-2 border-border hover:border-foreground transition-colors relative group"
                    style={{ backgroundColor: color }}
                    onClick={() => handlePresetClick(color)}
                    aria-label={`Select ${color} color`}
                  >
                    {tempColor.toLowerCase() === color.toLowerCase() && (
                      <Check className="h-3 w-3 text-white absolute inset-0 m-auto drop-shadow-xs" />
                    )}
                  </button>
                ))}
              </div>
            </div>
            
            {/* Actions */}
            <div className="flex gap-2 pt-2 border-t border-border">
              <Button 
                size="sm" 
                onClick={handleApply}
                className="flex-1"
              >
                Apply
              </Button>
              <Button 
                size="sm" 
                variant="ghost" 
                onClick={handleCancel}
              >
                Cancel
              </Button>
            </div>
          </Card>
        </PopoverContent>
      </Popover>
    </div>
  );
}