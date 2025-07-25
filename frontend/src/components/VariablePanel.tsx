import { useState } from "react";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { StatusBadge } from "@/components/ui/status-badge";
import { VariableStateIcon } from "@/components/ui/variable-state-icon";
import { Hash, X, FileText } from "lucide-react";

interface Variable {
  name: string;
  value?: string;
  defaultValue?: string;
  state: 'provided' | 'default' | 'missing';
}

interface VariablePanelProps {
  variables: Variable[];
  onVariableChange: (name: string, value: string) => void;
  snippets: string[];
}

export function VariablePanel({ variables, onVariableChange, snippets }: VariablePanelProps) {
  const [values, setValues] = useState<Record<string, string>>({});

  const handleValueChange = (name: string, value: string) => {
    setValues(prev => ({ ...prev, [name]: value }));
    onVariableChange(name, value);
  };

  const handleClearVariable = (name: string) => {
    setValues(prev => {
      const newValues = { ...prev };
      delete newValues[name];
      return newValues;
    });
    onVariableChange(name, '');
  };

  const getVariableState = (variable: Variable) => {
    if (values[variable.name]) return 'provided';
    if (variable.defaultValue) return 'default';
    return 'missing';
  };



  const getInputClasses = (state: string) => {
    switch (state) {
      case 'provided':
        return "border-variable-provided/50 focus:border-variable-provided focus:ring-variable-provided/20";
      case 'default':
        return "border-variable-default/50 focus:border-variable-default focus:ring-variable-default/20";
      case 'missing':
        return "border-variable-missing/50 focus:border-variable-missing focus:ring-variable-missing/20";
      default:
        return "";
    }
  };

  // Count variables by state for summary
  const stateCounts = variables.reduce((acc, variable) => {
    const state = getVariableState(variable);
    acc[state] = (acc[state] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  return (
    <Card className="h-full flex flex-col">
      <div className="p-4 border-b border-border">
        <div className="flex items-center gap-2 mb-3">
          <Hash className="h-5 w-5 text-primary" />
          <h3 className="font-semibold">Variables</h3>
          <Badge variant="secondary" className="ml-auto border-primary/20">
            {variables.length}
          </Badge>
        </div>
        
        {/* Status Summary */}
        {variables.length > 0 && (
          <div className="flex gap-3 text-sm">
            <div className="flex items-center gap-1">
              <VariableStateIcon state="provided" size="sm" />
              <span className="text-variable-provided">{stateCounts.provided || 0}</span>
            </div>
            <div className="flex items-center gap-1">
              <VariableStateIcon state="default" size="sm" />
              <span className="text-variable-default">{stateCounts.default || 0}</span>
            </div>
            <div className="flex items-center gap-1">
              <VariableStateIcon state="missing" size="sm" />
              <span className="text-variable-missing">{stateCounts.missing || 0}</span>
            </div>
          </div>
        )}
      </div>

      <div className="flex-1 overflow-y-auto p-4">
        {variables.length === 0 ? (
          <div className="text-center text-muted-foreground py-8">
            <Hash className="h-8 w-8 mx-auto mb-2 opacity-50" />
            <p>No variables found</p>
            <p className="text-xs">Use {`{{variable_name}}`} syntax to add variables</p>
          </div>
        ) : (
          <div className="space-y-3">
            {/* Render variables in editor order to maintain stability */}
            {variables.map((variable, index) => {
              const state = getVariableState(variable);
              const hasValue = !!values[variable.name];
              
              return (
                <div key={`${variable.name}-${index}`} className="space-y-2">
                  <div className="flex items-center justify-between">
                    <Label htmlFor={variable.name} className="text-sm font-medium">
                      {variable.name}
                    </Label>
                    <div className="flex items-center gap-2">
                      <StatusBadge status={state}>
                        {state === 'provided' ? 'Set' : state === 'default' ? 'Default' : 'Required'}
                      </StatusBadge>
                      <VariableStateIcon state={state} />
                    </div>
                  </div>
                  
                  <div className="flex gap-2">
                    <Input
                      id={variable.name}
                      value={values[variable.name] || ''}
                      onChange={(e) => handleValueChange(variable.name, e.target.value)}
                      placeholder={variable.defaultValue || 'Enter value...'}
                      className={`flex-1 ${getInputClasses(state)}`}
                    />
                    {hasValue && (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleClearVariable(variable.name)}
                        className="px-2 h-10 text-muted-foreground hover:text-foreground"
                        title={`Clear ${variable.name}`}
                      >
                        <X className="h-4 w-4" />
                      </Button>
                    )}
                  </div>
                  
                  {variable.defaultValue && state !== 'provided' && (
                    <p className="text-xs text-muted-foreground">
                      Default: {variable.defaultValue}
                    </p>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Referenced Snippets Section */}
      {snippets.length > 0 && (
        <>
          <Separator />
          <div className="p-3 bg-muted/30">
            <div className="flex items-center gap-2 text-xs text-muted-foreground mb-2">
              <FileText className="h-3 w-3" />
              Referenced Snippets ({snippets.length})
            </div>
            <div className="flex flex-wrap gap-1">
              {snippets.map((snippet) => (
                <Badge
                  key={snippet}
                  variant="outline"
                  className="text-xs px-2 py-0.5 bg-primary/10 text-primary border-primary/30"
                >
                  @{snippet}
                </Badge>
              ))}
            </div>
          </div>
        </>
      )}

      {variables.length > 0 && (
        <div className="p-3 border-t border-border bg-muted/50">
          <div className="flex justify-between items-center">
            <div className="text-xs text-muted-foreground">
              {(stateCounts.missing || 0) > 0 
                ? `${stateCounts.missing} required` 
                : 'All variables set'
              }
            </div>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setValues({})}
              className="text-xs h-6 px-2"
              disabled={Object.keys(values).length === 0}
            >
              Clear All
            </Button>
          </div>
        </div>
      )}
    </Card>
  );
}