
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { Link } from "react-router-dom";
import { StudyActivity } from "@/types";
import { ArrowUpRight, Info } from "lucide-react";
import { AspectRatio } from "@/components/ui/aspect-ratio";

interface StudyActivityCardProps {
  activity: StudyActivity;
  className?: string;
}

export default function StudyActivityCard({
  activity,
  className,
}: StudyActivityCardProps) {
  const { id, name, thumbnail_url, description } = activity;

  return (
    <Card className={cn("overflow-hidden transition-all card-hover", className)}>
      <AspectRatio ratio={16 / 9}>
        <img
          src={thumbnail_url}
          alt={name}
          className="object-cover w-full h-full rounded-t-lg"
        />
      </AspectRatio>
      
      <CardContent className="p-4">
        <h3 className="text-lg font-semibold mb-2">{name}</h3>
        <p className="text-sm text-muted-foreground line-clamp-2">{description}</p>
      </CardContent>
      
      <CardFooter className="px-4 pb-4 pt-0 flex justify-between">
        <Button variant="outline" size="sm" asChild>
          <Link to={`/study_activities/${id}`} className="flex items-center gap-1">
            <Info className="w-4 h-4" />
            <span>Details</span>
          </Link>
        </Button>
        
        <Button asChild>
          <Link to={`/study_activities/${id}/launch`} className="flex items-center gap-1">
            <span>Launch</span>
            <ArrowUpRight className="w-4 h-4" />
          </Link>
        </Button>
      </CardFooter>
    </Card>
  );
}
