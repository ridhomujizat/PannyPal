import { Sidebar } from "@/components/dashboard/Sidebar";
import { Header } from "@/components/dashboard/Header";
import { StatCard } from "@/components/dashboard/StatCard";
import { ExpenseCategories } from "@/components/dashboard/ExpenseCategories";
import { SavingsGoals } from "@/components/dashboard/SavingsGoals";
import { SpendingChart } from "@/components/dashboard/SpendingChart";
import { RecentTransactions } from "@/components/dashboard/RecentTransactions";
import { QuickActions } from "@/components/dashboard/QuickActions";
import { Wallet, TrendingUp, PiggyBank, CreditCard } from "lucide-react";
import { useState } from "react";
import { format } from "date-fns";
import { useDashboardAnalytics } from "@/lib/api";

const Index = () => {
  const today = new Date();
  const startOfMonth = new Date(today.getFullYear(), today.getMonth(), 1);

  const [startDate, setStartDate] = useState<Date>(startOfMonth);
  const [endDate, setEndDate] = useState<Date>(today);

  // Format dates for API (YYYY-MM-DD)
  const startDateString = format(startDate, "yyyy-MM-dd");
  const endDateString = format(endDate, "yyyy-MM-dd");
  const selectedYear = String(today.getFullYear());

  const { data: dashboardData, isLoading, error } = useDashboardAnalytics(
    startDateString,
    endDateString
  );

  const handleDateChange = (start: Date | undefined, end: Date | undefined) => {
    if (start) setStartDate(start);
    if (end) setEndDate(end);
  };

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(value);
  };

  const stats = dashboardData?.data || {
    total_balance: 0,
    monthly_income: 0,
    monthly_income_change: 0,
    monthly_expense: 0,
    monthly_expense_change: 0,
  };

  return (
    <div className="min-h-screen w-full">
      {/* Floating Sidebar */}
      <Sidebar />

      {/* Main Content - offset for floating sidebar on desktop */}
      <main className="min-h-screen p-4 pt-20 sm:p-6 sm:pt-6 lg:pl-24 lg:p-8 overflow-auto">
        <div className="max-w-7xl mx-auto lg:ml-4">
          {/* Header with Date Range Filter */}
          <Header
            startDate={startDate}
            endDate={endDate}
            onDateChange={handleDateChange}
          />

          {error && (
            <div className="mb-8 p-4 bg-red-50 border border-red-200 rounded-lg text-red-800">
              Failed to load dashboard data. Please try again.
            </div>
          )}

          {/* Quick Actions */}
          <div className="mb-8">
            <QuickActions />
          </div>

          {/* Stats Grid */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
            <StatCard
              title="Total Balance"
              value={formatCurrency(stats.total_balance)}
              subtitle="Across all accounts"
              icon={Wallet}
              variant="sky"
              isLoading={isLoading}
            />
            <StatCard
              title="Monthly Income"
              value={formatCurrency(stats.monthly_income)}
              icon={TrendingUp}
              trend={{
                value: `${Math.abs(stats.monthly_income_change)}%`,
                positive: stats.monthly_income_change >= 0,
              }}
              variant="mint"
              isLoading={isLoading}
            />
            <StatCard
              title="Monthly Expenses"
              value={formatCurrency(stats.monthly_expense)}
              icon={CreditCard}
              trend={{
                value: `${Math.abs(stats.monthly_expense_change)}%`,
                positive: !(stats.monthly_expense_change >= 0),
              }}
              variant="peach"
              isLoading={isLoading}
            />
            <StatCard
              title="Total Savings"
              value={formatCurrency(
                0
              )}
              subtitle="Available for goals"
              icon={PiggyBank}
              variant="lavender"
              isLoading={isLoading}
            />
          </div>

          {/* Charts & Categories Section */}
          <div className="grid grid-cols-1 xl:grid-cols-2 gap-6 mb-8">
            <SpendingChart year={selectedYear} />
            <ExpenseCategories startDate={startDateString} endDate={endDateString} />
          </div>

          {/* Goals & Transactions Section */}
          <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
            {/* <SavingsGoals /> */}
            <RecentTransactions startDate={startDateString} endDate={endDateString} />
          </div>
        </div>
      </main>
    </div>
  );
};

export default Index;
