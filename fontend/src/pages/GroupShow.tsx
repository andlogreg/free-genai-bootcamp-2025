
import { useEffect, useState } from "react";
import { useParams, Link, useSearchParams } from "react-router-dom";
import PageLayout from "@/components/layout/PageLayout";
import { ArrowLeft, BookOpen, Calendar } from "lucide-react";
import { getGroup, getGroupWords, getGroupStudySessions } from "@/lib/api";
import { GroupDetail, PaginatedResponse, StudySession, Word } from "@/types";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import DataTable, { Column } from "@/components/shared/DataTable";
import Pagination from "@/components/shared/Pagination";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { formatDistanceToNow, format } from "date-fns";
import { CheckCircle, XCircle } from "lucide-react";

export default function GroupShow() {
  const { id } = useParams<{ id: string }>();
  const groupId = id ? parseInt(id) : 0;
  
  const [searchParams, setSearchParams] = useSearchParams();
  const tab = searchParams.get("tab") || "words";
  const wordsPage = parseInt(searchParams.get("wordsPage") || "1");
  const sessionsPage = parseInt(searchParams.get("sessionsPage") || "1");
  
  const [group, setGroup] = useState<GroupDetail | null>(null);
  const [words, setWords] = useState<PaginatedResponse<Word> | null>(null);
  const [sessions, setSessions] = useState<PaginatedResponse<StudySession> | null>(null);
  const [loading, setLoading] = useState(true);
  const [loadingWords, setLoadingWords] = useState(true);
  const [loadingSessions, setLoadingSessions] = useState(true);

  useEffect(() => {
    const fetchGroupDetails = async () => {
      if (!groupId) return;
      
      try {
        setLoading(true);
        const data = await getGroup(groupId);
        setGroup(data);
      } catch (error) {
        console.error("Failed to fetch group details:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchGroupDetails();
  }, [groupId]);

  useEffect(() => {
    const fetchWords = async () => {
      if (!groupId) return;
      
      try {
        setLoadingWords(true);
        const data = await getGroupWords(groupId, wordsPage);
        setWords(data);
      } catch (error) {
        console.error("Failed to fetch group words:", error);
      } finally {
        setLoadingWords(false);
      }
    };

    if (tab === "words") {
      fetchWords();
    }
  }, [groupId, wordsPage, tab]);

  useEffect(() => {
    const fetchSessions = async () => {
      if (!groupId) return;
      
      try {
        setLoadingSessions(true);
        const data = await getGroupStudySessions(groupId, sessionsPage);
        setSessions(data);
      } catch (error) {
        console.error("Failed to fetch group sessions:", error);
      } finally {
        setLoadingSessions(false);
      }
    };

    if (tab === "sessions") {
      fetchSessions();
    }
  }, [groupId, sessionsPage, tab]);

  const handleTabChange = (value: string) => {
    setSearchParams({ tab: value });
  };

  const handleWordsPageChange = (page: number) => {
    setSearchParams({ tab: "words", wordsPage: String(page) });
  };

  const handleSessionsPageChange = (page: number) => {
    setSearchParams({ tab: "sessions", sessionsPage: String(page) });
  };

  const wordColumns: Column<Word>[] = [
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

  const formatDate = (dateString: string) => {
    try {
      return format(new Date(dateString), "MMM d, yyyy h:mm a");
    } catch (error) {
      return dateString;
    }
  };

  const sessionColumns: Column<StudySession>[] = [
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
        <div className="flex items-center gap-2">
          <Button asChild variant="link" className="p-0 h-auto font-medium">
            <Link to="/groups" className="flex items-center gap-1">
              <ArrowLeft className="h-4 w-4" />
              <span>Back to Groups</span>
            </Link>
          </Button>
        </div>

        {loading ? (
          <div className="h-16 bg-muted rounded animate-pulse" />
        ) : group ? (
          <div className="grid gap-6 md:grid-cols-3">
            <div className="md:col-span-2">
              <h1 className="text-3xl font-bold tracking-tight">{group.name}</h1>
              <p className="text-muted-foreground mt-1">
                A collection of thematically related words
              </p>
            </div>
            
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-lg flex items-center gap-2">
                  <BookOpen className="h-5 w-5 text-primary" />
                  <span>Group Statistics</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">Total Words</span>
                  <Badge variant="secondary" className="text-base">
                    {group.stats.total_word_count}
                  </Badge>
                </div>
              </CardContent>
            </Card>
          </div>
        ) : (
          <div className="text-center py-12">
            <p className="text-muted-foreground">Group not found</p>
          </div>
        )}

        <Tabs value={tab} onValueChange={handleTabChange} className="w-full">
          <TabsList className="grid w-full grid-cols-2 mb-6">
            <TabsTrigger value="words" className="flex items-center gap-2">
              <BookOpen className="h-4 w-4" />
              <span>Words</span>
            </TabsTrigger>
            <TabsTrigger value="sessions" className="flex items-center gap-2">
              <Calendar className="h-4 w-4" />
              <span>Study Sessions</span>
            </TabsTrigger>
          </TabsList>
          
          <TabsContent value="words" className="space-y-4">
            <DataTable
              columns={wordColumns}
              data={words?.items || []}
              keyExtractor={(word) => word.id}
              loading={loadingWords}
              onRowClick={(word) => window.location.href = `/words/${word.id}`}
              emptyMessage="No words in this group"
            />
            
            {words && (
              <div className="flex justify-center">
                <Pagination
                  pagination={words.pagination}
                  onPageChange={handleWordsPageChange}
                />
              </div>
            )}
          </TabsContent>
          
          <TabsContent value="sessions" className="space-y-4">
            <DataTable
              columns={sessionColumns}
              data={sessions?.items || []}
              keyExtractor={(session) => session.id}
              loading={loadingSessions}
              onRowClick={(session) => window.location.href = `/study_sessions/${session.id}`}
              emptyMessage="No study sessions for this group"
            />
            
            {sessions && (
              <div className="flex justify-center">
                <Pagination
                  pagination={sessions.pagination}
                  onPageChange={handleSessionsPageChange}
                />
              </div>
            )}
          </TabsContent>
        </Tabs>
      </div>
    </PageLayout>
  );
}
