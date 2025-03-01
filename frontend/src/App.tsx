
import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { ThemeProvider } from "@/components/ui/theme-provider";
import Dashboard from "./pages/Dashboard";
import NotFound from "./pages/NotFound";
import Words from "./pages/Words";
import WordShow from "./pages/WordShow";
import Groups from "./pages/Groups";
import GroupShow from "./pages/GroupShow";
import StudyActivities from "./pages/StudyActivities";
import StudyActivityShow from "./pages/StudyActivityShow";
import StudyActivityLaunch from "./pages/StudyActivityLaunch";
import StudySessions from "./pages/StudySessions";
import Settings from "./pages/Settings";

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <ThemeProvider defaultTheme="system">
      <TooltipProvider>
        <Toaster />
        <Sonner />
        <BrowserRouter>
          <Routes>
            {/* Redirect root to dashboard */}
            <Route path="/" element={<Navigate to="/dashboard" replace />} />
            
            {/* Main routes */}
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/words" element={<Words />} />
            <Route path="/words/:id" element={<WordShow />} />
            <Route path="/groups" element={<Groups />} />
            <Route path="/groups/:id" element={<GroupShow />} />
            <Route path="/study_activities" element={<StudyActivities />} />
            <Route path="/study_activities/:id" element={<StudyActivityShow />} />
            <Route path="/study_activities/:id/launch" element={<StudyActivityLaunch />} />
            <Route path="/study_sessions" element={<StudySessions />} />
            <Route path="/settings" element={<Settings />} />
            
            {/* Catch-all route */}
            <Route path="*" element={<NotFound />} />
          </Routes>
        </BrowserRouter>
      </TooltipProvider>
    </ThemeProvider>
  </QueryClientProvider>
);

export default App;
