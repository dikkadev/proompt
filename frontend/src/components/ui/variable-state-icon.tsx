import * as React from "react"
import { CheckCircle, Clock, AlertCircle } from "lucide-react"
import { cn } from "@/lib/utils"

export type VariableState = 'provided' | 'default' | 'missing'

interface VariableStateIconProps {
  state: VariableState;
  size?: 'sm' | 'md' | 'lg';
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

const getStateColorClass = (state: VariableState): string => {
  switch (state) {
    case 'provided':
      return "text-variable-provided";
    case 'default':
      return "text-variable-default";
    case 'missing':
      return "text-variable-missing";
    default:
      return "";
  }
};

function VariableStateIcon({ state, size = 'md', className }: VariableStateIconProps) {
  const iconClasses = cn(getSizeClasses(size), getStateColorClass(state), className);

  switch (state) {
    case 'provided':
      return <CheckCircle className={iconClasses} />;
    case 'default':
      return <Clock className={iconClasses} />;
    case 'missing':
      return <AlertCircle className={iconClasses} />;
    default:
      return null;
  }
}

export { VariableStateIcon } 