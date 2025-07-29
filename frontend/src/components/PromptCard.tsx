import { Prompt } from "@/lib/api";
import { cn } from "@/lib/utils";
import { BaseCard } from "./BaseCard";

interface PromptCardProps {
  prompt: Prompt;
  className?: string;
}

export function PromptCard({ prompt, className }: PromptCardProps) {
  return (
    <BaseCard
      title={prompt.title}
      description={prompt.use_case}
      tags={prompt.model_compatibility_tags}
      className={className}
      type="prompt"
    />
  );
} 