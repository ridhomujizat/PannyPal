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
import { ChartVisualization } from "./ChartVisualization";

interface ChatMessageProps {
    message: ChatMessageType;
}

// Parsed structured AI response
interface AIStructuredResponse {
    answer: string;
    insights: string[];
    recommendations: string[];
    optimization_opportunities?: string[];
    top_categories?: Array<{
        category: string;
        amount: number;
        percentage?: number;
        transaction_count?: number;
    }>;
    needs_visualization?: boolean;
    visualization_type?: string;
    visualization_hint?: string;
}

// Extracted chart data point
interface ExtractedDataPoint {
    label: string;
    value: number;
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
                recommendations: parsed.recommendations || parsed.optimization_opportunities || [],
                optimization_opportunities: parsed.optimization_opportunities || [],
                top_categories: parsed.top_categories || [],
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

/**
 * Extract chart data from AI response.
 * Priority 1: Use top_categories from parsed JSON (most reliable)
 * Priority 2: Use metadata visualization data from backend
 * Priority 3: Fall back to regex text extraction from answer/insights
 */
function extractChartData(
    structured: AIStructuredResponse | null,
    intro: string
): ExtractedDataPoint[] {
    // Priority 1: Use top_categories from JSON block
    if (structured?.top_categories && structured.top_categories.length >= 2) {
        return structured.top_categories.map((cat) => ({
            label: cat.category,
            value: cat.amount,
        }));
    }

    // Priority 2: Fall back to regex text extraction
    const points: ExtractedDataPoint[] = [];
    const allTexts = [
        intro,
        structured?.answer || "",
        ...(structured?.insights || []),
    ].filter(Boolean);
    const combined = allTexts.join("\n");
    let m;

    // Pattern 1: **Label:** Rp X.XXX or **Label:** Rp X.XXX.XXX
    const boldLabelPattern =
        /\*\*([^*]+?)\*\*[:\s]+Rp\s?([\d.]+(?:\.\d{3})*)/g;
    while ((m = boldLabelPattern.exec(combined)) !== null) {
        const label = m[1].trim().replace(/:$/, "");
        const value = parseRpValue(m[2]);
        if (value > 0) points.push({ label, value });
    }

    // Pattern 2: text **Rp X.XXX.XXX** (value wrapped in bold)
    const boldValuePattern =
        /([^.!?\n]{5,50}?)\s*\*\*Rp\s?([\d.]+(?:\.\d{3})*)\*\*/g;
    while ((m = boldValuePattern.exec(combined)) !== null) {
        const label = extractLabel(m[1]);
        const value = parseRpValue(m[2]);
        if (value > 0 && label.length > 2) points.push({ label, value });
    }

    // Pattern 3: text (Rp X.XXX.XXX) (value in parentheses)
    const parenPattern =
        /([^()\n]{5,60}?)\s*\(Rp\s?([\d.]+(?:\.\d{3})*)\)/g;
    while ((m = parenPattern.exec(combined)) !== null) {
        const label = extractLabel(m[1]);
        const value = parseRpValue(m[2]);
        if (value > 0 && label.length > 2) points.push({ label, value });
    }

    // Deduplicate by value (keep unique amounts, prefer first found)
    const seenValues = new Set<number>();
    const unique = points.filter((p) => {
        if (seenValues.has(p.value)) return false;
        seenValues.add(p.value);
        return true;
    });

    // Also deduplicate by label similarity
    const seenLabels = new Set<string>();
    return unique.filter((p) => {
        const key = p.label.toLowerCase().replace(/[^a-z]/g, "");
        if (seenLabels.has(key)) return false;
        seenLabels.add(key);
        return true;
    });
}

/** Parse "10.765.306" → 10765306 */
function parseRpValue(str: string): number {
    const cleaned = str.replace(/\./g, "");
    const val = parseInt(cleaned, 10);
    return isNaN(val) ? 0 : val;
}

/** Extract a short label from preceding context text */
function extractLabel(text: string): string {
    let label = text
        .replace(/\*\*/g, "")
        .replace(/[*_`]/g, "")
        .replace(/^[\s,;:-]+/, "")
        .trim();
    // Take last meaningful phrase (after comma, colon, etc.)
    const parts = label.split(/[,;]/);
    label = parts[parts.length - 1].trim();
    // Truncate to reasonable length
    if (label.length > 40) {
        label = label.substring(label.length - 40).trim();
        // Find first word boundary
        const spaceIdx = label.indexOf(" ");
        if (spaceIdx > 0) label = label.substring(spaceIdx + 1);
    }
    return label;
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

    // Extract chart data from structured JSON or text
    const chartData = hasStructured
        ? extractChartData(parsed?.structured || null, parsed?.intro || "")
        : [];

    // Check for metadata-based visualization (if backend provides it)
    const hasMetadataViz =
        message.metadata?.visualization &&
        message.metadata.visualization.data?.labels?.length > 0;

    // Determine chart type: use visualization_type from JSON, or "pie" if we have top_categories
    const hasTopCategories =
        parsed?.structured?.top_categories &&
        parsed.structured.top_categories.length >= 2;
    const chartType = hasTopCategories
        ? "pie"
        : parsed?.structured?.visualization_type === "pie"
            ? "pie"
            : parsed?.structured?.visualization_type === "line"
                ? "line"
                : "bar";

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
                    "flex flex-col gap-2",
                    isUser
                        ? "max-w-[85%] md:max-w-[75%] items-end"
                        : "max-w-[95%] sm:max-w-[90%] md:max-w-[80%]"
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
                                    {renderFormattedText(parsed!.intro)}
                                </p>
                            </div>
                        )}

                        {/* Answer card */}
                        {parsed!.structured!.answer && (
                            <div className="rounded-2xl px-4 py-3 bg-gradient-to-br from-primary/10 to-primary/5 border border-primary/20 shadow-sm">
                                <div className="flex items-start gap-2">
                                    <Sparkles className="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
                                    <p className="text-sm leading-relaxed text-foreground">
                                        {renderFormattedText(
                                            parsed!.structured!.answer
                                        )}
                                    </p>
                                </div>
                            </div>
                        )}

                        {/* Chart Visualization from metadata */}
                        {hasMetadataViz &&
                            message.metadata!.visualization!.type !==
                            "table" && (
                                <div className="rounded-xl p-4 bg-gradient-to-br from-sky-500/5 to-sky-500/10 border border-sky-500/20 shadow-sm">
                                    <div className="flex items-center gap-2 mb-3">
                                        <BarChart3 className="w-4 h-4 text-sky-600" />
                                        <span className="text-xs font-semibold uppercase tracking-wider text-sky-700">
                                            Visualisasi Data
                                        </span>
                                    </div>
                                    <ChartVisualization
                                        type={
                                            message.metadata!.visualization!
                                                .type as
                                            | "bar"
                                            | "pie"
                                            | "line"
                                        }
                                        data={
                                            message.metadata!.visualization!
                                                .data
                                        }
                                        config={
                                            message.metadata!.visualization!
                                                .config
                                        }
                                    />
                                </div>
                            )}

                        {/* Chart Visualization from extracted data (top_categories or text regex) */}
                        {!hasMetadataViz &&
                            chartData.length >= 2 && (
                                <div className="rounded-xl p-4 bg-gradient-to-br from-sky-500/5 to-sky-500/10 border border-sky-500/20 shadow-sm">
                                    <div className="flex items-center gap-2 mb-3">
                                        <BarChart3 className="w-4 h-4 text-sky-600" />
                                        <span className="text-xs font-semibold uppercase tracking-wider text-sky-700">
                                            Visualisasi Data
                                        </span>
                                    </div>
                                    <ChartVisualization
                                        type={chartType}
                                        data={{
                                            labels: chartData.map(
                                                (d) => d.label
                                            ),
                                            values: chartData.map(
                                                (d) => d.value
                                            ),
                                        }}
                                        config={{
                                            y_label: "Jumlah (Rp)",
                                            format: "currency",
                                        }}
                                    />
                                </div>
                            )}

                        {/* Visualization hint (fallback if no chart data at all) */}
                        {!hasMetadataViz &&
                            chartData.length < 2 &&
                            parsed!.structured!.needs_visualization &&
                            parsed!.structured!.visualization_hint && (
                                <div className="rounded-xl px-4 py-3 bg-sky-500/5 border border-sky-500/15 shadow-sm">
                                    <div className="flex items-start gap-2.5">
                                        <BarChart3 className="w-4 h-4 text-sky-500 mt-0.5 flex-shrink-0" />
                                        <div>
                                            <p className="text-xs font-semibold text-sky-600 mb-0.5">
                                                Visualization Available
                                            </p>
                                            <p className="text-xs leading-relaxed text-muted-foreground">
                                                {
                                                    parsed!.structured!
                                                        .visualization_hint
                                                }
                                            </p>
                                        </div>
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
                                {parsed!.structured!.insights.map(
                                    (insight, idx) => (
                                        <div
                                            key={idx}
                                            className="rounded-xl px-4 py-3 bg-amber-500/5 border border-amber-500/15 shadow-sm"
                                        >
                                            <div className="flex items-start gap-2.5">
                                                <TrendingUp className="w-4 h-4 text-amber-500 mt-0.5 flex-shrink-0" />
                                                <p className="text-sm leading-relaxed text-foreground">
                                                    {renderFormattedText(
                                                        insight
                                                    )}
                                                </p>
                                            </div>
                                        </div>
                                    )
                                )}
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
                                {parsed!.structured!.recommendations.map(
                                    (rec, idx) => (
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
                                    )
                                )}
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
