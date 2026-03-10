import { createLazyFileRoute, Link } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'
import { ArrowRight, Shield, Zap, Layout } from 'lucide-react'

export const Route = createLazyFileRoute('/')({
  component: Index,
})

function Index() {
  return (
    <div className="flex flex-col min-h-[calc(100vh-3.5rem)]">
      {/* Hero Section */}
      <section className="py-20 px-4 flex flex-col items-center text-center space-y-6">
        <div className="inline-flex items-center rounded-lg bg-muted px-3 py-1 text-sm font-medium">
          🚀 Welcome to the future of task management
        </div>
        <h1 className="text-4xl md:text-6xl font-extrabold tracking-tight max-w-3xl">
          Manage your tasks with <span className="text-primary">Speed</span> and <span className="text-primary">Elegance</span>
        </h1>
        <p className="text-xl text-muted-foreground max-w-2xl">
          The all-in-one SaaS solution for individuals and teams. Built with Go, React, and Clean Architecture.
        </p>
        <div className="flex gap-4">
          <Link to="/register">
            <Button size="lg" className="gap-2 group font-semibold">
              Get Started <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-1" />
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
      <section className="bg-muted/50 py-20 px-4">
        <div className="max-w-6xl mx-auto grid md:grid-cols-3 gap-8">
          <FeatureCard 
            icon={<Zap className="h-10 w-10 text-primary" />}
            title="Fast Performance"
            description="Blazing fast task management powered by Go and optimized React."
          />
          <FeatureCard 
            icon={<Shield className="h-10 w-10 text-primary" />}
            title="Secure by Design"
            description="JWT-based authentication and protected routes for your data safety."
          />
          <FeatureCard 
            icon={<Layout className="h-10 w-10 text-primary" />}
            title="Clean UI"
            description="Beautifully crafted with Tailwind CSS and shadcn/ui components."
          />
        </div>
      </section>

      <footer className="mt-auto py-8 border-t px-4 text-center text-sm text-muted-foreground">
        © 2026 SaaS Task Management System. All rights reserved.
      </footer>
    </div>
  )
}

function FeatureCard({ icon, title, description }: { icon: React.ReactNode, title: string, description: string }) {
  return (
    <div className="bg-background p-8 rounded-xl border space-y-4 hover:shadow-lg transition-shadow">
      <div>{icon}</div>
      <h3 className="text-xl font-bold">{title}</h3>
      <p className="text-muted-foreground">{description}</p>
    </div>
  )
}
