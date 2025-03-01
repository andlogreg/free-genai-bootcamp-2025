
import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import PageLayout from "@/components/layout/PageLayout";
import { getGroups } from "@/lib/api";
import { Group, PaginatedResponse } from "@/types";
import DataTable, { Column } from "@/components/shared/DataTable";
import Pagination from "@/components/shared/Pagination";
import { Badge } from "@/components/ui/badge";
import { Layers } from "lucide-react";

export default function Groups() {
  const [searchParams, setSearchParams] = useSearchParams();
  const navigate = useNavigate();
  
  const currentPage = parseInt(searchParams.get("page") || "1");
  
  const [groups, setGroups] = useState<PaginatedResponse<Group> | null>(null);
  const [loading, setLoading] = useState(true);
  
  useEffect(() => {
    const fetchGroups = async () => {
      try {
        setLoading(true);
        const data = await getGroups(currentPage);
        setGroups(data);
      } catch (error) {
        console.error("Failed to fetch groups:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchGroups();
  }, [currentPage]);

  const handlePageChange = (page: number) => {
    setSearchParams({ page: String(page) });
  };

  const columns: Column<Group>[] = [
    {
      header: "Group Name",
      accessor: "name",
      className: "font-medium",
    },
    {
      header: "Word Count",
      accessor: (group) => (
        <Badge variant="outline" className="ml-auto">
          {group.word_count} words
        </Badge>
      ),
      className: "text-right",
    },
  ];

  return (
    <PageLayout>
      <div className="space-y-6">
        <section>
          <h1 className="text-3xl font-bold tracking-tight flex items-center gap-2">
            <Layers className="h-7 w-7 text-primary" />
            <span>Word Groups</span>
          </h1>
          <p className="text-muted-foreground mt-1">
            Browse thematic word groups for focused learning
          </p>
        </section>

        <DataTable
          columns={columns}
          data={groups?.items || []}
          keyExtractor={(group) => group.id}
          loading={loading}
          onRowClick={(group) => navigate(`/groups/${group.id}`)}
          emptyMessage="No word groups available"
        />
        
        {groups && (
          <div className="flex justify-center">
            <Pagination
              pagination={groups.pagination}
              onPageChange={handlePageChange}
            />
          </div>
        )}
      </div>
    </PageLayout>
  );
}
