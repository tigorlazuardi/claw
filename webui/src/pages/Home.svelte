<div class="page-content">
  <h1>Welcome to Claw</h1>
  <p>Image collection and management system</p>
  
  <div class="dashboard-cards">
    <div class="dashboard-card">
      <h3>📚 Total Images</h3>
      <p class="stat-number">1,247</p>
    </div>
    
    <div class="dashboard-card">
      <h3>📱 Active Devices</h3>
      <p class="stat-number">3</p>
    </div>
    
    <div class="dashboard-card">
      <h3>📡 Sources</h3>
      <p class="stat-number">5</p>
    </div>
    
    <div class="dashboard-card">
      <h3>⚙️ Running Jobs</h3>
      <p class="stat-number">2</p>
    </div>
  </div>

  <div class="recent-activity">
    <h2>Recent Activity</h2>
    <div class="activity-list">
      <div class="activity-item">
        <span class="activity-time">2 hours ago</span>
        <span class="activity-text">Downloaded 23 images from Reddit source</span>
      </div>
      <div class="activity-item">
        <span class="activity-time">5 hours ago</span>
        <span class="activity-text">Synced images to Mobile Device</span>
      </div>
      <div class="activity-item">
        <span class="activity-time">1 day ago</span>
        <span class="activity-text">Added new Booru source configuration</span>
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
  
  .dashboard-cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1.5rem;
    margin-bottom: 3rem;
  }
  
  .dashboard-card {
    background-color: #2d2d2d;
    border-radius: 8px;
    padding: 1.5rem;
    border-left: 4px solid #646cff;
  }
  
  .dashboard-card h3 {
    margin: 0 0 0.5rem 0;
    color: #ffffff;
    font-size: 1rem;
    font-weight: 600;
  }
  
  .stat-number {
    font-size: 2rem;
    font-weight: 700;
    color: #646cff;
    margin: 0;
  }
  
  .recent-activity h2 {
    color: #ffffff;
    margin: 0 0 1rem 0;
    font-size: 1.5rem;
    font-weight: 600;
  }
  
  .activity-list {
    background-color: #2d2d2d;
    border-radius: 8px;
    padding: 1.5rem;
  }
  
  .activity-item {
    display: flex;
    gap: 1rem;
    padding: 0.75rem 0;
    border-bottom: 1px solid #3d3d3d;
  }
  
  .activity-item:last-child {
    border-bottom: none;
  }
  
  .activity-time {
    color: #aaa;
    font-size: 0.9rem;
    min-width: 100px;
  }
  
  .activity-text {
    color: #ffffff;
  }
</style>