<div class="page-content">
  <h1>Sources</h1>
  <p>Configure and manage your image sources</p>
  
  <div class="sources-grid">
    <div class="source-card">
      <div class="source-header">
        <h3>🔴 Reddit</h3>
        <span class="source-status active">Active</span>
      </div>
      <div class="source-info">
        <p><strong>Subreddits:</strong> wallpapers, EarthPorn, SkyPorn</p>
        <p><strong>Schedule:</strong> Every 6 hours</p>
        <p><strong>Last Run:</strong> 2 hours ago</p>
        <p><strong>Images Collected:</strong> 1,247</p>
      </div>
      <div class="source-actions">
        <button class="btn-primary">Configure</button>
        <button class="btn-secondary">Run Now</button>
      </div>
    </div>

    <div class="source-card">
      <div class="source-header">
        <h3>🎨 Danbooru</h3>
        <span class="source-status active">Active</span>
      </div>
      <div class="source-info">
        <p><strong>Tags:</strong> landscape, scenery, nature</p>
        <p><strong>Schedule:</strong> Daily at 3:00 AM</p>
        <p><strong>Last Run:</strong> 8 hours ago</p>
        <p><strong>Images Collected:</strong> 432</p>
      </div>
      <div class="source-actions">
        <button class="btn-primary">Configure</button>
        <button class="btn-secondary">Run Now</button>
      </div>
    </div>

    <div class="source-card">
      <div class="source-header">
        <h3>📸 Unsplash</h3>
        <span class="source-status paused">Paused</span>
      </div>
      <div class="source-info">
        <p><strong>Keywords:</strong> minimal, abstract, architecture</p>
        <p><strong>Schedule:</strong> Every 12 hours</p>
        <p><strong>Last Run:</strong> 3 days ago</p>
        <p><strong>Images Collected:</strong> 89</p>
      </div>
      <div class="source-actions">
        <button class="btn-primary">Configure</button>
        <button class="btn-secondary">Resume</button>
      </div>
    </div>

    <div class="source-card">
      <div class="source-header">
        <h3>🌐 Custom RSS</h3>
        <span class="source-status inactive">Inactive</span>
      </div>
      <div class="source-info">
        <p><strong>Feed URL:</strong> example.com/wallpapers.rss</p>
        <p><strong>Schedule:</strong> Not configured</p>
        <p><strong>Last Run:</strong> Never</p>
        <p><strong>Images Collected:</strong> 0</p>
      </div>
      <div class="source-actions">
        <button class="btn-primary">Configure</button>
        <button class="btn-secondary">Enable</button>
      </div>
    </div>
  </div>

  <div class="add-source-section">
    <h2>Add New Source</h2>
    <div class="source-types">
      <button class="source-type-btn">
        <span class="icon">🔴</span>
        <span>Reddit</span>
      </button>
      <button class="source-type-btn">
        <span class="icon">🎨</span>
        <span>Booru</span>
      </button>
      <button class="source-type-btn">
        <span class="icon">📸</span>
        <span>Unsplash</span>
      </button>
      <button class="source-type-btn">
        <span class="icon">🌐</span>
        <span>RSS Feed</span>
      </button>
      <button class="source-type-btn">
        <span class="icon">📁</span>
        <span>Local Folder</span>
      </button>
    </div>
  </div>
</div>

<style>
  .page-content {
    padding: 2rem;
    max-width: 1200px;
  }
  
  .page-content h1 {
    margin: 0 0 1rem 0;
    font-size: 2rem;
    font-weight: 600;
    color: #ffffff;
  }
  
  .page-content > p {
    color: #aaa;
    font-size: 1.1rem;
    margin-bottom: 2rem;
  }
  
  .sources-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    gap: 1.5rem;
    margin-bottom: 3rem;
  }
  
  .source-card {
    background-color: #2d2d2d;
    border-radius: 8px;
    padding: 1.5rem;
    border-left: 4px solid #646cff;
  }
  
  .source-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }
  
  .source-header h3 {
    margin: 0;
    color: #ffffff;
    font-size: 1.1rem;
    font-weight: 600;
  }
  
  .source-status {
    padding: 0.25rem 0.75rem;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: 600;
    text-transform: uppercase;
  }
  
  .source-status.active {
    background-color: rgba(34, 197, 94, 0.2);
    color: #22c55e;
  }
  
  .source-status.paused {
    background-color: rgba(251, 191, 36, 0.2);
    color: #fbbf24;
  }
  
  .source-status.inactive {
    background-color: rgba(239, 68, 68, 0.2);
    color: #ef4444;
  }
  
  .source-info {
    margin-bottom: 1.5rem;
  }
  
  .source-info p {
    margin: 0.5rem 0;
    color: #aaa;
    font-size: 0.9rem;
  }
  
  .source-info strong {
    color: #ffffff;
  }
  
  .source-actions {
    display: flex;
    gap: 0.75rem;
  }
  
  .btn-primary, .btn-secondary {
    padding: 0.5rem 1rem;
    border-radius: 6px;
    border: none;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }
  
  .btn-primary {
    background-color: #646cff;
    color: #ffffff;
  }
  
  .btn-primary:hover {
    background-color: #5a5fcf;
  }
  
  .btn-secondary {
    background-color: #3d3d3d;
    color: #aaa;
    border: 1px solid #5d5d5d;
  }
  
  .btn-secondary:hover {
    background-color: #4d4d4d;
    color: #ffffff;
  }
  
  .add-source-section {
    background-color: #2d2d2d;
    border-radius: 8px;
    padding: 2rem;
    border: 2px dashed #5d5d5d;
  }
  
  .add-source-section h2 {
    margin: 0 0 1.5rem 0;
    color: #ffffff;
    font-size: 1.3rem;
    font-weight: 600;
    text-align: center;
  }
  
  .source-types {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
    gap: 1rem;
  }
  
  .source-type-btn {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    padding: 1rem;
    background-color: #3d3d3d;
    border: 1px solid #5d5d5d;
    border-radius: 8px;
    color: #aaa;
    cursor: pointer;
    transition: all 0.2s ease;
  }
  
  .source-type-btn:hover {
    background-color: #4d4d4d;
    color: #ffffff;
    border-color: #646cff;
  }
  
  .source-type-btn .icon {
    font-size: 1.5rem;
  }
  
  .source-type-btn span:last-child {
    font-size: 0.9rem;
    font-weight: 500;
  }
</style>