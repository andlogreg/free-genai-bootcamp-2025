
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { cn } from "@/lib/utils";
import { ReactNode, useState } from "react";

export interface Column<T> {
  header: string;
  accessor: keyof T | ((item: T) => ReactNode);
  className?: string;
}

interface DataTableProps<T> {
  columns: Column<T>[];
  data: T[];
  keyExtractor: (item: T) => string | number;
  className?: string;
  onRowClick?: (item: T) => void;
  loading?: boolean;
  emptyMessage?: string;
}

export default function DataTable<T>({
  columns,
  data,
  keyExtractor,
  className,
  onRowClick,
  loading = false,
  emptyMessage = "No data available",
}: DataTableProps<T>) {
  const [hoveredRow, setHoveredRow] = useState<string | number | null>(null);

  // Handle empty data
  if (!loading && data.length === 0) {
    return (
      <div className="bg-card rounded-lg border shadow-sm p-8 text-center">
        <p className="text-muted-foreground">{emptyMessage}</p>
      </div>
    );
  }

  return (
    <div className={cn("rounded-lg border shadow-sm overflow-hidden", className)}>
      <Table>
        <TableHeader>
          <TableRow>
            {columns.map((column, index) => (
              <TableHead
                key={`header-${index}`}
                className={cn("whitespace-nowrap", column.className)}
              >
                {column.header}
              </TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody>
          {loading ? (
            // Loading state with skeleton rows
            Array.from({ length: 5 }).map((_, rowIndex) => (
              <TableRow key={`skeleton-${rowIndex}`}>
                {columns.map((_, colIndex) => (
                  <TableCell key={`skeleton-cell-${rowIndex}-${colIndex}`}>
                    <div className="h-4 bg-muted rounded animate-pulse" />
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            // Actual data rows
            data.map((item) => {
              const rowKey = keyExtractor(item);
              return (
                <TableRow
                  key={rowKey}
                  className={cn(
                    onRowClick && "cursor-pointer transition-colors hover:bg-muted/50",
                    hoveredRow === rowKey && "bg-muted/50"
                  )}
                  onClick={() => onRowClick?.(item)}
                  onMouseEnter={() => setHoveredRow(rowKey)}
                  onMouseLeave={() => setHoveredRow(null)}
                >
                  {columns.map((column, colIndex) => {
                    const cellContent =
                      typeof column.accessor === "function"
                        ? column.accessor(item)
                        : (item[column.accessor] as ReactNode);
                    
                    return (
                      <TableCell
                        key={`cell-${rowKey}-${colIndex}`}
                        className={column.className}
                      >
                        {cellContent}
                      </TableCell>
                    );
                  })}
                </TableRow>
              );
            })
          )}
        </TableBody>
      </Table>
    </div>
  );
}
