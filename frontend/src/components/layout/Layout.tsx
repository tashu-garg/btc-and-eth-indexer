import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Search, LayoutDashboard, Coins, Github } from 'lucide-react';
import { BTCIcon, ETHIcon } from '../ui/Icons';

export const Navbar = () => {
  const [query, setQuery] = useState('');
  const navigate = useNavigate();

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (query.trim()) {
      navigate(`/search?q=${encodeURIComponent(query.trim())}`);
      setQuery('');
    }
  };

  return (
    <nav className="h-16 border-b border-border bg-surface/80 backdrop-blur-md sticky top-0 z-50 flex items-center px-6 justify-between">
      <div className="flex items-center gap-4">
        <Link to="/" className="flex items-center gap-2 group">
          <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center shadow-lg shadow-primary/20 group-hover:scale-110 transition-transform">
            <Coins className="text-white w-5 h-5" />
          </div>
          <span className="font-bold text-xl tracking-tight hidden md:block">
            Block<span className="text-primary">Pulse</span>
          </span>
        </Link>
      </div>

      <div className="flex-1 max-w-xl px-8">
        <form onSubmit={handleSearch} className="relative group">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted w-4 h-4 group-focus-within:text-primary transition-colors" />
          <input
            type="text"
            placeholder="Search height, hash, or address..."
            className="w-full bg-background border border-border rounded-xl py-2 pl-10 pr-4 outline-none focus:border-primary/50 focus:ring-4 focus:ring-primary/10 transition-all text-sm"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </form>
      </div>

      <div className="flex items-center gap-4">
        <a href="https://github.com" target="_blank" rel="noreferrer" className="text-text-muted hover:text-white transition-colors">
          <Github className="w-5 h-5" />
        </a>
      </div>
    </nav>
  );
};

export const Sidebar = () => {
  return (
    <aside className="w-64 border-r border-border h-[calc(100vh-64px)] overflow-y-auto p-4 flex flex-col gap-2">
      <NavLink to="/" icon={<LayoutDashboard size={18} />} label="Dashboard" />
      <div className="mt-4 px-3 text-[10px] font-bold text-text-muted uppercase tracking-wider">Explorers</div>
      <NavLink to="/btc" icon={<BTCIcon size={18} />} label="Bitcoin" />
      <NavLink to="/eth" icon={<ETHIcon size={18} />} label="Ethereum" />
    </aside>
  );
};

const NavLink = ({ to, icon, label }: { to: string, icon: React.ReactNode, label: string }) => {
  return (
    <Link
      to={to}
      className="flex items-center gap-3 px-4 py-2.5 rounded-xl text-text-muted hover:text-white hover:bg-surface-hover transition-all group"
    >
      <span className="group-hover:text-primary transition-colors">{icon}</span>
      <span className="text-sm font-medium">{label}</span>
    </Link>
  );
};
