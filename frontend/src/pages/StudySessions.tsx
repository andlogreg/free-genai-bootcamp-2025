
import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import PageLayout from "@/components/layout/PageLayout";
import { getStudySessions } from "@/lib/api";
import { PaginatedResponse, StudySession } from "@/types";
import DataTable, { Column } from "@/components/shared/DataTable";
import Pagination from "@/components/shared/Pagination";
import { format } from "date-fns";
import { Calendar } from "lucide-react";

export default function StudySessions() {
  const [searchParams, setSearchParams] = useSearchParams();
  const navigate = useNavigate();
  
  const currentPage = parseInt(searchParams.get("page") || "1");
  
  const [sessions, setSessions] = useState<PaginatedResponse<StudySession> | null>(null);
  const [loading, setLoading] = useState(true);
  
  useEffect(() => {
    const fetchSessions = async () => {
      try {
        setLoading(true);
        const data = await getStudySessions(currentPage);
        setSessions(data);
      } catch (error) {
        console.error("Failed to fetch study sessions:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchSessions();
  }, [currentPage]);

  const handlePageChange = (page: number) => {
    setSearchParams({ page: String(page) });
  };

  const formatDate = (dateString: string) => {
    try {
      return format(new Date(dateString), "MMM d, yyyy h:mm a");
    } catch (error) {
      return dateString;
    }
  };

  const columns: Column<StudySession>[] = [
    {
      header: "ID",
      accessor: "id",
      className: "font-medium",
    },
    {
      header: "Activity",
      accessor: "activity_name",
    },
    {
      header: "Group",
      accessor: "group_name",
    },
    {
      header: "Start Time",
      accessor: (session) => formatDate(session.start_time),
    },
    {
      header: "End Time",
      accessor: (session) => formatDate(session.end_time),
    },
    {
      header: "Items",
      accessor: "review_items_count",
      className: "text-right",
    },
  ];

  return (
    <PageLayout>
      <div className="space-y-6">
        <section>
          <h1 className="text-3xl font-bold tracking-tight flex items-center gap-2">
            <Calendar className="h-7 w-7 text-primary" />
            <span>Study Sessions</span>
          </h1>
          <p className="text-muted-foreground mt-1">
            View your history of learning activities
          </p>
        </section>

        <DataTable
          columns={columns}
          data={sessions?.items || []}
          keyExtractor={(session) => session.id}
          loading={loading}
          onRowClick={(session) => navigate(`/study_sessions/${session.id}`)}
          emptyMessage="No study sessions available"
        />
        
        {sessions && (
          <div className="flex justify-center">
            <Pagination
              pagination={sessions.pagination}
              onPageChange={handlePageChange}
            />
          </div>
        )}
      </div>
    </PageLayout>
  );
}
