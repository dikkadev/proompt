import { Snippet } from "@/lib/api";
import { cn } from "@/lib/utils";
import { BaseCard } from "./BaseCard";

interface SnippetCardProps {
  snippet: Snippet;
  onInsert: (name: string) => void;
  className?: string;
}

export function SnippetCard({
  snippet,
  onInsert,
  className,
}: SnippetCardProps) {
  return (
    <BaseCard
      title={snippet.title}
      description={snippet.description}
      tags={snippet.tags}
      className={className}
      type="snippet"
    >
      <span className="text-muted-foreground ml-0.5 break-all">
        {/* Assuming snippet has a category for now */}
        {(snippet as any).category}
      </span>
    </BaseCard>
  );
} 