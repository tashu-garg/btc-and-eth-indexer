import { Copy } from 'lucide-react';

export const CopyButton = ({ text }: { text: string }) => {
  const copy = () => {
    navigator.clipboard.writeText(text);
  };

  return (
    <button
      onClick={copy}
      className="p-1 hover:bg-surface-hover rounded transition-colors text-text-muted hover:text-white"
    >
      <Copy className="w-3.5 h-3.5" />
    </button>
  );
};

export const AddressDisplay = ({ address }: { address: string }) => {
  if (address === 'coinbase') {
    return (
      <span className="bg-orange-500/10 text-orange-500 text-[10px] font-bold px-2 py-0.5 rounded ring-1 ring-orange-500/20 uppercase">
        Mining / Coinbase
      </span>
    );
  }

  if (address === 'unknown') {
    return <span className="text-text-muted italic text-xs">Unknown / Encrypted</span>;
  }

  if (address === 'non-standard') {
    return <span className="text-text-muted italic text-xs">Non-standard Script</span>;
  }

  // Handle +X others
  if (address.includes(',+')) {
    const [main, count] = address.split(',+');
    return (
      <div className="flex items-center gap-2">
        <span className="font-mono text-sm text-text-muted truncate max-w-[120px]">
          {main.substring(0, 6)}...{main.substring(main.length - 4)}
        </span>
        <span className="bg-primary/10 text-primary text-[9px] font-bold px-1.5 py-0.5 rounded ring-1 ring-primary/20">
          +{count}
        </span>
      </div>
    );
  }

  return (
    <div className="flex items-center gap-2">
      <span className="font-mono text-sm text-text-muted truncate max-w-[160px]">
        {address.length > 16 ? `${address.substring(0, 8)}...${address.substring(address.length - 6)}` : address}
      </span>
      <CopyButton text={address} />
    </div>
  );
};

export const Table = ({ headers, children }: { headers: string[], children: React.ReactNode }) => (
  <div className="glass rounded-2xl overflow-hidden border border-border">
    <table className="w-full text-left border-collapse">
      <thead>
        <tr className="bg-surface/50 border-b border-border">
          {headers.map((h) => (
            <th key={h} className="px-6 py-4 text-xs font-bold text-text-muted uppercase tracking-wider">
              {h}
            </th>
          ))}
        </tr>
      </thead>
      <tbody className="divide-y divide-border">
        {children}
      </tbody>
    </table>
  </div>
);

export const Pagination = ({ current, total, limit, onPageChange }: { 
  current: number, 
  total: number, 
  limit: number, 
  onPageChange: (p: number) => void 
}) => {
  const totalPages = Math.ceil(total / limit);
  
  return (
    <div className="flex items-center justify-between mt-6 px-2">
      <span className="text-sm text-text-muted">
        Showing Page {current} of {totalPages || 1}
      </span>
      <div className="flex gap-2">
        <button
          disabled={current === 1}
          onClick={() => onPageChange(current - 1)}
          className="px-4 py-2 text-sm font-medium border border-border rounded-xl disabled:opacity-30 hover:bg-surface-hover transition-colors"
        >
          Previous
        </button>
        <button
          disabled={current >= totalPages}
          onClick={() => onPageChange(current + 1)}
          className="px-4 py-2 text-sm font-medium border border-border rounded-xl disabled:opacity-30 hover:bg-surface-hover transition-colors bg-surface-hover"
        >
          Next
        </button>
      </div>
    </div>
  );
};
