import { useQuery } from "@tanstack/react-query";

const BASE_URL = import.meta.env.VITE_APP_BASE_URL;

export interface DashboardData {
  data: {
    total_balance: number;
    monthly_income: number;
    monthly_income_change: number;
    monthly_expense: number;
    monthly_expense_change: number;
    month: number;
    year: number;
    previous_month: number;
    previous_year: number;
  };
  message: string;
  debug: {
    requestId: string;
    version: string;
    error: null | string;
    startTime: string;
    endTime: string;
    runtimeMs: number;
  };
}

export interface MonthlyAnalyticsItem {
  month: number;
  month_name: string;
  year: number;
  total_income: number;
  total_expense: number;
  balance: number;
  transaction_count: number;
}

export interface MonthlyAnalyticsData {
  data: {
    data: MonthlyAnalyticsItem[];
  };
  message: string;
  debug: {
    requestId: string;
    version: string;
    error: null | string;
    startTime: string;
    endTime: string;
    runtimeMs: number;
  };
}

export interface CategoryItem {
  category_id: number;
  category_name: string;
  type: "EXPENSE" | "INCOME";
  total_amount: number;
  count: number;
  percentage: number;
  average_amount: number;
}

export interface CategoriesAnalyticsData {
  data: {
    data: CategoryItem[];
    total_amount: number;
    transaction_count: number;
    period: {
      start_date: string;
      end_date: string;
    };
    top_category: CategoryItem;
  };
  message: string;
  debug: {
    requestId: string;
    version: string;
    error: null | string;
    startTime: string;
    endTime: string;
    runtimeMs: number;
  };
}

export interface TransactionCategory {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: null | string;
  name: string;
}

export interface Transaction {
  id: number;
  user_id: number;
  category_id: number;
  category: TransactionCategory;
  amount: number;
  description: string;
  transaction_date: string;
  type: "EXPENSE" | "INCOME";
  created_at: string;
  updated_at: string;
}

export interface TransactionsPagination {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface TransactionsData {
  data: {
    transactions: Transaction[];
    pagination: TransactionsPagination;
  };
  message: string;
  debug: {
    requestId: string;
    version: string;
    error: null | string;
    startTime: string;
    endTime: string;
    runtimeMs: number;
  };
}
  

export const useTransactions = (
  startDate: string,
  endDate: string,
  page: number = 1,
  limit: number = 10
) => {
  return useQuery<TransactionsData>({
    queryKey: ["transactions", startDate, endDate, page, limit],
    queryFn: async () => {
      const response = await fetch(
        `${BASE_URL}/api/transactions?start_date=${startDate}&end_date=${endDate}&page=${page}&limit=${limit}`
      );

      if (!response.ok) {
        throw new Error("Failed to fetch transactions");
      }

      return response.json();
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
};

export const useDashboardAnalytics = (startDate: string, endDate: string) => {
  return useQuery<DashboardData>({
    queryKey: ["dashboardAnalytics", startDate, endDate],
    queryFn: async () => {
      const response = await fetch(
        `${BASE_URL}/api/analytics/dashboard?start_date=${startDate}&end_date=${endDate}`
      );

      if (!response.ok) {
        throw new Error("Failed to fetch dashboard analytics");
      }

      return response.json();
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useCategoriesAnalytics = (startDate: string, endDate: string) => {
  return useQuery<CategoriesAnalyticsData>({
    queryKey: ["categoriesAnalytics", startDate, endDate],
    queryFn: async () => {
      const response = await fetch(
        `${BASE_URL}/api/analytics/categories?start_date=${startDate}&end_date=${endDate}`
      );

      if (!response.ok) {
        throw new Error("Failed to fetch categories analytics");
      }

      return response.json();
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useMonthlyAnalytics = (year: string) => {
  return useQuery<MonthlyAnalyticsData>({
    queryKey: ["monthlyAnalytics", year],
    queryFn: async () => {
      const response = await fetch(
        `${BASE_URL}/api/analytics/monthly?year=${year}`
      );

      if (!response.ok) {
        throw new Error("Failed to fetch monthly analytics");
      }

      return response.json();
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useCategoryTransactions = (
  categoryId: number | null,
  page: number = 1,
  limit: number = 10
) => {
  return useQuery<TransactionsData>({
    queryKey: ["categoryTransactions", categoryId, page, limit],
    queryFn: async () => {
      const response = await fetch(
        `${BASE_URL}/api/transactions?category_id=${categoryId}&page=${page}&limit=${limit}`
      );

      if (!response.ok) {
        throw new Error("Failed to fetch category transactions");
      }

      return response.json();
    },
    enabled: categoryId !== null,
    staleTime: 2 * 60 * 1000,
  });
};
