import { motion } from 'framer-motion';
import type { LucideProps } from 'lucide-react';
import React from 'react';

interface StatCardProps {
  title: string;
  value: string | number;
  icon: React.ComponentType<LucideProps>;
  subValue?: string;
  synced?: boolean;
  loading?: boolean;
}

export const StatCard = ({ title, value, icon: Icon, subValue, synced, loading }: StatCardProps) => {
  if (loading) return <SkeletonCard />;

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="glass p-6 rounded-2xl flex flex-col gap-4 relative overflow-hidden group hover:border-primary/30 transition-colors"
    >
      <div className="flex items-center justify-between">
        <span className="text-text-muted text-sm font-medium">{title}</span>
        <div className="w-10 h-10 rounded-xl bg-surface-hover flex items-center justify-center border border-border group-hover:bg-primary/10 group-hover:border-primary/20 transition-all">
          <Icon className="w-5 h-5 text-text-muted group-hover:text-primary transition-colors" />
        </div>
      </div>

      <div className="flex flex-col">
        <span className="text-3xl font-bold tracking-tight">{value}</span>
        {subValue && <span className="text-text-muted text-xs mt-1">{subValue}</span>}
      </div>

      {synced !== undefined && (
        <div className="flex items-center gap-1.5 mt-2">
          <div className={`w-2 h-2 rounded-full ${synced ? 'bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.5)]' : 'bg-yellow-500 shadow-[0_0_8px_rgba(234,179,8,0.5)]'}`} />
          <span className="text-[10px] font-bold uppercase tracking-wider text-text-muted">
            {synced ? 'Synced' : 'Syncing...'}
          </span>
        </div>
      )}

      {/* Background Glow */}
      <div className="absolute -right-4 -bottom-4 w-24 h-24 bg-primary/5 blur-3xl rounded-full" />
    </motion.div>
  );
};

const SkeletonCard = () => (
  <div className="glass p-6 rounded-2xl flex flex-col gap-4 animate-pulse">
    <div className="flex items-center justify-between">
      <div className="h-4 w-24 bg-border rounded" />
      <div className="w-10 h-10 rounded-xl bg-border" />
    </div>
    <div className="h-8 w-32 bg-border rounded" />
    <div className="h-3 w-16 bg-border rounded" />
  </div>
);

export const Skeleton = ({ className }: { className?: string }) => (
  <div className={`animate-pulse bg-border rounded ${className}`} />
);
