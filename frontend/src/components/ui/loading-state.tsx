import * as React from "react"
import { RefreshCw } from "lucide-react"
import { cn } from "@/lib/utils"

interface LoadingStateProps {
  message?: string;
  size?: 'sm' | 'md' | 'lg';
  inline?: boolean;
  className?: string;
}

const getSizeClasses = (size: 'sm' | 'md' | 'lg'): string => {
  switch (size) {
    case 'sm':
      return "h-3 w-3";
    case 'md':
      return "h-4 w-4";
    case 'lg':
      return "h-5 w-5";
    default:
      return "h-4 w-4";
  }
};

function LoadingState({ message = "Loading...", size = 'md', inline = false, className }: LoadingStateProps) {
  const iconClasses = cn(getSizeClasses(size), "animate-spin loading-spinner");
  
  if (inline) {
    return (
      <div className={cn("flex items-center gap-2 text-sm text-muted-foreground", className)}>
        <RefreshCw className={iconClasses} />
        {message && <span>{message}</span>}
      </div>
    );
  }

  return (
    <div className={cn("text-center py-8", className)}>
      <RefreshCw className={cn(iconClasses, "mx-auto mb-2")} />
      {message && <div className="text-sm text-muted-foreground">{message}</div>}
    </div>
  );
}

export { LoadingState } 