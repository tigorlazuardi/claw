<div class="page-content">
  <h1>Devices</h1>
  <p>Manage your target devices and their image filtering preferences</p>
  
  <div class="devices-grid">
    <div class="device-card">
      <div class="device-header">
        <h3>📱 Mobile Phone</h3>
        <span class="device-status active">Active</span>
      </div>
      <div class="device-info">
        <p><strong>Resolution:</strong> 1080 × 2340</p>
        <p><strong>Aspect Ratio:</strong> 9:19.5</p>
        <p><strong>Images:</strong> 342</p>
        <p><strong>Storage Used:</strong> 2.1 GB</p>
      </div>
      <div class="device-actions">
        <button class="btn-primary">Configure</button>
        <button class="btn-secondary">Sync Now</button>
      </div>
    </div>

    <div class="device-card">
      <div class="device-header">
        <h3>💻 Desktop</h3>
        <span class="device-status active">Active</span>
      </div>
      <div class="device-info">
        <p><strong>Resolution:</strong> 2560 × 1440</p>
        <p><strong>Aspect Ratio:</strong> 16:9</p>
        <p><strong>Images:</strong> 856</p>
        <p><strong>Storage Used:</strong> 5.7 GB</p>
      </div>
      <div class="device-actions">
        <button class="btn-primary">Configure</button>
        <button class="btn-secondary">Sync Now</button>
      </div>
    </div>

    <div class="device-card">
      <div class="device-header">
        <h3>📺 TV</h3>
        <span class="device-status inactive">Inactive</span>
      </div>
      <div class="device-info">
        <p><strong>Resolution:</strong> 3840 × 2160</p>
        <p><strong>Aspect Ratio:</strong> 16:9</p>
        <p><strong>Images:</strong> 123</p>
        <p><strong>Storage Used:</strong> 1.3 GB</p>
      </div>
      <div class="device-actions">
        <button class="btn-primary">Configure</button>
        <button class="btn-secondary">Enable</button>
      </div>
    </div>
  </div>

  <div class="add-device-section">
    <h2>Add New Device</h2>
    <button class="btn-primary">➕ Add Device</button>
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
  
  .devices-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 1.5rem;
    margin-bottom: 3rem;
  }
  
  .device-card {
    background-color: #2d2d2d;
    border-radius: 8px;
    padding: 1.5rem;
    border-left: 4px solid #646cff;
  }
  
  .device-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }
  
  .device-header h3 {
    margin: 0;
    color: #ffffff;
    font-size: 1.1rem;
    font-weight: 600;
  }
  
  .device-status {
    padding: 0.25rem 0.75rem;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: 600;
    text-transform: uppercase;
  }
  
  .device-status.active {
    background-color: rgba(34, 197, 94, 0.2);
    color: #22c55e;
  }
  
  .device-status.inactive {
    background-color: rgba(239, 68, 68, 0.2);
    color: #ef4444;
  }
  
  .device-info {
    margin-bottom: 1.5rem;
  }
  
  .device-info p {
    margin: 0.5rem 0;
    color: #aaa;
    font-size: 0.9rem;
  }
  
  .device-info strong {
    color: #ffffff;
  }
  
  .device-actions {
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
  
  .add-device-section {
    background-color: #2d2d2d;
    border-radius: 8px;
    padding: 2rem;
    text-align: center;
    border: 2px dashed #5d5d5d;
  }
  
  .add-device-section h2 {
    margin: 0 0 1rem 0;
    color: #ffffff;
    font-size: 1.3rem;
    font-weight: 600;
  }
</style>