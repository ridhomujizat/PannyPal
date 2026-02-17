import { ShoppingCart, CupSoda, Fuel, Wifi, Film, Utensils, Heart, Gamepad2, Lightbulb, Car, ChevronLeft, ChevronRight, Wallet, CreditCard, GraduationCap, Home, PiggyBank, HelpCircle, Pill, Zap } from "lucide-react";
import { cn } from "@/lib/utils";
import { useTransactions } from "@/lib/api";
import { formatDistanceToNow, format } from "date-fns";
import { id as idLocale } from "date-fns/locale";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { useState } from "react";

// Category Icons mapped to their purpose:
// Daily Meals (FOOD) = Makan sehari-hari (makan berat)
// Salary = Salary, Side Income, Bonus/Allowances
// Other = NULL
// Transportation = Fuel, Parking, Vehicle Maintenance
// Loans & Fixed Payments = Loan, CC Payment, Fixed Bills (Paylater)
// Lifestyle & Entertainment = Clothing, Personal Care, Entertainment, Subscriptions
// Health/Medical = Medical Expenses, Medication, and Health Insurance
// Education = NULL
// Groceries = Belanjaan yang dibeli di Supermarket
// Miscellaneous = Unexpected Expenses
// Saving & Investments = Emergency Fund, General Savings, Investments
// Housing & Utilities = Rent, Electricity, Water, Internet/Wifi, Cooking Gas, Home Maintenance
// Snack & Beverages = Jajan (cilor, tahu, gorengan) sama Kopi
const categoryIcons: { [key: string]: typeof ShoppingCart } = {
  "Food": Utensils,
  "Salary": Wallet,
  "Other": HelpCircle,
  "Transportation": Fuel,
  "Loans & Fixed Payments": CreditCard,
  "Lifestyle & Entertainment": Gamepad2,
  "Health/Medical": Pill,
  "Education": GraduationCap,
  "Groceries": ShoppingCart,
  "Miscellaneous": Zap,
  "Saving & Investments": PiggyBank,
  "Housing & Utilities": Home,
  "⁠Snack & Beverages": CupSoda,
};

const categoryColors: { [key: string]: "peach" | "mint" | "lavender" | "lemon" | "sky" | "coral" | "rose" | "amber" | "indigo" | "teal" | "purple" | "orange" } = {
  "Food": "peach",       // Food - warm peach
  "Salary": "mint",                     // Income - fresh green/mint
  "Other": "lavender",                  // Neutral - soft lavender
  "Transportation": "indigo",           // Transportation - deep blue
  "Loans & Fixed Payments": "coral",    // Debt - warning coral/red
  "Lifestyle & Entertainment": "purple", // Fun - vibrant purple
  "Health/Medical": "rose",             // Health - caring rose/pink
  "Education": "sky",                   // Education - bright sky blue
  "Groceries": "teal",                  // Groceries - fresh teal
  "Miscellaneous": "amber",             // Unexpected - attention amber
  "Saving & Investments": "mint",       // Savings - positive green/mint
  "Housing & Utilities": "lemon",       // Utilities - electric yellow
  "⁠Snack & Beverages": "orange",        // Snacks - playful orange
};

const iconStyles = {
  peach: "bg-peach text-peach-dark",
  mint: "bg-mint text-mint-dark",
  lavender: "bg-lavender text-lavender-dark",
  lemon: "bg-lemon text-lemon-dark",
  sky: "bg-sky text-sky-dark",
  coral: "bg-coral text-coral-dark",
  rose: "bg-rose-100 text-rose-700",
  amber: "bg-amber-100 text-amber-700",
  indigo: "bg-indigo-100 text-indigo-700",
  teal: "bg-teal-100 text-teal-700",
  purple: "bg-purple-100 text-purple-700",
  orange: "bg-orange-100 text-orange-700",
};

interface RecentTransactionsProps {
  startDate: string;
  endDate: string;
}

export function RecentTransactions({ startDate, endDate }: RecentTransactionsProps) {
  const [currentPage, setCurrentPage] = useState(1);
  const { data: transactionsData, isLoading, error } = useTransactions(startDate, endDate, currentPage, 10);

  const transactions = transactionsData?.data?.transactions || [];
  const pagination = transactionsData?.data?.pagination;

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(Math.abs(value));
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffInHours = (now.getTime() - date.getTime()) / (1000 * 60 * 60);

    if (diffInHours < 24) {
      return formatDistanceToNow(date, { addSuffix: true, locale: idLocale });
    }
    return format(date, "d MMM yyyy", { locale: idLocale });
  };

  if (isLoading) {
    return (
      <div className="stat-card animate-fade-in" style={{ animationDelay: "0.5s" }}>
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-lg font-bold text-foreground">Recent Transactions</h3>
        </div>

        <div className="space-y-3">
          {[...Array(6)].map((_, i) => (
            <Skeleton key={i} className="h-16 rounded-lg" />
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="stat-card animate-fade-in" style={{ animationDelay: "0.5s" }}>
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-lg font-bold text-foreground">Recent Transactions</h3>
        </div>
        <div className="flex items-center justify-center h-32">
          <p className="text-red-500">Failed to load transactions</p>
        </div>
      </div>
    );
  }

  return (
    <div className="stat-card animate-fade-in" style={{ animationDelay: "0.5s" }}>
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg font-bold text-foreground">Recent Transactions</h3>
        <span className="text-xs text-muted-foreground">
          Page {pagination?.page || 1} of {pagination?.total_pages || 1}
        </span>
      </div>

      <div className="space-y-3">
        {transactions.map((transaction, index) => {
          const IconComponent = categoryIcons[transaction.category.name] || ShoppingCart;
          const color = categoryColors[transaction.category.name] || "peach";

          return (
            <div
              key={transaction.id}
              className="flex items-center gap-3 p-3 rounded-xl hover:bg-muted/50 transition-colors duration-200 cursor-pointer"
              style={{ animationDelay: `${0.05 * (index + 1)}s` }}
            >
              <div className={cn(
                "w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0",
                iconStyles[color]
              )}>
                <IconComponent className="w-5 h-5" />
              </div>

              <div className="flex-1 min-w-0">
                <h4 className="font-semibold text-foreground truncate">
                  {transaction.description}
                </h4>
                <p className="text-xs text-muted-foreground">
                  {transaction.category.name} • {formatDate(transaction.transaction_date)}
                </p>
              </div>

              <div className="text-right flex-shrink-0">
                <p className={cn(
                  "font-semibold text-sm",
                  transaction.type === "EXPENSE" ? "text-red-500" : "text-mint-dark"
                )}>
                  {transaction.type === "EXPENSE" ? "-" : "+"}{formatCurrency(transaction.amount)}
                </p>
              </div>
            </div>
          );
        })}
      </div>

      {pagination && pagination.total_pages > 1 && (
        <div className="flex items-center justify-between mt-6 pt-4 border-t border-border">
          <div className="text-xs text-muted-foreground">
            Showing {(pagination.page - 1) * pagination.limit + 1} - {Math.min(pagination.page * pagination.limit, pagination.total)} of {pagination.total}
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
              disabled={pagination.page === 1 || isLoading}
              className="bg-card border-border hover:bg-muted"
            >
              <ChevronLeft className="w-4 h-4" />
            </Button>

            <Button
              variant="outline"
              size="sm"
              onClick={() => setCurrentPage(p => Math.min(pagination.total_pages, p + 1))}
              disabled={pagination.page === pagination.total_pages || isLoading}
              className="bg-card border-border hover:bg-muted"
            >
              <ChevronRight className="w-4 h-4" />
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
