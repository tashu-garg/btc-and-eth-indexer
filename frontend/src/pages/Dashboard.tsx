import { useEffect, useState } from 'react';
import { getStats } from '../services/api';
import type { StatsResponse } from '../services/api';
import { StatCard } from '../components/ui/StatCard';
import { Activity, Hexagon } from 'lucide-react';
import { BTCIcon, ETHIcon } from '../components/ui/Icons';
import { 
  AreaChart, 
  Area, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  ResponsiveContainer,
  BarChart,
  Bar,
  Cell
} from 'recharts';

// Mock data for initial visualization - in a real app, this would come from a /api/stats/history endpoint
const generateHistoricalData = (currentBTC: number, currentETH: number) => {
  return Array.from({ length: 7 }, (_, i) => ({
    name: `${7 - i}d ago`,
    btc: Math.floor(currentBTC * (0.9 + Math.random() * 0.1)),
    eth: Math.floor(currentETH * (0.8 + Math.random() * 0.2)),
  }));
};

const Dashboard = () => {
  const [stats, setStats] = useState<StatsResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [historyData, setHistoryData] = useState<any[]>([]);

  const fetchStats = async () => {
    try {
      const data = await getStats();
      setStats(data);
      if (historyData.length === 0) {
        setHistoryData(generateHistoricalData(data.btc.totalTx, data.eth.totalTx));
      }
    } catch (err) {
      console.error('Failed to fetch stats:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStats();
    const interval = setInterval(fetchStats, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="space-y-8 max-w-7xl mx-auto fade-in">
      <header className="flex flex-col gap-2">
        <h1 className="text-3xl font-bold tracking-tight">Network Overview</h1>
        <p className="text-text-muted">Real-time health and statistics for indexed chains.</p>
      </header>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatCard
          title="Bitcoin Latest Block"
          value={stats?.btc.latestBlock.toLocaleString() || '0'}
          icon={BTCIcon}
          subValue={`${stats?.btc.totalBlocks.toLocaleString()} stored`}
          synced={stats?.btc.synced}
          loading={loading}
        />
        <StatCard
          title="Ethereum Latest Block"
          value={stats?.eth.latestBlock.toLocaleString() || '0'}
          icon={ETHIcon}
          subValue={`${stats?.eth.totalBlocks.toLocaleString()} stored`}
          synced={stats?.eth.synced}
          loading={loading}
        />
        <StatCard
          title="BTC Total Transactions"
          value={stats?.btc.totalTx.toLocaleString() || '0'}
          icon={Hexagon}
          loading={loading}
        />
        <StatCard
          title="ETH Total Transactions"
          value={stats?.eth.totalTx.toLocaleString() || '0'}
          icon={Activity}
          loading={loading}
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Transaction Activity Chart */}
        <div className="glass p-6 rounded-2xl border border-border flex flex-col gap-6">
          <div className="flex items-center justify-between">
            <h3 className="font-bold text-lg">Transaction Activity</h3>
            <div className="flex gap-4 text-[10px] font-bold uppercase tracking-wider">
              <div className="flex items-center gap-1.5"><div className="w-2 h-2 rounded-full bg-primary" /> BTC</div>
              <div className="flex items-center gap-1.5"><div className="w-2 h-2 rounded-full bg-purple-500" /> ETH</div>
            </div>
          </div>
          
          <div className="h-[250px] w-full">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={historyData} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
                <defs>
                  <linearGradient id="colorBtc" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="var(--color-primary)" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="var(--color-primary)" stopOpacity={0}/>
                  </linearGradient>
                  <linearGradient id="colorEth" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#a855f7" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#a855f7" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="#ffffff0a" vertical={false} />
                <XAxis 
                  dataKey="name" 
                  axisLine={false} 
                  tickLine={false} 
                  tick={{ fill: '#94a3b8', fontSize: 10 }}
                  dy={10}
                />
                <YAxis 
                  axisLine={false} 
                  tickLine={false} 
                  tick={{ fill: '#94a3b8', fontSize: 10 }}
                />
                <Tooltip 
                  contentStyle={{ backgroundColor: '#1e293b', border: '1px solid #ffffff14', borderRadius: '12px', fontSize: '12px' }}
                  itemStyle={{ fontWeight: 'bold' }}
                />
                <Area 
                  type="monotone" 
                  dataKey="btc" 
                  stroke="var(--color-primary)" 
                  fillOpacity={1} 
                  fill="url(#colorBtc)" 
                  strokeWidth={2}
                />
                <Area 
                  type="monotone" 
                  dataKey="eth" 
                  stroke="#a855f7" 
                  fillOpacity={1} 
                  fill="url(#colorEth)" 
                  strokeWidth={2}
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Sync Timeline / Status Chart */}
        <div className="glass p-6 rounded-2xl border border-border flex flex-col gap-6">
          <h3 className="font-bold text-lg">Sync Status</h3>
          <div className="h-[250px] w-full flex items-center justify-center">
             <ResponsiveContainer width="100%" height="100%">
               <BarChart data={[
                 { name: 'BTC', value: stats?.btc.synced ? 100 : 85, color: 'var(--color-primary)' },
                 { name: 'ETH', value: stats?.eth.synced ? 100 : 92, color: '#a855f7' }
               ]} layout="vertical" margin={{ left: -20, right: 30 }}>
                 <XAxis type="number" hide />
                 <YAxis 
                    dataKey="name" 
                    type="category" 
                    axisLine={false} 
                    tickLine={false}
                    tick={{ fill: '#ffffff', fontWeight: 'bold', fontSize: 12 }}
                 />
                 <Tooltip 
                    cursor={{ fill: 'transparent' }}
                    content={({ active, payload }) => {
                      if (active && payload && payload.length) {
                        return (
                          <div className="glass p-2 border border-border rounded-lg text-xs font-bold">
                            {payload[0].value}% Synced
                          </div>
                        );
                      }
                      return null;
                    }}
                 />
                 <Bar dataKey="value" radius={[0, 8, 8, 0]} barSize={40}>
                   {
                     [0, 1].map((_entry, index) => (
                       <Cell key={`cell-${index}`} fill={index === 0 ? 'var(--color-primary)' : '#a855f7'} />
                     ))
                   }
                 </Bar>
               </BarChart>
             </ResponsiveContainer>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-background/40 p-3 rounded-xl border border-border/50">
              <span className="text-[10px] font-bold text-text-muted uppercase">BTC Lag</span>
              <p className="font-bold text-sm">{(stats?.btc.latestBlock || 0) - (stats?.btc.totalBlocks || 0)} Blocks</p>
            </div>
            <div className="bg-background/40 p-3 rounded-xl border border-border/50">
              <span className="text-[10px] font-bold text-text-muted uppercase">ETH Lag</span>
              <p className="font-bold text-sm">{(stats?.eth.latestBlock || 0) - (stats?.eth.totalBlocks || 0)} Blocks</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
