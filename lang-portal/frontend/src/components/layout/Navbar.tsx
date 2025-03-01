
import { Link, useLocation } from "react-router-dom";
import { cn } from "@/lib/utils";
import { ThemeToggle } from "@/components/ui/theme-toggle";
import { useState, useEffect } from "react";

interface NavLink {
  path: string;
  label: string;
}

const navLinks: NavLink[] = [
  { path: "/dashboard", label: "Dashboard" },
  { path: "/study_activities", label: "Study Activities" },
  { path: "/words", label: "Words" },
  { path: "/groups", label: "Word Groups" },
  { path: "/study_sessions", label: "Sessions" },
  { path: "/settings", label: "Settings" },
];

export default function Navbar() {
  const location = useLocation();
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 10);
    };

    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  return (
    <header
      className={cn(
        "fixed top-0 left-0 right-0 z-50 transition-all duration-300",
        scrolled
          ? "py-2 bg-background/80 backdrop-blur-md border-b shadow-sm"
          : "py-4 bg-transparent"
      )}
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 flex justify-between items-center">
        <div className="flex items-center">
          <Link
            to="/dashboard"
            className="text-2xl font-semibold tracking-tight transition-colors duration-300 hover:opacity-80"
          >
            Language Learning Portal
          </Link>
        </div>

        <nav className="hidden md:flex items-center space-x-1">
          {navLinks.map((link) => (
            <Link
              key={link.path}
              to={link.path}
              className={cn(
                "nav-link",
                location.pathname === link.path && "active"
              )}
            >
              {link.label}
            </Link>
          ))}
        </nav>

        <ThemeToggle />
      </div>
    </header>
  );
}
