import { createLazyFileRoute, Link } from "@tanstack/react-router";
import { ArrowRight, Layout, Shield, Zap } from "lucide-react";
import { Button } from "@/components/ui/button";

export const Route = createLazyFileRoute("/")({
  component: Index,
});

function Index() {
  return (
    <div className="flex min-h-[calc(100vh-3.5rem)] flex-col">
      {/* Hero Section */}
      <section className="flex flex-col items-center space-y-6 px-4 py-20 text-center">
        <div className="bg-muted inline-flex items-center rounded-lg px-3 py-1 text-sm font-medium">
          🚀 Welcome to the future of task management
        </div>
        <h1 className="max-w-3xl text-4xl font-extrabold tracking-tight md:text-6xl">
          Manage your tasks with <span className="text-primary">Speed</span> and{" "}
          <span className="text-primary">Elegance</span>
        </h1>
        <p className="text-muted-foreground max-w-2xl text-xl">
          The all-in-one SaaS solution for individuals and teams. Built with Go,
          React, and Clean Architecture.
        </p>
        <div className="flex gap-4">
          <Link to="/register">
            <Button size="lg" className="group gap-2 font-semibold">
              Get Started{" "}
              <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-1" />
            </Button>
          </Link>
          <Link to="/login">
            <Button variant="outline" size="lg" className="font-semibold">
              Login
            </Button>
          </Link>
        </div>
      </section>

      {/* Features Section */}
      <section className="bg-muted/50 px-4 py-20">
        <div className="mx-auto grid max-w-6xl gap-8 md:grid-cols-3">
          <FeatureCard
            icon={<Zap className="text-primary h-10 w-10" />}
            title="Fast Performance"
            description="Blazing fast task management powered by Go and optimized React."
          />
          <FeatureCard
            icon={<Shield className="text-primary h-10 w-10" />}
            title="Secure by Design"
            description="JWT-based authentication and protected routes for your data safety."
          />
          <FeatureCard
            icon={<Layout className="text-primary h-10 w-10" />}
            title="Clean UI"
            description="Beautifully crafted with Tailwind CSS and shadcn/ui components."
          />
        </div>
      </section>

      <footer className="text-muted-foreground mt-auto border-t px-4 py-8 text-center text-sm">
        © 2026 SaaS Task Management System. All rights reserved.
      </footer>
    </div>
  );
}

function FeatureCard({
  icon,
  title,
  description,
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
}) {
  return (
    <div className="bg-background space-y-4 rounded-xl border p-8 transition-shadow hover:shadow-lg">
      <div>{icon}</div>
      <h3 className="text-xl font-bold">{title}</h3>
      <p className="text-muted-foreground">{description}</p>
    </div>
  );
}
