import { LucideIcon } from "lucide-react";
import { cn } from "@/lib/utils";
import { Skeleton } from "@/components/ui/skeleton";

interface StatCardProps {
  title: string;
  value: string;
  subtitle?: string;
  icon: LucideIcon;
  trend?: {
    value: string;
    positive: boolean;
  };
  variant?: "default" | "peach" | "lavender" | "mint" | "sky";
  isLoading?: boolean;
}

const variantStyles = {
  default: "bg-card",
  peach: "bg-peach",
  lavender: "bg-lavender",
  mint: "bg-mint",
  sky: "bg-sky",
};

const iconVariantStyles = {
  default: "bg-muted text-muted-foreground",
  peach: "bg-peach-dark/20 text-peach-dark",
  lavender: "bg-lavender-dark/20 text-lavender-dark",
  mint: "bg-mint-dark/20 text-mint-dark",
  sky: "bg-sky-dark/20 text-sky-dark",
};

export function StatCard({ 
  title, 
  value, 
  subtitle, 
  icon: Icon, 
  trend,
  variant = "default",
  isLoading = false
}: StatCardProps) {
  if (isLoading) {
    return (
      <div 
        className={cn(
          "stat-card animate-fade-in",
          variantStyles[variant]
        )}
        style={{ animationDelay: "0.1s" }}
      >
        <div className="flex items-start justify-between mb-4">
          <Skeleton className="w-11 h-11 rounded-xl" />
          <Skeleton className="w-16 h-6 rounded-lg" />
        </div>
        
        <Skeleton className="h-4 w-20 mb-2" />
        <Skeleton className="h-8 w-32 mb-2" />
        {subtitle && <Skeleton className="h-4 w-24" />}
      </div>
    );
  }

  return (
    <div 
      className={cn(
        "stat-card animate-fade-in",
        variantStyles[variant]
      )}
      style={{ animationDelay: "0.1s" }}
    >
      <div className="flex items-start justify-between mb-4">
        <div 
          className={cn(
            "w-11 h-11 rounded-xl flex items-center justify-center",
            iconVariantStyles[variant]
          )}
        >
          <Icon className="w-5 h-5" />
        </div>
        {trend && (
          <span 
            className={cn(
              "text-sm font-semibold px-2 py-1 rounded-lg",
              trend.positive 
                ? "bg-mint text-mint-dark" 
                : "bg-coral text-coral-dark"
            )}
          >
            {trend.positive ? "↑" : "↓"} {trend.value}
          </span>
        )}
      </div>
      
      <p className="text-sm font-medium text-muted-foreground mb-1">
        {title}
      </p>
      <p className="text-2xl font-bold text-foreground">
        {value}
      </p>
      {subtitle && (
        <p className="text-sm text-muted-foreground mt-1">
          {subtitle}
        </p>
      )}
    </div>
  );
}
