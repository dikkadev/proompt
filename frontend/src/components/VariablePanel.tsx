import { useState } from "react";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Hash, AlertCircle, CheckCircle, Clock } from "lucide-react";

interface Variable {
  name: string;
  value?: string;
  defaultValue?: string;
  state: 'provided' | 'default' | 'missing';
}

interface VariablePanelProps {
  variables: Variable[];
  onVariableChange: (name: string, value: string) => void;
}

export function VariablePanel({ variables, onVariableChange }: VariablePanelProps) {
  const [values, setValues] = useState<Record<string, string>>({});

  const handleValueChange = (name: string, value: string) => {
    setValues(prev => ({ ...prev, [name]: value }));
    onVariableChange(name, value);
  };

  const getVariableState = (variable: Variable) => {
    if (values[variable.name]) return 'provided';
    if (variable.defaultValue) return 'default';
    return 'missing';
  };

  const getStateIcon = (state: string) => {
    switch (state) {
      case 'provided':
        return <CheckCircle className="h-4 w-4 text-variable-provided" />;
      case 'default':
        return <Clock className="h-4 w-4 text-variable-default" />;
      case 'missing':
        return <AlertCircle className="h-4 w-4 text-variable-missing" />;
      default:
        return null;
    }
  };

  const getStateBadge = (state: string) => {
    const baseClasses = "text-xs px-2 py-1";
    switch (state) {
      case 'provided':
        return `${baseClasses} bg-variable-provided/10 text-variable-provided border-variable-provided`;
      case 'default':
        return `${baseClasses} bg-variable-default/10 text-variable-default border-variable-default`;
      case 'missing':
        return `${baseClasses} bg-variable-missing/10 text-variable-missing border-variable-missing`;
      default:
        return baseClasses;
    }
  };

  const providedVars = variables.filter(v => getVariableState(v) === 'provided');
  const defaultVars = variables.filter(v => getVariableState(v) === 'default');
  const missingVars = variables.filter(v => getVariableState(v) === 'missing');

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
        <div className="flex gap-2 text-sm">
          <div className="flex items-center gap-1">
            <CheckCircle className="h-3 w-3 text-variable-provided" />
            <span className="text-variable-provided">{providedVars.length}</span>
          </div>
          <div className="flex items-center gap-1">
            <Clock className="h-3 w-3 text-variable-default" />
            <span className="text-variable-default">{defaultVars.length}</span>
          </div>
          <div className="flex items-center gap-1">
            <AlertCircle className="h-3 w-3 text-variable-missing" />
            <span className="text-variable-missing">{missingVars.length}</span>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {variables.length === 0 ? (
          <div className="text-center text-muted-foreground py-8">
            <Hash className="h-8 w-8 mx-auto mb-2 opacity-50" />
            <p>No variables found</p>
            <p className="text-xs">Use {`{{variable_name}}`} syntax to add variables</p>
          </div>
        ) : (
          <>
            {/* Missing Variables (Priority) */}
            {missingVars.length > 0 && (
              <div className="space-y-3">
                <div className="flex items-center gap-2">
                  <AlertCircle className="h-4 w-4 text-variable-missing" />
                  <h4 className="text-sm font-medium text-variable-missing">Required Variables</h4>
                </div>
                {missingVars.map((variable) => (
                  <div key={variable.name} className="space-y-2">
                    <div className="flex items-center justify-between">
                      <Label htmlFor={variable.name} className="text-sm font-medium">
                        {variable.name}
                      </Label>
                      <Badge className={getStateBadge('missing')}>
                        Required
                      </Badge>
                    </div>
                    <Input
                      id={variable.name}
                      value={values[variable.name] || ''}
                      onChange={(e) => handleValueChange(variable.name, e.target.value)}
                      placeholder="Enter value..."
                      className="border-variable-missing focus:border-variable-missing"
                    />
                  </div>
                ))}
              </div>
            )}

            {missingVars.length > 0 && (defaultVars.length > 0 || providedVars.length > 0) && (
              <Separator />
            )}

            {/* Default Variables */}
            {defaultVars.length > 0 && (
              <div className="space-y-3">
                <div className="flex items-center gap-2">
                  <Clock className="h-4 w-4 text-variable-default" />
                  <h4 className="text-sm font-medium text-variable-default">Optional Variables</h4>
                </div>
                {defaultVars.map((variable) => (
                  <div key={variable.name} className="space-y-2">
                    <div className="flex items-center justify-between">
                      <Label htmlFor={variable.name} className="text-sm font-medium">
                        {variable.name}
                      </Label>
                      <Badge className={getStateBadge('default')}>
                        Default
                      </Badge>
                    </div>
                    <Input
                      id={variable.name}
                      value={values[variable.name] || ''}
                      onChange={(e) => handleValueChange(variable.name, e.target.value)}
                      placeholder={variable.defaultValue || 'Enter value...'}
                      className="border-variable-default focus:border-variable-default"
                    />
                    {variable.defaultValue && (
                      <p className="text-xs text-muted-foreground">
                        Default: {variable.defaultValue}
                      </p>
                    )}
                  </div>
                ))}
              </div>
            )}

            {defaultVars.length > 0 && providedVars.length > 0 && (
              <Separator />
            )}

            {/* Provided Variables */}
            {providedVars.length > 0 && (
              <div className="space-y-3">
                <div className="flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 text-variable-provided" />
                  <h4 className="text-sm font-medium text-variable-provided">Completed Variables</h4>
                </div>
                {providedVars.map((variable) => (
                  <div key={variable.name} className="space-y-2">
                    <div className="flex items-center justify-between">
                      <Label htmlFor={variable.name} className="text-sm font-medium">
                        {variable.name}
                      </Label>
                      <Badge className={getStateBadge('provided')}>
                        Set
                      </Badge>
                    </div>
                    <Input
                      id={variable.name}
                      value={values[variable.name] || ''}
                      onChange={(e) => handleValueChange(variable.name, e.target.value)}
                      className="border-variable-provided focus:border-variable-provided"
                    />
                  </div>
                ))}
              </div>
            )}
          </>
        )}
      </div>

      {variables.length > 0 && (
        <div className="p-4 border-t border-border bg-muted/50">
          <div className="flex justify-between items-center">
            <div className="text-xs text-muted-foreground">
              {missingVars.length > 0 ? `${missingVars.length} required` : 'All variables set'}
            </div>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setValues({})}
              className="text-xs h-6 px-2"
            >
              Clear All
            </Button>
          </div>
        </div>
      )}
    </Card>
  );
}