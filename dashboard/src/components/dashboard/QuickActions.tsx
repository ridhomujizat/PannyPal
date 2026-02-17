import { Plus, ArrowUpRight, ArrowDownLeft, Repeat } from "lucide-react";
import { cn } from "@/lib/utils";

const actions = [
  { label: "Add Income", icon: ArrowDownLeft, color: "mint" as const },
  { label: "Add Expense", icon: ArrowUpRight, color: "coral" as const },
  { label: "Transfer", icon: Repeat, color: "lavender" as const },
  { label: "New Goal", icon: Plus, color: "sky" as const },
];

const buttonStyles = {
  mint: "bg-mint hover:bg-mint-dark hover:text-card text-mint-dark",
  coral: "bg-coral hover:bg-coral-dark hover:text-card text-coral-dark",
  lavender: "bg-lavender hover:bg-lavender-dark hover:text-card text-lavender-dark",
  sky: "bg-sky hover:bg-sky-dark hover:text-card text-sky-dark",
};

export function QuickActions() {
  return (
    <div className="flex flex-wrap gap-2 animate-fade-in">
      {actions.map((action) => (
        <button
          key={action.label}
          className={cn(
            "flex items-center gap-2 px-4 py-2.5 rounded-xl font-semibold text-sm transition-all duration-200",
            buttonStyles[action.color]
          )}
        >
          <action.icon className="w-4 h-4" />
          <span className="hidden sm:inline">{action.label}</span>
        </button>
      ))}
    </div>
  );
}
