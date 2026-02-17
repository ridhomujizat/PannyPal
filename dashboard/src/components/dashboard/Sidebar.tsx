import { 
  Home, 
  Wallet, 
  PiggyBank, 
  TrendingUp, 
  Receipt, 
  Menu,
  X
} from "lucide-react";
import { useState } from "react";
import { cn } from "@/lib/utils";

const navItems = [
  { icon: Home, label: "Dashboard", active: true },
  { icon: Wallet, label: "Transactions", active: false },
  { icon: PiggyBank, label: "Savings", active: false },
  { icon: TrendingUp, label: "Investments", active: false },
  { icon: Receipt, label: "Bills", active: false },
];

export function Sidebar() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      {/* Mobile menu button */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="fixed top-4 left-4 z-50 p-2 rounded-xl bg-primary text-primary-foreground lg:hidden shadow-lg"
      >
        {isOpen ? <X size={24} /> : <Menu size={24} />}
      </button>

      {/* Overlay */}
      {isOpen && (
        <div
          className="fixed inset-0 bg-foreground/20 backdrop-blur-sm z-40 lg:hidden"
          onClick={() => setIsOpen(false)}
        />
      )}

      {/* Sidebar - Floating on desktop */}
      <aside
        className={cn(
          "fixed left-0 top-0 h-full bg-sidebar z-40 transition-all duration-300 lg:translate-x-0",
          "w-64 lg:w-[70px] flex flex-col py-6",
          // Floating on desktop
          "lg:left-4 lg:top-4 lg:h-[calc(100vh-2rem)] lg:rounded-2xl lg:shadow-xl",
          isOpen ? "translate-x-0" : "-translate-x-full"
        )}
      >
        {/* Logo */}
        <div className="px-4 mb-6 flex justify-center">
          <div className="w-10 h-10 rounded-xl bg-sidebar-primary flex items-center justify-center">
            <PiggyBank className="w-6 h-6 text-sidebar-primary-foreground" />
          </div>
        </div>

        {/* Navigation */}
        <nav className="flex-1 px-2">
          <ul className="space-y-1">
            {navItems.map((item) => (
              <li key={item.label}>
                <a
                  href="#"
                  className={cn(
                    "flex items-center justify-center lg:justify-center gap-3 p-3 rounded-xl transition-all duration-200",
                    "hover:bg-sidebar-accent group",
                    item.active
                      ? "bg-sidebar-primary text-sidebar-primary-foreground"
                      : "text-sidebar-foreground"
                  )}
                  title={item.label}
                >
                  <item.icon 
                    className={cn(
                      "w-5 h-5 flex-shrink-0",
                      item.active ? "" : "group-hover:scale-110 transition-transform"
                    )} 
                  />
                  <span className="font-medium lg:hidden">
                    {item.label}
                  </span>
                </a>
              </li>
            ))}
          </ul>
        </nav>
      </aside>
    </>
  );
}
