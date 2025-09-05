<div class="page-content">
  <h1>Jobs</h1>
  <p>Monitor and manage background tasks and scheduled operations</p>
  
  <div class="job-controls">
    <button class="btn-primary">🔄 Refresh</button>
    <button class="btn-secondary">⏸️ Pause All</button>
  </div>

  <div class="jobs-section">
    <h2>Running Jobs</h2>
    <div class="jobs-list">
      <div class="job-card running">
        <div class="job-header">
          <h3>Reddit Image Download</h3>
          <span class="job-status running">Running</span>
        </div>
        <div class="job-info">
          <p><strong>Source:</strong> r/wallpapers</p>
          <p><strong>Started:</strong> 5 minutes ago</p>
          <p><strong>Progress:</strong> 45 / 100 images</p>
        </div>
        <div class="job-progress">
          <div class="progress-bar">
            <div class="progress-fill" style="width: 45%"></div>
          </div>
          <span class="progress-text">45%</span>
        </div>
        <div class="job-actions">
          <button class="btn-danger">Stop</button>
          <button class="btn-secondary">View Logs</button>
        </div>
      </div>

      <div class="job-card running">
        <div class="job-header">
          <h3>Device Sync</h3>
          <span class="job-status running">Running</span>
        </div>
        <div class="job-info">
          <p><strong>Target:</strong> Mobile Device</p>
          <p><strong>Started:</strong> 2 minutes ago</p>
          <p><strong>Progress:</strong> 23 / 67 files</p>
        </div>
        <div class="job-progress">
          <div class="progress-bar">
            <div class="progress-fill" style="width: 34%"></div>
          </div>
          <span class="progress-text">34%</span>
        </div>
        <div class="job-actions">
          <button class="btn-danger">Stop</button>
          <button class="btn-secondary">View Logs</button>
        </div>
      </div>
    </div>
  </div>

  <div class="jobs-section">
    <h2>Recent Jobs</h2>
    <div class="jobs-list">
      <div class="job-card completed">
        <div class="job-header">
          <h3>Danbooru Image Download</h3>
          <span class="job-status completed">Completed</span>
        </div>
        <div class="job-info">
          <p><strong>Source:</strong> danbooru.donmai.us</p>
          <p><strong>Completed:</strong> 2 hours ago</p>
          <p><strong>Result:</strong> 34 images downloaded</p>
        </div>
        <div class="job-actions">
          <button class="btn-secondary">View Results</button>
          <button class="btn-secondary">View Logs</button>
        </div>
      </div>

      <div class="job-card failed">
        <div class="job-header">
          <h3>RSS Feed Processing</h3>
          <span class="job-status failed">Failed</span>
        </div>
        <div class="job-info">
          <p><strong>Source:</strong> example.com/feed.rss</p>
          <p><strong>Failed:</strong> 6 hours ago</p>
          <p><strong>Error:</strong> Connection timeout</p>
        </div>
        <div class="job-actions">
          <button class="btn-primary">Retry</button>
          <button class="btn-secondary">View Logs</button>
        </div>
      </div>

      <div class="job-card completed">
        <div class="job-header">
          <h3>Image Processing</h3>
          <span class="job-status completed">Completed</span>
        </div>
        <div class="job-info">
          <p><strong>Task:</strong> Generate thumbnails</p>
          <p><strong>Completed:</strong> 1 day ago</p>
          <p><strong>Result:</strong> 156 thumbnails created</p>
        </div>
        <div class="job-actions">
          <button class="btn-secondary">View Results</button>
          <button class="btn-secondary">View Logs</button>
        </div>
      </div>
    </div>
  </div>

  <div class="jobs-section">
    <h2>Scheduled Jobs</h2>
    <div class="jobs-list">
      <div class="job-card scheduled">
        <div class="job-header">
          <h3>Daily Reddit Sync</h3>
          <span class="job-status scheduled">Scheduled</span>
        </div>
        <div class="job-info">
          <p><strong>Next Run:</strong> Today at 6:00 PM</p>
          <p><strong>Frequency:</strong> Daily</p>
          <p><strong>Last Result:</strong> 45 images downloaded</p>
        </div>
        <div class="job-actions">
          <button class="btn-primary">Run Now</button>
          <button class="btn-secondary">Edit Schedule</button>
        </div>
      </div>

      <div class="job-card scheduled">
        <div class="job-header">
          <h3>Weekly Cleanup</h3>
          <span class="job-status scheduled">Scheduled</span>
        </div>
        <div class="job-info">
          <p><strong>Next Run:</strong> Sunday at 2:00 AM</p>
          <p><strong>Frequency:</strong> Weekly</p>
          <p><strong>Last Result:</strong> 12 old files removed</p>
        </div>
        <div class="job-actions">
          <button class="btn-primary">Run Now</button>
          <button class="btn-secondary">Edit Schedule</button>
        </div>
      </div>
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
  
  .job-controls {
    display: flex;
    gap: 1rem;
    margin-bottom: 2rem;
  }
  
  .jobs-section {
    margin-bottom: 3rem;
  }
  
  .jobs-section h2 {
    color: #ffffff;
    margin: 0 0 1rem 0;
    font-size: 1.5rem;
    font-weight: 600;
  }
  
  .jobs-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .job-card {
    background-color: #2d2d2d;
    border-radius: 8px;
    padding: 1.5rem;
    border-left: 4px solid #646cff;
  }
  
  .job-card.running {
    border-left-color: #3b82f6;
  }
  
  .job-card.completed {
    border-left-color: #22c55e;
  }
  
  .job-card.failed {
    border-left-color: #ef4444;
  }
  
  .job-card.scheduled {
    border-left-color: #f59e0b;
  }
  
  .job-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }
  
  .job-header h3 {
    margin: 0;
    color: #ffffff;
    font-size: 1.1rem;
    font-weight: 600;
  }
  
  .job-status {
    padding: 0.25rem 0.75rem;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: 600;
    text-transform: uppercase;
  }
  
  .job-status.running {
    background-color: rgba(59, 130, 246, 0.2);
    color: #3b82f6;
  }
  
  .job-status.completed {
    background-color: rgba(34, 197, 94, 0.2);
    color: #22c55e;
  }
  
  .job-status.failed {
    background-color: rgba(239, 68, 68, 0.2);
    color: #ef4444;
  }
  
  .job-status.scheduled {
    background-color: rgba(245, 158, 11, 0.2);
    color: #f59e0b;
  }
  
  .job-info {
    margin-bottom: 1rem;
  }
  
  .job-info p {
    margin: 0.5rem 0;
    color: #aaa;
    font-size: 0.9rem;
  }
  
  .job-info strong {
    color: #ffffff;
  }
  
  .job-progress {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-bottom: 1rem;
  }
  
  .progress-bar {
    flex: 1;
    height: 8px;
    background-color: #3d3d3d;
    border-radius: 4px;
    overflow: hidden;
  }
  
  .progress-fill {
    height: 100%;
    background-color: #3b82f6;
    transition: width 0.3s ease;
  }
  
  .progress-text {
    color: #aaa;
    font-size: 0.9rem;
    font-weight: 600;
    min-width: 40px;
  }
  
  .job-actions {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
  }
  
  .btn-primary, .btn-secondary, .btn-danger {
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
  
  .btn-danger {
    background-color: #ef4444;
    color: #ffffff;
  }
  
  .btn-danger:hover {
    background-color: #dc2626;
  }
</style>