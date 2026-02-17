import {
    BarChart,
    Bar,
    PieChart,
    Pie,
    LineChart,
    Line,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    Legend,
    ResponsiveContainer,
    Cell,
} from "recharts";

interface ChartVisualizationProps {
    type: "bar" | "pie" | "line";
    data: {
        labels: string[];
        values: number[];
        colors?: string[];
    };
    config?: {
        x_label?: string;
        y_label?: string;
        format?: string;
    };
}

const DEFAULT_COLORS = [
    "#FF6384",
    "#36A2EB",
    "#FFCE56",
    "#4BC0C0",
    "#9966FF",
    "#FF9F40",
    "#FF6384",
    "#C9CBCF",
];

export function ChartVisualization({
    type,
    data,
    config,
}: ChartVisualizationProps) {
    // Transform data to recharts format
    const chartData = data.labels.map((label, index) => ({
        name: label,
        value: data.values[index],
    }));

    const colors = data.colors || DEFAULT_COLORS;

    // Format currency for tooltips
    const formatValue = (value: number) => {
        if (config?.format === "currency") {
            return new Intl.NumberFormat("id-ID", {
                style: "currency",
                currency: "IDR",
                minimumFractionDigits: 0,
            }).format(value);
        }
        return value.toLocaleString("id-ID");
    };

    if (type === "bar") {
        return (
            <div className="w-full h-48 sm:h-56 lg:h-64 mt-2">
                <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={chartData}>
                        <CartesianGrid strokeDasharray="3 3" opacity={0.3} />
                        <XAxis
                            dataKey="name"
                            tick={{ fontSize: 12 }}
                            angle={-45}
                            textAnchor="end"
                            height={80}
                        />
                        <YAxis
                            tick={{ fontSize: 12 }}
                            tickFormatter={(value) => {
                                if (config?.format === "currency") {
                                    return `Rp ${(value / 1000000).toFixed(1)}M`;
                                }
                                return value.toLocaleString();
                            }}
                        />
                        <Tooltip
                            formatter={(value: number) => formatValue(value)}
                            contentStyle={{
                                backgroundColor: "rgba(255, 255, 255, 0.95)",
                                border: "1px solid #e5e7eb",
                                borderRadius: "8px",
                                fontSize: "12px",
                            }}
                        />
                        <Bar dataKey="value" radius={[8, 8, 0, 0]}>
                            {chartData.map((_, index) => (
                                <Cell
                                    key={`cell-${index}`}
                                    fill={colors[index % colors.length]}
                                />
                            ))}
                        </Bar>
                    </BarChart>
                </ResponsiveContainer>
            </div>
        );
    }

    if (type === "pie") {
        return (
            <div className="w-full mt-2 overflow-hidden">
                <div className="h-52 sm:h-64 lg:h-72">
                    <ResponsiveContainer width="100%" height="100%">
                        <PieChart>
                            <Pie
                                data={chartData}
                                cx="50%"
                                cy="50%"
                                labelLine={true}
                                label={({ percent }) =>
                                    `${(percent * 100).toFixed(0)}%`
                                }
                                outerRadius="75%"
                                fill="#8884d8"
                                dataKey="value"
                            >
                                {chartData.map((_, index) => (
                                    <Cell
                                        key={`cell-${index}`}
                                        fill={colors[index % colors.length]}
                                    />
                                ))}
                            </Pie>
                            <Tooltip
                                formatter={(value: number) => formatValue(value)}
                                contentStyle={{
                                    backgroundColor: "rgba(255, 255, 255, 0.95)",
                                    border: "1px solid #e5e7eb",
                                    borderRadius: "8px",
                                    fontSize: "12px",
                                }}
                            />
                        </PieChart>
                    </ResponsiveContainer>
                </div>
                {/* Custom legend - grid on mobile, flex-wrap on desktop */}
                <div className="grid grid-cols-2 sm:flex sm:flex-wrap gap-x-4 gap-y-1.5 sm:gap-3 justify-center mt-2 px-1 sm:px-2">
                    {chartData.map((item, index) => (
                        <div key={index} className="flex items-center gap-1.5 min-w-0">
                            <div
                                className="w-2.5 h-2.5 sm:w-3 sm:h-3 rounded-sm flex-shrink-0"
                                style={{
                                    backgroundColor:
                                        colors[index % colors.length],
                                }}
                            />
                            <span className="text-[10px] sm:text-xs text-muted-foreground truncate">
                                {item.name}
                            </span>
                            <span className="text-[10px] sm:text-xs font-medium whitespace-nowrap">
                                {formatValue(item.value)}
                            </span>
                        </div>
                    ))}
                </div>
            </div>
        );
    }

    if (type === "line") {
        return (
            <div className="w-full h-64 mt-2">
                <ResponsiveContainer width="100%" height="100%">
                    <LineChart data={chartData}>
                        <CartesianGrid strokeDasharray="3 3" opacity={0.3} />
                        <XAxis
                            dataKey="name"
                            tick={{ fontSize: 12 }}
                            angle={-45}
                            textAnchor="end"
                            height={80}
                        />
                        <YAxis
                            tick={{ fontSize: 12 }}
                            tickFormatter={(value) => {
                                if (config?.format === "currency") {
                                    return `Rp ${(value / 1000000).toFixed(1)}M`;
                                }
                                return value.toLocaleString();
                            }}
                        />
                        <Tooltip
                            formatter={(value: number) => formatValue(value)}
                            contentStyle={{
                                backgroundColor: "rgba(255, 255, 255, 0.95)",
                                border: "1px solid #e5e7eb",
                                borderRadius: "8px",
                                fontSize: "12px",
                            }}
                        />
                        <Legend />
                        <Line
                            type="monotone"
                            dataKey="value"
                            stroke={colors[0]}
                            strokeWidth={2}
                            dot={{ fill: colors[0], r: 4 }}
                            activeDot={{ r: 6 }}
                        />
                    </LineChart>
                </ResponsiveContainer>
            </div>
        );
    }

    return null;
}
