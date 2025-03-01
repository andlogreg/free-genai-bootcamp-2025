
import Navbar from "./Navbar";
import { ReactNode } from "react";
import { cn } from "@/lib/utils";

interface PageLayoutProps {
  children: ReactNode;
  className?: string;
  containerClassName?: string;
  fullWidth?: boolean;
}

export default function PageLayout({
  children,
  className,
  containerClassName,
  fullWidth = false,
}: PageLayoutProps) {
  return (
    <div className="min-h-screen flex flex-col">
      <Navbar />
      <main
        className={cn(
          "flex-1 pt-24 pb-12 px-4 sm:px-6 animate-in",
          className
        )}
      >
        <div
          className={cn(
            fullWidth ? "w-full" : "max-w-7xl mx-auto",
            containerClassName
          )}
        >
          {children}
        </div>
      </main>
    </div>
  );
}
