import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Navbar, Sidebar } from './components/layout/Layout';
import Dashboard from './pages/Dashboard';
import Explorer from './pages/Explorer';
import BlockDetails from './pages/BlockDetails';
import SearchResults from './pages/SearchResults';

function App() {
  return (
    <Router>
      <div className="flex flex-col min-h-screen">
        <Navbar />
        <div className="flex flex-1">
          <Sidebar />
          <main className="flex-1 bg-background p-6 overflow-y-auto max-h-[calc(100vh-64px)]">
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/btc" element={<Explorer chain="btc" />} />
              <Route path="/eth" element={<Explorer chain="eth" />} />
              <Route path="/btc/block/:height" element={<BlockDetails chain="btc" />} />
              <Route path="/eth/block/:height" element={<BlockDetails chain="eth" />} />
              <Route path="/search" element={<SearchResults />} />
            </Routes>
          </main>
        </div>
      </div>
    </Router>
  );
}

export default App;
