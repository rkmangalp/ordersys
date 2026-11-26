import { useState, useCallback } from 'react'

// API base URL - uses Vite proxy in dev, relative in prod
const API_BASE = '/api/v1'

// Helper to format JSON nicely
const formatJson = (obj) => {
  try {
    return JSON.stringify(obj, null, 2)
  } catch {
    return String(obj)
  }
}

// Helper to get current timestamp
const now = () => new Date().toLocaleTimeString()

function App() {
  const [healthResponse, setHealthResponse] = useState(null)
  const [productsResponse, setProductsResponse] = useState(null)
  const [loading, setLoading] = useState({ health: false, products: false })
  const [logs, setLogs] = useState([])
  const [serverStatus, setServerStatus] = useState('unknown') // 'healthy' | 'unhealthy' | 'unknown'

  const addLog = useCallback((method, path, status, duration) => {
    setLogs(prev => [
      { method, path, status, duration, time: now() },
      ...prev.slice(0, 19) // Keep last 20
    ])
  }, [])

  const callApi = async (endpoint, setter, loadingKey, base = API_BASE) => {
    setLoading(prev => ({ ...prev, [loadingKey]: true }))
    const start = performance.now()
    try {
      const res = await fetch(`${base}${endpoint}`)
      const duration = Math.round(performance.now() - start)
      const data = await res.json()
      setter(data)
      addLog(endpoint.includes('health') ? 'GET' : 'GET', endpoint, res.status, `${duration}ms`)
      if (endpoint.includes('health')) {
        setServerStatus(res.ok ? 'healthy' : 'unhealthy')
      }
      return { ok: res.ok, data }
    } catch (err) {
      const duration = Math.round(performance.now() - start)
      addLog(endpoint.includes('health') ? 'GET' : 'GET', endpoint, 'ERR', `${duration}ms`)
      setter({ error: err.message })
      if (endpoint.includes('health')) {
        setServerStatus('unhealthy')
      }
      return { ok: false, error: err.message }
    } finally {
      setLoading(prev => ({ ...prev, [loadingKey]: false }))
    }
  }

  // /health is not under /api/v1 — the Go server registers it at the root.
  const testHealth = () => callApi('/health', setHealthResponse, 'health', '')
  const testProducts = () => callApi('/products', setProductsResponse, 'products')

  const testAll = async () => {
    await testHealth()
    await testProducts()
  }

  const clearLogs = () => setLogs([])

  return (
    <div className="container">
      <header>
        <h1>OrderSys API Test UI</h1>
        <p className="subtitle">
          Visual test harness for Phase 1 HTTP server endpoints
        </p>
        <div className="status-indicator" style={{ marginTop: '1rem' }}>
          <span className={`status-dot ${serverStatus === 'healthy' ? 'status-healthy' : serverStatus === 'unhealthy' ? 'status-error' : 'status-unknown'}`}></span>
          Server: {serverStatus === 'healthy' ? 'Healthy ✓' : serverStatus === 'unhealthy' ? 'Unhealthy ✗' : 'Unknown'}
        </div>
      </header>

      {/* Health Endpoint */}
      <section className="endpoint-card">
        <div className="endpoint-header">
          <span className="method-badge method-get">GET</span>
          <span className="endpoint-path">/health</span>
        </div>
        <p style={{ marginBottom: '1rem', color: '#666' }}>
          Health check endpoint — returns server status and timestamp
        </p>
        <div style={{ display: 'flex', gap: '0.75rem', flexWrap: 'wrap', marginBottom: '1rem' }}>
          <button
            className="btn btn-primary"
            onClick={testHealth}
            disabled={loading.health}
          >
            {loading.health && <span className="spinner" />}
            Test Health
          </button>
        </div>
        {healthResponse && (
          <div className="response-area">
            <span className="response-timestamp">[{now()}] Response:</span>
            <span className={healthResponse.error ? 'response-error' : 'response-success'}>
{formatJson(healthResponse)}
            </span>
          </div>
        )}
      </section>

      {/* Products Endpoint */}
      <section className="endpoint-card">
        <div className="endpoint-header">
          <span className="method-badge method-get">GET</span>
          <span className="endpoint-path">/api/v1/products</span>
        </div>
        <p style={{ marginBottom: '1rem', color: '#666' }}>
          List products — currently returns empty array (placeholder)
        </p>
        <div style={{ display: 'flex', gap: '0.75rem', flexWrap: 'wrap', marginBottom: '1rem' }}>
          <button
            className="btn btn-primary"
            onClick={testProducts}
            disabled={loading.products}
          >
            {loading.products && <span className="spinner" />}
            Test Products
          </button>
        </div>
        {productsResponse && (
          <div className="response-area">
            <span className="response-timestamp">[{now()}] Response:</span>
            <span className={productsResponse.error ? 'response-error' : 'response-success'}>
{formatJson(productsResponse)}
            </span>
          </div>
        )}
      </section>

      {/* Bulk Actions */}
      <section className="endpoint-card">
        <div style={{ display: 'flex', gap: '0.75rem', flexWrap: 'wrap', marginBottom: '1rem' }}>
          <button
            className="btn btn-secondary"
            onClick={testAll}
            disabled={loading.health || loading.products}
          >
            Test All Endpoints
          </button>
        </div>
      </section>

      {/* Request Logs */}
      <section className="logs-section">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
          <h3 style={{ fontSize: '1rem' }}>Request Logs</h3>
          <button className="btn btn-secondary" onClick={clearLogs} style={{ padding: '0.375rem 0.75rem', fontSize: '0.8rem' }}>
            Clear
          </button>
        </div>
        <div style={{ maxHeight: '300px', overflow: 'auto' }}>
          {logs.length === 0 ? (
            <p style={{ color: '#999', textAlign: 'center', padding: '2rem' }}>No requests yet</p>
          ) : (
            logs.map((log, i) => (
              <div key={i} className="log-entry">
                <span className="log-method">{log.method}</span>
                <span className="log-path">{log.path}</span>
                <span className={`log-status ${log.status === 200 || log.status === '200' ? 'log-status-200' : log.status >= 400 && log.status < 500 ? 'log-status-4xx' : log.status >= 500 ? 'log-status-5xx' : ''}`}>
                  {log.status}
                </span>
                <span className="log-duration">{log.duration}</span>
                <span className="log-duration">{log.time}</span>
              </div>
            ))
          )}
        </div>
      </section>

      {/* Footer Info */}
      <footer style={{ marginTop: '2rem', paddingTop: '1rem', borderTop: '1px solid #eee', color: '#999', fontSize: '0.8rem' }}>
        <p>Make sure the Go server is running on <code>http://localhost:8080</code></p>
        <p>Vite proxies <code>/api/*</code> to the Go server (see <code>vite.config.js</code>)</p>
      </footer>
    </div>
  )
}

export default App