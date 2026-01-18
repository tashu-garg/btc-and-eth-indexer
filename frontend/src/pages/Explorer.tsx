import { useEffect, useState } from 'react';
import { getBlocks } from '../services/api';
import type { PaginatedBlocksResponse } from '../services/api';
import { Table, Pagination, AddressDisplay } from '../components/ui/DataDisplay';
import { formatDistanceToNow } from 'date-fns';
import { Link } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { BTCIcon, ETHIcon } from '../components/ui/Icons';

const Explorer = ({ chain }: { chain: 'btc' | 'eth' }) => {
  const [data, setData] = useState<PaginatedBlocksResponse | null>(null);
  const [prevLatestHash, setPrevLatestHash] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);

  const fetchData = async (p: number) => {
    if (!data) setLoading(true);
    try {
      const resp = await getBlocks(chain, p, 15);
      if (p === 1 && data && resp.blocks.length > 0) {
        if (resp.blocks[0].hash !== data.blocks[0].hash) {
          setPrevLatestHash(data.blocks[0].hash);
        }
      }
      setData(resp);
    } catch (err) {
      console.error('Failed to fetch blocks:', err);
    } finally {
      if (!data) setLoading(false);
    }
  };

  useEffect(() => {
    fetchData(page);
    let interval: any;
    if (page === 1) {
      interval = setInterval(() => fetchData(1), 5000);
    }
    return () => clearInterval(interval);
  }, [chain, page]);

  const Icon = chain === 'btc' ? BTCIcon : ETHIcon;

  return (
    <div className="space-y-8 max-w-7xl mx-auto">
      <header className="flex items-center justify-between">
        <div className="flex flex-col gap-2">
          <div className="flex items-center gap-3">
             <div className="w-10 h-10 bg-surface rounded-xl flex items-center justify-center border border-border">
                <Icon className="w-6 h-6 text-primary" />
             </div>
             <h1 className="text-3xl font-bold tracking-tight capitalize">{chain === 'btc' ? 'Bitcoin' : 'Ethereum'} Explorer</h1>
          </div>
          <p className="text-text-muted">Real-time block production and history for {chain.toUpperCase()}.</p>
        </div>
        <div className="flex items-center gap-4 text-xs font-bold text-text-muted bg-surface p-2 rounded-xl border border-border">
          <div className="flex items-center gap-1">
            <div className="w-1.5 h-1.5 rounded-full bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.5)]" /> Live
          </div>
          <span>Total: {data?.total?.toLocaleString() || '0'} Blocks</span>
        </div>
      </header>

      <div className="relative">
        <AnimatePresence mode="wait">
          {!data && loading ? (
            <div className="space-y-4">
              {[...Array(5)].map((_, i) => (
                <div key={i} className="h-16 w-full bg-surface animate-pulse rounded-xl" />
              ))}
            </div>
          ) : (
            <motion.div
              key={chain + page}
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
            >
              <Table headers={['Height', 'Hash', 'Transactions', 'Time Ago']}>
                {data?.blocks.map((block, idx) => {
                  const isNew = idx === 0 && prevLatestHash && block.hash !== prevLatestHash;
                  return (
                    <motion.tr 
                      key={block.hash} 
                      initial={isNew ? { backgroundColor: '#3b82f633' } : false}
                      animate={isNew ? { backgroundColor: 'transparent' } : false}
                      transition={{ duration: 2 }}
                      className="hover:bg-surface-hover/50 transition-colors group"
                    >
                      <td className="px-6 py-4">
                        <Link to={`/${chain}/block/${block.height}`} className="font-mono text-primary font-bold hover:underline">
                          {block.height.toLocaleString()}
                        </Link>
                      </td>
                      <td className="px-6 py-4">
                        <AddressDisplay address={block.hash} />
                      </td>
                      <td className="px-6 py-4">
                        <span className="bg-surface-hover px-2 py-1 rounded-lg text-xs font-bold ring-1 ring-border">
                          {block.txCount} txs
                        </span>
                      </td>
                      <td className="px-6 py-4 text-sm text-text-muted">
                        {formatDistanceToNow(new Date(block.timestamp * 1000))} ago
                      </td>
                    </motion.tr>
                  );
                })}
              </Table>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      <Pagination
        current={page}
        total={data?.total || 0}
        limit={15}
        onPageChange={setPage}
      />
    </div>
  );
};

export default Explorer;
