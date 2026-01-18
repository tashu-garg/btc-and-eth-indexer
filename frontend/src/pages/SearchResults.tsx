import { useEffect, useState } from 'react';
import { useSearchParams, useNavigate, Link } from 'react-router-dom';
import { search } from '../services/api';
import type { SearchResult as ISearchResult } from '../services/api';
import { AddressDisplay } from '../components/ui/DataDisplay';
import { Loader2, AlertCircle, ArrowRight, Box, ArrowRightLeft } from 'lucide-react';
import { BTCIcon, ETHIcon } from '../components/ui/Icons';

const SearchResults = () => {
  const [searchParams] = useSearchParams();
  const query = searchParams.get('q');
  const [results, setResults] = useState<ISearchResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const performSearch = async () => {
      if (!query) return;
      setLoading(true);
      setError(null);
      try {
        const resp = await search(query);
        if (resp.type === 'block') {
          // Auto-redirect if it's a block
          navigate(`/${resp.chain}/block/${resp.result.height}`);
        } else {
          setResults(resp);
        }
      } catch (err: any) {
        setError(err.response?.data?.error || 'Search failed. Please try again.');
      } finally {
        setLoading(false);
      }
    };
    performSearch();
  }, [query, navigate]);

  if (loading) return (
    <div className="flex flex-col items-center justify-center py-20 gap-4">
      <Loader2 className="w-8 h-8 text-primary animate-spin" />
      <p className="text-text-muted">Searching the blockchain...</p>
    </div>
  );

  if (error) return (
    <div className="flex flex-col items-center justify-center py-20 gap-4 text-center">
      <AlertCircle className="w-12 h-12 text-red-500/50" />
      <div className="space-y-1">
        <h2 className="text-xl font-bold">Search Error</h2>
        <p className="text-text-muted max-w-md">{error}</p>
      </div>
      <Link to="/" className="text-primary font-bold hover:underline">Back to Dashboard</Link>
    </div>
  );

  return (
    <div className="max-w-4xl mx-auto space-y-8">
      <header>
        <h1 className="text-3xl font-bold tracking-tight">Search Results</h1>
        <p className="text-text-muted">Showing results for: <span className="text-white font-mono">{query}</span></p>
      </header>

      {!results ? (
        <div className="glass p-12 rounded-2xl text-center space-y-4">
          <p className="text-text-muted">No exact matches found for this query.</p>
          <Link to="/" className="text-primary font-bold hover:underline inline-block">Return to Dashboard</Link>
        </div>
      ) : (
        <div className="space-y-6">
          <div className="glass p-6 rounded-2xl border border-border space-y-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 bg-surface rounded-xl flex items-center justify-center border border-border">
                  {results.chain === 'btc' ? <BTCIcon className="text-primary w-6 h-6" /> : <ETHIcon className="text-primary w-6 h-6" />}
                </div>
                <div>
                  <h3 className="font-bold capitalize">{results.type} Found</h3>
                  <p className="text-xs text-text-muted capitalize">{results.chain} Network</p>
                </div>
              </div>
              <span className="bg-green-500/10 text-green-500 text-[10px] font-bold px-2 py-1 rounded ring-1 ring-green-500/20 uppercase">
                Confirmed
              </span>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="bg-background/50 p-4 rounded-xl border border-border space-y-1">
                <div className="flex items-center gap-2 text-text-muted">
                  <Box className="w-3 h-3" />
                  <span className="text-[10px] font-bold uppercase">Block Height</span>
                </div>
                <p className="font-bold">{results.result.height.toLocaleString()}</p>
              </div>
              <div className="bg-background/50 p-4 rounded-xl border border-border space-y-1">
                <div className="flex items-center gap-2 text-text-muted">
                  <ArrowRightLeft className="w-3 h-3" />
                  <span className="text-[10px] font-bold uppercase">Value</span>
                </div>
                <p className="font-bold text-primary">{results.result.value} {results.chain === 'btc' ? 'BTC' : 'ETH'}</p>
              </div>
            </div>

            <div className="space-y-4 pt-4 border-t border-border">
               <div className="flex items-start gap-4">
                  <div className="flex-1 flex flex-col gap-1">
                    <span className="text-[10px] font-bold text-text-muted uppercase tracking-wider">From</span>
                    <AddressDisplay address={results.result.from} />
                  </div>
                  <div className="flex items-center justify-center pt-6">
                    <ArrowRight className="text-text-muted w-4 h-4 opacity-30" />
                  </div>
                  <div className="flex-1 flex flex-col gap-1">
                    <span className="text-[10px] font-bold text-text-muted uppercase tracking-wider">To</span>
                    <AddressDisplay address={results.result.to} />
                  </div>
               </div>
               <div className="space-y-1">
                  <span className="text-[10px] font-bold text-text-muted uppercase">Transaction Hash</span>
                  <p className="font-mono text-xs break-all">{results.result.hash}</p>
               </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default SearchResults;
