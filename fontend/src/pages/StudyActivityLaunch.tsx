
import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import PageLayout from "@/components/layout/PageLayout";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ArrowLeft, ArrowUpRight, CheckCircle } from "lucide-react";
import { getStudyActivity, getGroups, launchStudyActivity } from "@/lib/api";
import { Group, StudyActivity } from "@/types";
import { Link } from "react-router-dom";
import { toast } from "sonner";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";

export default function StudyActivityLaunch() {
  const { id } = useParams<{ id: string }>();
  const activityId = id ? parseInt(id) : 0;
  const navigate = useNavigate();
  
  const [activity, setActivity] = useState<StudyActivity | null>(null);
  const [groups, setGroups] = useState<Group[]>([]);
  const [selectedGroupId, setSelectedGroupId] = useState<string>("");
  const [loading, setLoading] = useState(true);
  const [launching, setLaunching] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      if (!activityId) return;
      
      try {
        setLoading(true);
        const [activityData, groupsData] = await Promise.all([
          getStudyActivity(activityId),
          getGroups()
        ]);
        
        setActivity(activityData);
        setGroups(groupsData.items || []);
        
        if (groupsData.items.length > 0) {
          setSelectedGroupId(String(groupsData.items[0].id));
        }
      } catch (error) {
        console.error("Failed to fetch data:", error);
        toast.error("Failed to load required data");
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [activityId]);

  const handleLaunch = async () => {
    if (!activity || !selectedGroupId) {
      toast.error("Please select a word group");
      return;
    }
    
    try {
      setLaunching(true);
      
      const result = await launchStudyActivity(
        activity.id,
        parseInt(selectedGroupId)
      );
      
      // Successful launch
      toast.success("Study session started successfully!");
      
      // Simulate launching in a new tab (in a real app, this would have a proper URL)
      window.open(
        `https://example.com/activities/${activity.id}?session=${result.id}&group=${result.group_id}`,
        "_blank"
      );
      
      // Navigate to the study session page
      navigate(`/study_sessions/${result.id}`);
    } catch (error) {
      console.error("Failed to launch study activity:", error);
      toast.error("Failed to launch the study activity");
    } finally {
      setLaunching(false);
    }
  };

  return (
    <PageLayout>
      <div className="max-w-2xl mx-auto">
        <div className="mb-4">
          <Button asChild variant="link" className="p-0 h-auto font-medium">
            <Link to={`/study_activities/${activityId}`} className="flex items-center gap-1">
              <ArrowLeft className="h-4 w-4" />
              <span>Back to Activity</span>
            </Link>
          </Button>
        </div>
        
        <Card className="animate-in">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <ArrowUpRight className="h-5 w-5 text-primary" />
              {loading ? (
                <Skeleton className="h-7 w-48" />
              ) : (
                <span>Launch {activity?.name}</span>
              )}
            </CardTitle>
            <CardDescription>
              {loading ? (
                <Skeleton className="h-4 w-full mt-2" />
              ) : (
                "Select a word group to begin your study session"
              )}
            </CardDescription>
          </CardHeader>
          
          <CardContent className="space-y-4">
            {loading ? (
              <Skeleton className="h-10 w-full" />
            ) : (
              <div className="space-y-2">
                <label className="text-sm font-medium">
                  Word Group
                </label>
                <Select
                  value={selectedGroupId}
                  onValueChange={setSelectedGroupId}
                  disabled={launching}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select a word group" />
                  </SelectTrigger>
                  <SelectContent>
                    {groups.map((group) => (
                      <SelectItem key={group.id} value={String(group.id)}>
                        {group.name} ({group.word_count} words)
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            )}
          </CardContent>
          
          <CardFooter>
            <Button
              onClick={handleLaunch}
              disabled={loading || launching || !selectedGroupId}
              className="w-full"
            >
              {launching ? (
                <span className="flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 animate-spin" />
                  Launching...
                </span>
              ) : (
                <span className="flex items-center gap-2">
                  Launch Now
                  <ArrowUpRight className="h-4 w-4" />
                </span>
              )}
            </Button>
          </CardFooter>
        </Card>
      </div>
    </PageLayout>
  );
}
