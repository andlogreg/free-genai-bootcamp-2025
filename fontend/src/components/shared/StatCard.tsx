
import { Card, CardContent } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { ReactNode } from "react";

interface StatCardProps {
  title: string;
  value: ReactNode;
  icon?: ReactNode;
  description?: string;
  className?: string;
  trending?: "up" | "down" | "neutral";
  trendingValue?: string;
}

export default function StatCard({
  title,
  value,
  icon,
  description,
  className,
  trending,
  trendingValue,
}: StatCardProps) {
  return (
    <Card className={cn("overflow-hidden transition-all duration-300 card-hover", className)}>
      <CardContent className="p-6">
        <div className="flex items-start justify-between">
          <div>
            <p className="text-sm font-medium text-muted-foreground mb-1">{title}</p>
            <h3 className="text-2xl font-bold tracking-tight">{value}</h3>
            
            {description && (
              <p className="text-xs text-muted-foreground mt-1">{description}</p>
            )}
            
            {trending && trendingValue && (
              <div className="flex items-center mt-2">
                <span
                  className={cn(
                    "text-xs font-medium flex items-center",
                    trending === "up" && "text-green-500",
                    trending === "down" && "text-red-500",
                    trending === "neutral" && "text-muted-foreground"
                  )}
                >
                  {trendingValue}
                </span>
              </div>
            )}
          </div>
          
          {icon && (
            <div className="p-2 bg-primary/10 rounded-full">
              {icon}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
