import { Plane, Home, GraduationCap, Car } from "lucide-react";
import { cn } from "@/lib/utils";

const goals = [
  {
    name: "Family Vacation",
    target: 5000000,
    current: 3200000,
    icon: Plane,
    color: "peach" as const,
    deadline: "Aug 2025"
  },
  {
    name: "Home Renovation",
    target: 15000000,
    current: 8500000,
    icon: Home,
    color: "mint" as const,
    deadline: "Dec 2025"
  },
  {
    name: "Kids Education",
    target: 20000000,
    current: 12000000,
    icon: GraduationCap,
    color: "lavender" as const,
    deadline: "Sep 2026"
  },
  {
    name: "New Car Fund",
    target: 25000000,
    current: 5000000,
    icon: Car,
    color: "sky" as const,
    deadline: "Jan 2027"
  },
];

const iconStyles = {
  peach: "bg-peach text-peach-dark",
  mint: "bg-mint text-mint-dark",
  lavender: "bg-lavender text-lavender-dark",
  sky: "bg-sky text-sky-dark",
};

const progressStyles = {
  peach: "bg-peach-dark",
  mint: "bg-mint-dark",
  lavender: "bg-lavender-dark",
  sky: "bg-sky-dark",
};

export function SavingsGoals() {
  return (
    <div className="stat-card animate-fade-in" style={{ animationDelay: "0.3s" }}>
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg font-bold text-foreground">Savings Goals</h3>
        <button className="text-sm font-medium text-ring hover:underline">
          + Add Goal
        </button>
      </div>

      <div className="space-y-4">
        {goals.map((goal, index) => {
          const percentage = Math.round((goal.current / goal.target) * 100);
          
          return (
            <div 
              key={goal.name} 
              className="p-4 rounded-xl bg-muted/50 hover:bg-muted transition-colors duration-200"
              style={{ animationDelay: `${0.1 * (index + 1)}s` }}
            >
              <div className="flex items-center gap-3 mb-3">
                <div className={cn(
                  "w-10 h-10 rounded-xl flex items-center justify-center",
                  iconStyles[goal.color]
                )}>
                  <goal.icon className="w-5 h-5" />
                </div>
                <div className="flex-1 min-w-0">
                  <h4 className="font-semibold text-foreground truncate">
                    {goal.name}
                  </h4>
                  <p className="text-xs text-muted-foreground">
                    Target: {goal.deadline}
                  </p>
                </div>
                <div className="text-right">
                  <p className="font-bold text-foreground">
                    Rp {goal.current.toLocaleString('id-ID')}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    of Rp {goal.target.toLocaleString('id-ID')}
                  </p>
                </div>
              </div>
              
              <div className="progress-bar">
                <div 
                  className={cn("progress-fill animate-progress-fill", progressStyles[goal.color])}
                  style={{ 
                    "--progress-width": `${percentage}%`,
                    width: `${percentage}%`
                  } as React.CSSProperties}
                />
              </div>
              <p className="text-xs text-muted-foreground mt-2 text-right">
                {percentage}% complete
              </p>
            </div>
          );
        })}
      </div>
    </div>
  );
}
