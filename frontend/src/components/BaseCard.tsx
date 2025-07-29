import { Card } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { GripVertical } from "lucide-react";

interface BaseCardProps {
  title: string;
  description?: string;
  tags?: string[];
  children?: React.ReactNode;
  className?: string;
  type: "prompt" | "snippet"; // To differentiate styling
}

export function BaseCard({
  title,
  description,
  tags,
  children,
  className,
  type,
}: BaseCardProps) {
  const isSnippet = type === "snippet";

  return (
    <Card 
      className={cn(
        "p-1.5 cursor-pointer transition-all duration-200 w-full relative",
        isSnippet ? "border border-transparent hover:border-accent bg-background cursor-grab active:cursor-grabbing" : "border bg-background hover:border-accent",
        className
      )}
    >
      <div className="flex items-start gap-0.5">
        <div className="flex-1 min-w-0 flex flex-col">
          <h4 className={cn("text-sm font-medium break-all mb-0", isSnippet ? "font-mono" : "")}>
            {isSnippet && <span className="text-accent mr-0.5">@</span>}{title}
          </h4>
          
          <div className="flex justify-between items-start gap-2 mt-1">
            {description && (
              <p className="text-xs text-muted-foreground line-clamp-2 flex-grow min-w-0">
                {description}
              </p>
            )}
            {isSnippet && (
              <GripVertical className="h-5 w-5 text-muted-foreground/60 hover:text-accent flex-shrink-0" />
            )}
          </div>
          
          <div className="flex items-end justify-between text-xs mt-auto pt-1">
            <div className="flex flex-wrap gap-0.5 min-w-0">
              {tags && tags.slice(0, 2).map((tag, index) => (
                <span 
                  key={`${tag}-${index}`}
                  className="px-1 py-0.5 bg-secondary/60 text-secondary-foreground rounded text-xs break-all"
                >
                  #{tag}
                </span>
              ))}
              {tags && tags.length > 2 && (
                <span className="text-muted-foreground">
                  +{tags.length - 2}
                </span>
              )}
            </div>
            
            <div className="flex items-center">
                {children}
            </div>
          </div>
        </div>
      </div>
    </Card>
  );
}
