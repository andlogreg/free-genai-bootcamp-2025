
import { useEffect, useState } from "react";
import PageLayout from "@/components/layout/PageLayout";
import { getStudyActivities } from "@/lib/api";
import { StudyActivity } from "@/types";
import StudyActivityCard from "@/components/shared/StudyActivityCard";

export default function StudyActivities() {
  const [activities, setActivities] = useState<StudyActivity[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchActivities = async () => {
      try {
        setLoading(true);
        const data = await getStudyActivities();
        setActivities(data);
      } catch (error) {
        console.error("Failed to fetch study activities:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchActivities();
  }, []);

  return (
    <PageLayout>
      <div className="space-y-6">
        <section className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Study Activities</h1>
            <p className="text-muted-foreground mt-1">
              Choose an activity to practice your language skills
            </p>
          </div>
        </section>

        {loading ? (
          <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 3 }).map((_, index) => (
              <div
                key={index}
                className="bg-card rounded-lg border shadow-sm h-[280px] animate-pulse"
              />
            ))}
          </div>
        ) : (
          <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3">
            {activities.map((activity, index) => (
              <StudyActivityCard
                key={activity.id}
                activity={activity}
                className={`staggered-item staggered-delay-${index + 1}`}
              />
            ))}
          </div>
        )}
      </div>
    </PageLayout>
  );
}
