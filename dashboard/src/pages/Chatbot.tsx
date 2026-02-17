import { useState, useEffect, useRef } from "react";
import { Sidebar } from "@/components/dashboard/Sidebar";
import { ChatMessage } from "@/components/chatbot/ChatMessage";
import { ChatInput } from "@/components/chatbot/ChatInput";
import { TypingIndicator } from "@/components/chatbot/TypingIndicator";
import { ConversationList } from "@/components/chatbot/ConversationList";
import {
    useConversations,
    useConversation,
    useSendMessage,
    useClearConversation,
    ChatMessage as ChatMessageType
} from "@/lib/api";
import { useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { MessageSquare, Loader2, PanelLeftClose, PanelLeftOpen } from "lucide-react";
import { cn } from "@/lib/utils";

const Chatbot = () => {
    const [activeSessionId, setActiveSessionId] = useState<string | null>(null);
    const [messages, setMessages] = useState<ChatMessageType[]>([]);
    const [isTyping, setIsTyping] = useState(false);
    const [showSidebar, setShowSidebar] = useState(false);
    const messagesEndRef = useRef<HTMLDivElement>(null);
    const queryClient = useQueryClient();

    // API Hooks
    const { data: conversations = [], refetch: refetchConversations } = useConversations(20);
    const { data: conversationData, isLoading: isLoadingConversation } = useConversation(activeSessionId);
    const sendMessage = useSendMessage();
    const clearConversation = useClearConversation();

    // Update messages when conversation data changes
    useEffect(() => {
        if (conversationData?.messages) {
            setMessages(conversationData.messages.reverse()); // Backend returns newest first, we want oldest first
        }
    }, [conversationData]);

    // Auto-scroll to bottom when new messages arrive
    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }, [messages, isTyping]);

    const handleSendMessage = async (messageText: string) => {
        if (!messageText.trim()) return;

        // Add user message optimistically
        const userMessage: ChatMessageType = {
            session_id: activeSessionId || "",
            role: "user",
            content: messageText,
            token_used: 0,
            response_time: 0,
            created_at: new Date().toISOString(),
        };

        setMessages((prev) => [...prev, userMessage]);
        setIsTyping(true);

        try {
            const response = await sendMessage({
                session_id: activeSessionId || undefined,
                message: messageText,
            });

            // If this was a new conversation, set the session ID
            if (!activeSessionId) {
                setActiveSessionId(response.session_id);
            }

            // Add assistant response
            setMessages((prev) => [...prev, response]);

            // Refetch conversations to update the list
            refetchConversations();

            // Invalidate conversation cache
            queryClient.invalidateQueries({ queryKey: ["conversation", response.session_id] });
        } catch (error) {
            toast.error("Failed to send message. Please try again.");
            console.error("Error sending message:", error);
            // Remove optimistic user message on error
            setMessages((prev) => prev.slice(0, -1));
        } finally {
            setIsTyping(false);
        }
    };

    const handleSelectConversation = (sessionId: string) => {
        setActiveSessionId(sessionId);
        setMessages([]);
        setShowSidebar(false); // Close sidebar on mobile after selecting
    };

    const handleDeleteConversation = async (sessionId: string) => {
        try {
            await clearConversation(sessionId);
            toast.success("Conversation deleted");

            // If we're deleting the active conversation, reset
            if (sessionId === activeSessionId) {
                setActiveSessionId(null);
                setMessages([]);
            }

            // Refetch conversations
            refetchConversations();
        } catch (error) {
            toast.error("Failed to delete conversation");
            console.error("Error deleting conversation:", error);
        }
    };

    const handleNewChat = () => {
        setActiveSessionId(null);
        setMessages([]);
        setShowSidebar(false); // Close sidebar on mobile
    };

    return (
        <div className="min-h-screen w-full">
            {/* Floating Sidebar */}
            <Sidebar />

            {/* Main Content */}
            <main className="min-h-screen p-3 pt-16 sm:p-4 sm:pt-6 lg:pl-24 lg:p-8 overflow-hidden">
                <div className="max-w-7xl mx-auto lg:ml-4 h-[calc(100vh-4rem)] sm:h-[calc(100vh-2rem)] flex flex-col">

                    {/* Chat Container */}
                    <div className="flex-1 flex gap-4 lg:gap-6 overflow-hidden min-h-0">
                        {/* Conversation List Sidebar - Desktop always visible, Mobile toggle */}
                        {/* Mobile overlay backdrop */}
                        {showSidebar && (
                            <div
                                className="lg:hidden fixed inset-0 bg-black/40 z-20 animate-in fade-in duration-200"
                                onClick={() => setShowSidebar(false)}
                            />
                        )}

                        {/* Conversation list panel */}
                        <div className={cn(
                            // Desktop styles
                            "hidden lg:flex lg:w-80 lg:flex-shrink-0 bg-card rounded-2xl border border-sidebar shadow-sm overflow-hidden flex-col",
                            // Mobile styles - slide-in overlay (z-30, below app sidebar z-40/z-50)
                            showSidebar && "!flex fixed inset-y-0 left-0 w-[85%] max-w-sm z-30 rounded-none border-r lg:relative lg:inset-auto lg:w-80 lg:rounded-2xl lg:border animate-in slide-in-from-left duration-300"
                        )}>
                            <ConversationList
                                conversations={conversations}
                                activeSessionId={activeSessionId}
                                onSelectConversation={handleSelectConversation}
                                onDeleteConversation={handleDeleteConversation}
                                onNewChat={handleNewChat}
                            />
                        </div>

                        {/* Chat Area */}
                        <div className="flex-1 bg-card rounded-2xl border border-sidebar shadow-sm flex flex-col overflow-hidden min-w-0">
                            {/* Chat Header */}
                            <div className="p-3 sm:p-4 border-b border-sidebar flex items-center gap-2">
                                {/* Mobile sidebar toggle */}
                                <button
                                    onClick={() => setShowSidebar(!showSidebar)}
                                    className="lg:hidden w-9 h-9 rounded-xl bg-muted flex items-center justify-center hover:bg-muted/80 transition-colors flex-shrink-0"
                                >
                                    {showSidebar ? (
                                        <PanelLeftClose className="w-4 h-4 text-muted-foreground" />
                                    ) : (
                                        <PanelLeftOpen className="w-4 h-4 text-muted-foreground" />
                                    )}
                                </button>

                                <div className="w-9 h-9 sm:w-10 sm:h-10 rounded-xl bg-gradient-to-br from-primary to-primary/80 flex items-center justify-center flex-shrink-0">
                                    <MessageSquare className="w-4 h-4 sm:w-5 sm:h-5 text-primary-foreground" />
                                </div>
                                <div className="min-w-0 flex-1">
                                    <h2 className="font-semibold text-foreground text-sm sm:text-base truncate">
                                        {conversationData?.title || "New Conversation"}
                                    </h2>
                                    <p className="text-xs text-muted-foreground">
                                        AI Financial Assistant
                                    </p>
                                </div>
                            </div>

                            {/* Messages Container */}
                            <div className="flex-1 overflow-y-auto p-3 sm:p-4 lg:p-6 space-y-4">
                                {isLoadingConversation ? (
                                    <div className="flex items-center justify-center h-full">
                                        <Loader2 className="w-8 h-8 animate-spin text-primary" />
                                    </div>
                                ) : messages.length === 0 && !activeSessionId ? (
                                    <div className="flex flex-col items-center justify-center h-full text-center px-4">
                                        <div className="w-16 h-16 sm:w-20 sm:h-20 rounded-2xl bg-gradient-to-br from-primary/20 to-primary/10 flex items-center justify-center mb-4">
                                            <MessageSquare className="w-8 h-8 sm:w-10 sm:h-10 text-primary" />
                                        </div>
                                        <h3 className="text-base sm:text-lg font-semibold mb-2">Start a Conversation</h3>
                                        <p className="text-muted-foreground text-xs sm:text-sm max-w-md">
                                            Ask me anything about your finances! I can help you analyze spending,
                                            track income, set goals, and provide financial insights.
                                        </p>
                                    </div>
                                ) : (
                                    <>
                                        {messages.map((message, index) => (
                                            <ChatMessage key={`${message.session_id}-${index}`} message={message} />
                                        ))}
                                        {isTyping && <TypingIndicator />}
                                        <div ref={messagesEndRef} />
                                    </>
                                )}
                            </div>

                            {/* Input Area */}
                            <div className="p-3 sm:p-4 border-t border-sidebar bg-background/50">
                                <ChatInput
                                    onSend={handleSendMessage}
                                    disabled={isTyping}
                                    placeholder="Ask me about your finances..."
                                />
                            </div>
                        </div>
                    </div>
                </div>
            </main>
        </div>
    );
};

export default Chatbot;
