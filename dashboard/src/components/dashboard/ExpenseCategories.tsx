import { useState } from "react";
import { ShoppingCart, Car, Utensils, Fuel, Heart, Gamepad2, Wallet, CreditCard, GraduationCap, Home, PiggyBank, HelpCircle, CupSoda, Pill, Zap } from "lucide-react";
import { cn } from "@/lib/utils";
import { useCategoriesAnalytics, CategoryItem } from "@/lib/api";
import { Skeleton } from "@/components/ui/skeleton";
import { CategoryTransactionsModal } from "./CategoryTransactionsModal";

// Category Icons mapped to their purpose:
// Food = Makan sehari-hari (makan berat)
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

const categoryColors: { [key: string]: "peach" | "lavender" | "mint" | "lemon" | "coral" | "sky" | "rose" | "amber" | "indigo" | "teal" | "purple" | "orange" | "emerald" | "cyan" | "fuchsia" | "lime" } = {
  "Food": "orange",                     // Food - warm orange
  "Salary": "emerald",                  // Income - fresh emerald green
  "Other": "lavender",                  // Neutral - soft lavender
  "Transportation": "indigo",           // Transportation - deep blue
  "Loans & Fixed Payments": "coral",    // Debt - warning coral/red
  "Lifestyle & Entertainment": "fuchsia", // Fun - vibrant fuchsia
  "Health/Medical": "rose",             // Health - caring rose/pink
  "Education": "cyan",                  // Education - bright cyan
  "Groceries": "teal",                  // Groceries - fresh teal
  "Miscellaneous": "purple",            // Unexpected - attention purple
  "Saving & Investments": "lime",       // Savings - positive lime green
  "Housing & Utilities": "sky",         // Utilities - sky blue
  "⁠Snack & Beverages": "amber",         // Snacks - playful amber
};

const colorStyles = {
  peach: "bg-peach text-peach-dark",
  lavender: "bg-lavender text-lavender-dark",
  mint: "bg-mint text-mint-dark",
  lemon: "bg-lemon text-lemon-dark",
  coral: "bg-coral text-coral-dark",
  sky: "bg-sky-100 text-sky-700",
  rose: "bg-rose-100 text-rose-700",
  amber: "bg-amber-100 text-amber-700",
  indigo: "bg-indigo-100 text-indigo-700",
  teal: "bg-teal-100 text-teal-700",
  purple: "bg-purple-100 text-purple-700",
  orange: "bg-orange-100 text-orange-700",
  emerald: "bg-emerald-100 text-emerald-700",
  cyan: "bg-cyan-100 text-cyan-700",
  fuchsia: "bg-fuchsia-100 text-fuchsia-700",
  lime: "bg-lime-100 text-lime-700",
};

const progressColors = {
  peach: "bg-peach-dark",
  lavender: "bg-lavender-dark",
  mint: "bg-mint-dark",
  lemon: "bg-lemon-dark",
  coral: "bg-coral-dark",
  sky: "bg-sky-500",
  rose: "bg-rose-500",
  amber: "bg-amber-500",
  indigo: "bg-indigo-500",
  teal: "bg-teal-500",
  purple: "bg-purple-500",
  orange: "bg-orange-500",
  emerald: "bg-emerald-500",
  cyan: "bg-cyan-500",
  fuchsia: "bg-fuchsia-500",
  lime: "bg-lime-500",
};

interface ExpenseCategoriesProps {
  startDate: string;
  endDate: string;
}

export function ExpenseCategories({ startDate, endDate }: ExpenseCategoriesProps) {
  const { data: categoriesData, isLoading, error } = useCategoriesAnalytics(startDate, endDate);
  const [selectedCategory, setSelectedCategory] = useState<CategoryItem | null>(null);
  const [modalOpen, setModalOpen] = useState(false);

  const categories = categoriesData?.data?.data || [];

  const handleCategoryClick = (category: CategoryItem) => {
    setSelectedCategory(category);
    setModalOpen(true);
  };

  if (isLoading) {
    return (
      <div className="stat-card animate-fade-in" style={{ animationDelay: "0.2s" }}>
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-lg font-bold text-foreground">Expense Categories</h3>
          <span className="text-sm text-muted-foreground">Loading...</span>
        </div>

        <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
          {[...Array(6)].map((_, i) => (
            <Skeleton key={i} className="h-24 rounded-lg" />
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="stat-card animate-fade-in" style={{ animationDelay: "0.2s" }}>
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-lg font-bold text-foreground">Expense Categories</h3>
        </div>
        <div className="flex items-center justify-center h-32">
          <p className="text-red-500">Failed to load categories</p>
        </div>
      </div>
    );
  }

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(value);
  };

  return (
    <div className="stat-card animate-fade-in" style={{ animationDelay: "0.2s" }}>
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg font-bold text-foreground">Expense Categories</h3>
        <span className="text-sm text-muted-foreground">{categories.length} categories</span>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
        {categories.map((category, index) => {
          const IconComponent = categoryIcons[category.category_name] || ShoppingCart;
          const color = categoryColors[category.category_name] || "peach";

          return (
            <div
              key={category.category_id}
              className={cn(
                "category-card cursor-pointer",
                colorStyles[color]
              )}
              style={{ animationDelay: `${0.1 * (index + 1)}s` }}
              onClick={() => handleCategoryClick(category)}
            >
              <div className="flex items-center gap-2 mb-3">
                <IconComponent className="w-4 h-4" />
                <span className="text-xs font-semibold truncate">{category.category_name}</span>
              </div>
              <p className="text-sm font-bold mb-2">{formatCurrency(category.total_amount)}</p>
              <div className="progress-bar bg-foreground/10">
                <div
                  className={cn("progress-fill", progressColors[color])}
                  style={{
                    width: `${category.percentage}%`,
                    "--progress-width": `${category.percentage}%`
                  } as React.CSSProperties}
                />
              </div>
              <p className="text-xs mt-1 opacity-80">{category.percentage.toFixed(1)}% of total</p>
            </div>
          );
        })}
      </div>

      <CategoryTransactionsModal
        category={selectedCategory}
        open={modalOpen}
        onOpenChange={setModalOpen}
      />
    </div>
  );
}
