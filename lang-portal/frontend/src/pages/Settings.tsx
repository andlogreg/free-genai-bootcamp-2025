
import { useState } from "react";
import PageLayout from "@/components/layout/PageLayout";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Settings as SettingsIcon, Trash2, RefreshCw, AlertTriangle } from "lucide-react";
import { toast } from "sonner";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useTheme } from "@/components/ui/theme-provider";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { resetHistory, resetFull } from "@/lib/api";

export default function Settings() {
  const { theme, setTheme } = useTheme();
  const [isResettingHistory, setIsResettingHistory] = useState(false);
  const [isResettingFull, setIsResettingFull] = useState(false);

  const handleResetHistory = async () => {
    try {
      setIsResettingHistory(true);
      await resetHistory();
      toast.success("Study history has been reset successfully");
    } catch (error) {
      toast.error("Failed to reset study history");
      console.error(error);
    } finally {
      setIsResettingHistory(false);
    }
  };

  const handleFullReset = async () => {
    try {
      setIsResettingFull(true);
      await resetFull();
      toast.success("System has been fully reset");
    } catch (error) {
      toast.error("Failed to perform full reset");
      console.error(error);
    } finally {
      setIsResettingFull(false);
    }
  };

  return (
    <PageLayout>
      <div className="space-y-6 max-w-3xl mx-auto">
        <section>
          <h1 className="text-3xl font-bold tracking-tight flex items-center gap-2">
            <SettingsIcon className="h-7 w-7 text-primary" />
            <span>Settings</span>
          </h1>
          <p className="text-muted-foreground mt-1">
            Configure your vocabulary learning portal
          </p>
        </section>

        <Card>
          <CardHeader>
            <CardTitle>Appearance</CardTitle>
            <CardDescription>
              Customize how the application looks
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <label htmlFor="theme" className="text-sm font-medium">
                Theme
              </label>
              <Select
                value={theme}
                onValueChange={(value) => setTheme(value as "light" | "dark" | "system")}
              >
                <SelectTrigger id="theme" className="w-full sm:w-[200px]">
                  <SelectValue placeholder="Select theme" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="light">Light</SelectItem>
                  <SelectItem value="dark">Dark</SelectItem>
                  <SelectItem value="system">System Default</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Reset Data</CardTitle>
            <CardDescription>
              Reset your study history or perform a full system reset
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <h3 className="text-sm font-medium">Reset Study History</h3>
              <p className="text-sm text-muted-foreground">
                This will delete all study sessions and word review items, but keep your words and groups intact.
              </p>
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button variant="outline" className="mt-2">
                    <RefreshCw className="h-4 w-4 mr-2" />
                    Reset History
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Reset Study History</AlertDialogTitle>
                    <AlertDialogDescription>
                      This will delete all your study sessions and progress. This action cannot be undone.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction
                      onClick={handleResetHistory}
                      disabled={isResettingHistory}
                    >
                      {isResettingHistory ? "Resetting..." : "Reset"}
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>

            <div className="space-y-2 pt-4">
              <h3 className="text-sm font-medium">Full Reset</h3>
              <p className="text-sm text-muted-foreground">
                This will drop all tables and re-create them with seed data. All your customizations will be lost.
              </p>
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button variant="destructive" className="mt-2">
                    <Trash2 className="h-4 w-4 mr-2" />
                    Full Reset
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle className="flex items-center gap-2">
                      <AlertTriangle className="h-5 w-5 text-destructive" />
                      Full System Reset
                    </AlertDialogTitle>
                    <AlertDialogDescription>
                      This will reset ALL data in the system, including words, groups, and study history.
                      Everything will be reverted to initial seed data. This action cannot be undone.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction
                      onClick={handleFullReset}
                      disabled={isResettingFull}
                      className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                    >
                      {isResettingFull ? "Resetting..." : "Yes, Reset Everything"}
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>
          </CardContent>
        </Card>
      </div>
    </PageLayout>
  );
}
