import * as React from "react"
import { cn } from "@/lib/utils"
import { Badge } from "./badge"

export type StatusType = 'provided' | 'default' | 'missing' | 'warning'

interface StatusBadgeProps extends React.HTMLAttributes<HTMLDivElement> {
  status: StatusType;
  count?: number;
  children?: React.ReactNode;
}

const getStatusClasses = (status: StatusType): string => {
  const baseClasses = "gap-1";
  switch (status) {
    case 'provided':
      return `${baseClasses} border-variable-provided text-variable-provided`;
    case 'default':
      return `${baseClasses} border-variable-default text-variable-default`;
    case 'missing':
      return `${baseClasses} border-variable-missing text-variable-missing`;
    case 'warning':
      return `${baseClasses} border-warning text-warning`;
    default:
      return baseClasses;
  }
};

function StatusBadge({ status, count, children, className, ...props }: StatusBadgeProps) {
  return (
    <Badge 
      variant="outline" 
      className={cn(getStatusClasses(status), className)} 
      {...props}
    >
      {children}
    </Badge>
  );
}

export { StatusBadge } 