import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getBlock } from '../services/api';
import type { BlockDetails as IBlockDetails } from '../services/api';
import { Table, AddressDisplay } from '../components/ui/DataDisplay';
import { format } from 'date-fns';
import { ChevronLeft, Box, Clock, Hash, ListChecks, ArrowRightLeft } from 'lucide-react';
import { BTCIcon, ETHIcon } from '../components/ui/Icons';

const BlockDetails = ({ chain }: { chain: 'btc' | 'eth' }) => {
  const { height } = useParams();
  const [block, setBlock] = useState<IBlockDetails | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchBlock = async () => {
      if (!height) return;
      setLoading(true);
      try {
        const data = await getBlock(chain, height);
        setBlock(data);
      } catch (err) {
        console.error('Failed to fetch block:', err);
      } finally {
        setLoading(false);
      }
    };
    fetchBlock();
  }, [chain, height]);

  if (loading) return (
    <div className="flex flex-col gap-8 animate-pulse">
      <div className="h-6 w-32 bg-surface rounded" />
      <div className="h-10 w-64 bg-surface rounded" />
      <div className="grid grid-cols-4 gap-6">
        {[...Array(4)].map((_, i) => <div key={i} className="h-20 bg-surface rounded-xl" />)}
      </div>
    </div>
  );

  if (!block) return (
    <div className="flex flex-col items-center justify-center py-20 gap-4">
      <h2 className="text-2xl font-bold">Block Details Not Found</h2>
      <p className="text-text-muted">We couldn't retrieve information for this block height.</p>
      <Link to={`/${chain}`} className="text-primary font-bold hover:underline">Return to Explorer</Link>
    </div>
  );

  const Icon = chain === 'btc' ? BTCIcon : ETHIcon;

  return (
    <div className="space-y-8 max-w-7xl mx-auto fade-in">
      <Link to={`/${chain}`} className="flex items-center gap-2 text-text-muted hover:text-white transition-colors group">
        <ChevronLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
        <span className="text-sm font-medium">Back to Explorer</span>
      </Link>

      <header className="flex flex-col gap-2">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-surface rounded-xl flex items-center justify-center border border-border">
            <Icon className="w-6 h-6 text-primary" />
          </div>
          <h1 className="text-3xl font-bold tracking-tight">Block #{block.height.toLocaleString()}</h1>
        </div>
        <p className="text-text-muted font-mono text-xs break-all">{block.hash}</p>
      </header>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <DetailCard label="Timestamp" value={format(new Date(block.timestamp * 1000), 'MMM dd, yyyy HH:mm:ss')} icon={Clock} />
        <DetailCard label="Transactions" value={block.txCount.toString()} icon={ListChecks} />
        <DetailCard label="Chain" value={chain === 'btc' ? 'Bitcoin' : 'Ethereum'} icon={Hash} />
        <DetailCard label="Status" value="Confirmed" icon={Box} />
      </div>

      <div className="space-y-4">
        <div className="flex items-center gap-2 text-lg font-bold">
          <ArrowRightLeft className="w-5 h-5 text-primary" />
          <h2>Transactions</h2>
        </div>
        
        <Table headers={['Hash', 'From', 'To', 'Value']}>
          {block.transactions.map((tx) => (
            <tr key={tx.hash} className="hover:bg-surface-hover/50 transition-colors">
              <td className="px-6 py-4">
                <AddressDisplay address={tx.hash} />
              </td>
              <td className="px-6 py-4 font-mono text-sm text-text-muted">
                <AddressDisplay address={tx.from} />
              </td>
              <td className="px-6 py-4 font-mono text-sm text-text-muted">
                <AddressDisplay address={tx.to} />
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <span className="font-bold text-sm">{tx.value} {chain === 'btc' ? 'BTC' : 'ETH'}</span>
              </td>
            </tr>
          ))}
        </Table>
        
        {block.transactions.length === 0 && (
          <div className="glass p-12 rounded-2xl text-center text-text-muted italic">
            No transactions found in this block.
          </div>
        )}
      </div>
    </div>
  );
};

const DetailCard = ({ label, value, icon: Icon }: { label: string, value: string, icon: any }) => (
  <div className="glass p-4 rounded-xl border border-border flex flex-col gap-2">
    <div className="flex items-center gap-2 text-text-muted">
      <Icon className="w-3.5 h-3.5" />
      <span className="text-[10px] font-bold uppercase tracking-wider">{label}</span>
    </div>
    <span className="font-bold text-sm tracking-tight">{value}</span>
  </div>
);

export default BlockDetails;
