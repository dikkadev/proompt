import { useState, useEffect, useRef, useCallback } from "react";
import { HexColorPicker } from "react-colorful";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Palette, Check, GripVertical, X } from "lucide-react";
import { applyThemeColor, getStoredAccentColor, getDefaultAccentColor } from "@/lib/colorUtils";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

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
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [isDragging, setIsDragging] = useState(false);
  const [dragOffset, setDragOffset] = useState({ x: 0, y: 0 });
  const debounceRef = useRef<NodeJS.Timeout>();
  const dialogRef = useRef<HTMLDivElement>(null);
  const dragHandleRef = useRef<HTMLDivElement>(null);
  const triggerRef = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    // Load stored color on mount
    const stored = getStoredAccentColor();
    if (stored) {
      setCurrentColor(stored);
      setTempColor(stored);
      applyThemeColor(stored);
    }
  }, []);

  // Initialize position when dialog opens for the first time
  useEffect(() => {
    if (isOpen && position.x === 0 && position.y === 0 && triggerRef.current) {
      const rect = triggerRef.current.getBoundingClientRect();
      const dialogWidth = 264; // As per comment in original code
      const dialogHeight = 400; // Approximate height

      let x = rect.right - dialogWidth;
      let y = rect.bottom + 8; // 8px below the button

      // Keep dialog within viewport bounds
      const maxX = window.innerWidth - dialogWidth;
      const maxY = window.innerHeight - dialogHeight;

      setPosition({
        x: Math.max(8, Math.min(x, maxX - 8)),
        y: Math.max(8, Math.min(y, maxY - 8)),
      });
    }
  }, [isOpen, position.x, position.y]);

  const handleCancel = useCallback(() => {
    setTempColor(currentColor);
    applyThemeColor(currentColor);
    setIsOpen(false);
  }, [currentColor]);

  // Handle click outside to close
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (isOpen && dialogRef.current && !dialogRef.current.contains(event.target as Node)) {
        handleCancel();
      }
    };

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
      return () => document.removeEventListener('mousedown', handleClickOutside);
    }
  }, [isOpen, handleCancel]);

  // Handle escape key to close
  useEffect(() => {
    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape' && isOpen) {
        handleCancel();
      }
    };

    if (isOpen) {
      document.addEventListener('keydown', handleEscape);
      return () => document.removeEventListener('keydown', handleEscape);
    }
  }, [isOpen, handleCancel]);

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

  const handleMouseDown = useCallback((event: React.MouseEvent) => {
    if (event.target === dragHandleRef.current || dragHandleRef.current?.contains(event.target as Node)) {
      setIsDragging(true);
      const rect = dialogRef.current?.getBoundingClientRect();
      if (rect) {
        setDragOffset({
          x: event.clientX - rect.left,
          y: event.clientY - rect.top,
        });
      }
    }
  }, []);

  const handleMouseMove = useCallback((event: MouseEvent) => {
    if (isDragging) {
      const newX = event.clientX - dragOffset.x;
      const newY = event.clientY - dragOffset.y;
      
      // Keep dialog within viewport bounds
      const dialogWidth = 264;
      const dialogHeight = 400;
      const maxX = window.innerWidth - dialogWidth;
      const maxY = window.innerHeight - dialogHeight;
      
      setPosition({
        x: Math.max(8, Math.min(newX, maxX - 8)),
        y: Math.max(8, Math.min(newY, maxY - 8)),
      });
    }
  }, [isDragging, dragOffset]);

  const handleMouseUp = useCallback(() => {
    setIsDragging(false);
  }, []);

  useEffect(() => {
    if (isDragging) {
      document.addEventListener('mousemove', handleMouseMove);
      document.addEventListener('mouseup', handleMouseUp);
      return () => {
        document.removeEventListener('mousemove', handleMouseMove);
        document.removeEventListener('mouseup', handleMouseUp);
      };
    }
  }, [isDragging, handleMouseMove, handleMouseUp]);

  return (
    <>
      {/* Trigger Button */}
      <Button
        ref={triggerRef}
        variant="ghost"
        size="sm"
        className="gap-2"
        onClick={() => setIsOpen(true)}
        aria-label="Choose accent color"
      >
        <Palette className="h-4 w-4" />
      </Button>

      {/* Custom Draggable Dialog */}
      {isOpen && (
        <>
          {/* Backdrop */}
          <div className="fixed inset-0 bg-black/20 z-40" />
          
          {/* Dialog */}
          <div
            ref={dialogRef}
            className="fixed z-50 w-64"
            style={{
              left: position.x,
              top: position.y,
              cursor: isDragging ? 'grabbing' : 'default',
            }}
            onMouseDown={handleMouseDown}
          >
            <Card className="bg-background/95 backdrop-blur-sm border-2 border-border shadow-lg">
              {/* Drag Handle Row */}
              <div 
                ref={dragHandleRef}
                className="flex items-center justify-between p-2 border-b border-border cursor-grab active:cursor-grabbing hover:bg-muted/50 transition-colors"
              >
                <GripVertical className="h-4 w-4 text-muted-foreground rotate-90" />
                <span className="text-xs font-medium text-muted-foreground select-none">Drag to move</span>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-6 w-6 p-0 hover:bg-destructive/10 hover:text-destructive"
                  onClick={handleCancel}
                  aria-label="Close color picker"
                >
                  <X className="h-3 w-3" />
                </Button>
              </div>

              {/* Main Content */}
              <div className="p-4 space-y-4">
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
                <div className="flex gap-2 pt-4 mt-4 border-t border-border">
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
              </div>
            </Card>
          </div>
        </>
      )}
    </>
  );
}
