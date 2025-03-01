
import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import PageLayout from "@/components/layout/PageLayout";
import { getWords } from "@/lib/api";
import { PaginatedResponse, Word } from "@/types";
import DataTable, { Column } from "@/components/shared/DataTable";
import Pagination from "@/components/shared/Pagination";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { CheckCircle, FileText, Search, XCircle } from "lucide-react";

export default function Words() {
  const [searchParams, setSearchParams] = useSearchParams();
  const navigate = useNavigate();
  
  const currentPage = parseInt(searchParams.get("page") || "1");
  const searchQuery = searchParams.get("q") || "";
  
  const [words, setWords] = useState<PaginatedResponse<Word> | null>(null);
  const [loading, setLoading] = useState(true);
  const [localSearch, setLocalSearch] = useState(searchQuery);
  
  useEffect(() => {
    const fetchWords = async () => {
      try {
        setLoading(true);
        const data = await getWords(currentPage);
        
        // Client-side filtering based on search query if provided
        if (searchQuery) {
          const filtered = {
            ...data,
            items: data.items.filter(
              (word) => 
                word.portuguese.toLowerCase().includes(searchQuery.toLowerCase()) ||
                word.english.toLowerCase().includes(searchQuery.toLowerCase())
            ),
          };
          setWords(filtered);
        } else {
          setWords(data);
        }
      } catch (error) {
        console.error("Failed to fetch words:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchWords();
  }, [currentPage, searchQuery]);

  const handlePageChange = (page: number) => {
    setSearchParams({ page: String(page), ...(searchQuery && { q: searchQuery }) });
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearchParams({ ...(localSearch && { q: localSearch }), page: "1" });
  };

  const columns: Column<Word>[] = [
    {
      header: "Portuguese",
      accessor: "portuguese",
      className: "font-medium cursor-pointer hover:text-primary transition-colors",
    },
    {
      header: "English",
      accessor: "english",
    },
    {
      header: "Correct",
      accessor: (word) => (
        <div className="flex items-center justify-end">
          <Badge variant="outline" className="bg-green-500/10 text-green-700 flex items-center gap-1">
            <CheckCircle className="h-3 w-3" />
            <span>{word.correct_count}</span>
          </Badge>
        </div>
      ),
      className: "text-right",
    },
    {
      header: "Wrong",
      accessor: (word) => (
        <div className="flex items-center justify-end">
          <Badge variant="outline" className="bg-red-500/10 text-red-700 flex items-center gap-1">
            <XCircle className="h-3 w-3" />
            <span>{word.wrong_count}</span>
          </Badge>
        </div>
      ),
      className: "text-right",
    },
  ];

  return (
    <PageLayout>
      <div className="space-y-6">
        <section className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold tracking-tight flex items-center gap-2">
              <FileText className="h-7 w-7 text-primary" />
              <span>Words</span>
            </h1>
            <p className="text-muted-foreground mt-1">
              Browse all vocabulary words in the database
            </p>
          </div>
          
          <form onSubmit={handleSearch} className="w-full sm:w-auto">
            <div className="relative">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                type="search"
                placeholder="Search words..."
                className="pl-9"
                value={localSearch}
                onChange={(e) => setLocalSearch(e.target.value)}
              />
            </div>
          </form>
        </section>

        <DataTable
          columns={columns}
          data={words?.items || []}
          keyExtractor={(word) => word.id}
          loading={loading}
          onRowClick={(word) => navigate(`/words/${word.id}`)}
          emptyMessage={
            searchQuery
              ? `No words found matching "${searchQuery}"`
              : "No words available"
          }
        />
        
        {words && (
          <div className="flex justify-center">
            <Pagination
              pagination={words.pagination}
              onPageChange={handlePageChange}
            />
          </div>
        )}
      </div>
    </PageLayout>
  );
}
