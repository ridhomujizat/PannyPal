import { Conversation } from "@/lib/api";
import { cn } from "@/lib/utils";
import { MessageSquare, Trash2 } from "lucide-react";
import { formatDistanceToNow } from "date-fns";
import { Button } from "@/components/ui/button";

interface ConversationListProps {
    conversations: Conversation[];
    activeSessionId: string | null;
    onSelectConversation: (sessionId: string) => void;
    onDeleteConversation: (sessionId: string) => void;
    onNewChat: () => void;
}

export function ConversationList({
    conversations,
    activeSessionId,
    onSelectConversation,
    onDeleteConversation,
    onNewChat,
}: ConversationListProps) {
    return (
        <div className="flex flex-col h-full">
            {/* New Chat Button */}
            <div className="p-4 border-b border-sidebar">
                <Button
                    onClick={onNewChat}
                    className="w-full bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl"
                >
                    <MessageSquare className="w-4 h-4 mr-2" />
                    New Chat
                </Button>
            </div>

            {/* Conversation List */}
            <div className="flex-1 overflow-y-auto p-2 space-y-1">
                {conversations.length === 0 ? (
                    <div className="text-center py-8 text-muted-foreground text-sm">
                        No conversations yet
                    </div>
                ) : (
                    conversations.map((conv) => {
                        const isActive = activeSessionId === conv.session_id;
                        return (
                            <div
                                key={conv.session_id}
                                className={cn(
                                    "group relative p-3 rounded-xl transition-all",
                                    isActive
                                        ? "bg-primary text-primary-foreground pointer-events-none shadow-md"
                                        : "cursor-pointer hover:bg-sidebar/50 border border-transparent"
                                )}
                                onClick={() => !isActive && onSelectConversation(conv.session_id)}
                            >
                                <div className="flex items-start gap-2 pr-8">
                                    <MessageSquare className={cn(
                                        "w-4 h-4 mt-1 flex-shrink-0",
                                        isActive ? "text-primary-foreground" : "text-primary"
                                    )} />
                                    <div className="flex-1 min-w-0">
                                        <p className={cn(
                                            "text-sm font-medium truncate",
                                            isActive && "text-primary-foreground"
                                        )}>{conv.title}</p>
                                        <p className={cn(
                                            "text-xs",
                                            isActive ? "text-primary-foreground/70" : "text-muted-foreground"
                                        )}>
                                            {conv.message_count} messages · {formatDistanceToNow(new Date(conv.last_message), { addSuffix: true })}
                                        </p>
                                    </div>
                                </div>

                                <Button
                                    variant="ghost"
                                    size="icon"
                                    className="absolute right-2 top-2 h-6 w-6 opacity-0 group-hover:opacity-100 transition-opacity"
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        onDeleteConversation(conv.session_id);
                                    }}
                                >
                                    <Trash2 className="w-3 h-3 text-destructive" />
                                </Button>
                            </div>
                        );
                    })
                )}
            </div>
        </div>
    );
}
