
import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import PageLayout from "@/components/layout/PageLayout";
import { getStudyActivity, getStudyActivitySessions } from "@/lib/api";
import { PaginatedResponse, StudyActivity, StudySession } from "@/types";
import { Button } from "@/components/ui/button";
import { ArrowLeft, ArrowUpRight, Calendar } from "lucide-react";
import DataTable, { Column } from "@/components/shared/DataTable";
import Pagination from "@/components/shared/Pagination";
import { formatDistanceToNow, format } from "date-fns";

export default function StudyActivityShow() {
  const { id } = useParams<{ id: string }>();
  const activityId = id ? parseInt(id) : 0;
  
  const [activity, setActivity] = useState<StudyActivity | null>(null);
  const [sessions, setSessions] = useState<PaginatedResponse<StudySession> | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [loadingSessions, setLoadingSessions] = useState(true);

  useEffect(() => {
    const fetchActivityDetails = async () => {
      if (!activityId) return;
      
      try {
        setLoading(true);
        const activityData = await getStudyActivity(activityId);
        setActivity(activityData);
      } catch (error) {
        console.error("Failed to fetch activity details:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchActivityDetails();
  }, [activityId]);

  useEffect(() => {
    const fetchSessions = async () => {
      if (!activityId) return;
      
      try {
        setLoadingSessions(true);
        const sessionsData = await getStudyActivitySessions(activityId, currentPage);
        setSessions(sessionsData);
      } catch (error) {
        console.error("Failed to fetch activity sessions:", error);
      } finally {
        setLoadingSessions(false);
      }
    };

    fetchSessions();
  }, [activityId, currentPage]);

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

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  return (
    <PageLayout>
      <div className="space-y-8">
        <div className="flex items-center gap-2">
          <Button asChild variant="link" className="p-0 h-auto font-medium">
            <Link to="/study_activities" className="flex items-center gap-1">
              <ArrowLeft className="h-4 w-4" />
              <span>Back to Activities</span>
            </Link>
          </Button>
        </div>

        {loading ? (
          <div className="grid gap-6">
            <div className="h-8 bg-muted rounded w-1/3 animate-pulse" />
            <div className="h-[200px] bg-muted rounded animate-pulse" />
            <div className="h-4 bg-muted rounded w-2/3 animate-pulse" />
          </div>
        ) : activity ? (
          <div className="grid gap-6 md:grid-cols-2">
            <div className="space-y-4">
              <h1 className="text-3xl font-bold tracking-tight">{activity.name}</h1>
              <p className="text-muted-foreground">{activity.description}</p>
              
              <div className="pt-4">
                <Button asChild>
                  <Link to={`/study_activities/${activityId}/launch`} className="flex items-center gap-1">
                    <span>Launch Activity</span>
                    <ArrowUpRight className="h-4 w-4" />
                  </Link>
                </Button>
              </div>
            </div>
            
            <div className="relative rounded-lg overflow-hidden aspect-video shadow-md">
              <img
                src={activity.thumbnail_url}
                alt={activity.name}
                className="absolute inset-0 w-full h-full object-cover"
              />
            </div>
          </div>
        ) : (
          <div className="text-center py-12">
            <p className="text-muted-foreground">Activity not found</p>
          </div>
        )}

        <div className="pt-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold flex items-center gap-2">
              <Calendar className="h-5 w-5 text-primary" />
              <span>Study Sessions</span>
            </h2>
          </div>
          
          <DataTable
            columns={columns}
            data={sessions?.items || []}
            keyExtractor={(session) => session.id}
            loading={loadingSessions}
            onRowClick={(session) => 
              window.location.href = `/study_sessions/${session.id}`
            }
            emptyMessage="No study sessions found for this activity"
          />
          
          {sessions && (
            <div className="mt-4 flex justify-center">
              <Pagination
                pagination={sessions.pagination}
                onPageChange={handlePageChange}
              />
            </div>
          )}
        </div>
      </div>
    </PageLayout>
  );
}
