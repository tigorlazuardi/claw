<script lang="ts">
  import Sidebar from './lib/Sidebar.svelte'
  
  let currentPage = 'images'
  let currentComponent: any = null
  let isLoading = false
  
  // Page component cache to avoid re-importing
  const pageCache = new Map<string, any>()
  
  // Lazy load page components
  async function loadPage(pageName: string) {
    if (pageCache.has(pageName)) {
      return pageCache.get(pageName)
    }
    
    isLoading = true
    
    try {
      let module
      switch (pageName) {
        case 'home':
          module = await import('./pages/Home.svelte')
          break
        case 'images':
          module = await import('./pages/Images.svelte')
          break
        case 'devices':
          module = await import('./pages/Devices.svelte')
          break
        case 'sources':
          module = await import('./pages/Sources.svelte')
          break
        case 'jobs':
          module = await import('./pages/Jobs.svelte')
          break
        default:
          module = await import('./pages/Home.svelte')
      }
      
      const component = module.default
      pageCache.set(pageName, component)
      return component
    } catch (error) {
      console.error(`Failed to load page: ${pageName}`, error)
      // Fallback to Home page
      const fallback = await import('./pages/Home.svelte')
      return fallback.default
    } finally {
      isLoading = false
    }
  }
  
  // Update current component when page changes
  $: {
    if (currentPage) {
      loadPage(currentPage).then(component => {
        currentComponent = component
      })
    }
  }
  
  // Simple URL-based routing
  function updatePageFromURL() {
    const path = window.location.pathname.replace('/', '') || 'home'
    currentPage = path
  }
  
  function handleNavigation(event) {
    currentPage = event.detail.page
    // Update URL without page reload
    window.history.pushState(null, '', `/${currentPage === 'home' ? '' : currentPage}`)
  }
  
  // Listen for browser back/forward navigation
  function handlePopState() {
    updatePageFromURL()
  }
  
  // Initialize page from URL on load
  updatePageFromURL()
</script>

<svelte:window on:popstate={handlePopState} />

<main class="app">
  <Sidebar {currentPage} on:navigate={handleNavigation} />
  <div class="content">
    {#if isLoading}
      <div class="loading-container">
        <div class="loading-spinner">
          <div class="spinner"></div>
        </div>
        <p class="loading-text">Loading {currentPage}...</p>
      </div>
    {:else if currentComponent}
      <svelte:component this={currentComponent} />
    {:else}
      <div class="error-container">
        <h2>Page not found</h2>
        <p>The requested page could not be loaded.</p>
      </div>
    {/if}
  </div>
</main>

<style>
  .app {
    display: flex;
    height: 100vh;
    width: 100vw;
    background-color: #1a1a1a;
    color: #ffffff;
    margin: 0;
    padding: 0;
  }
  
  .content {
    flex: 1;
    overflow: auto;
  }
  
  .loading-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    padding: 2rem;
  }
  
  .loading-spinner {
    margin-bottom: 1rem;
  }
  
  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid #3d3d3d;
    border-top: 4px solid #646cff;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }
  
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
  
  .loading-text {
    color: #aaa;
    font-size: 1rem;
    margin: 0;
    text-transform: capitalize;
  }
  
  .error-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    padding: 2rem;
    text-align: center;
  }
  
  .error-container h2 {
    color: #ffffff;
    margin: 0 0 1rem 0;
    font-size: 1.5rem;
  }
  
  .error-container p {
    color: #aaa;
    margin: 0;
    font-size: 1rem;
  }
  
</style>
