import { useState } from "react";
import { Calendar } from "lucide-react";
import { format } from "date-fns";
import { id } from "date-fns/locale";
import { Button } from "@/components/ui/button";
import { Calendar as CalendarComponent } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

interface HeaderProps {
  startDate: Date | undefined;
  endDate: Date | undefined;
  onDateChange: (startDate: Date | undefined, endDate: Date | undefined) => void;
}

export function Header({
  startDate,
  endDate,
  onDateChange,
}: HeaderProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [tempStartDate, setTempStartDate] = useState<Date | undefined>(startDate);
  const [tempEndDate, setTempEndDate] = useState<Date | undefined>(endDate);

  const handleApply = () => {
    onDateChange(tempStartDate, tempEndDate);
    setIsOpen(false);
  };

  const handleReset = () => {
    const today = new Date();
    const startOfMonth = new Date(today.getFullYear(), today.getMonth(), 1);
    setTempStartDate(startOfMonth);
    setTempEndDate(today);
  };

  const dateRangeLabel = startDate && endDate
    ? `${format(startDate, "d MMM yyyy", { locale: id })} - ${format(endDate, "d MMM yyyy", { locale: id })}`
    : "Pilih tanggal";

  return (
    <header className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 sm:gap-4 mb-6 sm:mb-8">
      <div className="min-w-0">
        <h1 className="text-xl sm:text-2xl lg:text-3xl font-bold text-foreground truncate">
          Hi Family! 👋
        </h1>
        <p className="text-xs sm:text-sm text-muted-foreground mt-0.5 sm:mt-1 truncate">
          Your financial overview for {dateRangeLabel}
        </p>
      </div>

      <Popover open={isOpen} onOpenChange={setIsOpen}>
        <PopoverTrigger asChild>
          <Button 
            variant="outline" 
            className="w-full sm:w-fit bg-card border-border hover:bg-muted flex-shrink-0 text-xs sm:text-sm"
          >
            <Calendar className="mr-2 h-3 w-3 sm:h-4 sm:w-4" />
            <span className="truncate sm:whitespace-nowrap">{dateRangeLabel}</span>
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-auto p-4 bg-card border-border" align="end">
          <div className="space-y-4">
            <div className="space-y-2">
              <label className="text-sm font-medium text-foreground">
                Tanggal Mulai
              </label>
              <CalendarComponent
                mode="single"
                selected={tempStartDate}
                onSelect={setTempStartDate}
                disabled={(date) =>
                  tempEndDate ? date > tempEndDate : false
                }
                className="rounded-md border-border"
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-foreground">
                Tanggal Akhir
              </label>
              <CalendarComponent
                mode="single"
                selected={tempEndDate}
                onSelect={setTempEndDate}
                disabled={(date) =>
                  tempStartDate ? date < tempStartDate : false
                }
                className="rounded-md border-border"
              />
            </div>

            <div className="flex gap-2 justify-end pt-4 border-t border-border">
              <Button
                variant="outline"
                size="sm"
                onClick={handleReset}
                className="bg-card border-border hover:bg-muted"
              >
                Reset
              </Button>
              <Button
                size="sm"
                onClick={handleApply}
                className="bg-mint-dark hover:bg-mint-dark/90 text-white"
              >
                Terapkan
              </Button>
            </div>
          </div>
        </PopoverContent>
      </Popover>
    </header>
  );
}
