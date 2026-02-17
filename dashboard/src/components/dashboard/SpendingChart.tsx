import { BarChart, Bar, XAxis, YAxis, CartesianGrid, ResponsiveContainer, Tooltip } from "recharts";
import { useMonthlyAnalytics } from "@/lib/api";
import { Skeleton } from "@/components/ui/skeleton";

const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    return (
      <div className="bg-card border border-border rounded-xl p-3 shadow-soft">
        <p className="font-semibold text-foreground mb-2">{label}</p>
        <p className="text-sm text-mint-dark">
          Income: Rp {payload[0]?.value?.toLocaleString('id-ID')}
        </p>
        <p className="text-sm text-peach-dark">
          Expenses: Rp {payload[1]?.value?.toLocaleString('id-ID')}
        </p>
      </div>
    );
  }
  return null;
};

interface SpendingChartProps {
  year?: string;
}

export function SpendingChart({ year }: SpendingChartProps) {
  const currentYear = year || String(new Date().getFullYear());
  const { data: analyticsData, isLoading, error } = useMonthlyAnalytics(currentYear);

  // Transform API data to chart format
  const chartData = analyticsData?.data?.data?.map((item) => ({
    month: item.month_name.slice(0, 3), // "Jan", "Feb", etc.
    income: item.total_income,
    expenses: item.total_expense,
  })) || [];

  if (isLoading) {
    return (
      <div className="stat-card animate-fade-in" style={{ animationDelay: "0.4s" }}>
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-lg font-bold text-foreground">Spending Trends</h3>
          <div className="flex items-center gap-4 text-sm">
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 rounded-full bg-mint-dark" />
              <span className="text-muted-foreground">Income</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 rounded-full bg-peach-dark" />
              <span className="text-muted-foreground">Expenses</span>
            </div>
          </div>
        </div>
        <div className="h-64">
          <Skeleton className="w-full h-full" />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="stat-card animate-fade-in" style={{ animationDelay: "0.4s" }}>
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-lg font-bold text-foreground">Spending Trends</h3>
        </div>
        <div className="h-64 flex items-center justify-center">
          <p className="text-red-500">Failed to load spending trends</p>
        </div>
      </div>
    );
  }

  return (
    <div className="stat-card animate-fade-in" style={{ animationDelay: "0.4s" }}>
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg font-bold text-foreground">Spending Trends</h3>
        <div className="flex items-center gap-4 text-sm">
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 rounded-full bg-mint-dark" />
            <span className="text-muted-foreground">Income</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 rounded-full bg-peach-dark" />
            <span className="text-muted-foreground">Expenses</span>
          </div>
        </div>
      </div>

      <div className="h-64">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={chartData} barGap={8}>
            <CartesianGrid 
              strokeDasharray="3 3" 
              vertical={false}
              stroke="hsl(var(--border))"
            />
            <XAxis 
              dataKey="month" 
              axisLine={false}
              tickLine={false}
              tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
            />
            <YAxis 
              axisLine={false}
              tickLine={false}
              tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
              tickFormatter={(value) => `${value / 1000000}jt`}
            />
            <Tooltip content={<CustomTooltip />} cursor={{ fill: 'hsl(var(--muted) / 0.5)' }} />
            <Bar 
              dataKey="income" 
              fill="hsl(var(--mint-dark))" 
              radius={[6, 6, 0, 0]}
              maxBarSize={40}
            />
            <Bar 
              dataKey="expenses" 
              fill="hsl(var(--peach-dark))" 
              radius={[6, 6, 0, 0]}
              maxBarSize={40}
            />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
