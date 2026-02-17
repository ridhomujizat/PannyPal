import { ChatMessage as ChatMessageType } from "@/lib/api";
import { cn } from "@/lib/utils";
import {
    User,
    Bot,
    Lightbulb,
    TrendingUp,
    BarChart3,
    CheckCircle2,
    Sparkles,
} from "lucide-react";
import { formatDistanceToNow } from "date-fns";

interface ChatMessageProps {
    message: ChatMessageType;
}

// Parsed structured AI response
interface AIStructuredResponse {
    answer: string;
    insights: string[];
    recommendations: string[];
    needs_visualization?: boolean;
    visualization_type?: string;
    visualization_hint?: string;
}

/**
 * Try to extract structured JSON from AI content.
 * The AI response format is:
 *   "Some intro text...\n\n```json\n{ ... }\n```"
 * We extract the intro and parse the JSON block separately.
 */
function parseAIContent(content: string): {
    intro: string;
    structured: AIStructuredResponse | null;
} {
    // Match ```json ... ``` block
    const jsonBlockRegex = /```json\s*\n([\s\S]*?)\n```/;
    const match = content.match(jsonBlockRegex);

    if (!match) {
        return { intro: content, structured: null };
    }

    // Everything before the json block is intro text
    const intro = content.substring(0, match.index).trim();

    try {
        const parsed = JSON.parse(match[1]);
        return {
            intro,
            structured: {
                answer: parsed.answer || "",
                insights: parsed.insights || [],
                recommendations: parsed.recommendations || [],
                needs_visualization: parsed.needs_visualization,
                visualization_type: parsed.visualization_type,
                visualization_hint: parsed.visualization_hint,
            },
        };
    } catch {
        // JSON parse failed, return raw content
        return { intro: content, structured: null };
    }
}

/** Render text with **bold** markdown support */
function renderFormattedText(text: string) {
    const parts = text.split(/(\*\*[^*]+\*\*)/g);
    return parts.map((part, i) => {
        if (part.startsWith("**") && part.endsWith("**")) {
            return (
                <strong key={i} className="font-semibold">
                    {part.slice(2, -2)}
                </strong>
            );
        }
        return <span key={i}>{part}</span>;
    });
}

export function ChatMessage({ message }: ChatMessageProps) {
    const isUser = message.role === "user";

    // Parse structured AI response
    const parsed = !isUser ? parseAIContent(message.content) : null;
    const hasStructured = parsed?.structured !== null;

    return (
        <div
            className={cn(
                "flex gap-3 mb-6 animate-in slide-in-from-bottom-4 duration-300",
                isUser ? "justify-end" : "justify-start"
            )}
        >
            {!isUser && (
                <div className="flex-shrink-0 w-8 h-8 rounded-xl bg-gradient-to-br from-primary to-primary/80 flex items-center justify-center shadow-lg">
                    <Bot className="w-5 h-5 text-primary-foreground" />
                </div>
            )}

            <div
                className={cn(
                    "flex flex-col gap-2 max-w-[85%] md:max-w-[75%]",
                    isUser && "items-end"
                )}
            >
                {/* User message or plain AI message */}
                {isUser || !hasStructured ? (
                    <div
                        className={cn(
                            "rounded-2xl px-4 py-3 shadow-sm",
                            isUser
                                ? "bg-primary text-primary-foreground rounded-tr-sm"
                                : "bg-white dark:bg-card border border-border/60 text-foreground rounded-tl-sm"
                        )}
                    >
                        <p className="text-sm leading-relaxed whitespace-pre-wrap">
                            {message.content}
                        </p>
                    </div>
                ) : (
                    /* Structured AI response */
                    <div className="space-y-3 w-full">
                        {/* Intro text */}
                        {parsed!.intro && (
                            <div className="rounded-2xl rounded-tl-sm px-4 py-3 shadow-sm bg-white dark:bg-card border border-border/60 text-foreground">
                                <p className="text-sm leading-relaxed whitespace-pre-wrap">
                                    {parsed!.intro}
                                </p>
                            </div>
                        )}

                        {/* Answer card */}
                        {parsed!.structured!.answer && (
                            <div className="rounded-2xl px-4 py-3 bg-gradient-to-br from-primary/10 to-primary/5 border border-primary/20 shadow-sm">
                                <div className="flex items-start gap-2">
                                    <Sparkles className="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
                                    <p className="text-sm leading-relaxed text-foreground">
                                        {renderFormattedText(parsed!.structured!.answer)}
                                    </p>
                                </div>
                            </div>
                        )}

                        {/* Insights */}
                        {parsed!.structured!.insights.length > 0 && (
                            <div className="space-y-2">
                                <div className="flex items-center gap-1.5 px-1">
                                    <Lightbulb className="w-3.5 h-3.5 text-amber-500" />
                                    <span className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                                        Insights
                                    </span>
                                </div>
                                {parsed!.structured!.insights.map((insight, idx) => (
                                    <div
                                        key={idx}
                                        className="rounded-xl px-4 py-3 bg-amber-500/5 border border-amber-500/15 shadow-sm"
                                    >
                                        <div className="flex items-start gap-2.5">
                                            <TrendingUp className="w-4 h-4 text-amber-500 mt-0.5 flex-shrink-0" />
                                            <p className="text-sm leading-relaxed text-foreground">
                                                {renderFormattedText(insight)}
                                            </p>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}

                        {/* Recommendations */}
                        {parsed!.structured!.recommendations.length > 0 && (
                            <div className="space-y-2">
                                <div className="flex items-center gap-1.5 px-1">
                                    <CheckCircle2 className="w-3.5 h-3.5 text-emerald-500" />
                                    <span className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                                        Recommendations
                                    </span>
                                </div>
                                {parsed!.structured!.recommendations.map((rec, idx) => (
                                    <div
                                        key={idx}
                                        className="rounded-xl px-4 py-3 bg-emerald-500/5 border border-emerald-500/15 shadow-sm"
                                    >
                                        <div className="flex items-start gap-2.5">
                                            <span className="flex-shrink-0 w-5 h-5 rounded-full bg-emerald-500/20 flex items-center justify-center text-xs font-bold text-emerald-600 mt-0.5">
                                                {idx + 1}
                                            </span>
                                            <p className="text-sm leading-relaxed text-foreground">
                                                {renderFormattedText(rec)}
                                            </p>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}

                        {/* Visualization hint */}
                        {parsed!.structured!.needs_visualization &&
                            parsed!.structured!.visualization_hint && (
                                <div className="rounded-xl px-4 py-3 bg-sky-500/5 border border-sky-500/15 shadow-sm">
                                    <div className="flex items-start gap-2.5">
                                        <BarChart3 className="w-4 h-4 text-sky-500 mt-0.5 flex-shrink-0" />
                                        <div>
                                            <p className="text-xs font-semibold text-sky-600 mb-0.5">
                                                Visualization Available
                                            </p>
                                            <p className="text-xs leading-relaxed text-muted-foreground">
                                                {parsed!.structured!.visualization_hint}
                                            </p>
                                        </div>
                                    </div>
                                </div>
                            )}
                    </div>
                )}

                <span className="text-xs text-muted-foreground px-1">
                    {formatDistanceToNow(new Date(message.created_at), {
                        addSuffix: true,
                    })}
                </span>
            </div>

            {isUser && (
                <div className="flex-shrink-0 w-8 h-8 rounded-xl bg-gradient-to-br from-sky to-sky/80 flex items-center justify-center shadow-lg">
                    <User className="w-5 h-5 text-white" />
                </div>
            )}
        </div>
    );
}
