import { useState, FormEvent } from "react";
import { Send } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";

interface ChatInputProps {
    onSend: (message: string) => void;
    disabled?: boolean;
    placeholder?: string;
}

export function ChatInput({ onSend, disabled = false, placeholder = "Ask me anything about your finances..." }: ChatInputProps) {
    const [message, setMessage] = useState("");

    const handleSubmit = (e: FormEvent) => {
        e.preventDefault();
        if (message.trim() && !disabled) {
            onSend(message.trim());
            setMessage("");
        }
    };

    const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
        if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault();
            handleSubmit(e);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="relative">
            <Textarea
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder={placeholder}
                disabled={disabled}
                className="min-h-[60px] max-h-[200px] pr-12 resize-none rounded-2xl border-sidebar bg-background focus-visible:ring-primary"
                rows={2}
            />
            <Button
                type="submit"
                size="icon"
                disabled={!message.trim() || disabled}
                className="absolute right-2 bottom-2 h-8 w-8 rounded-xl bg-primary hover:bg-primary/90 transition-all disabled:opacity-50"
            >
                <Send className="h-4 w-4" />
            </Button>
        </form>
    );
}
