import axios from 'axios';

const BASE_URL = 'http://localhost:8989/api';

const api = axios.create({
  baseURL: BASE_URL,
  timeout: 10000,
});

export interface ChainStats {
  latestBlock: number;
  totalBlocks: number;
  totalTx: number;
  synced: boolean;
}

export interface StatsResponse {
  btc: ChainStats;
  eth: ChainStats;
}

export interface Block {
  height: number;
  hash: string;
  txCount: number;
  timestamp: number;
}

export interface PaginatedBlocksResponse {
  page: number;
  limit: number;
  total: number;
  blocks: Block[];
}

export interface Transaction {
  hash: string;
  from: string;
  to: string;
  value: string;
  height: number;
  timestamp: number;
}

export interface BlockDetails extends Block {
  transactions: Transaction[];
}

export interface SearchResult {
  type: 'block' | 'transaction';
  chain: 'btc' | 'eth';
  result: any;
}

export const getStats = async (): Promise<StatsResponse> => {
  const { data } = await api.get<StatsResponse>('/stats');
  return data;
};

export const getBlocks = async (
  chain: 'btc' | 'eth',
  page: number = 1,
  limit: number = 20
): Promise<PaginatedBlocksResponse> => {
  const { data } = await api.get<PaginatedBlocksResponse>(`/${chain}/blocks`, {
    params: { page, limit },
  });
  return data;
};

export const getBlock = async (
  chain: 'btc' | 'eth',
  height: string | number
): Promise<BlockDetails> => {
  const { data } = await api.get<BlockDetails>(`/${chain}/blocks/${height}`);
  return data;
};

export const search = async (q: string): Promise<SearchResult> => {
  const { data } = await api.get<SearchResult>('/search', {
    params: { q },
  });
  return data;
};

export default api;
