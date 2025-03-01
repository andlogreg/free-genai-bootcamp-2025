
import {
  Pagination as PaginationUI,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { Pagination as PaginationType } from "@/types";

interface PaginationProps {
  pagination: PaginationType;
  onPageChange: (page: number) => void;
  className?: string;
}

export default function Pagination({
  pagination,
  onPageChange,
  className,
}: PaginationProps) {
  const { current_page, total_pages } = pagination;

  // Don't render pagination if there's only one page
  if (total_pages <= 1) {
    return null;
  }

  // Helper to generate page numbers to show
  const getPageNumbers = () => {
    const visiblePages = [];
    // Always show first page
    visiblePages.push(1);

    // Calculate range around current page
    let rangeStart = Math.max(2, current_page - 1);
    let rangeEnd = Math.min(total_pages - 1, current_page + 1);

    // Add ellipsis if necessary after page 1
    if (rangeStart > 2) {
      visiblePages.push("ellipsis1");
    }

    // Add pages around current page
    for (let i = rangeStart; i <= rangeEnd; i++) {
      visiblePages.push(i);
    }

    // Add ellipsis if necessary before last page
    if (rangeEnd < total_pages - 1) {
      visiblePages.push("ellipsis2");
    }

    // Always show last page if more than 1 page
    if (total_pages > 1) {
      visiblePages.push(total_pages);
    }

    return visiblePages;
  };

  const visiblePages = getPageNumbers();

  return (
    <PaginationUI className={className}>
      <PaginationContent>
        <PaginationItem>
          <PaginationPrevious
            href="#"
            onClick={(e) => {
              e.preventDefault();
              if (current_page > 1) {
                onPageChange(current_page - 1);
              }
            }}
            className={current_page <= 1 ? "pointer-events-none opacity-50" : ""}
          />
        </PaginationItem>
        
        {visiblePages.map((page, index) => (
          <PaginationItem key={`page-${page}-${index}`}>
            {page === "ellipsis1" || page === "ellipsis2" ? (
              <PaginationEllipsis />
            ) : (
              <PaginationLink
                href="#"
                onClick={(e) => {
                  e.preventDefault();
                  if (typeof page === "number") {
                    onPageChange(page);
                  }
                }}
                isActive={current_page === page}
              >
                {page}
              </PaginationLink>
            )}
          </PaginationItem>
        ))}
        
        <PaginationItem>
          <PaginationNext
            href="#"
            onClick={(e) => {
              e.preventDefault();
              if (current_page < total_pages) {
                onPageChange(current_page + 1);
              }
            }}
            className={
              current_page >= total_pages ? "pointer-events-none opacity-50" : ""
            }
          />
        </PaginationItem>
      </PaginationContent>
    </PaginationUI>
  );
}
