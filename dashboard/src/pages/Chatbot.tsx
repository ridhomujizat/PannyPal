import { useState, useEffect, useRef } from "react";
import { Sidebar } from "@/components/dashboard/Sidebar";
import { Header } from "@/components/dashboard/Header";
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
import { MessageSquare, Loader2 } from "lucide-react";
import { cn } from "@/lib/utils";

const Chatbot = () => {
    const [activeSessionId, setActiveSessionId] = useState<string | null>(null);
    const [messages, setMessages] = useState<ChatMessageType[]>([]);
    const [isTyping, setIsTyping] = useState(false);
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
    };

    return (
        <div className="min-h-screen w-full">
            {/* Floating Sidebar */}
            <Sidebar />

            {/* Main Content */}
            <main className="min-h-screen p-4 pt-20 sm:p-6 sm:pt-6 lg:pl-24 lg:p-8 overflow-hidden">
                <div className="max-w-7xl mx-auto lg:ml-4 h-[calc(100vh-2rem)] flex flex-col">
                    {/* Header */}
                    <div className="mb-6">
                        <Header
                            startDate={new Date()}
                            endDate={new Date()}
                            onDateChange={() => { }}
                        />
                    </div>

                    {/* Chat Container */}
                    <div className="flex-1 grid grid-cols-1 lg:grid-cols-4 gap-6 overflow-hidden">
                        {/* Conversation List Sidebar */}
                        <div className="lg:col-span-1 bg-card rounded-2xl border border-sidebar shadow-sm overflow-hidden">
                            <ConversationList
                                conversations={conversations}
                                activeSessionId={activeSessionId}
                                onSelectConversation={handleSelectConversation}
                                onDeleteConversation={handleDeleteConversation}
                                onNewChat={handleNewChat}
                            />
                        </div>

                        {/* Chat Area */}
                        <div className="lg:col-span-3 bg-card rounded-2xl border border-sidebar shadow-sm flex flex-col overflow-hidden">
                            {/* Chat Header */}
                            <div className="p-4 border-b border-sidebar flex items-center gap-2">
                                <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-primary to-primary/80 flex items-center justify-center">
                                    <MessageSquare className="w-5 h-5 text-primary-foreground" />
                                </div>
                                <div>
                                    <h2 className="font-semibold text-foreground">
                                        {conversationData?.title || "New Conversation"}
                                    </h2>
                                    <p className="text-xs text-muted-foreground">
                                        AI Financial Assistant
                                    </p>
                                </div>
                            </div>

                            {/* Messages Container */}
                            <div className="flex-1 overflow-y-auto p-6 space-y-4">
                                {isLoadingConversation ? (
                                    <div className="flex items-center justify-center h-full">
                                        <Loader2 className="w-8 h-8 animate-spin text-primary" />
                                    </div>
                                ) : messages.length === 0 && !activeSessionId ? (
                                    <div className="flex flex-col items-center justify-center h-full text-center">
                                        <div className="w-20 h-20 rounded-2xl bg-gradient-to-br from-primary/20 to-primary/10 flex items-center justify-center mb-4">
                                            <MessageSquare className="w-10 h-10 text-primary" />
                                        </div>
                                        <h3 className="text-lg font-semibold mb-2">Start a Conversation</h3>
                                        <p className="text-muted-foreground text-sm max-w-md">
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
                            <div className="p-4 border-t border-sidebar bg-background/50">
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
