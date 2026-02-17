import { cn } from "@/lib/utils";

export function TypingIndicator() {
    return (
        <div className="flex items-center space-x-2 p-4 bg-sidebar/30 rounded-2xl max-w-[100px]">
            <div className="flex space-x-1">
                <div className={cn(
                    "w-2 h-2 rounded-full bg-primary/60 animate-bounce",
                    "[animation-delay:-0.3s]"
                )} />
                <div className={cn(
                    "w-2 h-2 rounded-full bg-primary/60 animate-bounce",
                    "[animation-delay:-0.15s]"
                )} />
                <div className="w-2 h-2 rounded-full bg-primary/60 animate-bounce" />
            </div>
        </div>
    );
}
