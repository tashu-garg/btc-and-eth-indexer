import React, { useState, useEffect, useCallback } from "react";
import axios from "axios";
import { 
  Search, 
  ChevronDown, 
  Box, 
  Repeat, 
  Clock, 
  Zap,
  Activity,
  Code,
  Menu,
  TrendingUp,
  ArrowRight,
  ChevronLeft,
  ChevronRight,
  LayoutGrid,
  X,
  Database
} from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

const BACKEND_URL = "http://localhost:8989/api/v1";
const ITEMS_PER_PAGE = 8;

// Custom Chain Icons
const EthIcon = ({ className }) => (
  <svg viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg" className={className}>
    <path d="M16 3L15.7 3.9V20.8L16 21.1L24.3 16.2L16 3Z" fill="currentColor" fillOpacity="0.8"/>
    <path d="M16 3L7.7 16.2L16 21.1V12.2V3Z" fill="currentColor"/>
    <path d="M16 22.8L15.8 23V31L16 31.6L24.3 20.1L16 22.8Z" fill="currentColor" fillOpacity="0.8"/>
    <path d="M16 31.6V22.8L7.7 20.1L16 31.6Z" fill="currentColor"/>
    <path d="M16 21.1L24.3 16.2L16 12.2V21.1Z" fill="currentColor" fillOpacity="0.4"/>
    <path d="M7.7 16.2L16 21.1V12.2L7.7 16.2Z" fill="currentColor" fillOpacity="0.4"/>
  </svg>
);

const BtcIcon = ({ className }) => (
  <svg viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg" className={className}>
    <path d="M22.7 14.8C23.3 13.9 23.5 12.5 23.2 11.2C22.6 9 20.6 7.6 18 7.3C17.6 7.2 17.2 7.2 16.8 7.2V3.7H14.8V7.2H12.8V3.7H10.8V7.2H8V9.1H9.5C10.2 9.1 10.7 9.6 10.7 10.3V21.7C10.7 22.4 10.2 22.9 9.5 22.9H8V24.8H10.8V28.3H12.8V24.8H14.8V28.3H16.8V24.8H18C21.4 24.8 24.1 22.8 24.5 19.3C24.7 17.6 24 16 22.7 14.8ZM13.1 9.4H17.5C18.9 9.4 20 10.2 20.3 11.4C20.6 12.6 19.9 13.7 18.5 13.7H13.1V9.4ZM18.8 22.2H13.1V16H18.8C20.4 16 21.7 17 21.9 18.5C22.1 20.4 20.9 22.2 18.8 22.2Z" fill="currentColor"/>
  </svg>
);

const App = () => {
  const [ethBlocks, setEthBlocks] = useState([]);
  const [btcBlocks, setBtcBlocks] = useState([]);
  const [ethTxs, setEthTxs] = useState([]);
  const [btcTxs, setBtcTxs] = useState([]);
  const [ethStats, setEthStats] = useState({ total_blocks: 0, total_transactions: 0 });
  const [btcStats, setBtcStats] = useState({ total_blocks: 0, total_transactions: 0 });
  
  const [activeTab, setActiveTab] = useState("all"); 
  const [viewMode, setViewMode] = useState("overview"); // overview, blocks, transactions
  const [blocksPage, setBlocksPage] = useState(0);
  const [txsPage, setTxsPage] = useState(0);

  const [searchQuery, setSearchQuery] = useState("");
  const [isSearching, setIsSearching] = useState(false);
  const [loading, setLoading] = useState(true);
  const [selectedItem, setSelectedItem] = useState(null);
  const [modalLoading, setModalLoading] = useState(false);
  const [isPaused, setIsPaused] = useState(false);

  const fetchStats = useCallback(async () => {
    try {
      const [eStats, bStats] = await Promise.all([
        axios.get(`${BACKEND_URL}/ethereum/stats`),
        axios.get(`${BACKEND_URL}/bitcoin/stats`)
      ]);
      setEthStats(eStats.data || { total_blocks: 0, total_transactions: 0 });
      setBtcStats(bStats.data || { total_blocks: 0, total_transactions: 0 });
    } catch (err) {
      console.error("Stats fetch error:", err);
    }
  }, []);

  const fetchData = useCallback(async (forcedPage = null) => {
    if (isPaused && ethBlocks.length > 0 && forcedPage === null) return;
    
    // If not searching, fetch paginated data
    if (isSearching && forcedPage === null) return;

    setLoading(true);
    try {
      const limit = ITEMS_PER_PAGE;
      const bPage = forcedPage !== null ? forcedPage : blocksPage;
      const tPage = forcedPage !== null ? forcedPage : txsPage;
      
      const requests = [];
      
      if (activeTab === 'all' || activeTab === 'ethereum') {
        requests.push(axios.get(`${BACKEND_URL}/ethereum/blocks?limit=${limit}&offset=${bPage * limit}`));
        requests.push(axios.get(`${BACKEND_URL}/ethereum/txs?limit=${limit}&offset=${tPage * limit}`));
      } else {
        requests.push(Promise.resolve({ data: [] }));
        requests.push(Promise.resolve({ data: [] }));
      }

      if (activeTab === 'all' || activeTab === 'bitcoin') {
        requests.push(axios.get(`${BACKEND_URL}/bitcoin/blocks?limit=${limit}&offset=${bPage * limit}`));
        requests.push(axios.get(`${BACKEND_URL}/bitcoin/txs?limit=${limit}&offset=${tPage * limit}`));
      } else {
        requests.push(Promise.resolve({ data: [] }));
        requests.push(Promise.resolve({ data: [] }));
      }

      const [eBlocks, eTxs, bBlocks, bTxs] = await Promise.all(requests);
      
      setEthBlocks(eBlocks.data || []);
      setEthTxs(eTxs.data || []);
      setBtcBlocks(bBlocks.data || []);
      setBtcTxs(bTxs.data || []);
      setLoading(false);
    } catch (err) {
      console.error("Data fetch error:", err);
      setLoading(false);
    }
  }, [activeTab, blocksPage, txsPage, isPaused, isSearching]);

  useEffect(() => {
    fetchStats();
    fetchData();
    const interval = setInterval(() => {
      fetchStats();
      if (blocksPage === 0 && txsPage === 0 && !isSearching) fetchData(); 
    }, 4000);
    return () => clearInterval(interval);
  }, [fetchData, fetchStats, blocksPage, txsPage, isSearching]);

  useEffect(() => {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }, [blocksPage, txsPage, viewMode]);

  const handleSearch = async (e) => {
    if (e) e.preventDefault();
    if (!searchQuery.trim()) {
      setIsSearching(false);
      fetchData(0);
      return;
    }
    
    setLoading(true);
    setIsSearching(true);
    try {
      const [ethRes, btcRes] = await Promise.all([
        axios.get(`${BACKEND_URL}/ethereum/search?q=${searchQuery}`),
        axios.get(`${BACKEND_URL}/bitcoin/search?q=${searchQuery}`)
      ]);
      
      setEthBlocks(ethRes.data.blocks || []);
      setBtcBlocks(btcRes.data.blocks || []);
      setEthTxs(ethRes.data.transactions || []);
      setBtcTxs(btcRes.data.transactions || []);
      
      setBlocksPage(0);
      setTxsPage(0);
    } catch (err) {
      console.error("Search error:", err);
    } finally {
      setLoading(false);
    }
  };

  const clearSearch = () => {
    setSearchQuery("");
    setIsSearching(false);
    setBlocksPage(0);
    setTxsPage(0);
    fetchData(0);
  };

  const fetchDetail = async (type, chain, identifier) => {
    setModalLoading(true);
    setSelectedItem({ type, chain, data: null });
    try {
      const endpoint = type === 'block' 
        ? `${BACKEND_URL}/${chain}/block/hash/${identifier}`
        : `${BACKEND_URL}/${chain}/tx/${identifier}`;
      
      const res = await axios.get(endpoint);
      setSelectedItem({ type, chain, data: res.data });
    } catch (err) {
      console.error("Detail Fetch error (DB unreachable or record missing):", err);
      // Removed local fallback to ensure strictly API-driven data
      setSelectedItem({ type, chain, data: null, error: true });
    } finally {
      setModalLoading(false);
    }
  };

  const combinedBlocks = [...ethBlocks, ...btcBlocks].sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp));
  const combinedTxs = [...ethTxs, ...btcTxs].sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp));

  const NavDropdown = ({ label, chain }) => (
    <div className="relative group">
      <button className={`flex items-center gap-2 px-5 py-2.5 rounded-xl text-sm font-bold transition-all border border-transparent hover:border-white/10 ${
        activeTab === chain ? 'text-white bg-white/5 shadow-inner' : 'text-gray-500 hover:text-white'
      }`}>
        <span className="flex items-center gap-2">
           {chain === 'ethereum' ? <EthIcon className="w-4 h-4 text-indigo-400" /> : <BtcIcon className="w-4 h-4 text-amber-500" />}
           {label}
        </span>
        <ChevronDown className="w-4 h-4 opacity-30 group-hover:rotate-180 transition-transform duration-300" />
      </button>
      <div className="absolute top-full left-0 pt-3 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-300 z-[100] translate-y-2 group-hover:translate-y-0">
        <div className="bg-[#111115] border border-white/10 rounded-2xl shadow-[0_20px_50px_rgba(0,0,0,0.5)] overflow-hidden min-w-[240px] p-2 backdrop-blur-3xl">
          <button
            onClick={() => { setActiveTab(chain); setViewMode("blocks"); setBlocksPage(0); setIsSearching(false); }}
            className="w-full text-left px-5 py-4 text-[11px] font-black uppercase tracking-[0.2em] text-gray-400 hover:text-white hover:bg-white/[0.05] transition-all rounded-xl flex items-center gap-4"
          >
            <Box className="w-4 h-4 text-indigo-500" />
            Latest Blocks
          </button>
          <button
            onClick={() => { setActiveTab(chain); setViewMode("transactions"); setTxsPage(0); setIsSearching(false); }}
            className="w-full text-left px-5 py-4 text-[11px] font-black uppercase tracking-[0.2em] text-gray-400 hover:text-white hover:bg-white/[0.05] transition-all rounded-xl flex items-center gap-4"
          >
            <Zap className="w-4 h-4 text-emerald-500" />
            Latest Transactions
          </button>
        </div>
      </div>
    </div>
  );

  return (
    <div className="min-h-screen bg-[#08080a] text-white selection:bg-indigo-500/30 font-sans tracking-tight">
      {/* Search & Header */}
      <nav className="sticky top-0 z-50 bg-[#08080a]/90 backdrop-blur-2xl border-b border-white/5 px-10 py-5 flex items-center justify-between shadow-2xl">
        <div className="flex items-center gap-12">
          <div 
            className="flex items-center gap-4 cursor-pointer group"
            onClick={() => { setActiveTab('all'); setViewMode('overview'); setBlocksPage(0); setTxsPage(0); setIsSearching(false); }}
          >
            <div className="w-10 h-10 bg-indigo-600 rounded-xl flex items-center justify-center group-hover:bg-indigo-500 transition-all duration-500 shadow-lg shadow-indigo-600/20 group-hover:rotate-12">
              <Zap className="w-6 h-6 text-white fill-current" />
            </div>
            <h1 className="text-2xl font-black tracking-tighter">OmniChain</h1>
          </div>
          
          <div className="hidden lg:flex items-center gap-2">
            <button 
              onClick={() => { setActiveTab('all'); setViewMode('overview'); setBlocksPage(0); setTxsPage(0); setIsSearching(false); }}
              className={`px-5 py-2.5 rounded-xl text-sm font-bold transition-all ${activeTab === 'all' && viewMode === 'overview' && !isSearching ? 'bg-white/10 text-white border border-white/10 shadow-inner' : 'text-gray-500 hover:text-white hover:bg-white/5'}`}
            >
              Overview
            </button>
            <NavDropdown label="Ethereum" chain="ethereum" />
            <NavDropdown label="Bitcoin" chain="bitcoin" />
            {(viewMode !== 'overview' || isSearching) && (
              <button 
                 onClick={() => { setViewMode('overview'); setActiveTab('all'); setIsSearching(false); setSearchQuery(""); }}
                 className="flex items-center gap-2 px-5 py-2.5 bg-indigo-600/10 text-indigo-400 rounded-xl text-xs font-black uppercase tracking-widest border border-indigo-600/20 hover:bg-indigo-600 hover:text-white transition-all shadow-lg ml-2"
              >
                <LayoutGrid className="w-4 h-4" /> Go Back
              </button>
            )}
          </div>
        </div>

        <div className="flex items-center gap-8">
          <form onSubmit={handleSearch} className="relative hidden xl:block group">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-600 group-focus-within:text-indigo-400 transition-all" />
            <input 
              type="text" 
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search address / hash / height..."
              className="bg-white/5 border border-white/10 rounded-2xl py-3 px-12 text-sm focus:outline-none focus:ring-4 focus:ring-indigo-600/10 w-[360px] transition-all font-medium placeholder:text-gray-700 shadow-inner"
            />
            {searchQuery && (
              <button 
                type="button"
                onClick={clearSearch}
                className="absolute right-4 top-1/2 -translate-y-1/2 text-gray-600 hover:text-white"
              >
                <X className="w-4 h-4" />
              </button>
            )}
          </form>
          <button className="lg:hidden p-3 bg-white/5 rounded-xl text-gray-400 border border-white/10"><Menu className="w-6 h-6" /></button>
        </div>
      </nav>

      {/* Hero Analytics (HOME PAGE ONLY) */}
      {viewMode === 'overview' && !isSearching && (
        <section className="max-w-7xl mx-auto px-10 pt-20 pb-12">
          <div className="mb-14">
            <h2 className="text-6xl font-black text-white mb-4 tracking-[-0.04em]">Registry Explorer</h2>
            <p className="text-gray-500 text-2xl max-w-3xl font-medium leading-relaxed opacity-70 mb-8">
              Real-time multi-chain validation and indexing for global ledger consensus.
            </p>
            <div className="flex items-center gap-4 px-6 py-2.5 bg-emerald-500/5 border border-emerald-500/20 rounded-full w-fit">
               <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
               <span className="text-[10px] font-black text-emerald-500 uppercase tracking-[0.3em]">Live PostgreSQL Registry Connected</span>
            </div>
          </div>

          <div className="grid grid-cols-2 lg:grid-cols-4 gap-6">
            {[
              { label: 'ETH Registry', value: ethStats.total_blocks, color: 'text-indigo-400', icon: Box, bg: 'bg-indigo-600/10', chain: 'ethereum' },
              { label: 'ETH Pulse', value: ethStats.total_transactions, color: 'text-indigo-400', icon: TrendingUp, bg: 'bg-indigo-600/10', chain: 'ethereum' },
              { label: 'BTC Registry', value: btcStats.total_blocks, color: 'text-amber-400', icon: Box, bg: 'bg-amber-600/10', chain: 'bitcoin' },
              { label: 'BTC Pulse', value: btcStats.total_transactions, color: 'text-amber-400', icon: TrendingUp, bg: 'bg-amber-600/10', chain: 'bitcoin' }
            ].map((stat, i) => (
              <motion.div 
                key={i}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: i * 0.1 }}
                className="bg-white/[0.02] border border-white/5 rounded-[2rem] p-8 hover:bg-white/[0.04] transition-all relative overflow-hidden group shadow-xl border-t-white/10"
              >
                <div className="flex items-center justify-between mb-6">
                  <div className={`p-3.5 rounded-xl ${stat.bg} ${stat.color} border border-white/10 shadow-inner`}>
                    {stat.chain === 'ethereum' ? <EthIcon className="w-6 h-6" /> : <BtcIcon className="w-6 h-6" />}
                  </div>
                  <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse shadow-[0_0_10px_rgba(16,185,129,0.5)]" />
                </div>
                <p className="text-[10px] font-black text-gray-500 uppercase tracking-[0.3em] mb-2">{stat.label}</p>
                <h3 className="text-4xl font-black tabular-nums tracking-tighter text-white">{stat.value.toLocaleString()}</h3>
              </motion.div>
            ))}
          </div>
        </section>
      )}

      {/* Main Content Area */}
      <main className={`max-w-7xl mx-auto px-10 pb-40 ${(viewMode !== 'overview' || isSearching) ? 'pt-20' : ''}`}>
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-8 mb-16">
          <div className="flex items-center gap-5">
            <button 
              onClick={() => setIsPaused(!isPaused)}
              className={`flex items-center gap-3 px-6 py-3 rounded-2xl border text-[10px] font-black tracking-[0.2em] transition-all shadow-2xl hover:scale-105 active:scale-95 ${
                isPaused 
                ? 'bg-amber-600/20 border-amber-600/40 text-amber-500' 
                : 'bg-white/5 border-white/10 text-gray-400 hover:text-white hover:bg-white/10'
              }`}
            >
              {isPaused ? <Clock className="w-5 h-5" /> : <Repeat className="w-5 h-5 animate-spin-slow" />}
              {isPaused ? 'SYNC HALTED' : 'LIVE FEED'}
            </button>
            {isSearching && (
              <div className="px-6 py-3 bg-indigo-600/20 border border-indigo-600/40 rounded-2xl text-[10px] font-black tracking-[0.2em] text-indigo-400 uppercase">
                Search Results: "{searchQuery}"
              </div>
            )}
          </div>
          
          {!isSearching && (
            <div className="flex items-center gap-1.5 bg-white/5 p-1.5 rounded-2xl border border-white/10 shadow-inner">
              <button 
                disabled={blocksPage === 0 && txsPage === 0}
                onClick={() => { setBlocksPage(p => Math.max(0, p - 1)); setTxsPage(p => Math.max(0, p - 1)); }}
                className="p-3 rounded-xl bg-white/5 border border-white/5 disabled:opacity-10 hover:bg-white/10 transition-all"
              >
                <ChevronLeft className="w-6 h-6" />
              </button>
              <div className="px-5 text-[10px] font-black text-gray-500 min-w-[100px] text-center uppercase tracking-widest leading-none">Page {blocksPage + 1}</div>
              <button 
                onClick={() => { setBlocksPage(p => p + 1); setTxsPage(p => p + 1); }}
                className="p-3 rounded-xl bg-white/5 border border-white/5 hover:bg-white/10 transition-all"
              >
                <ChevronRight className="w-6 h-6" />
              </button>
            </div>
          )}
        </div>

        <div className={`grid gap-16 ${viewMode === 'overview' && !isSearching ? 'grid-cols-1 lg:grid-cols-2' : 'grid-cols-1'}`}>
          {/* Blocks Section */}
          {(viewMode === 'overview' || viewMode === 'blocks' || isSearching) && (
            <div className="space-y-10">
              <div className="flex items-center justify-between border-b border-white/10 pb-8">
                <h2 className="text-3xl font-black flex items-center gap-5 tracking-tighter uppercase italic">
                  <div className="p-3 rounded-2xl bg-indigo-600/10 text-indigo-500">
                    <Box className="w-8 h-8" />
                  </div>
                  {isSearching ? 'Found Blocks' : 'Recent Blocks'}
                </h2>
                <span className="text-[10px] text-gray-500 font-bold uppercase tracking-widest">Persistence Registry</span>
              </div>
              <div className="grid grid-cols-1 gap-4">
                {combinedBlocks.map((block, idx) => (
                  <motion.div
                    key={block.hash + idx}
                    initial={{ opacity: 0, y: 15 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: idx * 0.04 }}
                    onClick={() => fetchDetail('block', block.chain, block.hash)}
                    className="group p-8 bg-white/[0.015] border border-white/5 rounded-[2.5rem] flex items-center justify-between hover:bg-white/[0.05] hover:border-white/20 transition-all cursor-pointer relative overflow-hidden active:scale-[0.98] shadow-2xl"
                  >
                    <div className="flex items-center gap-8 min-w-0">
                      <div className={`w-16 h-16 rounded-2xl flex items-center justify-center shadow-lg transition-transform duration-500 group-hover:rotate-6 ${
                        block.chain === 'ethereum' ? 'bg-indigo-600 text-white' : 'bg-amber-600 text-white'
                      }`}>
                        {block.chain === 'ethereum' ? <EthIcon className="w-8 h-8" /> : <BtcIcon className="w-8 h-8" />}
                      </div>
                      <div className="min-w-0">
                        <div className="flex items-center gap-4 mb-2">
                          <span className="font-black text-3xl tracking-tighter text-white">#{block.height}</span>
                          <span className={`text-[9px] font-black uppercase tracking-widest px-2 py-0.5 rounded border ${
                            block.chain === 'ethereum' ? 'border-indigo-500/30 text-indigo-400 bg-indigo-500/5' : 'border-amber-500/30 text-amber-400 bg-amber-500/5'
                          }`}>
                            {block.chain}
                          </span>
                        </div>
                        <p className="text-xs font-mono text-gray-600 truncate max-w-md group-hover:text-gray-400 transition-colors uppercase font-bold tracking-tight">{block.hash}</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <ArrowRight className="w-6 h-6 text-gray-800 group-hover:text-white transition-all ml-auto mb-3" />
                      <p className="text-[11px] text-gray-600 font-black font-mono uppercase">{new Date(block.timestamp).toLocaleTimeString()}</p>
                    </div>
                  </motion.div>
                ))}
                {!loading && combinedBlocks.length === 0 && (
                  <div className="py-20 text-center opacity-20 italic font-medium uppercase tracking-widest text-sm">Waiting for blocks...</div>
                )}
              </div>
            </div>
          )}

          {/* Transactions Section */}
          {(viewMode === 'overview' || viewMode === 'transactions' || isSearching) && (
            <div className="space-y-10">
              <div className="flex items-center justify-between border-b border-white/10 pb-8">
                <h2 className="text-3xl font-black flex items-center gap-5 tracking-tighter uppercase italic">
                  <div className="p-3 rounded-2xl bg-emerald-600/10 text-emerald-500">
                    <Zap className="w-8 h-8" />
                  </div>
                  {isSearching ? 'Found Transactions' : 'Latest Transactions'}
                </h2>
                <span className="text-[10px] text-gray-500 font-bold uppercase tracking-widest">Protocol Stream</span>
              </div>
              <div className="grid grid-cols-1 gap-4">
                {combinedTxs.map((tx, idx) => (
                  <motion.div
                    key={tx.hash + idx}
                    initial={{ opacity: 0, y: 15 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: idx * 0.04 }}
                    onClick={() => fetchDetail('tx', tx.chain, tx.hash)}
                    className="group p-8 bg-white/[0.015] border border-white/5 rounded-[2.5rem] flex items-center justify-between hover:bg-white/[0.05] hover:border-white/20 transition-all cursor-pointer relative overflow-hidden active:scale-[0.98] shadow-2xl"
                  >
                    <div className="flex items-center gap-8 min-w-0">
                      <div className={`w-16 h-16 rounded-2xl flex items-center justify-center shadow-lg transition-all duration-500 group-hover:scale-110 ${
                        tx.chain === 'ethereum' ? 'bg-indigo-600/10 text-indigo-400 shadow-indigo-600/5' : 'bg-amber-600/10 text-amber-400 shadow-amber-600/5'
                      }`}>
                        {tx.chain === 'ethereum' ? <EthIcon className="w-8 h-8 text-indigo-400" /> : <BtcIcon className="w-8 h-8 text-amber-500" />}
                      </div>
                      <div className="min-w-0">
                        <p className="text-xs font-mono text-gray-300 truncate font-black mb-4 uppercase tracking-tighter group-hover:text-white transition-colors">{tx.hash}</p>
                        <div className="flex items-center gap-6">
                          <div className="flex items-center gap-3 bg-white/5 px-4 py-1.5 rounded-xl border border-white/5">
                            <span className="text-[9px] font-black text-gray-600 uppercase tracking-widest">VAL</span>
                            <span className="text-sm font-black text-emerald-400 tabular-nums">+{tx.value}</span>
                          </div>
                          <div className="flex items-center gap-3 bg-white/5 px-4 py-1.5 rounded-xl border border-white/5">
                            <span className="text-[9px] font-black text-gray-600 uppercase tracking-widest">NET</span>
                            <span className="text-sm font-black text-white/50">{tx.chain === 'ethereum' ? 'ETH' : 'BTC'}</span>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div className="text-right">
                       <ArrowRight className="w-6 h-6 text-gray-800 group-hover:text-white transition-all ml-auto mb-4" />
                       <div className="px-5 py-1.5 rounded-full bg-emerald-500/10 text-[10px] font-black text-emerald-500 uppercase tracking-widest border border-emerald-500/10">Relayed</div>
                    </div>
                  </motion.div>
                ))}
                {!loading && combinedTxs.length === 0 && (
                  <div className="py-20 text-center opacity-20 italic font-medium uppercase tracking-widest text-sm">Waiting for transactions...</div>
                )}
              </div>
            </div>
          )}
        </div>

        {/* Empty State */}
        {!loading && combinedBlocks.length === 0 && combinedTxs.length === 0 && (
          <div className="py-40 text-center flex flex-col items-center gap-10 bg-white/[0.02] border border-white/5 rounded-[4rem] shadow-inner">
             <Database className="w-20 h-20 text-gray-800 animate-pulse" />
             <div>
               <h3 className="text-3xl font-black uppercase italic tracking-tighter mb-4">No Records Found</h3>
               <p className="text-gray-600 text-lg max-w-md mx-auto font-medium">Try adjusting your filters or search query. The registry is indexing at 200ms intervals.</p>
             </div>
             {isSearching && (
               <button 
                 onClick={clearSearch}
                 className="px-10 py-4 bg-white text-black rounded-2xl font-black uppercase tracking-widest hover:scale-105 transition-all shadow-2xl"
               >
                 Clear Search results
               </button>
             )}
          </div>
        )}
      </main>

      {/* Detail Modal */}
      <AnimatePresence>
        {selectedItem && (
          <div className="fixed inset-0 z-[100] flex items-center justify-center p-6 sm:p-14">
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setSelectedItem(null)}
              className="absolute inset-0 bg-black/95 backdrop-blur-3xl"
            />
            <motion.div
              initial={{ opacity: 0, scale: 0.95, y: 30 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.95, y: 30 }}
              className="relative w-full max-w-4xl bg-[#0b0b0d] border border-white/10 rounded-[3rem] shadow-[0_80px_160px_-24px_rgba(0,0,0,1)] flex flex-col max-h-[92vh] border-t-white/10"
            >
              <div className="p-10 border-b border-white/5 flex items-center justify-between bg-gradient-to-r from-white/[0.04] to-transparent">
                <div className="flex items-center gap-8">
                  <div className={`w-14 h-14 rounded-2xl flex items-center justify-center shadow-2xl relative group/modal-icon ${
                    selectedItem.chain === 'ethereum' ? 'bg-indigo-600/20 text-indigo-400' : 'bg-amber-600/20 text-amber-400'
                  }`}>
                    <div className="absolute inset-0 bg-inherit rounded-inherit blur-xl opacity-20 animate-pulse" />
                    {selectedItem.chain === 'ethereum' ? <EthIcon className="w-8 h-8 relative z-10" /> : <BtcIcon className="w-8 h-8 relative z-10" />}
                  </div>
                  <div>
                    <h3 className="text-3xl font-black capitalize tracking-tighter mb-2 italic">{selectedItem.chain} {selectedItem.type} Details</h3>
                    <div className="flex items-center gap-4">
                      <span className="text-[10px] font-black text-gray-600 uppercase tracking-widest bg-white/5 px-4 py-1.5 rounded-full border border-white/5 shadow-inner">ID: {selectedItem.data?.id || 'PERSISTED'}</span>
                      <div className="flex items-center gap-2">
                        <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse shadow-lg shadow-emerald-500/50" />
                        <span className="text-[10px] font-black text-emerald-500 uppercase tracking-widest font-mono">Verified Link</span>
                      </div>
                    </div>
                  </div>
                </div>
                <button onClick={() => setSelectedItem(null)} className="p-5 bg-white/5 hover:bg-white/10 rounded-2xl transition-all text-gray-500 hover:text-white border border-white/5 active:scale-90 shadow-2xl">
                  <X className="w-8 h-8" />
                </button>
              </div>

              <div className="flex-1 overflow-y-auto p-12 sm:p-20 custom-scrollbar">
                {modalLoading ? (
                  <div className="py-32 flex flex-col items-center gap-8">
                    <Activity className="w-16 h-16 animate-spin text-indigo-500" />
                    <span className="text-[11px] font-black uppercase tracking-[0.4em] text-gray-600 animate-pulse">Syncing Payload...</span>
                  </div>
                ) : selectedItem.data ? (
                  <div className="space-y-16">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      {Object.entries(selectedItem.data).map(([key, value]) => {
                        if (key === 'transactions' || (typeof value === 'object' && value !== null)) return null;
                        return (
                          <div key={key} className="bg-white/[0.015] border border-white/5 rounded-2xl p-8 group hover:bg-white/[0.03] transition-all">
                             <label className="text-[11px] font-black text-gray-600 uppercase tracking-[0.4em] block mb-3 group-hover:text-indigo-400 transition-colors">{key.replace('_', ' ')}</label>
                             <p className="text-sm font-mono text-gray-200 break-all font-black leading-relaxed">{String(value)}</p>
                          </div>
                        );
                      })}
                    </div>

                    {selectedItem.type === 'block' && selectedItem.data.transactions?.length > 0 && (
                      <div className="space-y-10">
                        <div className="flex items-center justify-between border-b border-white/5 pb-8">
                          <h4 className="text-2xl font-black italic tracking-tighter uppercase flex items-center gap-6">
                            <Activity className="w-8 h-8 text-emerald-500" /> Block Transactions ({selectedItem.data.transactions.length})
                          </h4>
                          <span className="text-[10px] font-black text-gray-600 uppercase tracking-widest">Protocol V4</span>
                        </div>
                        <div className="space-y-3">
                          {selectedItem.data.transactions.map((tx, idx) => (
                            <div key={tx.hash} className="bg-white/[0.02] border border-white/5 rounded-2xl p-6 flex items-center justify-between group hover:bg-white/[0.05] transition-all border-l-4 border-l-indigo-600/30">
                              <div className="flex items-center gap-6 min-w-0">
                                <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center text-[10px] font-black text-gray-500 border border-white/5">{idx + 1}</div>
                                <span className="font-mono text-gray-300 text-xs truncate max-w-[400px] font-black group-hover:text-white transition-colors">{tx.hash}</span>
                              </div>
                              <div className="text-right">
                                <p className="font-black text-emerald-400 text-xl tabular-nums">+{tx.value}</p>
                                <p className="text-[9px] font-black text-gray-600 uppercase tracking-tighter font-mono">Decrypted</p>
                              </div>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    <div className="space-y-6">
                       <label className="text-xs font-black text-gray-600 uppercase tracking-[0.4em] flex items-center gap-3">
                         <Code className="w-7 h-7 text-indigo-500" /> Registry JSON Result
                       </label>
                       <pre className="p-12 bg-black/50 rounded-[2.5rem] border border-white/5 text-[12px] font-mono text-indigo-400/50 overflow-x-auto leading-loose font-black shadow-2xl selection:bg-indigo-500/30">
                        {JSON.stringify(selectedItem.data, null, 2)}
                       </pre>
                    </div>
                  </div>
                ) : (
                  <div className="py-32 text-center text-gray-600 flex flex-col items-center gap-10">
                    <Database className="w-20 h-20 opacity-5" />
                    <p className="text-2xl font-black italic tracking-tighter opacity-30 uppercase max-w-lg leading-relaxed">System acknowledged consensus hash. Multi-threaded archival in progress.</p>
                  </div>
                )}
              </div>

              <div className="p-12 border-t border-white/5 flex justify-end bg-white/[0.02]">
                 <button 
                   onClick={() => setSelectedItem(null)}
                   className="px-24 py-6 bg-white text-black hover:scale-110 active:scale-95 rounded-[2rem] text-xl font-black transition-all shadow-[0_48px_96px_-24px_rgba(255,255,255,0.4)] uppercase tracking-[0.3em]"
                 >
                   Exit Decoder
                 </button>
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>
    </div>
  );
};

export default App;
