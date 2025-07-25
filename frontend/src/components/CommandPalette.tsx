import { useState, useEffect } from "react";
import { Command, CommandInput, CommandList, CommandEmpty, CommandGroup, CommandItem } from "@/components/ui/command";
import { Dialog, DialogContent } from "@/components/ui/dialog";
import { Search, FileText, Hash, Settings, Keyboard } from "lucide-react";

interface CommandPaletteProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const commands = [
  {
    group: "Navigation",
    items: [
      { id: "new-prompt", label: "New Prompt", icon: FileText, shortcut: "Ctrl+N" },
      { id: "search-prompts", label: "Search Prompts", icon: Search, shortcut: "Ctrl+P" },
      { id: "variables", label: "Show Variables", icon: Hash, shortcut: "Alt+V" },
      { id: "snippets", label: "Toggle Snippets", icon: FileText, shortcut: "Alt+S" },
    ]
  },
  {
    group: "Actions",
    items: [
      { id: "save", label: "Save Current Prompt", icon: FileText, shortcut: "Ctrl+S" },
      { id: "preview", label: "Toggle Preview", icon: FileText, shortcut: "Alt+P" },
      { id: "settings", label: "Open Settings", icon: Settings, shortcut: "Ctrl+," },
      { id: "shortcuts", label: "Keyboard Shortcuts", icon: Keyboard, shortcut: "?" },
    ]
  }
];

export function CommandPalette({ open, onOpenChange }: CommandPaletteProps) {
  const [search, setSearch] = useState("");

  useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        onOpenChange(!open);
      }
    };

    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, [open, onOpenChange]);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="overflow-hidden p-0 shadow-lg max-w-[640px]" style={{ viewTransitionName: 'command-palette' }}>
        <Command className="[&_[cmdk-group-heading]]:px-2 [&_[cmdk-group-heading]]:font-medium [&_[cmdk-group-heading]]:text-muted-foreground [&_[cmdk-group]:not([hidden])_~[cmdk-group]]:pt-0 [&_[cmdk-group]]:px-2 [&_[cmdk-input-wrapper]_svg]:h-5 [&_[cmdk-input-wrapper]_svg]:w-5 [&_[cmdk-input]]:h-12 [&_[cmdk-item]]:px-2 [&_[cmdk-item]]:py-3 [&_[cmdk-item]_svg]:h-5 [&_[cmdk-item]_svg]:w-5">
          <CommandInput 
            placeholder="Type a command or search..." 
            value={search}
            onValueChange={setSearch}
          />
          <CommandList>
            <CommandEmpty>No results found.</CommandEmpty>
            {commands.map((group) => (
              <CommandGroup key={group.group} heading={group.group}>
                {group.items.map((item) => {
                  const Icon = item.icon;
                  return (
                    <CommandItem
                      key={item.id}
                      value={item.id}
                      onSelect={() => {
                        onOpenChange(false);
                        // Handle command execution here
                        console.log(`Executing command: ${item.id}`);
                      }}
                      className="flex items-center gap-2 px-3 py-2"
                    >
                      <Icon className="h-4 w-4" />
                      <span className="flex-1">{item.label}</span>
                      <kbd className="pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-xs font-medium text-muted-foreground opacity-100">
                        {item.shortcut}
                      </kbd>
                    </CommandItem>
                  );
                })}
              </CommandGroup>
            ))}
          </CommandList>
        </Command>
      </DialogContent>
    </Dialog>
  );
}