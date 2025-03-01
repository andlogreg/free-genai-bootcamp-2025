
import { useEffect, useState } from "react";
import PageLayout from "@/components/layout/PageLayout";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";
import { ArrowRight, Award, BookOpen, Calendar, CheckCircle, Clock, XCircle } from "lucide-react";
import { LastStudySession, QuickStats, StudyProgress } from "@/types";
import { getLastStudySession, getQuickStats, getStudyProgress } from "@/lib/api";
import StatCard from "@/components/shared/StatCard";
import { Progress } from "@/components/ui/progress";
import { formatDistanceToNow } from "date-fns";

export default function Dashboard() {
  const [lastSession, setLastSession] = useState<LastStudySession | null>(null);
  const [progress, setProgress] = useState<StudyProgress | null>(null);
  const [stats, setStats] = useState<QuickStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        setLoading(true);
        const [sessionData, progressData, statsData] = await Promise.all([
          getLastStudySession(),
          getStudyProgress(),
          getQuickStats()
        ]);
        
        setLastSession(sessionData);
        setProgress(progressData);
        setStats(statsData);
      } catch (error) {
        console.error("Failed to fetch dashboard data:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchDashboardData();
  }, []);

  const formatRelativeTime = (dateString: string) => {
    try {
      return formatDistanceToNow(new Date(dateString), { addSuffix: true });
    } catch (error) {
      return "recently";
    }
  };

  // Calculate progress percentage
  const progressPercentage = progress 
    ? Math.round((progress.total_words_studied / progress.total_available_words) * 100) 
    : 0;

  // Calculate mastery percentage (just an example formula)
  const masteryPercentage = stats 
    ? Math.round(stats.success_rate / 2) 
    : 0;

  return (
    <PageLayout>
      <div className="grid gap-6">
        <section className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
            <p className="text-muted-foreground mt-1">
              Welcome to your language learning portal
            </p>
          </div>
          
          <Button asChild size="lg" className="gap-2">
            <Link to="/study_activities">
              <span>Start Studying</span>
              <ArrowRight className="h-4 w-4" />
            </Link>
          </Button>
        </section>

        {lastSession && (
          <Card className="staggered-item staggered-delay-1">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Clock className="h-5 w-5 text-primary" />
                <span>Last Study Session</span>
              </CardTitle>
              <CardDescription>
                {formatRelativeTime(lastSession.created_at)}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <p className="text-sm font-medium text-muted-foreground mb-1">
                    Group
                  </p>
                  <p className="font-medium">{lastSession.group_name}</p>
                </div>
                <div className="flex items-center gap-6">
                  <div className="flex flex-col items-center">
                    <CheckCircle className="h-8 w-8 text-green-500 mb-1" />
                    <span className="text-sm font-medium">8 Correct</span>
                  </div>
                  <div className="flex flex-col items-center">
                    <XCircle className="h-8 w-8 text-red-500 mb-1" />
                    <span className="text-sm font-medium">2 Wrong</span>
                  </div>
                </div>
              </div>
              <Button variant="outline" size="sm" className="mt-4" asChild>
                <Link to={`/groups/${lastSession.group_id}`}>
                  View Group Details
                </Link>
              </Button>
            </CardContent>
          </Card>
        )}

        <div className="grid gap-6 sm:grid-cols-2">
          <Card className="staggered-item staggered-delay-2">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <BookOpen className="h-5 w-5 text-primary" />
                <span>Study Progress</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {loading ? (
                  <div className="h-[120px] flex items-center justify-center">
                    <div className="animate-pulse h-4 w-1/2 bg-muted rounded" />
                  </div>
                ) : progress ? (
                  <>
                    <div>
                      <div className="flex justify-between mb-1">
                        <p className="text-sm font-medium">
                          Words Studied
                        </p>
                        <p className="text-sm font-medium">
                          {progress.total_words_studied}/{progress.total_available_words}
                        </p>
                      </div>
                      <Progress value={progressPercentage} className="h-2" />
                      <p className="text-xs text-muted-foreground mt-1">
                        {progressPercentage}% of all available words
                      </p>
                    </div>
                    
                    <div>
                      <div className="flex justify-between mb-1">
                        <p className="text-sm font-medium">
                          Mastery
                        </p>
                        <p className="text-sm font-medium">
                          {masteryPercentage}%
                        </p>
                      </div>
                      <Progress value={masteryPercentage} className="h-2" />
                      <p className="text-xs text-muted-foreground mt-1">
                        Based on success rate across all sessions
                      </p>
                    </div>
                  </>
                ) : (
                  <p className="text-center text-muted-foreground">
                    No progress data available
                  </p>
                )}
              </div>
            </CardContent>
          </Card>

          <div className="grid gap-3 staggered-item staggered-delay-3">
            <h3 className="text-lg font-medium">Quick Stats</h3>
            
            <div className="grid grid-cols-2 gap-3">
              {loading ? (
                Array.from({ length: 4 }).map((_, i) => (
                  <div key={i} className="bg-card rounded-lg border shadow-sm p-4 animate-pulse">
                    <div className="h-4 bg-muted rounded w-2/3 mb-2" />
                    <div className="h-6 bg-muted rounded w-1/3" />
                  </div>
                ))
              ) : stats ? (
                <>
                  <StatCard
                    title="Success Rate"
                    value={`${stats.success_rate}%`}
                    icon={<Award className="h-5 w-5 text-primary" />}
                  />
                  
                  <StatCard
                    title="Study Sessions"
                    value={stats.total_study_sessions}
                    icon={<Calendar className="h-5 w-5 text-primary" />}
                  />
                  
                  <StatCard
                    title="Active Groups"
                    value={stats.total_active_groups}
                    icon={<BookOpen className="h-5 w-5 text-primary" />}
                  />
                  
                  <StatCard
                    title="Study Streak"
                    value={`${stats.study_streak_days} days`}
                    icon={<Award className="h-5 w-5 text-primary" />}
                  />
                </>
              ) : (
                <p className="col-span-2 text-center text-muted-foreground">
                  No stats available
                </p>
              )}
            </div>
          </div>
        </div>
      </div>
    </PageLayout>
  );
}
