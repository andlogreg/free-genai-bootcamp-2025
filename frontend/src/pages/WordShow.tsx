
import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import PageLayout from "@/components/layout/PageLayout";
import { ArrowLeft, CheckCircle, XCircle } from "lucide-react";
import { getWord } from "@/lib/api";
import { WordDetail } from "@/types";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";

export default function WordShow() {
  const { id } = useParams<{ id: string }>();
  const wordId = id ? parseInt(id) : 0;
  
  const [word, setWord] = useState<WordDetail | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchWordDetails = async () => {
      if (!wordId) return;
      
      try {
        setLoading(true);
        const data = await getWord(wordId);
        setWord(data);
      } catch (error) {
        console.error("Failed to fetch word details:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchWordDetails();
  }, [wordId]);

  return (
    <PageLayout>
      <div className="space-y-6">
        <div className="flex items-center gap-2">
          <Button asChild variant="link" className="p-0 h-auto font-medium">
            <Link to="/words" className="flex items-center gap-1">
              <ArrowLeft className="h-4 w-4" />
              <span>Back to Words</span>
            </Link>
          </Button>
        </div>

        {loading ? (
          <div className="space-y-6">
            <Skeleton className="h-12 w-1/3" />
            <Skeleton className="h-6 w-1/4" />
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Skeleton className="h-[200px] rounded-xl" />
              <Skeleton className="h-[200px] rounded-xl" />
            </div>
          </div>
        ) : word ? (
          <>
            <div>
              <h1 className="text-4xl font-bold tracking-tight">{word.portuguese}</h1>
              <p className="text-xl text-muted-foreground mt-1">{word.english}</p>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle className="text-xl">Study Statistics</CardTitle>
                  <CardDescription>
                    Performance history for this word
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-2 gap-4">
                    <div className="bg-green-50 dark:bg-green-900/20 p-4 rounded-lg flex flex-col items-center">
                      <CheckCircle className="h-10 w-10 text-green-500 mb-2" />
                      <span className="text-2xl font-bold">{word.stats.correct_count}</span>
                      <span className="text-sm text-muted-foreground">Correct</span>
                    </div>
                    
                    <div className="bg-red-50 dark:bg-red-900/20 p-4 rounded-lg flex flex-col items-center">
                      <XCircle className="h-10 w-10 text-red-500 mb-2" />
                      <span className="text-2xl font-bold">{word.stats.wrong_count}</span>
                      <span className="text-sm text-muted-foreground">Wrong</span>
                    </div>
                    
                    <div className="col-span-2 pt-2">
                      <div className="text-sm font-medium mb-1">Success Rate</div>
                      <div className="h-2 bg-muted rounded-full overflow-hidden">
                        {word.stats.correct_count + word.stats.wrong_count > 0 ? (
                          <div 
                            className="h-full bg-green-500 transition-all duration-500"
                            style={{ 
                              width: `${Math.round(
                                (word.stats.correct_count / 
                                (word.stats.correct_count + word.stats.wrong_count)) * 100
                              )}%` 
                            }}
                          />
                        ) : (
                          <div className="h-full bg-muted-foreground/20" />
                        )}
                      </div>
                      <div className="text-xs text-muted-foreground mt-1">
                        {word.stats.correct_count + word.stats.wrong_count > 0 ? (
                          `${Math.round(
                            (word.stats.correct_count / 
                            (word.stats.correct_count + word.stats.wrong_count)) * 100
                          )}%`
                        ) : (
                          "No attempts yet"
                        )}
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="text-xl">Word Groups</CardTitle>
                  <CardDescription>
                    Groups containing this word
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  {word.groups.length > 0 ? (
                    <div className="flex flex-wrap gap-2">
                      {word.groups.map((group) => (
                        <Badge
                          key={group.id}
                          variant="secondary"
                          className="cursor-pointer hover:bg-secondary/80 transition-colors"
                          onClick={() => window.location.href = `/groups/${group.id}`}
                        >
                          {group.name}
                        </Badge>
                      ))}
                    </div>
                  ) : (
                    <p className="text-muted-foreground">
                      This word is not part of any group
                    </p>
                  )}
                </CardContent>
              </Card>
            </div>
          </>
        ) : (
          <div className="text-center py-12">
            <p className="text-muted-foreground">Word not found</p>
          </div>
        )}
      </div>
    </PageLayout>
  );
}
