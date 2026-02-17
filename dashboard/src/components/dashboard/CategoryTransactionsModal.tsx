import { useState } from "react";
import { ShoppingCart, CupSoda, Fuel, Utensils, Gamepad2, Wallet, CreditCard, GraduationCap, Home, PiggyBank, HelpCircle, Pill, Zap, ChevronLeft, ChevronRight } from "lucide-react";
import { cn } from "@/lib/utils";
import { useCategoryTransactions, CategoryItem } from "@/lib/api";
import { formatDistanceToNow, format } from "date-fns";
import { id as idLocale } from "date-fns/locale";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

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
  "Food": "peach",
  "Salary": "mint",
  "Other": "lavender",
  "Transportation": "indigo",
  "Loans & Fixed Payments": "coral",
  "Lifestyle & Entertainment": "purple",
  "Health/Medical": "rose",
  "Education": "sky",
  "Groceries": "teal",
  "Miscellaneous": "amber",
  "Saving & Investments": "mint",
  "Housing & Utilities": "lemon",
  "⁠Snack & Beverages": "orange",
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

interface CategoryTransactionsModalProps {
  category: CategoryItem | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function CategoryTransactionsModal({ category, open, onOpenChange }: CategoryTransactionsModalProps) {
  const [currentPage, setCurrentPage] = useState(1);
  const { data: transactionsData, isLoading, error } = useCategoryTransactions(
    category?.category_id || null,
    currentPage,
    10
  );

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

  const IconComponent = category ? (categoryIcons[category.category_name] || ShoppingCart) : ShoppingCart;
  const color = category ? (categoryColors[category.category_name] || "peach") : "peach";

  // Reset page when category changes
  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      setCurrentPage(1);
    }
    onOpenChange(newOpen);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="max-w-lg max-h-[80vh] overflow-hidden flex flex-col">
        <DialogHeader className="pb-4 border-b border-border">
          <DialogTitle className="flex items-center gap-3">
            <div className={cn(
              "w-10 h-10 rounded-xl flex items-center justify-center",
              iconStyles[color]
            )}>
              <IconComponent className="w-5 h-5" />
            </div>
            <div>
              <span className="block">{category?.category_name}</span>
              <span className="text-sm font-normal text-muted-foreground">
                {formatCurrency(category?.total_amount || 0)} • {category?.count} transactions
              </span>
            </div>
          </DialogTitle>
        </DialogHeader>

        <div className="flex-1 overflow-y-auto">
          {isLoading ? (
            <div className="space-y-3">
              {[...Array(5)].map((_, i) => (
                <Skeleton key={i} className="h-16 rounded-lg" />
              ))}
            </div>
          ) : error ? (
            <div className="flex items-center justify-center h-32">
              <p className="text-red-500">Failed to load transactions</p>
            </div>
          ) : transactions.length === 0 ? (
            <div className="flex items-center justify-center h-32">
              <p className="text-muted-foreground">No transactions found</p>
            </div>
          ) : (
            <div className="space-y-3">
              {transactions.map((transaction) => (
                <div
                  key={transaction.id}
                  className="flex items-center gap-3 p-3 rounded-xl hover:bg-muted/50 transition-colors duration-200"
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
                      {formatDate(transaction.transaction_date)}
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
              ))}
            </div>
          )}
        </div>

        {pagination && pagination.total_pages > 1 && (
          <div className="flex items-center justify-between pt-4 border-t border-border">
            <div className="text-xs text-muted-foreground">
              Page {pagination.page} of {pagination.total_pages}
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
      </DialogContent>
    </Dialog>
  );
}